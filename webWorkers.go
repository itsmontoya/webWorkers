package webWorkers

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/missionMeteora/toolkit/errors"
)

const (
	// ErrEmptyWorkers is returned when workerCap is set to 0 (or ignored)
	ErrEmptyWorkers = errors.Error("worker capacity must be greater than zero")

	// ErrEmptyQueue is returned when queueLen is set to 0 (or ignored)
	ErrEmptyQueue = errors.Error("queue length must be greater than zero")

	// ErrIsListening is returned when Listen() is called on an instance of webWorkers already listening
	ErrIsListening = errors.Error("cannot listen on an instance already listening")

	// ErrIsClosed is returned when an action is attempted on a closed instance
	ErrIsClosed = errors.Error("cannot perform action on closed instance")

	// ErrHeadersSent is returned when header modifications are attempted after the headers have already been sent
	ErrHeadersSent = errors.Error("headers already sent")

	// ErrEmptyCerts is returned when no certificates are available and the TLS setting is enabled
	ErrEmptyCerts = errors.Error("number of certificates must be greater than zero if TLS is enabled")

	// ErrEmptyAddress is returned when an empty address is provided
	ErrEmptyAddress = errors.Error("address cannot be empty")

	// ErrInvalidHeaderStatus is returned when an invalid header status is provided
	ErrInvalidHeaderStatus = errors.Error("invalid header status")

	// ErrInvalidStatusCode is returned when an invalid status code is provided
	ErrInvalidStatusCode = errors.Error("invalid status code")
)

const (
	stateReady int32 = iota
	stateListening
	stateClosed
)

// New returns a new instance of Webworkers
func New(o Opts, fn Handler) (ww *Webworkers, err error) {
	if err = o.validate(); err != nil {
		return
	}

	ww = &Webworkers{
		w: make(workers, o.WorkerCap),
		l: log.New(o.ErrorOutput, "webWorkers ("+o.Address+"): ", log.Ldate|log.Ltime),

		addr: o.Address,
	}

	if o.TLS {
		if err = ww.initTLS(o.Certs); err != nil {
			return
		}
	}

	for i := range ww.w {
		ww.w[i] = newWorker(o.QueueLen, &ww.wg, ww.l, fn)
	}

	return
}

// Webworkers is the manager the web workers service
type Webworkers struct {
	wg sync.WaitGroup

	w workers
	q queue
	l *log.Logger

	// TLS configuration
	tc *tls.Config
	// Listening address
	addr string

	// TODO: Decide if I want to pull this out into it's own "router" struct
	// Worker mux
	wm sync.Mutex
	// Next worker index
	nwi int

	// Closed state
	cs int32
}

func (ww *Webworkers) nextWorker() (w *worker) {
	// Double you double you double you.. this needs a visual refactor. So many w's
	ww.wm.Lock()
	w = ww.w[ww.nwi]
	ww.nwi++
	// Our next worker index equals the length of our workers slice, reset to zero
	if ww.nwi == len(ww.w) {
		ww.nwi = 0
	}
	ww.wm.Unlock()
}

// isListening will return whether or not an instance is listening
func (ww *Webworkers) isListening() bool {
	return atomic.LoadInt32(&ww.cs) == stateListening
}

// isClosed will return whether or not an instance is closed
func (ww *Webworkers) isClosed() bool {
	return atomic.LoadInt32(&ww.cs) == stateClosed
}

func (ww *Webworkers) initTLS(tps []TLSPair) (err error) {
	var crt tls.Certificate
	ww.tc = &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            x509.NewCertPool(),
		Certificates:       make([]tls.Certificate, 0, 2),
	}

	for _, tp := range tps {
		if crt, err = tls.LoadX509KeyPair(tp.CRT, tp.Key); err != nil {
			ww.l.Println(err)
			continue
		}

		ww.tc.Certificates = append(ww.tc.Certificates, crt)
	}

	if len(ww.tc.Certificates) == 0 {
		return ErrEmptyCerts
	}

	ww.tc.BuildNameToCertificate()
	return nil
}

func (ww *Webworkers) newListener() (lst net.Listener, err error) {
	if ww.tc == nil {
		return net.Listen("tcp", ww.addr)
	}

	return tls.Listen("tcp", ww.addr, ww.tc)
}

// Listen will begin the listening loop
func (ww *Webworkers) Listen() (err error) {
	if !atomic.CompareAndSwapInt32(&ww.cs, stateReady, stateListening) {
		return ErrIsListening
	}

	var lst net.Listener
	if lst, err = ww.newListener(); err != nil {
		return
	}

	for {
		var c net.Conn
		if c, err = lst.Accept(); err != nil {
			ww.l.Println(err)
			goto ITERATIONEND
		}

		ww.nextWorker().q <- c

	ITERATIONEND:
		if ww.isClosed() {
			break
		}
	}

	lst.Close()
	return
}

// Close will close an instance of web workers
func (ww *Webworkers) Close() (err error) {
	if atomic.SwapInt32(&ww.cs, stateClosed) == stateClosed {
		// Instance of webWorkers is already closed, return ErrIsClosed
		return ErrIsClosed
	}

	// Close queue channel
	close(ww.q)
	return
}

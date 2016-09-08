package webWorkers

import (
	"io"
	"net"
	"sync"
	"sync/atomic"

	"github.com/missionMeteora/toolkit/errors"
)

import "log"

const (
	// ErrEmptyWorkers is returned when workerCap is set to 0 (or ignored)
	ErrEmptyWorkers = errors.Error("worker capacity must be greater than zero")

	// ErrEmptyQueue is returned when queueLen is set to 0 (or ignored)
	ErrEmptyQueue = errors.Error("queue length must be greater than zero")

	// ErrIsClosed is returned when an action is attempted on a closed instance
	ErrIsClosed = errors.Error("cannot perform action on closed instance")
)

const (
	stateOpen int32 = iota
	stateClosed
)

// New returns a new instance of Webworkers
func New(o Opts, fn func(w io.Writer, r io.Reader)) (ww *Webworkers, err error) {
	if err = o.validate(); err != nil {
		return
	}

	ww = &Webworkers{
		w: make(workers, o.WorkerCap),
		q: make(queue, o.QueueLen),
	}

	for i := range ww.w {
		ww.w[i] = newWorker(ww.q, &ww.wg, fn)
	}

	return
}

// Webworkers is the manager the web workers service
type Webworkers struct {
	w workers
	q queue

	// Closed state
	cs int32
	wg sync.WaitGroup
}

// isClosed will return whether or not an instance is closed
func (ww *Webworkers) isClosed() bool {
	return atomic.LoadInt32(&ww.cs) == stateClosed
}

// Listen will begin the listening loop
func (ww *Webworkers) Listen() {
	var (
		lst net.Listener
		err error
	)

	if lst, err = net.Listen("tcp", ":8080"); err != nil {
		// handle err here
		return
	}

	for {
		var c net.Conn
		if c, err = lst.Accept(); err != nil {
			log.Println("Error?", err)
			goto ITERATIONEND
		}

		ww.q <- c

	ITERATIONEND:
		if ww.isClosed() {
			break
		}
	}

	lst.Close()
}

// Close will close an instance of web workers
func (ww *Webworkers) Close() (err error) {
	if !atomic.CompareAndSwapInt32(&ww.cs, stateOpen, stateClosed) {
		return ErrIsClosed
	}

	close(ww.q)

	return
}

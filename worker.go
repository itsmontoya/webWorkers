package webWorkers

import (
	"bytes"
	"io"
	"log"
	"sync"
	"time"
)

const (
	dateFmt = time.RFC1123
)

var (
	httpType = []byte("HTTP/1.1")
	server   = []byte("Server: PandaNet/0.0.1\n")
)

const (
	stateKey state = iota
	stateVal
)

const (
	// ContentTypeHTML is the html content type
	ContentTypeHTML = "text/html"
	// ContentTypeJS is the javascript content type
	ContentTypeJS = "application/js"
)

type state uint8

// newWorker returns a new worker
func newWorker(queueLen int, wg *sync.WaitGroup, l *log.Logger, fn Handler) (w *worker) {
	w = &worker{
		q:  make(queue, queueLen),
		wg: wg,
		l:  l,
		fn: fn,
	}

	wg.Add(1)
	go w.listen()
	return
}

// worker is an independent web worker, processes one request at a time
type worker struct {
	q  queue
	wg *sync.WaitGroup
	l  *log.Logger

	fn Handler
}

// listen will listen to an inbound queue to process net.Conn's
func (w *worker) listen() {
	var (
		req Request
		res Response

		buf [1024 * 8]byte
		n   int
		hn  int // Header length
		err error

		brdr = bytes.NewBuffer(nil)
	)

	req.Cookies = newCookies()
	res.Cookies = newCookies()

	for c := range w.q {
		if n, err = c.Read(buf[:]); err != nil && err != io.EOF {
			log.Println("Error?!", err)
			continue
		}

		if hn, err = req.processHeader(buf[:n]); err != nil {
			w.l.Println(err)
			goto ITEREND
		}

		brdr.Write(buf[hn:n])

		if req.contentLength > n-hn {
			req.Body = io.MultiReader(brdr, c)
		} else {
			req.Body = brdr
		}

		res.conn = c
		w.fn(&res, &req)

	ITEREND:
		c.Close()
		req.clean()
		res.clean()
		brdr.Reset()
	}

	w.wg.Done()
}

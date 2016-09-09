package webWorkers

import (
	"bytes"
	"io"
	"sync"
	"time"
)

import "log"

const (
	dateFmt = time.RFC1123
)

var (
	statusOK = []byte("HTTP/1.1 200 OK\n")
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
func newWorker(in queue, wg *sync.WaitGroup, fn Handler) (w *worker) {
	w = &worker{
		in: in,
		wg: wg,
		fn: fn,
	}

	wg.Add(1)
	go w.listen()
	return
}

// worker is an independent web worker, processes one request at a time
type worker struct {
	in queue
	wg *sync.WaitGroup

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

	for c := range w.in {
		if n, err = c.Read(buf[:]); err != nil {
			log.Println("Error?!", err)
			continue
		}

		hn = req.processHeader(buf[:n])
		brdr.Write(buf[hn:n])

		if req.contentLength > n-hn {
			req.Body = io.MultiReader(brdr, c)
		} else {
			req.Body = brdr
		}

		res.conn = c
		w.fn(&res, &req)
		c.Close()

		req.clean()
		res.clean()
		brdr.Reset()
	}

	w.wg.Done()
}

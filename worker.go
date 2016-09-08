package webWorkers

import (
	"bytes"
	"io"
	"sync"
	"time"

	"strconv"
)

import "log"
import "fmt"

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

type state uint8

// newWorker returns a new worker
func newWorker(in queue, wg *sync.WaitGroup, fn func(w io.Writer, r io.Reader)) (w *worker) {
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

	fn func(w io.Writer, r io.Reader)
}

// listen will listen to an inbound queue to process net.Conn's
func (w *worker) listen() {
	var (
		buf [1028]byte
		n   int
		hn  int // Header length
		err error
	)

	for c := range w.in {
		var (
			req request
			res response
		)

		if n, err = c.Read(buf[:]); err != nil {
			log.Println("Error?!", err)
			continue
		}

		hn = req.processHeader(buf[:n])

		fmt.Print(string(buf[hn:n]))
		var rdr io.Reader
		wtr := bytes.NewBuffer(nil)
		rdr = bytes.NewReader(buf[hn:n])
		if req.ContentLength > n-hn {
			rdr = io.MultiReader(rdr, c)
		}

		w.fn(wtr, rdr)

		res.StatusCode = 200
		res.ContentLength = wtr.Len()
		c.Write(res.bytes())
		io.Copy(c, wtr)
		c.Close()
	}

	w.wg.Done()
}

type request struct {
	Method         string
	Host           string
	Connection     string
	UserAgent      string
	Accept         string
	AcceptEncoding string
	AcceptLanguage string
	Cookie         *Cookie

	ContentLength int
	ContentType   string
}

func (r *request) processHeader(bs []byte) (n int) {
	var (
		s   state
		key []byte
		val []byte
	)

	fmt.Println("About to process header", string(bs))

	for i, b := range bs {
		switch s {
		case stateKey:
			if b == ':' {
				s = stateVal
			}

			if b == '\n' {
				n = i + 1
				break
			}

			key = append(key, b)

		case stateVal:
			if b == '\n' {
				switch string(key) {
				case "Host":
					r.Host = string(val)
				case "Connection":
					r.Connection = string(val)
				case "User-Agent":
					r.UserAgent = string(val)
				case "Accept":
					r.Accept = string(val)
				case "AcceptEncoding":
					r.AcceptEncoding = string(val)
				case "AcceptLanguage":
					r.AcceptLanguage = string(val)
				case "Cookie":
					r.Cookie = newCookie(4)
					fmt.Println("About to parse cookie", string(val))
					r.Cookie.parse(string(val))
				case "Content-Length":
					r.ContentLength, _ = strconv.Atoi(string(val))
				case "Content-Type":
					r.ContentType = string(val)
				}

				s = stateKey
				key = key[:0]
				val = val[:0]
			}

			val = append(val, b)
		}
	}

	return
}

type response struct {
	StatusCode    int
	ContentType   string
	ContentLength int
	Connection    string
	Date          string
	ETag          string
	Server        string
	LastModified  string
	Cookie        *Cookie

	Body []byte
}

func (r *response) bytes() (out []byte) {
	now := time.Now().Format(dateFmt)

	switch r.StatusCode {
	case 200:
		out = append(out, statusOK...)
	}

	out = append(out, server...)
	out = append(out, "Content-Type: "+r.ContentType+"\n"...)
	out = append(out, "Content-Length: "+strconv.Itoa(len(r.Body))+"\n"...)
	out = append(out, "Connection: close\n"...)
	out = append(out, "Date: "+now+"\n"...)
	out = append(out, "Last-Modified: "+now+"\n"...)
	out = append(r.Body)
	return
}

const (
	// ContentTypeHTML is the html content type
	ContentTypeHTML = "text/html"
	// ContentTypeJS is the javascript content type
	ContentTypeJS = "application/js"
)

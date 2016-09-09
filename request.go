package webWorkers

import (
	"io"
	"strconv"

	"bytes"
)

// Request is an HTTP request
type Request struct {
	host           []byte
	method         []byte
	path           []byte
	httpType       []byte
	connection     []byte
	userAgent      []byte
	accept         []byte
	acceptEncoding []byte
	acceptLanguage []byte
	contentLength  int
	contentType    []byte

	Body    io.Reader
	Cookies *Cookies
}

func (r *Request) clean() {
	r.host = r.host[:0]
	r.method = r.method[:0]
	r.path = r.path[:0]
	r.httpType = r.httpType[:0]
	r.connection = r.connection[:0]
	r.userAgent = r.userAgent[:0]
	r.accept = r.accept[:0]
	r.acceptEncoding = r.acceptEncoding[:0]
	r.acceptLanguage = r.acceptLanguage[:0]
	r.contentLength = 0
	r.contentType = r.contentType[:0]

	r.Body = nil

	r.Cookies.release()
	r.Cookies = nil
}

// Host will return the host
func (r *Request) Host() string {
	return string(r.host)
}

// Method will return the method
func (r *Request) Method() string {
	return string(r.method)
}

// Path will return the path
func (r *Request) Path() string {
	return string(r.path)
}

// HTTPType will return the http type
func (r *Request) HTTPType() string {
	return string(r.httpType)
}

// Connection will return the connection
func (r *Request) Connection() string {
	return string(r.connection)
}

// UserAgent will return the user agent
func (r *Request) UserAgent() string {
	return string(r.userAgent)
}

// Accept will return the accept
func (r *Request) Accept() string {
	return string(r.accept)
}

// AcceptEncoding will return the accept encoding
func (r *Request) AcceptEncoding() string {
	return string(r.acceptEncoding)
}

// AcceptLanguage will return the accept language
func (r *Request) AcceptLanguage() string {
	return string(r.acceptLanguage)
}

// ContentLength will return the content length
func (r *Request) ContentLength() int {
	return r.contentLength
}

// ContentType will return the content type
func (r *Request) ContentType() string {
	return string(r.contentType)
}

func (r *Request) processHeader(bs []byte) (n int) {
	var (
		s   state
		key []byte
		val []byte
	)

	r.Cookies = newCookies()

	for i, b := range bs {
		if b == '\n' {
			n = i + 1

			spl := bytes.Split(key, []byte{' '})
			if len(spl) < 3 {
				// TODO: Handle error
				return
			}

			r.method = append(r.method, spl[0]...)
			r.path = append(r.path, spl[1]...)
			r.httpType = append(r.httpType, spl[2]...)
			key = key[:0]
			break
		}

		key = append(key, b)
	}

	for i, b := range bs[n:] {
		switch s {
		case stateKey:
			if b == ':' {
				s = stateVal
				continue
			}

			if b == '\n' {
				n = i + 1
				break
			}

			key = append(key, b)

		case stateVal:
			if b == '\n' {
				val = trimPrefix(val)
				switch string(key) {
				case "Connection":
					r.connection = append(r.connection, val...)
				case "User-Agent":
					r.userAgent = append(r.userAgent, val...)
				case "Accept":
					r.accept = append(r.accept, val...)
				case "Accept-Encoding":
					r.acceptEncoding = append(r.acceptEncoding, val...)
				case "Accept-Language":
					r.acceptLanguage = append(r.acceptLanguage, val...)
				case "Content-Length":
					r.contentLength, _ = strconv.Atoi(string(val))
				case "Content-Type":
					r.contentType = append(r.contentType, val...)

				case "Cookie":
					r.Cookies.parse(string(val))
				}

				s = stateKey
				key = key[:0]
				val = val[:0]
				continue
			}

			val = append(val, b)
		}
	}

	return
}

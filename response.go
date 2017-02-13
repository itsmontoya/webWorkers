package webWorkers

import (
	"bytes"
	"net"
	"time"
)

var (
	contentTypeBytes  = []byte("Content-Type: ")
	connectionBytes   = []byte("Connection: close\n")
	dateBytes         = []byte("Date: ")
	lastModifiedBytes = []byte("Last-Modified: ")
	setCookieBytes    = []byte("Set-Cookie: ")
)

// Response is an http response
type Response struct {
	headersSent bool
	conn        net.Conn

	// Header buffer
	hbuf *bytes.Buffer

	statusCode    []byte
	contentType   []byte
	connection    []byte
	date          []byte
	server        []byte
	lastModified  []byte
	contentLength int

	Cookies *Cookies
}

func (r *Response) bytes() []byte {
	now := []byte(time.Now().Format(dateFmt))

	r.hbuf.Write(httpType)
	r.hbuf.WriteByte(' ')
	r.hbuf.Write(r.statusCode)
	r.hbuf.Write(server)

	// Write content type
	r.hbuf.Write(contentTypeBytes)
	r.hbuf.Write(r.contentType)
	r.hbuf.WriteByte('\n')

	r.hbuf.Write(connectionBytes)

	// Write date
	r.hbuf.Write(dateBytes)
	r.hbuf.Write(now)
	r.hbuf.WriteByte('\n')

	// Write last modified
	r.hbuf.Write(lastModifiedBytes)
	r.hbuf.Write(now)
	r.hbuf.WriteByte('\n')

	for _, ck := range r.Cookies.cks {
		// Write cookie
		r.hbuf.Write(setCookieBytes)
		r.hbuf.Write(ck.Bytes())
		r.hbuf.WriteByte('\n')
	}

	r.hbuf.WriteByte('\n')
	return r.hbuf.Bytes()
}

func (r *Response) clean() {
	r.headersSent = false
	r.conn = nil

	r.hbuf.Reset()

	r.statusCode = r.statusCode[:0]
	r.contentType = r.contentType[:0]
	r.connection = r.connection[:0]
	r.date = r.date[:0]
	r.server = r.server[:0]
	r.lastModified = r.lastModified[:0]
	r.contentLength = 0

	r.Cookies.clean()
}

func (r *Response) Write(b []byte) (err error) {
	if !r.headersSent {
		r.conn.Write(r.bytes())
		r.headersSent = true
	}

	_, err = r.conn.Write(b)
	return
}

// StatusCode will set the status code
func (r *Response) StatusCode(sc int) (err error) {
	if r.headersSent {
		return ErrHeadersSent
	}

	var b []byte
	// Get the byteslice representation of the status code text
	if b, err = getStatusBytes(sc); err != nil {
		return
	}

	// Clear current status code
	r.statusCode = r.statusCode[:0]
	// Set status code bytes
	r.statusCode = append(r.statusCode, b...)
	return
}

// ContentType will set the content type
func (r *Response) ContentType(ct string) (err error) {
	if r.headersSent {
		return ErrHeadersSent
	}

	if len(r.contentType) > 0 {
		r.contentType = r.contentType[:0]
	}

	r.contentType = append(r.contentType, ct...)
	return
}

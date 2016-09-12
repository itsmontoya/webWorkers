package webWorkers

import (
	"net"
	"time"
)

// Response is an http response
type Response struct {
	headersSent bool
	conn        net.Conn

	statusCode    []byte
	contentType   []byte
	connection    []byte
	date          []byte
	server        []byte
	lastModified  []byte
	contentLength int

	Cookies *Cookies
}

func (r *Response) bytes() (out []byte) {
	now := time.Now().Format(dateFmt)
	out = make([]byte, 0, 256)

	out = append(out, httpType...)
	out = append(out, ' ')
	out = append(out, r.statusCode...)
	out = append(out, server...)
	out = append(out, "Content-Type: "+string(r.contentType)+"\n"...)
	out = append(out, "Connection: close\n"...)
	out = append(out, "Date: "+now+"\n"...)
	out = append(out, "Last-Modified: "+now+"\n"...)

	for _, ck := range r.Cookies.cks {
		out = append(out, "Set-Cookie: "+ck.String()+"\n"...)

	}

	out = append(out, '\n')
	return
}

func (r *Response) clean() {
	r.headersSent = false
	r.conn = nil

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

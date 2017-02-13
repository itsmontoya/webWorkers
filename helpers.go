package webWorkers

import (
	"net"
	"unsafe"
)

// queue is a queue of net.Conn's
type queue chan net.Conn

// Handler is the func used for handling http requests
type Handler func(*Response, *Request)

// TLSPair is a crt/key pair for TLS
type TLSPair struct {
	CRT string
	Key string
}

// mapASCII will return an ASCII-friendly version of a provided byteslice
func mapASCII(in []byte) (out []byte) {
	if isASCII(in) {
		// Provided byteslice is valid ASCII, return early
		return in
	}

	// Set out with a max capacity as the length of our inbound byteslice
	out = make([]byte, 0, len(in))
	for _, b := range in {
		if !isASCIIByte(b) {
			// Byte is not valid ASCII, continue
			continue
		}

		// Byte is valid ASCII, append it to the outbound byteslice
		out = append(out, b)
	}

	return
}

// isASCII will return if a provided byteslice is valid ASCII
func isASCII(in []byte) bool {
	for _, b := range in {
		if !isASCIIByte(b) {
			return false
		}
	}

	return true
}

// isASCIIByte will return if a provided byte is valid ASCII
func isASCIIByte(b byte) bool {
	return b >= 32 && b <= 127
}

// trimPrefix will remove all spaces and newlines preceeding characters within a provided byteslice
func trimPrefix(bs []byte) []byte {
	for i, b := range bs {
		if b == '\n' || b == ' ' {
			// Byte is either a newline or a space, continue
			continue
		}

		// Byte is a non-whitespace character, return byteslice starting at current index
		return bs[i:]
	}

	return nil
}

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

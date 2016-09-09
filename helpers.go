package webWorkers

import (
	"net"

	"github.com/missionMeteora/toolkit/errors"
)

type queue chan net.Conn

// Opts are the options used to configure an instance of web workers
type Opts struct {
	WorkerCap int    `ini:"workerCap"`
	QueueLen  int    `ini:"queueLen"`
	Address   string `ini:"address"`
}

// validate will return any errors (if any) with the set of Opts
func (o *Opts) validate() (err error) {
	var errs errors.ErrorList
	if o.WorkerCap == 0 {
		errs.Append(ErrEmptyWorkers)
	}

	if o.QueueLen == 0 {
		errs.Append(ErrEmptyQueue)
	}

	return errs.Err()
}

func mapASCII(in []byte) (out []byte) {
	if isASCII(in) {
		return in
	}

	out = make([]byte, 0, len(in))

	for _, b := range in {
		if !isASCIIByte(b) {
			continue
		}

		out = append(out, b)
	}

	return
}

func isASCII(in []byte) bool {
	for _, b := range in {
		if !isASCIIByte(b) {
			return false
		}
	}

	return true
}

func isASCIIByte(b byte) bool {
	return b >= 32 && b <= 127
}

func trimPrefix(bs []byte) []byte {
	for i, b := range bs {
		if b == '\n' || b == ' ' {
			continue
		}

		return bs[i:]
	}

	return nil
}

// Handler is the func used for handling http requests
type Handler func(*Response, *Request)

package webWorkers

import (
	"strings"

	"github.com/go-ini/ini"
)

// NewOpts returns new options given a provided source
// Please see the go-ini/ini docu (https://godoc.org/github.com/go-ini/ini#Load) for more information on the source argument
func NewOpts(src interface{}) (o Opts, err error) {
	var srcF *ini.File
	// Attempt to load ini file from provided source
	if srcF, err = ini.Load(src); err != nil {
		return
	}

	// Map ini file data to Opts struct
	if err = srcF.MapTo(&o); err != nil {
		return
	}

	if !o.TLS {
		// This is a non-TLS configuration, return early
		return
	}

	// Load TLS pairs
	err = o.loadTLSPairs(srcF)
	return
}

// Opts are the options used to configure an instance of web workers
type Opts struct {
	// Number of workers (Note: will maintain a running goroutine for each worker)
	WorkerCap int `ini:"workerCap"`
	// Length of requests queue (IE: Requests sitting in memory, rather than waiting on disk using epoll/kqueue)
	QueueLen int `ini:"queueLen"`
	// Address to be serving from
	// TODO: Consider changing this to port, and making it a uint16
	Address string `ini:"address"`

	// Whether or not TLS is enabled
	TLS bool `ini:"tls"`
	// List of TLS certifications (Only needed if TLS is set to true)
	Certs []TLSPair
}

func (o *Opts) loadTLSPairs(srcF *ini.File) (err error) {
	var (
		ik *ini.Key
		tp TLSPair
	)

	for _, sec := range srcF.Sections() {
		// If section is not a "certification" block, move along
		if strings.Index(sec.Name(), "certification") != 0 {
			continue
		}

		// Retrieve the key matching "key"
		if ik, err = sec.GetKey("key"); err != nil {
			return
		}

		// Set the value for Key
		tp.Key = ik.Value()

		// Retrieve the key matching "crt"
		if ik, err = sec.GetKey("crt"); err != nil {
			return
		}

		// Set the value for CRT
		tp.CRT = ik.Value()
		// Append TLSPair to o.Certs
		o.Certs = append(o.Certs, tp)
	}

	if len(o.Certs) == 0 {
		// Our list of certifications is empty, return ErrEmptyCerts
		err = ErrEmptyCerts
	}

	return
}

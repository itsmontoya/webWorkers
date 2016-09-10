package webWorkers

import (
	"strings"

	"log"

	"github.com/go-ini/ini"
)

// NewOpts returns new options
func NewOpts(src interface{}) (o Opts, err error) {
	var srcF *ini.File
	if srcF, err = ini.Load(src); err != nil {
		return
	}

	if err = srcF.MapTo(&o); err != nil {
		return
	}

	if o.TLS {
		if err = loadTLSPairs(&o, srcF); err != nil {
			return
		}
	}

	return
}

// Opts are the options used to configure an instance of web workers
type Opts struct {
	WorkerCap int    `ini:"workerCap"`
	QueueLen  int    `ini:"queueLen"`
	Address   string `ini:"address"`

	TLS   bool `ini:"tls"`
	Certs []TLSPair
}

func loadTLSPairs(o *Opts, srcF *ini.File) (err error) {
	log.Println("Sections", srcF.Sections())
	for _, sec := range srcF.Sections() {
		var (
			ik *ini.Key
			tp TLSPair
		)

		log.Println("Load section!", sec.Name())

		if strings.Index(sec.Name(), "certification") != 0 {
			continue
		}

		if ik, err = sec.GetKey("key"); err != nil {
			return
		}

		tp.Key = ik.Value()

		if ik, err = sec.GetKey("crt"); err != nil {
			return
		}

		tp.CRT = ik.Value()
		o.Certs = append(o.Certs, tp)
	}

	if len(o.Certs) == 0 {
		err = ErrEmptyCerts
	}

	return
}

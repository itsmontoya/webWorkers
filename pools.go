package webWorkers

import (
	"bytes"
	"sync"
)

func newPools() *pools {
	var p pools

	p.bufs = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(nil)
		},
	}

	return &p
}

type pools struct {
	bufs sync.Pool
}

func (p *pools) acquireBuffer() (b *bytes.Buffer) {
	var ok bool
	if b, ok = p.bufs.Get().(*bytes.Buffer); !ok {
		panic("invalid pool type (buffer pool)")
	}

	return
}

func (p *pools) releaseBuffer(b *bytes.Buffer) {
	b.Reset()
	p.bufs.Put(b)
}

package webWorkers

import (
	"strings"
	"sync"
	"time"
)

const (
	cookiesDefaultLen = 4
)

var cksPool = sync.Pool{
	New: func() interface{} {
		return &Cookies{
			cks: make([]*Cookie, 0, cookiesDefaultLen),
		}
	},
}

// newCookies will return a new Cookies with the requested length
// mmm.. Cookiess.
func newCookies() (c *Cookies) {
	var ok bool
	if c, ok = cksPool.Get().(*Cookies); !ok {
		// This should never happen
		panic(ok)
	}

	return
}

// Cookies represents of a set of cookies
type Cookies struct {
	cks []*Cookie
}

// parse will parse an inbound string and populate the Cookies
func (c *Cookies) parse(in string) {
	for _, kv := range strings.Split(in, "; ") {
		var (
			k, v string
			spl  = strings.Split(kv, "=")
		)

		if len(spl) < 2 {
			continue
		}

		if k = spl[0]; k == "" {
			continue
		}

		if v = spl[1]; v == "" {
			continue
		}

		c.cks = append(c.cks, &Cookie{
			Key: k,
			Val: string(mapASCII([]byte(v))),
		})
	}
}

func (c *Cookies) release() {
	c.clean()
	cksPool.Put(c)
}

// clean deletes all the keys within the internal map, so the Cookies can be re-used
func (c *Cookies) clean() {
	c.cks = c.cks[:0]
}

// Get will return the value of the cookie matching the provided key
func (c *Cookies) Get(key string) (val string) {
	for _, ck := range c.cks {
		if ck.Key != key {
			continue
		}

		val = ck.Val
	}

	return
}

// Set will set the value of the cookie matching the provided key
func (c *Cookies) Set(key, val, path string, expires int64) {
	var match bool
	for _, ck := range c.cks {
		if ck.Key != key {
			continue
		}

		ck.Val = val
		ck.Path = path
		ck.Exp = expires
		match = true
	}

	if match {
		return
	}

	c.cks = append(c.cks, &Cookie{
		Key:  key,
		Val:  val,
		Path: path,
		Exp:  expires,
	})
}

// Dup will return a copy of the Cookies, so it can be used after the HTTP request/response process has completed
func (c *Cookies) Dup() (nc *Cookies) {
	nc = newCookies()

	for _, v := range c.cks {
		nc.cks = append(nc.cks, v)
	}

	return
}

type workers []*worker

// String will return a string representation of our Cookies
func (c *Cookies) String() string {
	return string(c.Bytes())
}

// Bytes will return a byteslice representation of our Cookies
func (c *Cookies) Bytes() (bs []byte) {
	first := true
	bs = make([]byte, 0, 64)

	for _, ck := range c.cks {
		if first {
			first = false
		} else {
			bs = append(bs, "; "...)
		}

		bs = append(bs, ck.Key...)
		bs = append(bs, '=')
		bs = append(bs, ck.Val...)
	}

	return
}

// Cookie represents a cookie
type Cookie struct {
	Key  string
	Val  string
	Path string
	Exp  int64
}

func (c *Cookie) String() string {
	bs := make([]byte, 0, 64)
	bs = append(bs, c.Key...)
	bs = append(bs, '=')
	bs = append(bs, c.Val...)

	if c.Path != "" {
		bs = append(bs, "; "...)
		bs = append(bs, "Path"...)
		bs = append(bs, '=')
		bs = append(bs, c.Path...)
	}

	if c.Exp > 0 {
		bs = append(bs, "; "...)
		bs = append(bs, "Expires"...)
		bs = append(bs, '=')
		bs = append(bs, time.Unix(c.Exp, 0).Format(time.RFC1123)...)
	}

	return string(bs)
}

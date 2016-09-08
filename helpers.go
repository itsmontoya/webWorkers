package webWorkers

import (
	"net"
	"strings"

	"github.com/missionMeteora/toolkit/errors"
)

import "fmt"

type queue chan net.Conn

// Opts are the options used to configure an instance of web workers
type Opts struct {
	WorkerCap int `ini:"workerCap"`
	QueueLen  int `ini:"queueLen"`
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

// newCookie will return a new Cookie with the requested length
// mmm.. Cookies.
func newCookie(length int) (c *Cookie) {
	return &Cookie{
		m: make(map[string]string, length),
	}
}

// Cookie represents a cookie
type Cookie struct {
	m map[string]string
}

// parse will parse an inbound string and populate the Cookie
func (c *Cookie) parse(in string) {
	for _, kv := range strings.Split(in, ";") {
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

		c.m[k] = v
	}

	fmt.Println("Cookie!", c.m)
}

// clear deletes all the keys within the internal map, so the Cookie can be re-used
func (c *Cookie) clear() {
	for k := range c.m {
		delete(c.m, k)
	}
}

// Get will return the value of the cookie matching the provided key
func (c *Cookie) Get(key string) string {
	return c.m[key]
}

// Set will set the value of the cookie matching the provided key
func (c *Cookie) Set(key, val string) {
	c.m[key] = val
}

// Dup will return a copy of the Cookie, so it can be used after the HTTP request/response process has completed
func (c *Cookie) Dup() (nc *Cookie) {
	nc = newCookie(len(c.m))

	for k, v := range c.m {
		nc.m[k] = v
	}

	return
}

type workers []*worker

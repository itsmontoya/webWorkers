package webWorkers

import (
	"bytes"
	"time"
)

// cookiesDefaultLen is the default length a cookies internal slice will be set at
const cookiesDefaultLen = 4

// newCookies will return a new Cookies with the requested length
// mmm.. Cookiess.
func newCookies() (c *Cookies) {
	return &Cookies{
		cks: make([]*Cookie, 0, cookiesDefaultLen),
	}
}

// Cookies represents of a set of cookies
type Cookies struct {
	// cb represents our cookies in the form of a byteslice
	cb []byte
	// parsed is whether or not the cookie bytes have been parsed yet
	parsed bool
	// cks is a list of Cookie's
	cks []*Cookie
}

// set will set the Cookies cb (cookie bytes) value
// Note: This copies the bytes within the provided argument and does not use it after that. The argument slice is safe to use later in external code
func (c *Cookies) set(in []byte) {
	if len(c.cb) > 0 {
		// Cookie bytes already exist, set the slice to empty before writing
		c.cb = c.cb[:0]
	}

	// Append the bytes within "in" to our c.cb
	c.cb = append(c.cb, in...)
}

// parse will parse an inbound string and populate the Cookies
func (c *Cookies) parse() {
	if len(c.cb) == 0 {
		// Cannot parse an empty set of cookie bytes
		return
	}

	// Our key and value
	var k, v []byte
	for _, kv := range bytes.Split(c.cb, []byte{';', ' '}) {
		// spl is our key/value split into key and value
		spl := bytes.Split(kv, []byte{'='})
		k = k[:0]
		v = v[:0]

		if len(spl) < 2 {
			// Our split has less than two values, continue on
			continue
		}

		if k = spl[0]; len(k) == 0 {
			// Our key is empty, continue on
			continue
		}

		if v = spl[1]; len(v) == 0 {
			// Our value is empty, continue on
			continue
		}

		// Append our parsed cookie to our list of cookies
		c.cks = append(c.cks, &Cookie{
			Key: string(k),
			Val: string(mapASCII(v)),
		})
	}

	// Set parsed to true
	c.parsed = true
}

// clean deletes all the keys within the internal map, so the Cookies can be re-used
func (c *Cookies) clean() {
	// Set parsed to false
	c.parsed = false
	// Reset cookie bytes slice
	c.cb = c.cb[:0]
	// Reset cookies list
	c.cks = c.cks[:0]
}

// Get will return the value of the cookie matching the provided key
func (c *Cookies) Get(key string) (val string) {
	if !c.parsed {
		// We haven't yet parsed, do so before attempting to get value
		c.parse()
	}

	for _, ck := range c.cks {
		if ck.Key != key {
			// Cookie key does not match provided key, continue
			continue
		}

		// Set val as Cookie value
		val = ck.Val
		// We have found our match, break out of loop
		break
	}

	return
}

// Set will set the value of the cookie matching the provided key
func (c *Cookies) Set(key, val, path string, expires int64) {
	var match bool
	// Iterate through list of cookies
	for _, ck := range c.cks {
		if ck.Key != key {
			// Cookie key does not match the provided key, continue
			continue
		}

		// Set cookie value, path, and expiration
		ck.Val = val
		ck.Path = path
		ck.Exp = expires
		// Set match to true
		match = true
	}

	if match {
		// We edited an existing entry, no need to create something new
		return
	}

	// Append new cookie to cookies list
	c.cks = append(c.cks, &Cookie{
		Key:  key,
		Val:  val,
		Path: path,
		Exp:  expires,
	})
}

// Dup will return a copy of the Cookies, so it can be used after the HTTP request/response process has completed
func (c *Cookies) Dup() (nc *Cookies) {
	// Create a fresh set of cookies
	nc = newCookies()

	// Iterate through cookies list
	for _, v := range c.cks {
		// Append each entry in c.cks to nc.cks
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

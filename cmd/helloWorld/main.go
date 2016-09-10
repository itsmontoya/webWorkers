package main

import (
	"fmt"

	"github.com/itsmontoya/webWorkers"
)

const (
	cfgLoc = "./config.ini"
)

func main() {
	var (
		ww  *webWorkers.Webworkers
		o   webWorkers.Opts
		err error
	)

	if o, err = webWorkers.NewOpts(cfgLoc); err != nil {
		panic(err)
	}

	if ww, err = webWorkers.New(o, Handle); err != nil {
		panic(err)
	}

	fmt.Println("About to listen")
	fmt.Println("Err?", ww.Listen())
}

// Handle likes to handle
func Handle(res *webWorkers.Response, req *webWorkers.Request) {
	res.StatusCode(200)
	res.ContentType("application/json")
	res.Write([]byte(`{ "greeting" : "Hello world!" }`))
}

package main

import (
	"fmt"
	"time"

	"github.com/itsmontoya/webWorkers"
)

func main() {
	var (
		ww  *webWorkers.Webworkers
		err error
	)

	if ww, err = webWorkers.New(webWorkers.Opts{
		WorkerCap: 1,
		QueueLen:  128,
	}, Handle); err != nil {
		panic(err)
	}

	fmt.Println("About to listen")
	ww.Listen()
}

// Handle likes to handle
func Handle(res *webWorkers.Response, req *webWorkers.Request) {
	res.StatusCode(200)
	res.ContentType("application/json")
	res.Write([]byte(`{ "greeting" : "Hello world!" }`))
}

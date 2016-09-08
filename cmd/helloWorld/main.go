package main

import (
	"fmt"
	"io"

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

	ww.Listen()
}

// Handle likes to handle
func Handle(w io.Writer, r io.Reader) {
	var buf [32]byte
	n, _ := r.Read(buf[:])
	fmt.Println(string(buf[:n]))
}

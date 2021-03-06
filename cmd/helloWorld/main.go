package main

import (
	"fmt"

	"github.com/itsmontoya/webWorkers"
)

// cfgLoc is the location of our config ini file
// Note: An example config is provided at ./config.ini.example
const cfgLoc = "./config.ini"

func main() {
	var (
		ww  *webWorkers.Webworkers
		o   webWorkers.Opts
		err error
	)

	// Get a set of webWorkers options using our configuration file location
	if o, err = webWorkers.NewOpts(cfgLoc); err != nil {
		panic(err)
	}

	// Request a new instance of webWorkers with our options
	if ww, err = webWorkers.New(o, Handle); err != nil {
		panic(err)
	}

	fmt.Println("About to listen at:", o.Address)
	fmt.Println("Error listening", ww.Listen())
}

// Handle likes to handle
func Handle(res *webWorkers.Response, req *webWorkers.Request) {
	// Set status code to 200 (Status OK)
	res.StatusCode(200)
	// Set our content type to application/json
	res.ContentType("application/json")
	// Return a static []byteslice of a JSON object
	res.Write([]byte(`{ "greeting" : "Hello world!" }`))
}

package main

import (
	"fmt"
	"net/http"

	"github.com/itsmontoya/webworkers"
	"github.com/missionMeteora/toolkit/closer"
	"github.com/pkg/profile"
	"github.com/valyala/fasthttp"
	"time"
)

// cfgLoc is the location of our config ini file
// Note: An example config is provided at ./config.ini.example
const cfgLoc = "./config.ini"

func main() {
	p := profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	cl := closer.New()

	go initWW()
	go initStdlib()
	go initFastHTTP()

	cl.Wait()
	fmt.Println("Closing")
	p.Stop()

}

func initWW() {
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
	time.Sleep(time.Millisecond * 2)

	// Return a static []byteslice of a JSON object
	res.Write([]byte(`{ "greeting" : "Hello world!" }`))
}

func initStdlib() {
	//	srv := &stdlibSrv{}
	//	http.ListenAndServe(":8081", srv)

	s := &stdlibSrv{}
	srv := &http.Server{Addr: ":8081", Handler: s}
	srv.SetKeepAlivesEnabled(false)

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

type stdlibSrv struct {
}

func (s *stdlibSrv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-type", "application/json")
	time.Sleep(time.Millisecond * 2)
	w.Write([]byte(`{ "greeting" : "Hello world!" }`))
}

func initFastHTTP() {
	srv := fasthttp.Server{
		DisableKeepalive: true,
		Handler:          HandleFastHTTP,
	}

	srv.ListenAndServe(":8082")
	//	fasthttp.ListenAndServe(":8082", HandleFastHTTP)
}

func HandleFastHTTP(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(200)
	ctx.SetContentType("application/json")
	time.Sleep(time.Millisecond * 2)
	ctx.Write([]byte(`{ "greeting" : "Hello world!" }`))
}

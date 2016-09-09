package webWorkers

import (
	//	"fmt"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	statusCode      = 200
	jsonContentType = "application/json"
	jsonStr         = `{ "greeting" : "Hello world!" }`
)

var (
	jsonB = []byte(jsonStr)

	errInvalidResponse = errors.New("invalid response")
)

func TestMain(m *testing.M) {
	s := &srv{}
	go initWW(s)
	go initStdLib(s)

	time.Sleep(time.Second)
	sc := m.Run()

	os.Exit(sc)
}

func TestBasic(t *testing.T) {
	if err := httpReq("http://localhost:11110"); err != nil {
		t.Error(err)
	}
}

func BenchmarkWWBasic(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		if err = httpReq("http://localhost:11110"); err != nil {
			b.Error(err)
		}
	}

	b.ReportAllocs()
}

func BenchmarkStdlibBasic(b *testing.B) {
	var err error
	for i := 0; i < b.N; i++ {
		if err = httpReq("http://localhost:11111"); err != nil {
			b.Error(err)
		}
	}

	b.ReportAllocs()
}

func BenchmarkParaWWBasic(b *testing.B) {
	b.SetParallelism(4)
	b.RunParallel(func(pb *testing.PB) {
		var err error
		for pb.Next() {
			if err = httpReq("http://localhost:11110"); err != nil {
				b.Error(err)
			}
		}
	})

	b.ReportAllocs()
}

func BenchmarkParaStdlibBasic(b *testing.B) {
	b.SetParallelism(4)
	b.RunParallel(func(pb *testing.PB) {
		var err error
		for pb.Next() {
			if err = httpReq("http://localhost:11111"); err != nil {
				b.Error(err)
			}
		}
	})

	b.ReportAllocs()
}

func httpReq(loc string) (err error) {
	var (
		resp *http.Response
		buf  [len(jsonStr)]byte
	)

	if resp, err = http.Get(loc); err != nil {
		return
	}

	if resp.StatusCode != 200 {
		goto END
	}

	if _, err = resp.Body.Read(buf[:]); err != nil {
		if err != io.EOF {
			goto END
		}

		err = nil
	}

	if str := string(buf[:]); str != jsonStr {
		err = errInvalidResponse
		goto END
	}

END:
	resp.Body.Close()
	return
}

type srv struct{}

func (s *srv) wwHandler(res *Response, req *Request) {
	res.StatusCode(statusCode)
	res.ContentType(jsonContentType)
	res.Write(jsonB)
}

func (s *srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", jsonContentType)
	w.Write(jsonB)
}

func initWW(s *srv) {
	var (
		ww  *Webworkers
		err error
	)

	opts := Opts{
		WorkerCap: 8,
		QueueLen:  1024,
		Address:   ":11110",
	}

	if ww, err = New(opts, s.wwHandler); err != nil {
		panic(err)
	}

	ww.Listen()
}

func initStdLib(s *srv) {
	var err error
	if err = http.ListenAndServe(":11111", s); err != nil {
		panic(err)
	}
}

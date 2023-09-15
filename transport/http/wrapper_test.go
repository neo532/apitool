package http

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/neo532/apitool/transport"
	"github.com/neo532/apitool/transport/http/xhttp"

	"github.com/neo532/apitool/transport/http/xhttp/header"
	"github.com/neo532/apitool/transport/http/xhttp/queryargs"
)

type Logger struct {
}

func (l *Logger) Info(c context.Context, msg string) {
	fmt.Println(fmt.Sprintf("%+v", msg))
}
func (l *Logger) Error(c context.Context, msg string) {
	fmt.Println(fmt.Sprintf("%+v", msg))
}

func clt() (clt xhttp.Client) {
	env, _ := transport.String2Env("dev")
	clt = xhttp.New(
		xhttp.WithLogger(&Logger{}),
		xhttp.WithEnv(env),
		xhttp.WithMap("aa", "aaa"),
		xhttp.WithMap("bb", "bbb"),
	)
	return
}

func TestWrapT(t *testing.T) {
	count := 1
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			TWrap()
		}()
	}
	wg.Wait()
}

func TWrap() {

	//header["X-Auth-AppId"] = []string{"application/json"}
	c := context.Background()
	c = header.AppendToContext(
		c,
		"Content-Type", "application/jsona",
		"User-Agent", "ssssssss",
	)
	c = queryargs.AppendToContext(c, "{user}", "b")

	url := "http://127.0.0.1:8500/user/{user}"
	method := "PUT"

	opts := make([]xhttp.Opt, 0, 6)
	opts = append(opts, xhttp.WithUrl(url))
	opts = append(opts, xhttp.WithMethod(method))
	opts = append(opts, xhttp.WithHeader(c))
	type P struct {
		A    string `json:"a" form:"a"`
		B    int    `json:"b" form:"b"`
		Name string `json:"name" form:"name"`
	}

	var err error
	c, err = xhttp.AppendUrlByStruct(c, &P{"aaaa", 1222, "ccccc"})
	if err != nil {
		panic(err)
	}

	if method == "POST" {
		b, _ := xhttp.ToJsonBody(&P{"aaaaa", 111, "ccc"})
		opts = append(opts, xhttp.WithBody(b))
	}

	client := NewWrapper(clt()).Call(
		c,
		url,
		opts...,
	)

	_, err = client.Body(c)
	if err != nil {
		fmt.Println(fmt.Sprintf("err:\t%+v", err))
	}

	//fmt.Println(string(body))
	fmt.Println(runtime.Caller(0))
}

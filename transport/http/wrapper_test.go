package http

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/neo532/apitool/transport"
	"github.com/neo532/apitool/transport/http/xhttp"
	"github.com/neo532/apitool/transport/http/xhttp/client"
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

func clt() (clt client.Client) {
	env, _ := transport.String2Env("dev")
	clt = client.New(
		client.WithLogger(&Logger{}),
		client.WithEnv(env),
		client.WithMap("aa", "aaa"),
		client.WithMap("bb", "bbb"),
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
		"Content-Type", "application/json",
		"User-Agent", "ssssssss",
		"traceId", "aaaa",
	)
	c = queryargs.AppendToContext(c, "{user}", "b")

	url := "http://127.0.0.1:8500/user/{user}"
	method := "POST"

	opts := make([]xhttp.Opt, 0, 6)
	opts = append(opts, xhttp.WithUrl(url))
	opts = append(opts, xhttp.WithMethod(method))
	type P struct {
		A       string `json:"a" form:"a"`
		B       int    `json:"b" form:"b"`
		Name    string `json:"name" form:"name"`
		Message string `json:"message" form:"a"`
		Code    int    `json:"code" form:"b"`
		RankId  int    `json:"rankId" form:"b"`
	}

	var err error
	c, err = xhttp.AppendUrlByStruct(c, &P{A: "aaaa", B: 1222, Name: "ccccc"})
	if err != nil {
		panic(err)
	}

	if method == "POST" {
		//b, _ := xhttp.ToJsonBody(&P{"aaaaa", 111, "ccc"})
	}

	reply := &P{}
	err = NewWrapper(clt()).Call(
		c,
		&P{A: "aaaa", B: 1222, Name: "ccccc", RankId: 222},
		reply,
		opts...,
	)
	if err != nil {
		fmt.Println(fmt.Sprintf("err:\t%+v", err))
	}

	fmt.Println(reply)
	fmt.Println(runtime.Caller(0))
}

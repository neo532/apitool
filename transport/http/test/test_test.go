package test

// In case of that importing huge amount of packages,I annotation this file.

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"testing"

// 	"github.com/neo532/apitool/transport"
// 	"github.com/neo532/apitool/transport/http/proto"
// 	"github.com/neo532/apitool/transport/http/xhttp/client"
// 	"github.com/neo532/apitool/transport/http/xhttp/middleware"
// )

// func clt() (clt client.Client, err error) {
// 	var env transport.Env
// 	env, err = transport.String2Env("dev")
// 	clt = client.New(
// 		client.WithLogger(&transport.LoggerDefault{}),
// 		client.WithEnv(env),
// 		client.WithMaxConnsPerHost(1),
// 	)
// 	return
// }

// func TestXClient(t *testing.T) {
// 	clt, err := clt()
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	demo1 := proto.NewProtoXHttpClient(clt)
// 	demo1.WithMiddleware(mdw1("demo1"))

// 	demo2 := proto.NewProtoXHttpClient(clt)
// 	demo2.WithMiddleware(mdw2("demo2"))

// 	c := context.Background()
// 	count := 1
// 	var wg sync.WaitGroup
// 	wg.Add(count)
// 	for i := 1; i <= count; i++ {
// 		go func(i int) {
// 			defer wg.Done()
// 			demo1.Get(c, &proto.GetRequest{})
// 			demo2.GetById(c, &proto.GetByIdRequest{Id: uint64(i)})
// 			demo1.Get(c, &proto.GetRequest{})
// 		}(i)
// 	}
// 	wg.Wait()
// }

// func mdw1(flag string) middleware.Middleware {
// 	return func(handler middleware.Handler) middleware.Handler {
// 		return func(c context.Context, req, reply interface{}) (err error) {
// 			fmt.Println(flag, "mdw1")
// 			return handler(c, req, reply)
// 		}
// 	}
// }
// func mdw2(flag string) middleware.Middleware {
// 	return func(handler middleware.Handler) middleware.Handler {
// 		return func(c context.Context, req, reply interface{}) (err error) {
// 			fmt.Println(flag, "mdw2")
// 			return handler(c, req, reply)
// 		}
// 	}
// }

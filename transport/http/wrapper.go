package http

/*
 * @abstract 传输协议http的客户端包装
 * @mail neo532@126.com
 * @date 2023-09-12
 */

import (
	"context"

	"github.com/neo532/apitool/transport/http/xhttp"
	"github.com/neo532/apitool/transport/http/xhttp/client"
)

type Wrapper struct {
	clt client.Client
}

func NewWrapper(clt client.Client) *Wrapper {
	return &Wrapper{
		clt: clt,
	}
}

func (w *Wrapper) Call(
	c context.Context,
	req interface{},
	reply interface{},
	opts ...xhttp.Opt,
) (err error) {
	err = xhttp.New(w.clt, opts...).Do(c, req, reply)
	return
}

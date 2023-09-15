package http

/*
 * @abstract 传输协议http的客户端包装
 * @mail neo532@126.com
 * @date 2023-09-12
 */

import (
	"context"

	"github.com/neo532/apitool/transport/http/xhttp"
)

type Wrapper struct {
	clt xhttp.Client
}

func NewWrapper(clt xhttp.Client) *Wrapper {
	return &Wrapper{
		clt: clt,
	}
}

func (w *Wrapper) Call(
	c context.Context,
	domain string,
	opts ...xhttp.Opt,
) (clt *xhttp.Client) {

	clt = w.clt.Do(c, opts...)
	return
}

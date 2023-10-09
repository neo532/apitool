package http

/*
 * @abstract 传输协议http的客户端包装
 * @mail neo532@126.com
 * @date 2023-09-12
 */
import (
	"github.com/neo532/apitool/transport/http/xhttp/client"
	"github.com/neo532/apitool/transport/http/xhttp/middleware"
)

type XClient struct {
	Domain string
	Client client.Client
}

func NewXClient(clt client.Client) (xc *XClient) {
	return &XClient{
		Client: clt,
	}
}

func (xc *XClient) WithDomain(domain string) {
	xc.Domain = domain
	return
}

func (xc *XClient) WithMiddleware(mds ...middleware.Middleware) {
	for _, mw := range mds {
		xc.Client = xc.Client.AddMiddleware(mw)
	}
	return
}

package http

/*
 * @abstract 传输协议http的客户端包装
 * @mail neo532@126.com
 * @date 2023-09-12
 */
import (
	"github.com/neo532/apitool/transport/http/xhttp"
	"github.com/neo532/apitool/transport/http/xhttp/client"
	"github.com/neo532/apitool/transport/http/xhttp/middleware"
)

type XClient struct {
	Client client.Client

	Domain          string
	RequestEncoder  xhttp.EncodeRequestFunc
	ResponseDecoder xhttp.DecodeResponseFunc
	ErrorDecoder    xhttp.DecodeErrorFunc
}

func NewXClient(clt client.Client) (xc *XClient) {
	return &XClient{
		Client: clt.CopyMiddleware(),
	}
}

func (s *XClient) WithMiddleware(mds ...middleware.Middleware) {
	for _, mw := range mds {
		s.Client = s.Client.AddMiddleware(mw)
	}
	return
}

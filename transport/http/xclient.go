package http

/*
 * @abstract 传输协议http的客户端包装
 * @mail neo532@126.com
 * @date 2023-09-12
 */
import (
	"github.com/neo532/apitool/encoding"
	"github.com/neo532/apitool/encoding/json"
	"github.com/neo532/apitool/encoding/xml"
)

type XClient struct {
	Domain string
	Codecs map[string]encoding.Codec
}

func NewXClient() (xc *XClient) {
	return &XClient{
		Codecs: map[string]encoding.Codec{
			"json": json.NewCodec(),
			"xml":  xml.NewCodec(),
		},
	}
}

func (xc *XClient) WithDomain(domain string) *XClient {
	xc.Domain = domain
	return xc
}

func (xc *XClient) Codec(encoding string) encoding.Codec {
	return xc.Codecs[encoding]
}

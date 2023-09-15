package xml

import (
	"encoding/xml"
)

type opt func(cc *Codec)

// Codec is a Codec implementation with xml.
type Codec struct {
	name string
}

func NewCodec(opts ...opt) (cc *Codec) {
	cc = &Codec{
		name: Name,
	}
	for _, o := range opts {
		o(cc)
	}
	return
}

func (cc *Codec) Marshal(v interface{}) ([]byte, error) {
	return xml.Marshal(v)
}

func (cc *Codec) Unmarshal(data []byte, v interface{}) error {
	return xml.Unmarshal(data, v)
}

func (cc *Codec) Name() string {
	return cc.name
}

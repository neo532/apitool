package json

import (
	"encoding/json"
	"reflect"

	gtproto "github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/neo532/apitool/encoding"
)

func init() {
	encoding.RegisterCodec(NewCodec())
}

type opt func(cc *Codec)

func WithMarshalEmitUnpopulated(b bool) opt {
	return func(cc *Codec) {
		cc.marshalOptions.EmitUnpopulated = b
	}
}

// Codec is a Codec implementation with json.
type Codec struct {
	marshalOptions   protojson.MarshalOptions
	unmarshalOptions protojson.UnmarshalOptions
}

func NewCodec(opts ...opt) (cc *Codec) {
	cc = &Codec{
		marshalOptions:   MarshalOptions,
		unmarshalOptions: UnmarshalOptions,
	}
	for _, o := range opts {
		o(cc)
	}
	return
}

func (cc *Codec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case json.Marshaler:
		return m.MarshalJSON()
	case proto.Message:
		return cc.marshalOptions.Marshal(m)
	case gtproto.Message:
		mv := gtproto.MessageV2(m)
		return cc.marshalOptions.Marshal(mv)
	default:
		return json.Marshal(m)
	}
}

func (cc *Codec) Unmarshal(data []byte, v interface{}) error {
	switch m := v.(type) {
	case json.Unmarshaler:
		return m.UnmarshalJSON(data)
	case proto.Message:
		return cc.unmarshalOptions.Unmarshal(data, m)
	case gtproto.Message:
		mv := gtproto.MessageV2(m)
		return cc.unmarshalOptions.Unmarshal(data, mv)
	default:

		rv := reflect.ValueOf(v)
		for rv := rv; rv.Kind() == reflect.Ptr; {
			if rv.IsNil() {
				rv.Set(reflect.New(rv.Type().Elem()))
			}
			rv = rv.Elem()
		}
		if m, ok := reflect.Indirect(rv).Interface().(proto.Message); ok {
			return cc.unmarshalOptions.Unmarshal(data, m)
		}
		return json.Unmarshal(data, m)
	}
}

func (cc *Codec) Name() string {
	return "json"
}

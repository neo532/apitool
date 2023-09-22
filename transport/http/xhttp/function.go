package xhttp

/*
 * @abstract 传输协议http的客户端的一些通用方法
 * @mail neo532@126.com
 * @date 2022-05-30
 */

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/neo532/apitool/transport/http/xhttp/queryargs"
)

const (
	FORM_AND       = "&"
	FORM_ASSIGN    = "="
	FORM_KEY_SLICE = "%5B%5D"
)

var (
	TagName = "form"
	// ErrNotSupportType is a type of error that means invaild type.
	ErrNotSupportType error = errors.New("Invaild type,within string,int,int64,uint64,float64,[]string,[]int,[]int64,[]uint64,[]float64!")
	// ErrMustBeStruct is a type of error that means the type must be a struct.
	ErrMustBeStruct error = errors.New("QueryArgs must be a struct or a struct pointer!")
)

func ToJsonBody(v interface{}) (b []byte, err error) {
	b, err = json.Marshal(v)
	return
}

func ToQueryArgs(v interface{}) (s string, err error) {
	return
}

func AppendUrlByStruct(c context.Context, param interface{}) (ctx context.Context, err error) {
	// unify type
	T := reflect.TypeOf(param)
	V := reflect.ValueOf(param)
	switch {
	case T.Kind() == reflect.Struct:
	case T.Kind() == reflect.Ptr && T.Elem().Kind() == reflect.Struct:
		T = T.Elem()
		V = V.Elem()
	default:
		err = ErrMustBeStruct
		return
	}
	qa := queryargs.New()

	for i := 0; i < T.NumField(); i++ {
		field := T.Field(i)
		if field.PkgPath != "" && !field.Anonymous { // unexported
			continue
		}
		value := V.Field(i)

		name := field.Tag.Get(TagName)
		// don't parse that the name is -.
		if name == "-" {
			continue
		}
		nameSlice := name + FORM_KEY_SLICE

		// check whether if empty,in case of escape to heap,use strings.
		emptyIndex := strings.Index(name, ",omitempty")
		if emptyIndex != -1 {
			if value.IsZero() {
				continue
			}
			name = name[0:emptyIndex]
		}

		// identify type
		switch value.Kind() {
		case reflect.String:
			qa = qa.Add(name, url.QueryEscape(value.String()))
		case reflect.Int, reflect.Int64, reflect.Int32:
			qa = qa.Add(name, strconv.FormatInt(value.Int(), 10))
		case reflect.Uint64, reflect.Uint32:
			qa = qa.Add(name, strconv.FormatUint(value.Uint(), 10))
		case reflect.Float64:
			qa = qa.Add(name, strconv.FormatFloat(value.Float(), 'f', -1, 64))
		case reflect.Bool:
			qa = qa.Add(name, strconv.FormatBool(value.Bool()))
		case reflect.Slice, reflect.Array:
			o := value
			for i, lenS := 0, o.Len(); i < lenS; i++ {
				v := o.Index(i)
				switch v.Kind() {
				case reflect.String:
					qa = qa.Add(nameSlice, url.QueryEscape(v.String()))
				case reflect.Int, reflect.Int64, reflect.Int32:
					qa = qa.Add(nameSlice, strconv.FormatInt(v.Int(), 10))
				case reflect.Uint64, reflect.Uint32:
					qa = qa.Add(nameSlice, strconv.FormatUint(v.Uint(), 10))
				case reflect.Float64:
					qa = qa.Add(nameSlice, strconv.FormatFloat(v.Float(), 'f', -1, 64))
				case reflect.Bool:
					qa = qa.Add(nameSlice, strconv.FormatBool(v.Bool()))
				default:
					err = ErrNotSupportType
					return
				}
			}
		default:
			err = ErrNotSupportType
			return
		}
	}
	ctx = queryargs.MergeToContext(c, qa)
	return
}

func AppendUrlByKV(url, k, v string) (s string) {
	switch {
	case k == "":
	case string(k[0]) == "{":
		url = strings.Replace(url, k, v, -1)
	case strings.Index(url, "?") == -1:
		url += "?" + k + FORM_ASSIGN + v
	default:
		url += FORM_AND + k + FORM_ASSIGN + v
	}
	return url
}

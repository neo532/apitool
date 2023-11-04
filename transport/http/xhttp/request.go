package xhttp

/*
 * @abstract 传输协议http的客户端的操作方法
 * @mail neo532@126.com
 * @date 2022-05-30
 */

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/neo532/apitool/encoding"
	"github.com/neo532/apitool/transport"
	"github.com/neo532/apitool/transport/http/xhttp/client"
	"github.com/neo532/apitool/transport/http/xhttp/header"
	"github.com/neo532/apitool/transport/http/xhttp/middleware"
	"github.com/neo532/apitool/transport/http/xhttp/queryargs"
)

type Request struct {
	clt client.Client

	url                 string
	method              string
	contentType         string
	contentTypeResponse string

	retryTimes       int
	retryMaxDuration time.Duration
	retryDuration    time.Duration
	timeLimit        time.Duration

	encoder      EncodeRequestFunc
	decoder      DecodeResponseFunc
	errorDecoder DecodeErrorFunc
}

// ========== Opt ==========
type Opt func(o *Request)

func WithTimeLimit(d time.Duration) Opt {
	return func(o *Request) {
		o.timeLimit = d
	}
}
func WithUrl(s string) Opt {
	return func(o *Request) {
		o.url = s
	}
}
func WithMethod(m string) Opt {
	return func(o *Request) {
		o.method = m
	}
}
func WithContentType(ct string) Opt {
	return func(o *Request) {
		o.contentType = ct
	}
}
func WithContentTypeResponse(ct string) Opt {
	return func(o *Request) {
		o.contentTypeResponse = ct
	}
}
func WithRetryTimes(times int) Opt {
	return func(o *Request) {
		o.retryTimes = times
	}
}
func WithRetryDuration(d time.Duration) Opt {
	return func(o *Request) {
		o.retryDuration = d
	}
}
func WithRetryMaxDuration(d time.Duration) Opt {
	return func(o *Request) {
		o.retryMaxDuration = d
	}
}
func WithRequestEncoder(encoder EncodeRequestFunc) Opt {
	return func(o *Request) {
		o.encoder = encoder
	}
}
func WithResponseDecoder(decoder DecodeResponseFunc) Opt {
	return func(o *Request) {
		o.decoder = decoder
	}
}
func WithErrorDecoder(errorDecoder DecodeErrorFunc) Opt {
	return func(o *Request) {
		o.errorDecoder = errorDecoder
	}
}

// ========== /Opt ==========

func New(clt client.Client, opts ...Opt) (req *Request) {
	req = &Request{
		retryTimes:       clt.RetryTime(),
		retryDuration:    time.Microsecond,
		retryMaxDuration: 20 * time.Microsecond,
		timeLimit:        3 * time.Second,

		errorDecoder: DefaultErrorDecoder,
		encoder:      DefaultRequestEncoder,
		decoder:      DefaultResponseDecoder,
		clt:          clt,
	}
	for _, o := range opts {
		o(req)
	}
	return
}

func (r *Request) Do(c context.Context, req interface{}, reply interface{}) (err error) {

	h := func(c context.Context, req interface{}, reply interface{}) (err error) {

		url := r.FmtQueryArgs(c, r.url)

		var reqBody []byte
		if reqBody, err = r.encoder(c, r.contentType, req); err != nil {
			return
		}

		reqHeader, headerBCurl := r.FmtHeader(c)

		client := &http.Client{Timeout: r.timeLimit}

		retryDuration := r.retryDuration
		var er error
		for i := 0; i <= r.retryTimes; i++ {

			var param *http.Request
			if param, err = http.NewRequest(r.method, url, bytes.NewReader(reqBody)); err != nil {
				return
			}
			param.Header = reqHeader

			// request
			var resp *http.Response
			start := time.Now()
			resp, err = client.Do(param)
			cost := time.Now().Sub(start)

			var respCode int
			var respBody []byte
			var cancelRetry bool
			for j := 0; j < 1; j++ {
				if err != nil {
					break
				}
				if resp != nil {
					respCode = resp.StatusCode
				}
				if cancelRetry, err = r.errorDecoder(c, resp); err != nil {
					break
				}
				if resp != nil {
					if r.contentTypeResponse != "" {
						resp.Header.Set(ContentTypeHeaderKey, r.contentTypeResponse)
					}
					respBody, er = r.decoder(c, resp, reply)
				}
			}

			r.log(c, url, headerBCurl, reqBody, respCode, respBody, cost, err)
			c = header.AppendToContext(c,
				"url", r.url,
				"cost", cost.String(),
				"limit", r.timeLimit.String(),
				"httpcode", strconv.Itoa(respCode),
			)

			if cancelRetry || err == nil {
				break
			}

			time.Sleep(r.retryDuration)
			if retryDuration < r.retryMaxDuration {
				retryDuration = retryDuration + retryDuration
			}
		}
		if err == nil {
			err = er
		}
		return
	}

	if len(r.clt.Middlewares()) > 0 {
		h = middleware.Chain(r.clt.Middlewares()...)(h)
	}
	return h(c, req, reply)
}

func (r *Request) log(c context.Context,
	url string, header strings.Builder, reqBody []byte,
	respCode int, respBody []byte, cost time.Duration, err error) {
	respStr := string(respBody)
	if r.clt.Env() == transport.EnvProd && utf8.RuneCountInString(respStr) > r.clt.ResponseMaxLength() {
		respStr = string([]rune(respStr)[:r.clt.ResponseMaxLength()]) + "..."
	}

	reqBodyS := string(reqBody)
	if len(reqBodyS) != 0 && reqBodyS != "{}" {
		reqBodyS = " -d '" + string(reqBodyS) + "'"
	} else {
		reqBodyS = ""
	}
	msg := fmt.Sprintf("[code:%d] [limit:%s] [cost:%s] [curl -X '%s' '%s'%s%s] [rst:%s]",
		respCode,
		r.timeLimit.String(),
		cost.String(),
		r.method,
		url,
		header.String(),
		reqBodyS,
		respStr,
	)
	if err != nil {
		r.clt.Logger().Error(c, fmt.Sprintf("[err:%s] %s", err.Error(), msg))
		return
	}
	r.clt.Logger().Info(c, msg)
}

func (r *Request) FmtQueryArgs(c context.Context, url string) string {
	if qa, ok := queryargs.FromContext(c); ok {
		qa.Range(func(k, v string) (b bool) {
			url = AppendUrlByKV(url, k, v)
			return true
		})
	}
	return url
}

func (r *Request) FmtHeader(c context.Context) (h http.Header, curl strings.Builder) {
	h = http.Header{}
	if md, ok := header.FromContext(c); ok {
		md.Range(func(k string, vs []string) (b bool) {
			for _, v := range vs {
				h.Set(k, v)
				curl.WriteString(" -H '" + k + ":" + v + "'")
			}
			return true
		})
	}
	if r.contentType != "" {
		h.Set(ContentTypeHeaderKey, r.contentType)
		curl.WriteString(" -H '" + ContentTypeHeaderKey + ":" + r.contentType + "'")
		return
	}
	if h.Get(ContentTypeHeaderKey) == "" && HasBody(r.method) {
		h.Set(ContentTypeHeaderKey, ContentTypeHeaderDefaultValue)
		curl.WriteString(" -H '" + ContentTypeHeaderKey + ":" + ContentTypeHeaderDefaultValue + "'")
	}
	return
}

// DecodeErrorFunc is decode error func.
type DecodeErrorFunc func(c context.Context, res *http.Response) (cancelRetry bool, err error)

// EncodeRequestFunc is request encode func.
type EncodeRequestFunc func(c context.Context, contentType string, in interface{}) (body []byte, err error)

// DecodeResponseFunc is response decode func.
type DecodeResponseFunc func(c context.Context, res *http.Response, out interface{}) (body []byte, err error)

// DefaultRequestEncoder is an HTTP request encoder.
func DefaultRequestEncoder(c context.Context, contentType string, in interface{}) (body []byte, err error) {
	subContentType := ContentSubtype(contentType)
	if subContentType == "" {
		subContentType = DefaultContentType
	}
	codec := encoding.GetCodec(subContentType)
	return codec.Marshal(in)
}

// DefaultResponseDecoder is an HTTP response decoder.
func DefaultResponseDecoder(c context.Context, res *http.Response, v interface{}) (body []byte, err error) {
	if v == nil {
		return
	}
	subContentType := ContentSubtype(res.Header.Get("Content-Type"))
	if subContentType == "" {
		subContentType = DefaultContentType
	}
	codec := encoding.GetCodec(subContentType)
	if codec == nil {
		err = errors.New("wrong Content-Type from header")
		return
	}

	defer res.Body.Close()
	if body, err = io.ReadAll(res.Body); err != nil {
		return
	}
	err = codec.Unmarshal(body, v)
	return
}

// DefaultErrorDecoder is an HTTP error decoder.
func DefaultErrorDecoder(c context.Context, resp *http.Response) (cancelRetry bool, err error) {
	if resp == nil {
		err = errors.New("nil *http.Response")
		return
	}
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return
	}
	if resp.StatusCode >= 400 && resp.StatusCode <= 407 {
		cancelRetry = true
		err = errors.New(resp.Status)
		return
	}
	err = errors.New(resp.Status)
	return
}

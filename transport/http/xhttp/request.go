package xhttp

/*
 * @abstract 传输协议http的客户端的操作方法
 * @mail neo532@126.com
 * @date 2022-05-30
 */

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/neo532/apitool/transport"
	"github.com/neo532/apitool/transport/http/xhttp/header"
	"github.com/neo532/apitool/transport/http/xhttp/queryargs"
)

type Client struct {
	logger transport.Logger
	req    *request
	resp   *response
	env    transport.Env
}
type response struct {
	httpResp *http.Response
	body     []byte
	cost     time.Duration
	err      error
}

type request struct {
	timeLimit time.Duration
	url       string
	method    string
	body      *bytes.Reader
	header    http.Header

	mapValue sync.Map

	responseMaxLength int
	bodyCurl          string
	headerCurl        string
	retryTimes        int
	retryMaxDuration  time.Duration
	retryDuration     time.Duration
	succFunc          func(statusCode int) (b bool)
}

// ========== Opt ==========
type Opt func(o *Client)

func WithEnv(env transport.Env) Opt {
	return func(o *Client) {
		o.env = env
	}
}
func WithTimeLimit(d time.Duration) Opt {
	return func(o *Client) {
		o.req.timeLimit = d
	}
}
func WithUrl(s string) Opt {
	return func(o *Client) {
		o.req.url = s
	}
}
func WithMethod(m string) Opt {
	return func(o *Client) {
		o.req.method = m
	}
}
func WithBody(param []byte) Opt {
	return func(o *Client) {
		o.req.body = bytes.NewReader(param)
		if len(param) > 0 {
			o.req.bodyCurl = " -d " + "'" + string(param) + "'"
		}
	}
}
func WithHeader(c context.Context) Opt {
	return func(o *Client) {

		var bHeader strings.Builder
		if md, ok := header.FromContext(c); ok {
			md.Range(func(k string, vs []string) (b bool) {
				for _, v := range vs {
					o.req.header.Set(k, v)
					bHeader.WriteString(" -H '" + k + ":" + v + "'")
				}
				return true
			})
		}

		o.req.headerCurl = bHeader.String()
	}
}
func WithRetryTimes(times int) Opt {
	return func(o *Client) {
		o.req.retryTimes = times
	}
}
func WithRetryDuration(d time.Duration) Opt {
	return func(o *Client) {
		o.req.retryDuration = d
	}
}
func WithRetryMaxDuration(d time.Duration) Opt {
	return func(o *Client) {
		o.req.retryMaxDuration = d
	}
}
func WithLogger(l transport.Logger) Opt {
	return func(o *Client) {
		o.logger = l
	}
}
func WithResponseMaxLength(l int) Opt {
	return func(o *Client) {
		o.req.responseMaxLength = l
	}
}
func WithSuccFunc(fn func(statusCode int) (b bool)) Opt {
	return func(o *Client) {
		o.req.succFunc = fn
	}
}
func WithMap(key, value string) Opt {
	return func(o *Client) {
		o.req.mapValue.Store(key, value)
	}
}

// ========== /Opt ==========

func New(opts ...Opt) (client Client) {
	client = Client{
		env: transport.EnvProd,
		req: &request{
			retryTimes:        2,
			retryDuration:     time.Microsecond,
			retryMaxDuration:  20 * time.Microsecond,
			responseMaxLength: 512,
			timeLimit:         3 * time.Second,
			body:              bytes.NewReader([]byte("")),
			header:            http.Header{},
			succFunc: func(statusCode int) (b bool) {
				return statusCode == http.StatusOK
			},
		},
		resp:   &response{},
		logger: &transport.LoggerDefault{},
	}
	for _, o := range opts {
		o(&client)
	}
	return
}

func (r Client) Do(c context.Context, opts ...Opt) (clt *Client) {
	req := r
	for _, o := range opts {
		o(&req)
	}
	if qa, ok := queryargs.FromContext(c); ok {
		qa.Range(func(k, v string) (b bool) {
			r.req.url = AppendUrlByKV(r.req.url, k, v)
			return true
		})
	}

	for i := 0; i < req.req.retryTimes; i++ {

		clt = req.do(c)
		if clt.resp.httpResp != nil &&
			req.req.succFunc(
				clt.resp.httpResp.StatusCode) {
			break
		}

		time.Sleep(req.req.retryDuration)
		if req.req.retryDuration < req.req.retryMaxDuration {
			req.req.retryDuration = req.req.retryDuration + req.req.retryDuration
		}
	}
	return
}

func (r Client) do(c context.Context) *Client {

	var req *http.Request
	req, r.resp.err = http.NewRequest(r.req.method, r.req.url, r.req.body)
	if r.resp.err != nil {
		return &r
	}
	req.Header = r.req.header

	// request
	client := &http.Client{Timeout: r.req.timeLimit}
	start := time.Now()
	r.resp.httpResp, r.resp.err = client.Do(req)
	r.resp.cost = time.Now().Sub(start)

	defer func() {
		var respCode int
		if r.resp.err == nil &&
			r.resp.httpResp != nil &&
			r.resp.httpResp.Body != nil {

			r.resp.body, r.resp.err = ioutil.ReadAll(r.resp.httpResp.Body)
			respCode = r.resp.httpResp.StatusCode
			r.resp.httpResp.Body.Close()
		}

		c = header.AppendToContext(c,
			"url", r.req.url,
			"cost", r.resp.cost.String(),
			"limit", r.req.timeLimit.String(),
			"httpcode", strconv.Itoa(respCode),
		)

		// cut response
		responseStr := string(r.resp.body)
		if r.env == transport.EnvProd && utf8.RuneCountInString(responseStr) > r.req.responseMaxLength {
			responseStr = string([]rune(responseStr)[:r.req.responseMaxLength]) + "..."
		}

		msg := fmt.Sprintf("[code:%d] [limit:%s] [cost:%s] [%s] [rst:%s]",
			respCode,
			r.req.timeLimit.String(),
			r.resp.cost.String(),
			r.fmtCurl(),
			responseStr,
		)

		if r.resp.err != nil {
			r.logger.Error(c, fmt.Sprintf("[err:%s] %s", r.resp.err.Error(), msg))
			return
		}

		r.logger.Info(c, msg)

		// body reader重置
		if r.req.body.Size() > 0 {
			r.req.body.Seek(0, io.SeekStart)
		}
	}()

	return &r
}

func (r Client) Cookie(c context.Context) (cookies []*http.Cookie) {
	cookies = r.resp.httpResp.Cookies()

	c = header.AppendToContext(c,
		"url", r.req.url,
		"cost", r.resp.cost.String(),
		"limit", r.req.timeLimit.String(),
		"httpcode", strconv.Itoa(r.resp.httpResp.StatusCode),
	)

	r.logger.Info(c, fmt.Sprintf("[code:%d] [limit:%s] [cost:%s] [%s] [cookie:%+v]",
		r.resp.httpResp.StatusCode,
		r.req.timeLimit.String(),
		r.resp.cost.String(),
		r.fmtCurl(),
		cookies,
	))

	return
}

func (r Client) Body(c context.Context) (body []byte, err error) {
	body = r.resp.body
	err = r.resp.err
	return
}

func (r Client) fmtCurl() string {
	var b strings.Builder
	b.WriteString("curl -X '")
	b.WriteString(r.req.method)
	b.WriteString("' '")
	b.WriteString(r.req.url)
	b.WriteString("'")
	b.WriteString(r.req.headerCurl)
	b.WriteString(r.req.bodyCurl)
	return b.String()
}

func (r Client) Env() transport.Env {
	return r.env
}

func (r Client) Get(key string) (value string) {
	if v, ok := r.req.mapValue.Load(key); ok {
		if s, ok := v.(string); ok {
			value = s
		}
	}
	return
}

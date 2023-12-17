package client

/*
 * @abstract 传输协议http的客户端的操作方法
 * @mail neo532@126.com
 * @date 2022-05-30
 */

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/neo532/apitool/transport"
	"github.com/neo532/apitool/transport/http/xhttp/middleware"
)

type Client struct {
	env               transport.Env
	logger            transport.Logger
	mapValue          sync.Map
	middlewares       []middleware.Middleware
	responseMaxLength int
	retryTimes        int

	transport *http.Transport

	httpClient *http.Client
}

// ========== Opt ==========
type Opt func(o *Client)

// ---------- xhttp ----------
func WithEnv(env transport.Env) Opt {
	return func(o *Client) {
		o.env = env
	}
}
func WithLogger(l transport.Logger) Opt {
	return func(o *Client) {
		o.logger = l
	}
}
func WithResponseMaxLength(l int) Opt {
	return func(o *Client) {
		o.responseMaxLength = l
	}
}
func WithMap(key, value string) Opt {
	return func(o *Client) {
		o.mapValue.Store(key, value)
	}
}

// ---------- request ----------
func WithMiddleware(ms ...middleware.Middleware) Opt {
	return func(o *Client) {
		o.middlewares = append(o.middlewares, ms...)
	}
}
func WithRetryTimes(times int) Opt {
	return func(o *Client) {
		o.retryTimes = times
	}
}
func WithTimeLimit(d time.Duration) Opt {
	return func(o *Client) {
		o.httpClient.Timeout = d
	}
}

// ---------- connect pool ----------
func WithIdleConnTimeout(d time.Duration) Opt {
	return func(o *Client) {
		o.transport.IdleConnTimeout = d
	}
}
func WithMaxConnsPerHost(n int) Opt {
	return func(o *Client) {
		o.transport.MaxConnsPerHost = n
	}
}
func WithMaxIdleConns(n int) Opt {
	return func(o *Client) {
		o.transport.MaxIdleConns = n
	}
}
func WithMaxIdleConnsPerHost(n int) Opt {
	return func(o *Client) {
		o.transport.MaxIdleConnsPerHost = n
	}
}

// ---------- tls ----------
func initTLS(o *Client) {
	if o.transport.TLSClientConfig == nil {
		o.transport.TLSClientConfig = &tls.Config{}
	}
}
func WithInsecureSkipVerify(b bool) Opt {
	return func(o *Client) {
		initTLS(o)
		o.transport.TLSClientConfig.InsecureSkipVerify = b
	}
}
func WithCertFile(crt, key string) (oR Opt, err error) {
	var cert tls.Certificate
	if cert, err = tls.LoadX509KeyPair(crt, key); err != nil {
		return
	}
	oR = func(o *Client) {
		initTLS(o)
		o.transport.TLSClientConfig.Certificates = []tls.Certificate{cert}
	}
	return
}
func WithCaCertFile(crt string) (oR Opt, err error) {
	var caCrt []byte
	if caCrt, err = os.ReadFile(crt); err != nil {
		return
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caCrt)

	oR = func(o *Client) {
		initTLS(o)
		o.transport.TLSClientConfig.RootCAs = pool
	}
	return
}
func WithCert(crt, key []byte) (oR Opt, err error) {
	var cert tls.Certificate
	if cert, err = tls.X509KeyPair(crt, key); err != nil {
		return
	}
	oR = func(o *Client) {
		initTLS(o)
		o.transport.TLSClientConfig.Certificates = []tls.Certificate{cert}
	}
	return
}

// ========== /Opt ==========

func New(opts ...Opt) (client Client) {

	client = Client{
		env:               transport.EnvProd,
		responseMaxLength: 1024,
		logger:            &transport.LoggerDefault{},
		middlewares:       make([]middleware.Middleware, 0, 10),
		retryTimes:        0,
		transport:         http.DefaultTransport.(*http.Transport).Clone(),
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
	client.transport.MaxConnsPerHost = runtime.NumCPU()*2 + 1
	for _, o := range opts {
		o(&client)
	}
	client.httpClient.Transport = client.transport
	return
}

func (r Client) Env() transport.Env {
	return r.env
}

func (r Client) Logger() transport.Logger {
	return r.logger
}

func (r Client) ResponseMaxLength() int {
	return r.responseMaxLength
}

func (r Client) RetryTime() int {
	return r.retryTimes
}

func (r Client) Value(key string) (value string) {
	if v, ok := r.mapValue.Load(key); ok {
		if s, ok := v.(string); ok {
			value = s
		}
	}
	return
}

func (r Client) CopyMiddleware() Client {
	mds := make([]middleware.Middleware, 0, 1)
	for _, mw := range mds {
		mds = append(mds, mw)
	}
	r.middlewares = mds
	return r
}

func (r Client) Copy() (clt Client) {
	clt = Client{
		env:               r.env,
		responseMaxLength: r.responseMaxLength,
		logger:            r.logger,
		middlewares:       r.middlewares,
		retryTimes:        r.retryTimes,
		httpClient: &http.Client{
			Timeout:   3 * time.Second,
			Transport: r.httpClient.Transport,
		},
	}
	return
}

func (r Client) AddMiddleware(mds ...middleware.Middleware) Client {
	for _, mw := range mds {
		r.middlewares = append(r.middlewares, mw)
	}
	return r
}

func (r Client) Middlewares() (ms []middleware.Middleware) {
	return r.middlewares
}

func (r Client) HttpClient() *http.Client {
	return r.httpClient
}

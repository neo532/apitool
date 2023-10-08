package client

/*
 * @abstract 传输协议http的客户端的操作方法
 * @mail neo532@126.com
 * @date 2022-05-30
 */

import (
	"sync"

	"github.com/neo532/apitool/transport"
	"github.com/neo532/apitool/transport/http/xhttp/middleware"
)

type Client struct {
	env               transport.Env
	logger            transport.Logger
	mapValue          sync.Map
	middlewares       []middleware.Middleware
	responseMaxLength int
}

// ========== Opt ==========
type Opt func(o *Client)

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

func WithMiddleware(ms ...middleware.Middleware) Opt {
	return func(o *Client) {
		o.middlewares = append(o.middlewares, ms...)
	}
}

// ========== /Opt ==========

func New(opts ...Opt) (client Client) {
	client = Client{
		env:               transport.EnvProd,
		responseMaxLength: 512,
		logger:            &transport.LoggerDefault{},
		middlewares:       make([]middleware.Middleware, 0, 1),
	}
	for _, o := range opts {
		o(&client)
	}
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

func (r Client) Get(key string) (value string) {
	if v, ok := r.mapValue.Load(key); ok {
		if s, ok := v.(string); ok {
			value = s
		}
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

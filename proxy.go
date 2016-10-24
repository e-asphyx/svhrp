package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type BufferPool sync.Pool

func (p *BufferPool) Get() []byte {
	return (*sync.Pool)(p).Get().([]byte)
}

func (p *BufferPool) Put(b []byte) {
	(*sync.Pool)(p).Put(b)
}

var bufferPool = BufferPool{
	New: func() interface{} {
		return make([]byte, 32*1024)
	},
}

func NewProxy(routes map[string]Route) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		route, ok := routes[req.URL.Host]
		if !ok {
			// Nothing to do
			return
		}

		u, err := url.Parse(route.Host)
		if err != nil {
			return
		}

		req.URL.Host = u.Host
		req.URL.Scheme = u.Scheme
	}

	return &httputil.ReverseProxy{
		Director:   director,
		BufferPool: &bufferPool,
	}
}

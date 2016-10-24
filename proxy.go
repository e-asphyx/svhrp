package main

import (
	"log"
	"net/http"
	"net/http/httputil"
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

func NewProxy(routes map[string]Route, defaultRoute string) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		log.Printf("%s %s%s\n", req.Method, req.Host, req.URL.String())

		route, ok := routes[req.Host]
		if !ok {
			req.Host = defaultRoute
			return
		}

		req.Host = route.Host
	}

	return &httputil.ReverseProxy{
		Director:   director,
		BufferPool: &bufferPool,
	}
}

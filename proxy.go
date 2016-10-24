package main

import (
	"log"
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

func NewProxy(routes map[string]Route, defaultRoute string) *httputil.ReverseProxy {
	defUrl, _ := url.Parse(defaultRoute)

	director := func(req *http.Request) {
		log.Printf("%s %s\n", req.Method, req.URL.String())
		route, ok := routes[req.URL.Host]
		if !ok {
			req.URL.Host = defUrl.Host
			req.URL.Scheme = defUrl.Scheme
			return
		}

		u, err := url.Parse(route.Host)
		if err != nil {
			log.Println(err)
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

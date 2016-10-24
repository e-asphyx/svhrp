package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
)

func main() {
	confFile := flag.String("c", "config.json", "Config")
	flag.Parse()

	conf, err := NewConfig(*confFile)
	if err != nil {
		log.Fatalln(err)
	}

	tlsConfig := tls.Config{}
	// Load certs
	for _, route := range conf.Routes {
		var (
			cert tls.Certificate
			err  error
		)

		if route.CertPEM != "" && route.KeyPEM != "" {
			cert, err = tls.X509KeyPair([]byte(route.CertPEM), []byte(route.KeyPEM))
		} else {
			cert, err = tls.LoadX509KeyPair(route.CertFile, route.KeyFile)
		}

		if err != nil {
			log.Fatalln(err)
		}

		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}
	tlsConfig.BuildNameToCertificate()

	server := http.Server{
		Addr:      conf.Listen,
		Handler:   NewProxy(conf.Routes, conf.DefaultHost),
		TLSConfig: &tlsConfig,
	}

	// Redirect to https
	go http.ListenAndServe(conf.HttpsRedirectorListen, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		url := *req.URL
		url.Scheme = "https"
		http.Redirect(w, req, url.String(), http.StatusMovedPermanently)
	}))

	// Main server
	log.Fatalln(server.ListenAndServeTLS("", ""))
}

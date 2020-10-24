package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var port *string
var backend *string

func main() {
	log.Printf("GoFE v0.0.1")
	const (
		defaultPort      = ":9090"
		defaultPortUsage = "The port the reverse proxy should listen on, default is ':9090'"
		urlUsage         = "The backend url to proxy"
	)

	port = flag.String("port", defaultPort, defaultPortUsage)
	backend = flag.String("backend", "http://10.76.42.1:8000/", urlUsage)

	flag.Parse()
	log.Printf("Proxy listening on port:", *port)
	log.Printf("Proxying backend", *backend)

	http.HandleFunc("/", proxy)
	http.ListenAndServe(*port, nil)
}

func proxy(res http.ResponseWriter, req *http.Request) {
	url, _ := url.Parse(*backend)
	proxy := httputil.NewSingleHostReverseProxy(url)

	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	//req.URL.Path = mux.Vars(req)["rest"]

	proxy.ServeHTTP(res, req)
}

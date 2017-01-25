package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

//to do: read available servers from config file
var r = roundRobbin{servers: []url.URL{
	{
		Scheme: "http",
		Host:   "localhost:9091",
	},
	{
		Scheme: "http",
		Host:   "localhost:9092",
	},
}}

func balanceLoad() *httputil.ReverseProxy {
	balancer := func(req *http.Request) {
		target := r.GetNext()
		log.Println(target)
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}
	return &httputil.ReverseProxy{Director: balancer}
}

func main() {
	//fmt.Println(r.GetNext())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	proxy := balanceLoad()

	log.Fatal(http.ListenAndServe(":9090", proxy))
}

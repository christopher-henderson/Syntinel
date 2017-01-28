package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
	"net"
	"strings"
)

//to do: read available servers from config file
var r = roundRobbin{servers: []url.URL{
	{
		Scheme: "http",
		Host:   "127.0.0.1:9091",
	},
	{
		Scheme: "http",
		Host:   "127.0.0.1:9092",
	},
}}

func UrlToString(url url.URL) string {
	temp := url.String()
	port := strings.Split(temp, ":")
	name := strings.Split(port[1],"//")
	temp = "["+name[1] + "]" +":"+ port[2]
	return temp
}

func balanceLoad() *httputil.ReverseProxy {
	balancer := func(req *http.Request) {
		target := r.GetNext()
		log.Println(target)
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
	}
	return &httputil.ReverseProxy{
		Director: balancer,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: func(network, addr string) (net.Conn, error) {
				for{
					conn, err := net.Dial(network, UrlToString(r.GetNext()))
					if err != nil {
						continue
					}
					return conn, nil
				}
				
			},
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}
}

func main() {
	//fmt.Println(r.GetNext())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	proxy := balanceLoad()

	log.Fatal(http.ListenAndServe(":9090", proxy))
}

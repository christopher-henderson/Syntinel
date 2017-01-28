package main

import (
	"log"
	"net/http"
	"net/http/httputil"
"net/url"
	"time"
	"net"
	"strings"
    "fmt"
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

func balanceLoad() (net.Conn, error) {
        fmt.Println("I am here")
        failed:
        conn, err := net.Dial("tcp", UrlToString(r.GetNext()))
        if err != nil{
            goto failed
        }else{
            return conn, nil
        }
            
    
    return nil, fmt.Errorf("Something broke!")
}


func GetReverseProxy() http.HandlerFunc{
    transport := &http.Transport{
        Proxy: http.ProxyFromEnvironment,
        Dial: func(network, addr string)(net.Conn, error){
            log.Println(addr)
            log.Println(network)
            return balanceLoad()
        },
        TLSHandshakeTimeout: 10 * time.Second,
    }
    return func(w http.ResponseWriter, req *http.Request){
        (&httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = req.Host
			},
			Transport: transport,
        }).ServeHTTP(w, req)
    }
}


func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
       proxy := GetReverseProxy()
       proxy.ServeHTTP(w,r)
    })
	
	log.Fatal(http.ListenAndServe(":9090",nil))
}

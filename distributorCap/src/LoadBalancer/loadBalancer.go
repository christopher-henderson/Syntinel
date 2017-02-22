//curl -H "Content-Type: application/json" -X POST -d '{"hostName":"something","somethingElse":"xyz"}' http://localhost:9090

package LoadBalancer

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var r = roundRobbin{servers: []url.URL{
	{
		Scheme: "http",
		Host:   "localhost:9092",
	},
	{
		Scheme: "http",
		Host:   "localhost:9091",
	},
}}

//ServerStruct... for all connected servers. HostName: string, port:string, Scheme:string
type ServerStruct struct {
	HostName string `json:"hostName"`
	Port     string `json:"port"`
	Scheme   string `json:"Scheme"`
}

//outputs to format [protocollhost]:port
func UrlToString(url url.URL) string {
	temp := url.String()
	port := strings.Split(temp, ":")
	name := strings.Split(port[1], "//")
	temp = "[" + name[1] + "]" + ":" + port[2]
	fmt.Println(temp)
	return temp
}

func balanceLoad() (net.Conn, error) {
failed:
	conn, err := net.Dial("tcp", UrlToString(r.GetNext()))
	if err != nil {
		goto failed
	} else {
		return conn, nil
	}
}

func GetReverseProxy() http.HandlerFunc {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network, addr string) (net.Conn, error) {
			log.Println(addr)
			log.Println(network)
			return balanceLoad()
		},
		TLSHandshakeTimeout: 10 * time.Second,
	}
	return func(w http.ResponseWriter, req *http.Request) {
		(&httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = "http"
				req.URL.Host = req.Host
			},
			Transport: transport,
		}).ServeHTTP(w, req)
	}
}

func AddToHosts(s ServerStruct) {

	newServer := url.URL{
		Scheme: s.Scheme,
		Host:   s.HostName + ":" + s.Port,
	}
	log.Println(newServer)
	r.servers = append(r.servers, newServer)
}

func ValidateServer(s ServerStruct) (valid bool) {
	log.Println(len(s.HostName))
	log.Println(len(s.Port))
	if len(s.HostName) < 1 || len(s.Port) < 1 {
		return false
	}
	return true
}

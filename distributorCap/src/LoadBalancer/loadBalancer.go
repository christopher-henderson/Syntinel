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

	"../DAO"
)

var r = roundRobbin{servers: []url.URL{
	{
		Scheme: "http",
		Host:   "localhost:9091",
	},
	{
		Scheme: "http",
		Host:   "localhost:9090",
	},
}}

//ServerStruct... for all connected servers. HostName: string, port:string, Scheme:string

func removeServer(url url.URL) {
	//needs implimentation
}

//outputs to format [protocollhost]:port
func UrlToString(url url.URL) string {
	temp := url.String()
	port := strings.Split(temp, ":")
	name := strings.Split(port[1], "//")
	temp = "[" + name[1] + "]" + ":" + port[2]
	return temp
}

func updateLastExecutor(ID int, url url.URL) {
	log.Println("Reaching updatelast")
	tmp := DAO.GetJob(ID)
	log.Println(tmp)
	if tmp.Canceled == false && tmp.Interval != 0 && tmp.Id != 0 {
		log.Println("Job exists, updating")
		tmp.LastExecutor = url
		fmt.Println(tmp)
		DAO.PutJob(tmp.Id, tmp)
		fmt.Println(tmp)
	} else {
		log.Println("Does not exist")
		DAO.RemoveJob(tmp.Id)
	}
}

func balanceLoad(ID int, doIt bool) (net.Conn, error) {
	log.Println("I am balancing")
failed:
	url := r.GetNext()
	if doIt {
		updateLastExecutor(ID, url)
	}
	conn, err := net.Dial("tcp", UrlToString(url))
	if err != nil {
		removeServer(url)
		goto failed
	} else {
		return conn, nil
	}
}

func GetReverseProxy(ID int, doIt bool) http.HandlerFunc {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network, addr string) (net.Conn, error) {
			log.Println(addr)
			log.Println(network)
			return balanceLoad(ID, doIt)
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

func AddToHosts(s []ServerStruct) {
	for _, server := range s {
		valid := ValidateServer(server)
		if valid {
			newServer := url.URL{
				Scheme: server.Scheme,
				Host:   server.HostName + ":" + server.Port,
			}
			log.Println(newServer)
			r.servers = append(r.servers, newServer)
		}
	}
	log.Println(r.servers)
}

func ValidateServer(s ServerStruct) (valid bool) {
	log.Println(len(s.HostName))
	log.Println(len(s.Port))
	if len(s.HostName) < 1 || len(s.Port) < 1 {
		return false
	}
	return true
}

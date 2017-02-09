//curl -H "Content-Type: application/json" -X POST -d '{"hostName":"something","somethingElse":"xyz"}' http://localhost:9090

package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var jobMap = NewJobMap()

//ServerStruct... for all connected servers. HostName: string, port:string, Scheme:string
type ServerStruct struct {
	HostName string `json:"hostName"`
	Port     string `json:"port"`
	Scheme   string `json:"Scheme"`
}

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

func urlToString(url url.URL) string {
	temp := url.String()
	port := strings.Split(temp, ":")
	name := strings.Split(port[1], "//")
	temp = "[" + name[1] + "]" + ":" + port[2]
	return temp
}

func balanceLoad() (net.Conn, error) {
failed:
	conn, err := net.Dial("tcp", urlToString(r.GetNext()))
	if err != nil {
		goto failed
	} else {
		return conn, nil
	}
}

func getReverseProxy() http.HandlerFunc {
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

func addToHosts(s ServerStruct) {

	newServer := url.URL{
		Scheme: s.Scheme,
		Host:   s.HostName + ":" + s.Port,
	}
	log.Println(newServer)
	r.servers = append(r.servers, newServer)
}

func validateServer(s ServerStruct) (valid bool) {
	log.Println(len(s.HostName))
	log.Println(len(s.Port))
	if len(s.HostName) < 1 || len(s.Port) < 1 {
		return false
	}
	return true
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//requests to register must be in format {"hostName":"localhost", "port": "9093", "Scheme": "http"}
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("error")
		}
		log.Println(string(body))
		var t ServerStruct
		err = json.Unmarshal(body, &t)
		if err != nil {
			log.Println("error")
		}
		valid := validateServer(t)
		if valid {
			addToHosts(t)
			//to do. Set Headers and response codes
			io.WriteString(w, "accepted, you are now registered")
		} else {
			io.WriteString(w, "registration rejected")
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy := getReverseProxy()
		proxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/schedule", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("error")
		}
		log.Println(string(body))
		var t job
		err = json.Unmarshal(body, &t)
		if err != nil {
			log.Println("error")
			io.WriteString(w, "There was a problem with your submission")
		} else {
			t.Canceled = false
			jobMap.Put(t.Id, t)
			go scheduleJob(t)
			io.WriteString(w, "Scheduled")
		}

	})

	http.HandleFunc("/cancel", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("error")
		}
		log.Println(string(body))
		var t job
		err = json.Unmarshal(body, &t)
		if err != nil {
			log.Println("error")
			io.WriteString(w, "There was a problem with your submission")
		} else {
			tmp := jobMap.Get(t.Id)
			tmp.Canceled = true
			jobMap.Put(tmp.Id, tmp)
			io.WriteString(w, "Job has been canceled")
		}
	})

	log.Fatal(http.ListenAndServe(":9090", nil))
}

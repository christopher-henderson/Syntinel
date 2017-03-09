package controller

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"../LoadBalancer"
)

func AddExecutor(w http.ResponseWriter, r *http.Request) {
	log.Println("Attemping to add executors")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error")
		io.WriteString(w, "There was a problem with your submission")
	}
	log.Println(string(body))
	var t []LoadBalancer.ServerStruct
	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Println("error")
		io.WriteString(w, "There was a problem with your submission")
	} else {
		log.Println("Accepted, adding")
		go LoadBalancer.AddToHosts(t)
		io.WriteString(w, "Accepted")
	}
}

func BalanceLoad(w http.ResponseWriter, r *http.Request) {
	log.Println("Balance Load")
	proxy := LoadBalancer.GetReverseProxy(0, false)
	proxy.ServeHTTP(w, r)
}

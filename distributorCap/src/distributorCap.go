package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"./LoadBalancer"
	"./Scheduler"
	"./controller"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("starting")
	//router := mux.NewRouter()
	//requests to register must be in format {"hostName":"localhost", "port": "9093", "Scheme": "http"}

	http.HandleFunc("/register", controller.AddExecutor)
	http.HandleFunc("/schedule", controller.ScheduleTest)
	http.HandleFunc("/kill", controller.Kill)
	http.HandleFunc("/cancel", controller.CancelTest)
	http.HandleFunc("/", mainRoute)
	log.Fatal(http.ListenAndServe(":9093", nil))
}
func mainRoute(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering default")
	buf, _ := ioutil.ReadAll(r.Body)
	rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
	rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))

	if r.URL.Path == "/test/run" {
		body, err := ioutil.ReadAll(rdr1)
		if err != nil {
			log.Println("error")
		}
		var t Scheduler.ExecutorRequestObj
		err = json.Unmarshal(body, &t)
		if err != nil {
			log.Println(t.ID)
			r.Body = rdr2
			proxy := LoadBalancer.GetReverseProxy(0, false)
			proxy.ServeHTTP(w, r)
		} else {
			r.Body = rdr2
			proxy := LoadBalancer.GetReverseProxy(t.ID, true)
			proxy.ServeHTTP(w, r)
		}
	} else {
		r.Body = rdr2
		proxy := LoadBalancer.GetReverseProxy(0, false)
		proxy.ServeHTTP(w, r)
	}

}

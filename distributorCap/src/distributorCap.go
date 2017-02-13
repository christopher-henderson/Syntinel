package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"./LoadBalancer"
	"./Scheduler"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	//requests to register must be in format {"hostName":"localhost", "port": "9093", "Scheme": "http"}
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("error")
		}
		log.Println(string(body))
		var t LoadBalancer.ServerStruct
		err = json.Unmarshal(body, &t)
		if err != nil {
			log.Println("error")
		}
		valid := LoadBalancer.ValidateServer(t)
		if valid {
			LoadBalancer.AddToHosts(t)
			//to do. Set Headers and response codes
			io.WriteString(w, "accepted, you are now registered")
		} else {
			io.WriteString(w, "registration rejected")
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy := LoadBalancer.GetReverseProxy()
		proxy.ServeHTTP(w, r)
	})

	http.HandleFunc("/schedule", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("error")
		}
		log.Println(string(body))
		var t Scheduler.Job
		err = json.Unmarshal(body, &t)
		if err != nil {
			log.Println("error")
			io.WriteString(w, "There was a problem with your submission")
		} else {
			t.Canceled = false
			Scheduler.ExportedjobMap.Put(t.Id, t)
			go Scheduler.ScheduleJob(t)
			io.WriteString(w, "Scheduled")
		}

	})

	http.HandleFunc("/cancel", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("error")
		}
		log.Println(string(body))
		var t Scheduler.Job
		err = json.Unmarshal(body, &t)
		if err != nil {
			log.Println("error")
			io.WriteString(w, "There was a problem with your submission")
		} else {
			tmp := Scheduler.ExportedjobMap.Get(t.Id)
			tmp.Canceled = true
			Scheduler.ExportedjobMap.Put(tmp.Id, tmp)
			io.WriteString(w, "Job has been canceled")
		}
	})

	log.Fatal(http.ListenAndServe(":9093", nil))
}

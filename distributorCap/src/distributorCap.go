package main

import (
	"bytes"
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
		var t []LoadBalancer.ServerStruct
		err = json.Unmarshal(body, &t)
		if err != nil {
			log.Println("error")
			io.WriteString(w, "registration rejected")
		} else {
			LoadBalancer.AddToHosts(t)
			io.WriteString(w, "accepted, you are now registered")
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.URL.Path)
		buf, _ := ioutil.ReadAll(r.Body)
		rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
		rdr2 := ioutil.NopCloser(bytes.NewBuffer(buf))
		r.Body = rdr2
		if r.URL.Path == "/test/run" {
			body, err := ioutil.ReadAll(rdr1)
			if err != nil {
				log.Println("error")
			}
			var t Scheduler.ExecutorRequestObj
			err = json.Unmarshal(body, &t)
			if err == nil {
				log.Println(t.ID)
				proxy := LoadBalancer.GetReverseProxy(t.ID, true)
				proxy.ServeHTTP(w, r)
			} else {
				proxy := LoadBalancer.GetReverseProxy(0, false)
				proxy.ServeHTTP(w, r)
			}
		} else {
			proxy := LoadBalancer.GetReverseProxy(0, false)
			proxy.ServeHTTP(w, r)
		}

	})

	http.HandleFunc("/cancelrunning", func(w http.ResponseWriter, r *http.Request) {
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

		}
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
			Scheduler.ExportedjobMap.Put(t.Id, t)
			Scheduler.ScheduleJob(t)
			io.WriteString(w, "scheduled")
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

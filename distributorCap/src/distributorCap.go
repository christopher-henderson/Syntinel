package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

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
		log.Println(reflect.TypeOf(t))
		j, worked := Scheduler.UnmarshalStruct(body, t)
		log.Println(j)
		if worked {
			log.Print("54")
			j, ok := j.(Scheduler.Job)
			log.Println(ok)
			log.Println(j)
			log.Println(reflect.TypeOf(j))
			if ok {
				log.Println("57")
				j.Canceled = false
				Scheduler.ExportedjobMap.Put(j.Id, j)
				go Scheduler.ScheduleJob(j)
				io.WriteString(w, "Scheduled")
			}
		} else {
			log.Println("63")
			io.WriteString(w, "There was an error with your submission")
		}
	})
	/*
		http.HandleFunc("/cancel", func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println("error")
			}
			log.Println(string(body))
			Scheduler.UnmarshalJob(body)
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
	*/
	log.Fatal(http.ListenAndServe(":9093", nil))
}

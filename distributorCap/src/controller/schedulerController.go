package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"../DAO"
	"../Scheduler"
)

func ScheduleTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to Schedule Test")
	var t DAO.Job
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Println(err)
		io.WriteString(w, "There was a problem with your submission")
	} else {
		log.Println("Accepted, Attempting to schedule Test")
		go Scheduler.ScheduleAndRunJob(t)

	}
}

func Kill(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to Kill Running Test")
	var t DAO.Job
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Println(err)
		io.WriteString(w, "There was a problem with your submission")
	} else {
		log.Println("Accepted, Attempting to Kill Running Test")
		go Scheduler.KillJob(t)
		w.WriteHeader(300)
		io.WriteString(w, "Accepted")
	}
}

func CancelTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to Cancel Test")
	var t DAO.Job
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Println(err)
		io.WriteString(w, err.Error())
	} else {
		log.Println("Accepted, Attempting to Cancel Test")
		go Scheduler.CancelJob(t)
		io.WriteString(w, "Accepted")
	}
}

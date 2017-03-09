package controller

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"../DAO"
	"../Scheduler"
)

func ScheduleTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to Schedule Test")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error")
		io.WriteString(w, "There was a problem with your submission")
	}
	log.Println(string(body))
	var t DAO.Job
	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Println("error")
		io.WriteString(w, "There was a problem with your submission")
	} else {
		log.Println("Accepted, Attemptin to schedule Test")
		go Scheduler.ScheduleAndRunJob(t)
		io.WriteString(w, "Accepted")
	}
}

func CancelTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting to Cancel Test")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("error")
		io.WriteString(w, "There was a problem with your submission")
	}
	log.Println(string(body))
	var t DAO.Job
	err = json.Unmarshal(body, &t)
	if err != nil {
		log.Println("error")
		io.WriteString(w, "There was a problem with your submission")
	} else {
		log.Println("Accepted, Attemping to Cancel Test")
		go Scheduler.CancelJob(t)
		io.WriteString(w, "Accepted")
	}
}

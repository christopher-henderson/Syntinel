package Scheduler

import (
	"fmt"
	"log"
	"time"
)

var StoredJobMap = NewJobMap()

func executeJob(j Job) {
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go
	/*
	   	dockerfile := `FROM docker.io/centos
	   MAINTAINER Christopher Henderson
	   RUN yum install -y go git wget
	   COPY script.sh $HOME/script.sh
	   CMD chmod +x script.sh && ./script.sh`
	   	script := `
	     !/usr/bin/env bash
	   git clone https://github.com/christopher-henderson/TestTheTester.git && cd TestTheTester/GoBeInGoodHands && go test . -v -cover
	   `
	   	obj := struct {
	   		ID                   int    `json:"id"`
	   		Dockerfile           string `json:"dockerfile"`
	   		Script               string `json:"script"`
	   		EnvironmentVariables string `json:"environmentVairables"`
	   	}{1, dockerfile, script, "a=b"}
	   	r, w := io.Pipe()
	   	log.Println(r)
	   	log.Println(w)
	   	go func() {
	   		json.NewEncoder(w).Encode(obj)
	   		w.Close()
	   	}()
	   	req, err := http.NewRequest("POST", "http://localhost:9090/test/run", r)
	   	r.Close()
	   	if err != nil {
	   		// handle err
	   	}
	   	req.Header.Set("Content-Type", "application/json")

	   	resp, err := http.DefaultClient.Do(req)
	   	if err != nil {
	   		// handle err
	   	}
	   	defer resp.Body.Close()
	*/
	fmt.Print("Executed: " + fmt.Sprint(j.Id))

}

func ScheduleAndRunJob(j Job) {
	StoredJobMap.Put(j.Id, j)
	go ScheduleJob(j)
}

func CancelJob(j Job) {
	tmp := StoredJobMap.Get(j.Id)
	tmp.Canceled = true
	StoredJobMap.Put(tmp.Id, tmp)
}

func ScheduleJob(j Job) {
	fmt.Println("scheduling job")
	time.Sleep(time.Duration(j.Interval) * time.Second)
	tmp := StoredJobMap.Get(j.Id)

	if tmp.Canceled == false && tmp.Interval != 0 && tmp.Id != 0 {
		go ScheduleJob(j)
		go executeJob(j)

	} else {
		log.Println("erased job")
		StoredJobMap.Delete(j.Id)
	}
}

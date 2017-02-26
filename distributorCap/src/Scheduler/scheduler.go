package Scheduler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var ExportedjobMap = NewJobMap()

func executeJob(j Job) {
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
	go func() {
		json.NewEncoder(w).Encode(obj)
		w.Close()
	}()
	req, err := http.NewRequest("POST", "http://localhost:9093/test/run", r)
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
	fmt.Print("Executed: ")
	fmt.Println(j.Id)

}

func ScheduleJob(j Job) {
	fmt.Println("scheduling job")
	time.Sleep(time.Duration(j.Interval) * time.Second)
	tmp := ExportedjobMap.Get(j.Id)

	if tmp.Canceled == false && tmp.Interval != 0 && tmp.Id != 0 {
		go ScheduleJob(j)
		go executeJob(j)

	} else {
		log.Println("erased job")
		ExportedjobMap.Delete(j.Id)
	}
}

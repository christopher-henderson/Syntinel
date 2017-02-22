package tests

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"../LoadBalancer"
	"../Scheduler"
)

func TestJobMap(t *testing.T) {
	jobMap := Scheduler.ExportedjobMap
	jobMap.Put(123, Scheduler.Job{Id: 123, Interval: 1, Canceled: false})
	myInput := jobMap.Get(123)
	if myInput.Id != 123 || myInput.Interval != 1 || myInput.Canceled != false {
		t.Errorf("inconsistant data in jobMap")
	}
	tmp := jobMap.Get(122)
	if tmp.Id != 0 && tmp.Interval != 0 && tmp.Canceled != false {
		t.Errorf("Job should have been removed")
	}
	fmt.Println(tmp)
}

func TestScheduler(t *testing.T) {
	jobMap := Scheduler.ExportedjobMap
	jobMap.Put(123, Scheduler.Job{Id: 123, Interval: 2, Canceled: false})
	Scheduler.ScheduleJob(jobMap.Get(123))
	tmp := jobMap.Get(123)
	tmp.Canceled = true
	jobMap.Put(tmp.Id, tmp)
	time.Sleep(time.Duration(3) * time.Second)
	tmp = jobMap.Get(123)
	if tmp.Id != 0 {
		t.Errorf("Job should have been removed from jobMap")
	}
}

func TestUrlToString(t *testing.T) {
	url := url.URL{
		Scheme: "http",
		Host:   "localhost:9092",
	}
	result := LoadBalancer.UrlToString(url)
	if result != "[localhost]:9092" {
		t.Errorf("UrlToString malfunction")
	}
}

func TestValidateServers(t *testing.T) {
	tmp := LoadBalancer.ServerStruct{HostName: "localhost", Port: "22", Scheme: "http"}
	result := LoadBalancer.ValidateServer(tmp)
	if result != true {
		t.Errorf("Valadate Server Malfunction")
	}
}

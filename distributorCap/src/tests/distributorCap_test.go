package tests

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"../DAO"
	"../LoadBalancer"
	"../Scheduler"
)

func TestJobMap(t *testing.T) {

	DAO.PutJob(123, DAO.Job{Id: 123, Interval: 1, Canceled: false})
	myInput := DAO.GetJob(123)
	if myInput.Id != 123 || myInput.Interval != 1 || myInput.Canceled != false {
		t.Errorf("inconsistant data in jobMap")
	}
	tmp := DAO.GetJob(122)
	if tmp.Id != 0 && tmp.Interval != 0 && tmp.Canceled != false {
		t.Errorf("Job should have been removed")
	}
	fmt.Println(tmp)
}

func TestScheduler(t *testing.T) {
	DAO.PutJob(123, DAO.Job{Id: 123, Interval: 2, Canceled: false})
	Scheduler.ScheduleJob(DAO.GetJob(123))
	tmp := DAO.GetJob(123)
	tmp.Canceled = true
	DAO.PutJob(tmp.Id, tmp)
	time.Sleep(time.Duration(3) * time.Second)
	tmp = DAO.GetJob(123)
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

func TestKillJob(t *testing.T) {
	var testUrl = url.URL{
		Scheme: "http",
		Host:   "localhost:9091",
	}
	tmp := DAO.Job{Id: 123, Interval: 2, Canceled: false, LastExecutor: testUrl}
	DAO.PutJob(tmp.Id, tmp)
	Scheduler.KillJob(tmp)
}

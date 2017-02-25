package Scheduler

import "net/url"

type Job struct {
	Id           int `json:"TestID"`
	Interval     int `json:"Interval"`
	Canceled     bool
	LastExecutor url.URL
}

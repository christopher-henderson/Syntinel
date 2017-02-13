package Scheduler

type Job struct {
	Id       int `json:"TestID"`
	Interval int `json:"Interval"`
	Canceled bool
}

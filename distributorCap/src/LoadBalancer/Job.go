package main

type job struct {
	Id       int `json:"TestID"`
	Interval int `json:"Interval"`
	Canceled bool
}

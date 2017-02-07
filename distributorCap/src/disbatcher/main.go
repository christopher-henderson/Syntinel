package main

import (
	"fmt"
	"time"
)

var jobMap = make(map[int]job)
var jobChannel = make(chan job)

func main() {

	testJob := job{id: 111, interval: 10, canceled: false}
	testJob1 := job{id: 112, interval: 11, canceled: true}
	testJob2 := job{id: 113, interval: 13, canceled: false}
	go addToJobMap(jobChannel, testJob, jobMap)
	x := <-jobChannel
	go addToJobMap(jobChannel, testJob1, jobMap)
	y := <-jobChannel
	go addToJobMap(jobChannel, testJob2, jobMap)
	z := <-jobChannel
	fmt.Println(x)
	fmt.Println(y)
	fmt.Println(z)
	for {
		fmt.Println("running")
		time.Sleep(time.Duration(30) * time.Second)
	}

}
func addToJobMap(c chan job, j job, m map[int]job) {
	value, ok := m[j.id]
	if ok {
		c <- j
		fmt.Println(value)
	} else {
		m[j.id] = j
		c <- j
	}

	go scheduleJob(j)
}

//returns a job, if the job does not exist, returns empty job to the channel
func getJob(c chan job, id int, m map[int]job) {
	value, tmp := m[id]
	if tmp {
		c <- value
		fmt.Println(value)
	} else {
		c <- job{id: 0, interval: 0, canceled: false}
	}
}
func executeJob(j job) {
	fmt.Print("Executed: ")
	fmt.Println(j.id)
}

func removeFromJobMap(c chan job, j job, m map[int]job) {
	_, ok := jobMap[j.id]
	if ok {
		delete(jobMap, j.id)
		fmt.Print("Removed: ")
		fmt.Println(j)
	}
	//return empty job for validation
	c <- job{id: 0, interval: 0, canceled: true}

}

func scheduleJob(j job) {
	fmt.Println("scheduling job")
	time.Sleep(time.Duration(j.interval) * time.Second)
	go getJob(jobChannel, j.id, jobMap)
	tmp := <-jobChannel
	if tmp.canceled == false {
		go scheduleJob(j)
		go executeJob(tmp)
	} else {
		go removeFromJobMap(jobChannel, j, jobMap)
		tmp := <-jobChannel
		if tmp.id == 0 && tmp.interval == 0 && tmp.canceled == true {
			fmt.Println("succesfuly removed")
		} else {
			fmt.Println("There was a problem removing the item")
		}
	}
}

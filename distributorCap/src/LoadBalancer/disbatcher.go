package main

import (
	"fmt"
	"time"
)

func executeJob(j job) {
	fmt.Print("Executed: ")
	fmt.Println(j.Id)
}

func scheduleJob(j job) {
	fmt.Println("scheduling job")
	time.Sleep(time.Duration(j.Interval) * time.Second)
	tmp := jobMap.Get(j.Id)

	if tmp.Canceled == false {
		go scheduleJob(j)
		go executeJob(tmp)
	} else {
		jobMap.Delete(j.Id)
	}
}

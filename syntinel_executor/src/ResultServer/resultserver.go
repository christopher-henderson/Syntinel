package ResultServer

import (
	"log"
	"sync"
	"syntinel_executor/PAO/process"
	"time"
)

type FinalResult struct {
	ID     int
	Result *process.TestRunResult
}

const maxBackoff = 20

var rs = newResultServer()

type resultServer struct {
	in      chan *FinalResult
	retry   chan *FinalResult
	mutex   sync.Mutex
	backoff int
}

func newResultServer() *resultServer {
	rs := &resultServer{make(chan *FinalResult), make(chan *FinalResult), sync.Mutex{}, 0}
	go rs.ListenForResults()
	go rs.ListenForRetries()
	return rs
}

func SendResult(result *FinalResult) {
	rs.SendResult(result)
}

func (rs *resultServer) SendResult(result *FinalResult) {
	select {
	default:
		rs.in <- result
	}
}

func (rs *resultServer) ListenForResults() {
	for {
		go rs.handle(<-rs.in)
	}
}

func (rs *resultServer) ListenForRetries() {
	for {
		go rs.handle(<-rs.retry)
		time.Sleep(time.Second * time.Duration(rs.backoff))
	}
}

func (rs *resultServer) handle(result *FinalResult) {
	log.Println("The Result Server got the following result: ", result.Result.Output)
	succeeded := true
	if !succeeded {
		rs.mutex.Lock()
		rs.backoff += 1
		if rs.backoff < maxBackoff {
			rs.backoff = maxBackoff
		}
		rs.mutex.Unlock()
		rs.enqueue(result)
	} else {
		rs.mutex.Lock()
		rs.backoff -= 1
		if rs.backoff < 0 {
			rs.backoff = 0
		}
		rs.mutex.Unlock()
	}
}

func (rs *resultServer) enqueue(result *FinalResult) {

}

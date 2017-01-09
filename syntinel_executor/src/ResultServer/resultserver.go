package ResultServer

import (
	"log"
	"sync"
	"syntinel_executor/PAO/process"
	"time"
)

const maxBackoff = 20

var rs = newResultServer()

type resultServer struct {
	in      chan *process.TestRunResult
	retry   chan *process.TestRunResult
	mutex   sync.Mutex
	backoff int
}

func newResultServer() *resultServer {
	rs := &resultServer{make(chan *process.TestRunResult), make(chan *process.TestRunResult), sync.Mutex{}, 0}
	go rs.ListenForResults()
	go rs.ListenForRetries()
	return rs
}

func SendResult(result *process.TestRunResult) {
	rs.SendResult(result)
}

func (rs *resultServer) SendResult(result *process.TestRunResult) {
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

func (rs *resultServer) handle(result *process.TestRunResult) {
	log.Println("The Result Server got the following result: ", result)
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

func (rs *resultServer) enqueue(result *process.TestRunResult) {

}

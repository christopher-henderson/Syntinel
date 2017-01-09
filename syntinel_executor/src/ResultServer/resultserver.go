package ResultServer

import (
	"log"
	"syntinel_executor/PAO/process"
)

var rs = newResultServer()

type resultServer struct {
	in chan *process.TestRunResult
}

func newResultServer() *resultServer {
	rs := &resultServer{make(chan *process.TestRunResult)}
	go rs.ListenForResults()
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

func (rs *resultServer) handle(result *process.TestRunResult) {
	log.Println("The Result Server got the following result: ", result)
}

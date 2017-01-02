package work

import (
	"log"
	"syntinel_executor/PAO/process"
)

const (
	created          = iota
	startingDocker   = iota
	dockerStarted    = iota
	sendingScripts   = iota
	scriptsSent      = iota
	executingScripts = iota
	resultsReceived  = iota
	done             = iota
)

type Work struct {
	Cancel chan uint8
}

func NewWork() *Work {
	return &Work{make(chan uint8)}
}

func (w *Work) awaitOutput(function func() (*process.Process, <-chan *process.WorkResult, chan<- uint8)) bool {
	proc, result, cancel := function()
	defer close(cancel)
	proc.Start()
	select {
	case <-w.Cancel:
		log.Println("Received kill request in work.(*Work).awaitOuput")
		cancel <- 1
		workResult := <-result
		log.Println(workResult.Err)
		log.Println(workResult.Output)
		return false
	case workResult := <-result:
		log.Println("Received finished result in work.(*Work).awaitOuput")
		if workResult.Err != nil {
			log.Println(workResult.Err)
			log.Println(workResult.Output)
			return false
		}
		return true
	}
}

func (w *Work) Run() {
	defer w.destroyDocker()
	if !w.awaitOutput(w.createDocker) {
		return
	}
	if !w.awaitOutput(w.scpScript) {
		return
	}
	w.awaitOutput(w.RunTest)
}

func (w *Work) createDocker() (*process.Process, <-chan *process.WorkResult, chan<- uint8) {
	return process.NewProcess("echo", "hello world from create Docker")
}

func (w *Work) scpScript() (*process.Process, <-chan *process.WorkResult, chan<- uint8) {
	return process.NewProcess("echo", "hello world from SCP Script")
}

func (w *Work) RunTest() (*process.Process, <-chan *process.WorkResult, chan<- uint8) {
	return process.NewProcess("echo", "hello world from RunTest")
}

func (w *Work) destroyDocker() {
	proc, result, cancel := process.NewProcess("echo", "hello world from Destroy Docker!")
	defer close(cancel)
	proc.Start()
	<-result
}

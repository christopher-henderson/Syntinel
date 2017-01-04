package PAO

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

type TestRun struct {
	ID         int
	Cancel     chan uint8
	dockerPath string
	scriptPath string
}

func NewTestRun(id int, dockerPath string, scriptPath string) *TestRun {
	return &TestRun{id, make(chan uint8), dockerPath, scriptPath}
}

func (w *TestRun) awaitOutput(function func() (*process.Process, <-chan *process.TestRunResult, chan<- uint8)) bool {
	proc, result, cancel := function()
	defer close(cancel)
	proc.Start()
	select {
	case <-w.Cancel:
		log.Println("Received kill request in TestRun.(*TestRun).awaitOuput")
		cancel <- 1
		testRunResult := <-result
		log.Println(testRunResult.Err)
		log.Println(testRunResult.Output)
		return false
	case testRunResult := <-result:
		log.Println("Received finished result in TestRun.(*TestRun).awaitOuput")
		log.Println(testRunResult.Err)
		log.Println(testRunResult.Output)
		if testRunResult.Err != nil {
			return false
		}
		return true
	}
}

func (w *TestRun) Run() {
	defer w.destroyDocker()
	if !w.awaitOutput(w.createDocker) {
		return
	}
	if !w.awaitOutput(w.scpScript) {
		return
	}
	w.awaitOutput(w.RunTest)
}

func (w *TestRun) createDocker() (*process.Process, <-chan *process.TestRunResult, chan<- uint8) {
	return process.NewProcess("echo", "hello world from "+w.dockerPath)
}

func (w *TestRun) scpScript() (*process.Process, <-chan *process.TestRunResult, chan<- uint8) {
	return process.NewProcess("echo", "hello world from "+w.scriptPath)
}

func (w *TestRun) RunTest() (*process.Process, <-chan *process.TestRunResult, chan<- uint8) {
	return process.NewProcess(w.scriptPath)
	// args := []string{"-c", "from time import sleep;print('...thinking');sleep(5);print('AH HA!');raise Exception('wut happun')"}
	// return process.NewProcess("python", args...)
}

func (w *TestRun) destroyDocker() {
	proc, result, cancel := process.NewProcess("echo", "hello world from Destroy Docker!")
	defer close(cancel)
	proc.Start()
	<-result
}

package PAO

import (
	"log"
	"sync"
	"syntinel_executor/PAO/process"
)

type TestRun struct {
	ID         int
	Cancel     chan uint8
	state      int
	dockerPath string
	scriptPath string
	mutex      sync.Mutex
}

func NewTestRun(id int, dockerPath string, scriptPath string) *TestRun {
	return &TestRun{id, make(chan uint8), Queued, dockerPath, scriptPath, sync.Mutex{}}
}

func (t *TestRun) Run() {
	t.setState(Starting)
	defer t.destroyDocker()
	defer t.setState(Done)
	if result := t.awaitOutput(t.createDocker); result.Err != nil {
		log.Fatalln(result.Err)
	}
	if result := t.awaitOutput(t.scpScript); result.Err != nil {
		log.Fatalln(result.Err)
	}
	result := t.awaitOutput(t.runTest)
	// ResultServer.Post(t.ID, result)
	log.Println(result)
}

func (t *TestRun) Query() int {
	return t.getState()
}

func (t *TestRun) createDocker() (*process.Process, <-chan *process.TestRunResult, chan<- uint8) {
	t.setState(StartingDocker)
	defer t.setState(DockerStarted)
	args := []string{"-c", "from time import sleep;print('...thinking');sleep(15);print('Docker started!')"}
	return process.NewProcess("python", args...)
	// return process.NewProcess("echo", "hello world from "+t.dockerPath)
}

func (t *TestRun) scpScript() (*process.Process, <-chan *process.TestRunResult, chan<- uint8) {
	t.setState(SendingScripts)
	defer t.setState(ScriptsSent)
	args := []string{"-c", "from time import sleep;print('...thinking');sleep(5);print('SCP Done!')"}
	return process.NewProcess("python", args...)
	// return process.NewProcess("echo", "hello world from "+t.scriptPath)
}

func (t *TestRun) runTest() (*process.Process, <-chan *process.TestRunResult, chan<- uint8) {
	t.setState(ExecutingScripts)
	defer t.setState(ResultsReceived)
	return process.NewProcess(t.scriptPath)
	// args := []string{"-c", "from time import sleep;print('...thinking');sleep(5);print('AH HA!');raise Exception('wut happun')"}
	// return process.NewProcess("python", args...)
}

func (t *TestRun) destroyDocker() {
	t.setState(DestroyingDocker)
	defer t.setState(DockerDestroyed)
	proc, result, cancel := process.NewProcess("echo", "hello world from Destroy Docker!")
	defer close(cancel)
	proc.Start()
	<-result
}

func (t *TestRun) setState(state int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.state = state
}

func (t *TestRun) getState() int {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.state
}

func (w *TestRun) awaitOutput(function func() (*process.Process, <-chan *process.TestRunResult, chan<- uint8)) *process.TestRunResult {
	proc, result, cancel := function()
	var testRunResult *process.TestRunResult
	defer close(cancel)
	proc.Start()
	select {
	case <-w.Cancel:
		log.Println("Received kill request in TestRun.(*TestRun).awaitOuput")
		cancel <- 1
		testRunResult = <-result
	case testRunResult = <-result:
		log.Println("Received finished result in TestRun.(*TestRun).awaitOuput")
	}
	return testRunResult
}

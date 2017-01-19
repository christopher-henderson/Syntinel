package PAO

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"syntinel_executor/PAO/process"
	"syntinel_executor/ResultServer"
)

type TestRun struct {
	testID     int
	ID         int
	Cancel     chan uint8
	state      int
	dockerPath string
	scriptPath string
	mutex      sync.RWMutex
}

func NewTestRun(testID int, id int, dockerPath string, scriptPath string) *TestRun {
	return &TestRun{testID, id, make(chan uint8), Queued, dockerPath, scriptPath, sync.RWMutex{}}
}

func (t *TestRun) ImageName() string {
	return "executor_" + strconv.Itoa(t.ID)
}

func (t *TestRun) DockerBuildDirectory() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(path)
	return path + string(os.PathSeparator) + "assets" +
		string(os.PathSeparator) + "build" + string(os.PathSeparator) + t.ImageName()
}

func (t *TestRun) Run() {
	t.setState(Starting)
	finalResult := &ResultServer.FinalResult{t.ID, nil}
	// defer t.destroyDocker()
	defer t.setState(Done)
	defer ResultServer.SendResult(finalResult)
	// if result := t.awaitOutput(t.createDocker); result.Err != nil {
	// 	log.Println(result.Err)
	// 	finalResult.Result = result
	// 	return
	// }
	// if result := t.awaitOutput(t.startDocker); result.Err != nil {
	// 	log.Println(result.Err)
	// 	finalResult.Result = result
	// 	return
	// }
	// if result := t.awaitOutput(t.scpScript); result.Err != nil {
	// 	log.Println(result.Err)
	// 	finalResult.Result = result
	// 	return
	// }
	// result := t.awaitOutput(t.runTest)
	if result := t.awaitOutput(t.buildDockerImage); result.Err != nil {
		log.Println(result.Err)
		finalResult.Result = result
		return
	}
	finalResult.Result = t.awaitOutput(t.runDockerImage)
}

func (t *TestRun) Query() int {
	return t.getState()
}

func (t *TestRun) buildDockerImage() (*process.Process, <-chan *process.TestRunResult, chan<- uint8) {
	if err := os.MkdirAll(t.DockerBuildDirectory(), os.ModeDir); err != nil {
		log.Fatalln(err)
	}
	// if err := os.Chdir(t.DockerBuildDirectory()); err != nil {
	// 	log.Fatalln(err)
	// }
	if err := copy(t.scriptPath, t.DockerBuildDirectory()+string(os.PathSeparator)+"script.sh"); err != nil {
		log.Fatalln(err)
	}
	if err := copy(t.dockerPath, t.DockerBuildDirectory()+string(os.PathSeparator)+"Dockerfile"); err != nil {
		log.Fatalln(err)
	}
	command := "docker"
	args := []string{"build", "-t", t.ImageName(), "--force-rm", t.DockerBuildDirectory()}
	return process.NewProcess(command, args...)
}

func (t *TestRun) runDockerImage() (*process.Process, <-chan *process.TestRunResult, chan<- uint8) {
	command := "docker"
	args := []string{"run", "--rm", t.ImageName()}
	return process.NewProcess(command, args...)
}

func (t *TestRun) createDocker() (*process.Process, <-chan *process.TestRunResult, chan<- uint8) {
	t.setState(StartingDocker)
	defer t.setState(DockerStarted)
	command := "docker"
	// args := []string{"build", "-t", 50, "--force-rm", "."}
	args := []string{"-c", "from time import sleep;print('...thinking');sleep(15);print('Docker started!')"}
	return process.NewProcess(command, args...)
	// return process.NewProcess("echo", "hello world from "+t.dockerPath)
}

func (t *TestRun) startDocker() (*process.Process, <-chan *process.TestRunResult, chan<- uint8) {
	t.setState(StartingDocker)
	defer t.setState(DockerStarted)
	command := "docker"
	args := []string{"start", "ab"}
	// args := []string{"-c", "from time import sleep;print('...thinking');sleep(15);print('Docker started!')"}
	return process.NewProcess(command, args...)
	// return process.NewProcess("echo", "hello world from "+t.dockerPath)
}

func (t *TestRun) scpScript() (*process.Process, <-chan *process.TestRunResult, chan<- uint8) {
	t.setState(SendingScripts)
	defer t.setState(ScriptsSent)
	command := "docker"
	args := []string{"cp", t.scriptPath, "ab"}
	return process.NewProcess(command, args...)
}

func (t *TestRun) runTest() (*process.Process, <-chan *process.TestRunResult, chan<- uint8) {
	t.setState(ExecutingScripts)
	defer t.setState(ResultsReceived)
	command := "docker"
	args := []string{"exec", "ab"}
	return process.NewProcess(command, args...)
	// return process.NewProcess(t.scriptPath)
	// args := []string{"-c", "from time import sleep;print('...thinking');sleep(5);print('AH HA!');raise Exception('wut happun')"}
	// return process.NewProcess("python", args...)
}

func (t *TestRun) destroyDocker() {
	t.setState(DestroyingDocker)
	defer t.setState(DockerDestroyed)
	command := "docker"
	args := []string{"rm", "ab"}
	proc, result, cancel := process.NewProcess(command, args...)
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
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.state
}

func (w *TestRun) awaitOutput(function func() (*process.Process, <-chan *process.TestRunResult, chan<- uint8)) *process.TestRunResult {
	proc, result, cancel := function()
	var testRunResult *process.TestRunResult
	defer close(cancel)
	proc.Start()
	select {
	case <-w.Cancel:
		log.Println("Received kill request.")
		cancel <- 1
		testRunResult = <-result
	case testRunResult = <-result:
		log.Println("Received finished result.")
	}
	return testRunResult
}

func copy(source string, destination string) error {
	if _, err := os.Stat(source); err != nil {
		return nil
	}
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return nil
}

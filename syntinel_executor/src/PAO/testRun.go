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

func (t *TestRun) buildDockerImage() *process.Process {
	t.setState(MakingBuildDirectory)
	if err := os.MkdirAll(t.DockerBuildDirectory(), os.ModeDir); err != nil {
		t.setState(Failed)
		log.Fatalln(err)
	}
	t.setState(CopyingScript)
	if err := copy(t.scriptPath, t.DockerBuildDirectory()+string(os.PathSeparator)+"script.sh"); err != nil {
		t.setState(Failed)
		log.Fatalln(err)
	}
	t.setState(CopyingDockerfile)
	if err := copy(t.dockerPath, t.DockerBuildDirectory()+string(os.PathSeparator)+"Dockerfile"); err != nil {
		t.setState(Failed)
		log.Fatalln(err)
	}
	command := "docker"
	args := []string{"build", "-t", t.ImageName(), "--force-rm", t.DockerBuildDirectory()}
	t.setState(BuildingImage)
	return process.NewProcess(command, args...)
}

func (t *TestRun) runDockerImage() *process.Process {
	command := "docker"
	args := []string{"run", "--rm", t.ImageName()}
	return process.NewProcess(command, args...)
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

func (w *TestRun) awaitOutput(function func() *process.Process) *process.TestRunResult {
	proc := function()
	var testRunResult *process.TestRunResult
	result := make(chan *process.TestRunResult)
	defer close(result)
	proc.Start()
	go func() {
		result <- proc.Wait()
	}()
	select {
	case <-w.Cancel:
		log.Println("Received kill request.")
		proc.Kill()
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

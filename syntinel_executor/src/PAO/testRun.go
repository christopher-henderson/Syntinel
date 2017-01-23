package PAO

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"syntinel_executor/ResultServer"
	"syntinel_executor/utils"
	"syntinel_executor/utils/process"
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
	return fmt.Sprintf("%v%v", utils.ContainerPrefix, strconv.Itoa(t.ID))
}

func (t *TestRun) DockerBuildDirectory() string {
	return fmt.Sprintf("%v%v%v", utils.BuildDirectory(), t.ImageName(), string(os.PathSeparator))
}

func (t *TestRun) Run() {
	t.setState(Starting)
	defer t.setState(Done)
	if err := t.awaitOutput(t.buildDockerImage); err != nil {
		log.Println(err)
		return
	}
	t.awaitOutput(t.runDockerImage)
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
	if err := utils.FileCopy(t.scriptPath, t.DockerBuildDirectory()+utils.DockerScriptName); err != nil {
		t.setState(Failed)
		log.Fatalln(err)
	}
	t.setState(CopyingDockerfile)
	if err := utils.FileCopy(t.dockerPath, t.DockerBuildDirectory()+utils.DockerFile); err != nil {
		t.setState(Failed)
		log.Fatalln(err)
	}
	args := []string{
		utils.DockerBuild,        // Build image
		utils.DockerBuildTag,     // Define the image tag.
		t.ImageName(),            // Using the tag built by "ImageName"
		utils.DockerBuildForceRM, // Delete the temporary container.
		t.DockerBuildDirectory()} // Build in this directory.
	t.setState(BuildingImage)
	return process.NewProcess(utils.DockerCommand, args...)
}

func (t *TestRun) runDockerImage() *process.Process {
	defer t.deleteContainer()
	args := []string{
		utils.DockerRun,     // Run image
		utils.DockerRunRM,   // Delete after completed.
		utils.DockerRunName, // Name the container.
		t.ImageName(),       // Name of the container.
		t.ImageName()}       // Name of the image to run
	return process.NewProcess(utils.DockerCommand, args...)
}

func (t *TestRun) deleteContainer() {
	argsStop := []string{utils.DockerStop, t.ImageName()}
	argsDelete := []string{utils.DockerRM, t.ImageName()}
	process.NewProcess(utils.DockerCommand, argsStop...).Start().Wait()
	process.NewProcess(utils.DockerCommand, argsDelete...).Start().Wait()
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

func (t *TestRun) awaitOutput(function func() *process.Process) error {
	proc := function()
	var testRunResult error
	result := make(chan error)
	defer close(result)
	ResultServer.Stream(t.ID, proc.OutputStream())
	proc.Start()
	go func() {
		result <- proc.Wait()
	}()
	select {
	case <-t.Cancel:
		log.Println("Received kill request.")
		proc.Kill()
		testRunResult = <-result
	case testRunResult = <-result:
		log.Println("Received finished result.")
	}
	return testRunResult
}

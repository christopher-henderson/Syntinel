package entities

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"syntinel_executor/ResultServer"
	"syntinel_executor/utils"
	"syntinel_executor/utils/process"
)

const (
	DockerCommand      = "docker"
	DockerBuild        = "build"
	DockerBuildForceRM = "--force-rm" // Force deletion of the temporary image.
	DockerBuildTag     = "-t"         // Give a tag to the built image.
	DockerRun          = "run"
	DockerRunRM        = "--rm"   // Delete container upon completion of run.
	DockerRunName      = "--name" // Name the container.
	DockerStop         = "stop"   // Stop the named container. Used if the container is killed prematurely.
	DockerRM           = "rm"     // Destroy the named container. Used if the container is killed prematurely.
	DockerScriptName   = "script.sh"
	DockerFile         = "Dockerfile"
	ContainerPrefix    = "executor_" // Namespace containers created by the executor.
)

const (
	NotFound             = iota
	Queued               = iota
	MakingBuildDirectory = iota
	CopyingDockerfile    = iota
	CopyingScript        = iota
	RunningTest          = iota
	BuildingImage        = iota
	Starting             = iota
	StartingDocker       = iota
	DockerStarted        = iota
	SendingScripts       = iota
	ScriptsSent          = iota
	ExecutingScripts     = iota
	ResultsReceived      = iota
	DestroyingDocker     = iota
	DockerDestroyed      = iota
	Done                 = iota
	Failed               = iota
)

var stateMap = map[int]string{
	NotFound:             "Not found.",
	Queued:               "Test awaiting execution.",
	MakingBuildDirectory: "Making build directory",
	BuildingImage:        "Building image.",
	Starting:             "Test is starting.",
	StartingDocker:       "Test Docker is initializing.",
	DockerStarted:        "Test Docker has started.",
	SendingScripts:       "Sending artifacts to the test Docker.",
	ScriptsSent:          "Artifacts have been sent to the test Docker.",
	ExecutingScripts:     "Test is executing.",
	DestroyingDocker:     "Test Docker is being torn down.",
	DockerDestroyed:      "Test Docker has been torn down.",
	Done:                 "Done.",
}

func TestStateToString(state int) string {
	var stateString string
	var ok bool
	if stateString, ok = stateMap[state]; !ok {
		log.Fatalln("Bad test state passed to mapping function.")
	}
	return stateString
}

type TestRunEntity struct {
	ID                   int    `json:"id"`
	TestID               int    `json:"testID"`
	EnvironmentVariables string `json:"environmentVariables"`
	Dockerfile           string `json:"dockerfile"`
	Script               string `json:"script"`

	Cancel chan uint8
	state  int
	mutex  sync.RWMutex
}

func (t *TestRunEntity) Run() {
	t.Cancel = make(chan uint8)
	t.mutex = sync.RWMutex{}
	t.setState(Starting)
	defer t.setState(Done)
	if err := t.awaitOutput(t.buildDockerImage); err != nil {
		log.Println(err)
		return
	}
	testError := t.awaitOutput(t.runDockerImage)
	ResultServer.Finalize(t.ID, testError)
}

func (t *TestRunEntity) ImageName() string {
	return fmt.Sprintf("%v%v", ContainerPrefix, strconv.Itoa(t.TestID))
}

func (t *TestRunEntity) DockerBuildDirectory() string {
	return fmt.Sprintf("%v%v%v", utils.BuildDirectory(), t.ImageName(), string(os.PathSeparator))
}

func (t *TestRunEntity) Query() int {
	return t.getState()
}

func (t *TestRunEntity) buildDockerImage() *process.Process {
	t.setState(MakingBuildDirectory)
	buildDirectory := t.DockerBuildDirectory()
	scriptPath := fmt.Sprintf("%v%v", buildDirectory, DockerScriptName)
	dockerfilePath := fmt.Sprintf("%v%v", buildDirectory, DockerFile)
	if err := os.MkdirAll(buildDirectory, os.ModeDir); err != nil {
		t.setState(Failed)
		log.Fatalln(err)
	}
	ioutil.WriteFile(dockerfilePath, []byte(t.Dockerfile), 0666)
	ioutil.WriteFile(scriptPath, []byte(t.Script), 0666)
	args := []string{
		DockerBuild,              // Build image
		DockerBuildTag,           // Define the image tag.
		t.ImageName(),            // Using the tag built by "ImageName"
		DockerBuildForceRM,       // Delete the temporary container.
		t.DockerBuildDirectory()} // Build in this directory.
	t.setState(BuildingImage)
	return process.NewProcess(DockerCommand, args...)
}

func (t *TestRunEntity) runDockerImage() *process.Process {
	defer t.deleteContainer()
	envVars := strings.Split(t.EnvironmentVariables, ",")
	args := []string{
		DockerRun,   // Run image
		DockerRunRM} // Delete after completed.
	for _, envVar := range envVars {
		args = append(args, fmt.Sprintf("-e %v", envVar)) // -e a='b'
	}
	args = append(args, DockerRunName)
	args = append(args, t.ImageName()) // Name of the container.
	args = append(args, t.ImageName()) // Name of the image to run
	return process.NewProcess(DockerCommand, args...)
}

func (t *TestRunEntity) deleteContainer() {
	argsStop := []string{DockerStop, t.ImageName()}
	argsDelete := []string{DockerRM, t.ImageName()}
	process.NewProcess(DockerCommand, argsStop...).Start().Wait()
	process.NewProcess(DockerCommand, argsDelete...).Start().Wait()
}

func (t *TestRunEntity) setState(state int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.state = state
}

func (t *TestRunEntity) getState() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	return t.state
}

func (t *TestRunEntity) awaitOutput(function func() *process.Process) error {
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

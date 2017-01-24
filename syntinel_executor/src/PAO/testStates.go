package PAO

import "log"

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
	NotFound:         "Not found.",
	Queued:           "Test awaiting execution.",
	Starting:         "Test is starting.",
	StartingDocker:   "Test Docker is initializing.",
	DockerStarted:    "Test Docker has started.",
	SendingScripts:   "Sending artifacts to the test Docker.",
	ScriptsSent:      "Artifacts have been sent to the test Docker.",
	ExecutingScripts: "Test is executing.",
	DestroyingDocker: "Test Docker is being torn down.",
	DockerDestroyed:  "Test Docker has been torn down.",
	Done:             "Done.",
}

func TestStateToString(state int) string {
	var stateString string
	var ok bool
	if stateString, ok = stateMap[state]; !ok {
		log.Fatalln("Bad test state passed to mapping function.")
	}
	return stateString
}

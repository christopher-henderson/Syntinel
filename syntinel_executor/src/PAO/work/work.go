package work

import (
  "syntinel_executor/process"
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
	cancel <-chan uint8
  currentProcess process.Process
}

func NewWork() chan<- uint8 {
  cancel := chan uint 8
	work := &Work{cancel}
  go work.Run()
  return cancel
}

func (*w Work) Run() {

}

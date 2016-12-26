package test

import "os/exec"

type Test struct {
	command         string
	process         exec.Cmd
	to_controller   chan<- bool
	from_controller <-chan bool
}

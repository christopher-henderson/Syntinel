package tests

import (
	"syntinel_executor/process"
	"testing"
	"time"
)

func TestProcess(t *testing.T) {
	command := "echo"
	args := "hello"
	proc, result, cancel := process.NewProcess(command, args)
	defer close(cancel)
	proc.Start()
	output := <-result
	if output.Output != args {
		t.Errorf("Expected %v, got %v", args, output.Output)
		t.Errorf("Error is %v", output.Err)
	}
}

func TestProcessKill(t *testing.T) {
	command := "python"
	args := []string{"-c", "while True: pass"}
	proc, result, cancel := process.NewProcess(command, args...)
	defer close(cancel)
	proc.Start()
	cancel <- 1
	output := <-result
	if output.Err == nil {
		t.Errorf("No error on SIGKILL")
	}
}

func TestProcessBadInvocation(t *testing.T) {
	command := "totallynotacommandecho"
	args := "hello"
	proc, result, cancel := process.NewProcess(command, args)
	defer close(cancel)
	proc.Start()
	output := <-result
	if output.Err == nil {
		t.Errorf("No error on bad invocation.")
	}
}

func TestProcessBadInvocationKill(t *testing.T) {
	command := "totallynotacommandecho"
	args := "hello"
	proc, result, cancel := process.NewProcess(command, args)
	defer close(cancel)
	proc.Start()
	cancel <- 1
	output := <-result
	if output.Err == nil {
		t.Errorf("No error on bad invocation.")
	}
}

// There's no assertion on this one. Rather it is intentionally trying to
// send a kill signal at the exact same time the process completes in the hopes
// it reveals a panic. Turns out that it is INCREDIBLY unlikely, if even possible.
//
// The typical results are one of either:
// 	process_test.go:84: <nil>
// 	process_test.go:85: Done.
//
// 	process_test.go:84: killed
// 	process_test.go:85: Done.
//
// 	process_test.go:84: killed
// 	process_test.go:85:
//
// All of which seem fine.
func TestProcessRace(t *testing.T) {
	command := "python"
	args := []string{"-c", "from time import sleep;sleep(.01);print('Done.')"}
	proc, result, cancel := process.NewProcess(command, args...)
	defer close(cancel)
	proc.Start()
	time.Sleep(time.Millisecond * 33)
	cancel <- 1
	output := <-result
	t.Log(output.Err)
	t.Log(output.Output)
}

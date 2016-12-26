package tests

import (
	"syntinel_executor/process"
	"testing"
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

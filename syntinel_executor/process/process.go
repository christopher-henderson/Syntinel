package process

import (
	"os"
	"os/exec"
	"strings"
)

type WorkResult struct {
	Err    error
	Output string
}

type Process struct {
	proc               *exec.Cmd
	resultMailbox      chan<- WorkResult
	cancellationSignal <-chan os.Signal
	done               chan WorkResult
}

func NewProcess(resultMailbox chan<- WorkResult, cancellationSignal <-chan os.Signal, command string, args ...string) Process {
	return Process{exec.Command(command, args...), resultMailbox, cancellationSignal, make(chan WorkResult, 1)}
}

func (p *Process) PID() int {
	return p.proc.ProcessState.Pid()
}

func (p *Process) Start() {
	go func() {
		output, err := p.proc.CombinedOutput()
		p.done <- WorkResult{err, strings.TrimSpace(string(output))}
	}()
	go func() {
		defer p.cleanup()
		var result WorkResult
		select {
		case result = <-p.done:
		case <-p.cancellationSignal:
			// Not portable to Windows.
			p.proc.Process.Kill()
			// Process termination will cause p.proc.Wait() to return,
			// so we have to empty the channel to avoid a leak.
			result = <-p.done
		}
		p.resultMailbox <- result
	}()
}

func (p *Process) cleanup() {
	close(p.done)
	close(p.resultMailbox)
}

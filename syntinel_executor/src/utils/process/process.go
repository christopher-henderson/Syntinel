package process

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"sync"
)

// The Process type lightly wraps the exec.Cmd type. Its intent is for a
// situations where you would want a long running background process that
// can also be cancelled at any time.
type Process struct {
	proc               *exec.Cmd
	resultMailbox      chan error
	cancellationSignal chan uint8
	done               chan error
	started            bool

	outputReader    *io.PipeReader
	outputWriter    *io.PipeWriter
	outputStreamSet bool

	// Access to completed MUST capture the lock.
	completed bool
	mutex     sync.RWMutex
}

// NewProcess returns a new process.(*Process). The returned process is not
// executing at this point. In order to begin the process, call the process'
// Start method.
func NewProcess(command string, args ...string) *Process {
	process := &Process{
		exec.Command(command, args...),
		make(chan error, 1),
		make(chan uint8, 1),
		make(chan error, 1),
		false,
		&io.PipeReader{},
		&io.PipeWriter{},
		false,
		false,
		sync.RWMutex{}}
	return process
}

// OutputStream Sets stdout and stderr and returns a *bufio.Scanner
// of the combined stdout and stderr.
//
// Calling this method twice will result will result in a panic.
// Calling this method after calling Start() will result in a panic.
func (p *Process) OutputStream() *bufio.Scanner {
	if p.outputStreamSet {
		panic("process.(*Process).OutputStream was called twice.")
	}
	if p.started {
		panic("process.(*Process).OutputStream called after the process was started..")
	}
	p.outputReader, p.outputWriter = io.Pipe()
	p.proc.Stderr = p.outputWriter
	p.proc.Stdout = p.outputWriter
	p.outputStreamSet = true
	return bufio.NewScanner(p.outputReader)
}

// Start execution of the process.
func (p *Process) Start() *Process {
	p.started = true
	if err := p.proc.Start(); err != nil {
		defer p.cleanup()
		p.resultMailbox <- err
	} else {
		go p.awaitOutput()
		go p.selectResultOrDie()
	}
	return p
}

// Wait synchronously waits for the output of the process.
func (p *Process) Wait() error {
	return <-p.resultMailbox
}

// Kill sends a kill signal to the process. If the process has already completed,
// then Kill is a no-op.
func (p *Process) Kill() {
	// Locks are necessary. If p.cleanup is the only funcion that both:
	//  1. closes p.cancellationSignal
	//  2. Sets p.complete
	// There, p.complete is a proxy for "p.cancellationSignal is closed"
	//
	// There is no error returned on this one since either way the desired effect
	// has occurred (either we killed or it's already dead).
	p.mutex.RLock()
	if !p.completed {
		p.cancellationSignal <- 1
	}
	p.mutex.RUnlock()
}

// awaitOutput waits for the process to complete, whether it be successfully or
// by failure. The combined output of the process' stdout and stderr, as well
// as any error, is used to construct a process.(*TestRunResult) which is then
// communicated to the selectResultOrDie method via the 'done' channel.
func (p *Process) awaitOutput() {
	p.done <- p.proc.Wait()
}

// Selects on either the result channel or the kill channel. If the result
// channel is selected then the select falls through and the (assumed to be)
// successful process.(*TestRunResult) is placed into the result channel. If the
// kill channel is selected, then the process is terminated and the result
// channel is read and the (assumed to be) failed process.(*TestRunResult) is
// placed into the result channel.
func (p *Process) selectResultOrDie() {
	defer p.cleanup()
	var result error
	select {
	case result = <-p.done:
	case <-p.cancellationSignal:
		// Not portable to Windows.
		// If the process completes the moment before this is called, then
		// the returned error will simply say "Process already finished".
		// Which, if it finished, then awaitOutput will be putting the result
		// in p.done soon anyways.
		if err := p.proc.Process.Kill(); err != nil {
			log.Println(err)
		}
		// Process termination will cause p.proc.CombinedOutput() to return.
		result = <-p.done
	}
	p.resultMailbox <- result
}

// Closes the channels that the process.Process is responsible for.
func (p *Process) cleanup() {
	p.mutex.Lock()
	p.completed = true
	p.mutex.Unlock()
	if p.outputStreamSet {
		p.outputReader.Close()
		p.outputWriter.Close()
	}
	close(p.done)
	close(p.resultMailbox)
	close(p.cancellationSignal)
}

package process

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"sync"
)

// TestRun is a PODO that holds the error, if any, as well as the output
// of a command executed via a process.Process execution.
type TestRunResult struct {
	Err    error
	Output string
}

// The Process type lightly wraps the exec.Cmd type. Its intent is for a
// situations where you would want a long running background process that
// can also be cancelled at any time.
type Process struct {
	proc               *exec.Cmd
	resultMailbox      chan error
	cancellationSignal chan uint8
	done               chan error

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
//
// A unidirectional, receive only, channel of process.(*TestRunResult) and a
// unidirectional, send only, channel of unint8 are returned.
//
// The first, receive only, channel is the result channel. That is, upon
// completion of the command (either due to success or failure) this channel
// will be populated with a process.(*TestRunResult) with information about the
// the process. This is a buffered channel of size one, and such MUST be
// received from only once as the process will close it immediately upon
// completion.
//
// The second, send only, channel is the kill signal channel. Sending any
// value over this channel will send a SIGKILL to the process. After the
// process has been killed, the result can again be found in the above
// result channel. This channel MUST be closed by the caller.
func NewProcess(command string, args ...string) *Process {
	resultMailbox := make(chan error, 1)
	cancellationSignal := make(chan uint8, 1)
	done := make(chan error)
	process := &Process{exec.Command(command, args...), resultMailbox, cancellationSignal, done, &io.PipeReader{}, &io.PipeWriter{}, false, false, sync.RWMutex{}}
	return process
}

// OutputStream Sets stdout and stderr and returns a *bufio.Scanner
// If the combined stdout and stderr.
func (p *Process) OutputStream() *bufio.Scanner {
	if !p.outputStreamSet {
		p.outputReader, p.outputWriter = io.Pipe()
		p.proc.Stderr = p.outputWriter
		p.proc.Stdout = p.outputWriter
		p.outputStreamSet = true
		return bufio.NewScanner(p.outputReader)
	}
	panic("asds")
}

// Start execution of the process.
func (p *Process) Start() {
	if err := p.proc.Start(); err != nil {
		defer p.cleanup()
		p.resultMailbox <- err
	} else {
		go p.awaitOutput()
		go p.selectResultOrDie()
	}
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

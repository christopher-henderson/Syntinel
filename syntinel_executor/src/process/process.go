package process

import (
	"log"
	"os/exec"
	"strings"
)

// The Process type lightly wraps the exec.Cmd type. Its intent is for a
// situations where you would want a long running background process that
// can also be cancelled at any time.
type Process struct {
	// Underlying process.
	proc *exec.Cmd
	// Channels used for sending the process results and receiving a kill signal.
	resultMailbox      chan<- *WorkResult
	cancellationSignal <-chan uint8
	// Internally used by awaitOutput to send results to selectResultOrDie
	done chan *WorkResult
	// Internally used to signal completion of the process.
	completed bool
}

// NewProcess returns a new process.(*Process). The returned process is not
// executing at this point. In order to begin the process, call the process'
// Start method.
//
// A unidirectional, receive only, channel of process.(*WorkResult) and a
// unidirectional, send only, channel of unint8 are returned.
//
// The first, receive only, channel is the result channel. That is, upon
// completion of the command (either due to success or failure) this channel
// will be populated with a process.(*WorkResult) with information about the
// the process. This is a buffered channel of size one, and such MUST be
// received from only once as the process will close it immediately upon
// completion.
//
// The second, send only, channel is the kill signal channel. Sending any
// value over this channel will send a SIGKILL to the process. After the
// process has been killed, the result can again be found in the above
// result channel. This channel MUST be closed by the caller.
func NewProcess(command string, args ...string) (*Process, <-chan *WorkResult, chan<- uint8) {
	resultMailbox := make(chan *WorkResult, 1)
	cancellationSignal := make(chan uint8, 1)
	done := make(chan *WorkResult)
	process := &Process{exec.Command(command, args...), resultMailbox, cancellationSignal, done, false}
	return process, resultMailbox, cancellationSignal
}

// Start execution of the process.
func (p *Process) Start() {
	go p.awaitOutput()
	go p.selectResultOrDie()
	for p.proc.Process == nil && !p.completed {
		// Wait for either the process to spin up successfully or for it
		// to instantly die for whatever reason (usually a bad invocation).
		//
		// If p.proc.Process is not nil, then the process is up and running.
		// If p.completed then it instantly died somehow (probably).
		//
		// This is crucial for protecting against panics on illegal access to
		// the process. The goroutines will fall through which means without this,
		// then it is possible for the caller to cause a panic by immediately
		// killing the process. If they do, then:
		//
		// 1. If the process has not started up then p.proc.Process will be nil,
		//		causin a panic in p.selectResultOrDie()
		// 2. If the process has already died then attempts to kill it will also
		//		panic.
	}
}

// awaitOutput waits for the process to complete, whether it be successfully or
// by failure. The combined output of the process' stdout and stderr, as well
// as any error, is used to construct a process.(*WorkResult) which is then
// communicated to the selectResultOrDie method via the 'done' channel.
func (p *Process) awaitOutput() {
	output, err := p.proc.CombinedOutput()
	log.Println("Received output from process.")
	p.done <- &WorkResult{err, strings.TrimSpace(string(output))}
}

// Selects on either the result channel or the kill channel. If the result
// channel is selected then the select falls through and the (assumed to be)
// successful process.(*WorkResult) is placed into the result channel. If the
// kill channel is selected, then the process is terminated and the result
// channel is read and the (assumed to be) failed process.(*WorkResult) is
// placed into the result channel.
func (p *Process) selectResultOrDie() {
	defer p.cleanup()
	var result *WorkResult
	select {
	case result = <-p.done:
	case <-p.cancellationSignal:
		log.Println("Received kill signal.")
		// Not portable to Windows.
		p.proc.Process.Kill()
		log.Println("Process killed.")
		// Process termination will cause p.proc.CombinedOutput() to return.
		result = <-p.done
	}
	p.resultMailbox <- result
	log.Println("Process completed.")
	p.completed = true
}

// Closes the channels that the process.Process is responsible for.
func (p *Process) cleanup() {
	close(p.done)
	close(p.resultMailbox)
}

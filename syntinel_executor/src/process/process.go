package process

import (
	"os/exec"
	"strings"
	"time"
)

// WorkResult is a PODO that holds the error, if any, as well as the output
// of a command executed via a process.Process execution.
type WorkResult struct {
	Err    error
	Output string
}

// The Process type lightly wraps the exec.Cmd type. Its intent is for a
// situations where you would want a long running background process that
// can also be cancelled at any time.
type Process struct {
	proc               *exec.Cmd
	resultMailbox      chan<- *WorkResult
	cancellationSignal <-chan uint8
	done               chan *WorkResult
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
	process := &Process{exec.Command(command, args...), resultMailbox, cancellationSignal, done}
	return process, resultMailbox, cancellationSignal
}

// Start execution of the process.
func (p *Process) Start() {
	go p.awaitOutput()
	go p.selectResultOrDie()
}

// awaitOutput waits for the process to complete, whether it be successfully or
// by failure. The combined output of the process' stdout and stderr, as well
// as any error, is used to construct a process.(*WorkResult) which is then
// communicated to the selectResultOrDie method via the 'done' channel.
func (p *Process) awaitOutput() {
	output, err := p.proc.CombinedOutput()
	p.done <- &WorkResult{err, strings.TrimSpace(string(output))}
}

// Selects on either the result channel or the kill channel. If the result
// channel is selected then the select falls through and the (assumed to be)
// successful process.(*WorkResult) is placed into the result channel. If the
// kill channel is selected, then the process is terminated and the result
// channel is read and the (assumed to be) failed process.(*WorkResult) is
// placed into the result channel.
func (p *Process) selectResultOrDie() {
	time.Sleep(time.Second * 1)
	// @TODO Make more rigorous.
	//
	// If this sleep looks like code smell, well that's because it reeks. In testing,
	// if the process is killed IMMEDIATELY after starting, then we can get a
	// panic while in os.(*Process).Kill. Now, killing a proc immediately after
	// beginning it is going to be incredibly rare, but a panic is unacceptable.
	// This one second sleep stopped any panics in testing.
	//
	// Since this sleep is at the beginning of this long running goroutine,
	// it is uncommon for requests to be delayed. The only people who will
	// notice are those whose command finishes in under one second or those who
	// try to kill the process in less than one second after starting it.
	defer p.cleanup()
	var result *WorkResult
	select {
	case result = <-p.done:
	case <-p.cancellationSignal:
		// Not portable to Windows.
		p.proc.Process.Kill()
		// Process termination will cause p.proc.CombinedOutput() to return.
		result = <-p.done
	}
	p.resultMailbox <- result
}

// Closes the channels that the process.Process is responsible for.
func (p *Process) cleanup() {
	close(p.done)
	close(p.resultMailbox)
}

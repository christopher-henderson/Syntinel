package process

import (
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
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
	resultMailbox      chan *TestRunResult
	cancellationSignal chan uint8
	done               chan *TestRunResult
	completed          bool
	sync.RWMutex
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
	resultMailbox := make(chan *TestRunResult, 1)
	cancellationSignal := make(chan uint8, 1)
	done := make(chan *TestRunResult)
	process := &Process{exec.Command(command, args...), resultMailbox, cancellationSignal, done, false, sync.RWMutex{}}
	return process
}

// Start execution of the process.
func (p *Process) Start() {
	stdout, _ := p.proc.StdoutPipe()
	if err := p.proc.Start(); err != nil {
		defer p.cleanup()
		p.Lock()
		p.completed = true
		p.Unlock()
		p.resultMailbox <- &TestRunResult{err, ""}
		return
	}
	go p.awaitOutput(stdout)
	go p.selectResultOrDie()
}

func (p *Process) Wait() *TestRunResult {
	return <-p.resultMailbox
}

func (p *Process) Kill() {
	p.RLock()
	if !p.completed {
		p.RUnlock()
		p.cancellationSignal <- 1
	} else {
		p.RUnlock()
	}
}

// awaitOutput waits for the process to complete, whether it be successfully or
// by failure. The combined output of the process' stdout and stderr, as well
// as any error, is used to construct a process.(*TestRunResult) which is then
// communicated to the selectResultOrDie method via the 'done' channel.
func (p *Process) awaitOutput(stdout io.ReadCloser) {
	output, _ := ioutil.ReadAll(stdout)
	err := p.proc.Wait()
	p.done <- &TestRunResult{err, strings.TrimSpace(string(output))}
}

// Selects on either the result channel or the kill channel. If the result
// channel is selected then the select falls through and the (assumed to be)
// successful process.(*TestRunResult) is placed into the result channel. If the
// kill channel is selected, then the process is terminated and the result
// channel is read and the (assumed to be) failed process.(*TestRunResult) is
// placed into the result channel.
func (p *Process) selectResultOrDie() {
	defer p.cleanup()
	var result *TestRunResult
	select {
	case result = <-p.done:
	case <-p.cancellationSignal:
		// Not portable to Windows.
		if err := p.proc.Process.Kill(); err != nil {
			log.Println(err)
		}
		// Process termination will cause p.proc.CombinedOutput() to return.
		result = <-p.done
	}
	p.resultMailbox <- result
	// This assignment MUST come after placing the result in the mailbox.
	// It is the guarantee
	// p.completed = true
}

// Closes the channels that the process.Process is responsible for.
func (p *Process) cleanup() {
	close(p.done)
	close(p.resultMailbox)
	close(p.cancellationSignal)
}

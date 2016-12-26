package process

// WorkResult is a PODO that holds the error, if any, as well as the output
// of a command executed via a process.Process execution.
type WorkResult struct {
	Err    error
	Output string
}

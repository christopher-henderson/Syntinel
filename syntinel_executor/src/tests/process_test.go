package tests

import (
	"bufio"
	"log"
	"math"
	"runtime"
	"syntinel_executor/PAO/process"
	"testing"
	"time"
	"unsafe"
)

func LogScanner(scanner *bufio.Scanner) {
	for scanner.Scan() {
		log.Print(string(scanner.Bytes()))
	}
	log.Print("\n")
}

func TestProcess(t *testing.T) {
	command := "echo"
	args := "hello"
	proc := process.NewProcess(command, args)
	stdout := proc.OutputStream()
	proc.Start()
	var stdoutResult []byte
	go func() {
		for stdout.Scan() {
			bytes := stdout.Bytes()
			stdoutResult = append(stdoutResult, bytes...)
		}
		if string(stdoutResult) != args {
			t.Errorf("Stdout is wrong, expected %v got %v", args, string(stdoutResult))
		}
	}()
	if err := proc.Wait(); err != nil {
		t.Errorf("Unexpected error occured: %v", err)
	}
}

func TestProcessKill(t *testing.T) {
	command := "python"
	args := []string{"-c", "while True: pass"}
	proc := process.NewProcess(command, args...)
	go LogScanner(proc.OutputStream())
	proc.Start()
	proc.Kill()
	if err := proc.Wait(); err == nil {
		t.Errorf("No error on SIGKILL")
	}
}

func TestProcessBadInvocation(t *testing.T) {
	command := "totallynotacommandecho"
	args := "hello"
	proc := process.NewProcess(command, args)
	go LogScanner(proc.OutputStream())
	proc.Start()
	proc.Kill()
	if err := proc.Wait(); err == nil {
		t.Errorf("No error on bad invocation.")
	}
}

func TestRetrieveStderr(t *testing.T) {
	command := "python"
	args := []string{"-c", "from sys import stderr;print('Printing to stderr.', file=stderr);raise Exception('also this')"}
	proc := process.NewProcess(command, args...)
	go LogScanner(proc.OutputStream())
	proc.Start()
	proc.Wait()
}

func TestProcessBadInvocationKill(t *testing.T) {
	command := "totallynotacommandecho"
	args := "hello"
	proc := process.NewProcess(command, args)
	go LogScanner(proc.OutputStream())
	proc.Start()
	proc.Kill()
	if err := proc.Wait(); err == nil {
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
	proc := process.NewProcess(command, args...)
	proc.Start()
	time.Sleep(time.Millisecond * 33)
	proc.Kill()
	proc.Wait()
}

func TestProcessGoroutineLeak(t *testing.T) {
	startingGoroutines := runtime.NumGoroutine()
	command := "echo"
	args := "hello"
	proc := process.NewProcess(command, args)
	proc.Start()
	proc.Wait()
	// Give the goroutines a moment to shutdown. If you don't, then every
	// so often you'll see a 2 != 3 error.
	time.Sleep(time.Millisecond * 100)
	endingGoroutines := runtime.NumGoroutine()
	if endingGoroutines != startingGoroutines {
		t.Errorf("Possible Goroutine leak. Started test with %v goroutines, ended with %v", startingGoroutines, endingGoroutines)
	}
}

func TestProcessMemoryLeak(t *testing.T) {
	m1 := &runtime.MemStats{}
	runtime.ReadMemStats(m1)
	command := "echo"
	args := "hello"
	proc := process.NewProcess(command, args)
	proc.Start()
	proc.Wait()
	runtime.GC()
	m2 := &runtime.MemStats{}
	runtime.ReadMemStats(m2)
	frees := m2.Frees - m1.Frees
	mallocs := m2.Mallocs - m1.Mallocs
	t.Logf("Mallocs: %v, Frees: %v", mallocs, frees)
	if mallocs > frees {
		t.Errorf("Possible memory leak. In the course of this test there %v mallocs, but only %v frees", mallocs, frees)
	}
}

// Let's try to see if memory leaks occur. In the end, memory leaks
// tend to be hunted down over time, rather than in unit testing.
func TestProcessMemoryLeak2(t *testing.T) {
	command := "echo"
	args := "hello"

	proc := process.NewProcess(command, args)
	proc.Start()
	proc.Wait()
	runtime.GC()
	m1 := &runtime.MemStats{}
	runtime.ReadMemStats(m1)

	proc = process.NewProcess(command, args)
	proc.Start()
	proc.Wait()
	runtime.GC()
	m2 := &runtime.MemStats{}
	runtime.ReadMemStats(m2)

	proc = process.NewProcess(command, args)
	proc.Start()
	proc.Wait()
	runtime.GC()
	m3 := &runtime.MemStats{}
	runtime.ReadMemStats(m3)

	m1Allocs := int64(m1.Alloc) - int64(unsafe.Sizeof(m1))
	m2Allocs := int64(m2.Alloc) - int64(unsafe.Sizeof(m2))
	m3Allocs := int64(m3.Alloc) - int64(unsafe.Sizeof(m3))
	m1m2Difference := math.Abs(float64(m1Allocs - m2Allocs))
	m1m3Difference := math.Abs(float64(m1Allocs - m3Allocs))
	m2m3Difference := math.Abs(float64(m2Allocs - m3Allocs))
	// This is stating a tolernace of 3kB to allow for variance in the
	// runtime. This may or may not be reasonable.
	tolerance := 3000.0
	if m1m2Difference > tolerance || m1m3Difference > tolerance || m2m3Difference > tolerance {
		t.Errorf("Possible memory leak. After running three times, the Allocs are (in bytes)\n"+
			"1: %v\n"+
			"2: %v\n"+
			"3: %v\n", m1Allocs, m2Allocs, m3Allocs)
		t.Errorf("The differences are %v, %v, and %v", m1m2Difference, m1m3Difference, m2m3Difference)
	}
}

func TestMemoryLeak3(t *testing.T) {
	command := "echo"
	args := "hello"
	proc := process.NewProcess(command, args)
	proc.Start()
	proc.Wait()
	startingMemory := &runtime.MemStats{}
	runtime.ReadMemStats(startingMemory)
	for i := 0; i < 100; i++ {
		proc = process.NewProcess(command, args)
		proc.Start()
		proc.Wait()
	}
	runtime.GC()
	endingMemory := &runtime.MemStats{}
	runtime.ReadMemStats(endingMemory)
	t.Log(startingMemory.Alloc)
	t.Log(endingMemory.Alloc)
}

package main

import (
	"fmt"
	"os"
	"syntinel_executor/process"
	"syscall"
)

// func ls(args string) {
// 	c := exec.CommandContext(context.Background(), "ls", args)
// 	out, err := c.CombinedOutput()
// 	if err != nil {
// 		fmt.Println(err)
// 		fmt.Println(string(out))
// 		panic(1)
// 	}
// 	fmt.Println(string(out))
// }

// func timer(out chan<- bool) {
// 	time.Sleep(time.Second * 10)
// 	out <- true
// }

func main() {
	result := make(chan process.WorkResult, 1)
	killSignal := make(chan os.Signal, 1)
	command := "python"
	arg := "/tmp/lol.py"
	proc := process.NewProcess(result, killSignal, command, arg)
	proc.Start()
	killSignal <- syscall.SIGKILL
	output := <-result
	fmt.Println("Out is", output.Output)
	fmt.Println("Err is", output.Err)
}

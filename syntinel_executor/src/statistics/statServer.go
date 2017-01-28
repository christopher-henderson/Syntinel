package statistics

import (
	"fmt"
	"log"
	"runtime"
	"syntinel_executor/utils/process"
)

func getSystemDependentTop() (Top, error) {
	switch runtime.GOOS {
	case "darwin":
		return NewDarwinTop(), nil
	case "linux":
		return NewLinuxTop(), nil
	default:
		return nil, fmt.Errorf("Unsupported OS: %v", runtime.GOOS)
	}
}

func Start() {
	top, err := getSystemDependentTop()
	if err != nil {
		log.Fatalln(err)
	}
	proc := process.NewProcess(top.Command(), top.Args()...)
	out := proc.OutputStream()
	proc.Start()
	top.Parse(out)
}

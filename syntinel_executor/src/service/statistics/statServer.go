package statistics

import (
	"log"
	"reflect"
	"runtime"
	"syntinel_executor/utils/process"
)

var top = getSystemDependentTop()

func getSystemDependentTop() Top {
	switch runtime.GOOS {
	case "darwin":
		return NewDarwinTop()
	case "linux":
		return NewLinuxTop()
	default:
		log.Fatalln("Unsupported OS: %v", runtime.GOOS)
		return nil
	}
}

func StartStatisticsServer() {
	log.Println("Starting statistics server using: " + reflect.TypeOf(top).Name())
	proc := process.NewProcess(top.Command(), top.Args()...)
	out := proc.OutputStream()
	proc.Start()
	go top.Parse(out)
}

func Statistics() Stats {
	return top.Statistics()
}

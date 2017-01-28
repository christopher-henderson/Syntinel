package statistics

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type LinuxTop struct {
	// 9
	command string
	args    []string
	CPU     float64 `json:"cpu"`
	MEM     float64 `json:"memory"`
	Swap    float64 `json:"swap"`
}

func NewLinuxTop() LinuxTop {
	return LinuxTop{
		"top",
		[]string{"-p", strconv.Itoa(os.Getppid()), "-b"},
		0.0, 0.0, 0.0}
}

func (lt LinuxTop) Command() string {
	return lt.command
}

func (lt LinuxTop) Args() []string {
	return lt.args
}

func (lt LinuxTop) Parse(out *bufio.Scanner) {
	for {
		for line := 0; line < 9; line++ {
			if !out.Scan() {
				panic("This should never have happened")
			}
			lt.Dispatch(line, out.Text())
		}
	}
}

func (lt LinuxTop) AverageLoad(line string) {

}

func (lt LinuxTop) Tasks(line string) {

}

func (lt LinuxTop) CPUUsage(line string) {
	pattern := regexp.MustCompile(".*(\\d+\\.\\d+) sy.*")
	matches := pattern.FindSubmatch([]byte(line))
	fmt.Println(string(matches[1]))
}

func (lt LinuxTop) Dispatch(line int, output string) {
	switch line {
	case 0:
		lt.AverageLoad(output)
	case 1:
		lt.Tasks(output)
	case 2:
		lt.CPUUsage(output)
	case 3:
		// lt.Mem(output)
	case 4:
		// lt.Swap(output)
	case 5:
		// Blank line
	case 6:
		// PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND
	case 7:
		// lt.ThisProcess(output)
	case 8:
		// Blank line
	default:
		panic("failed to understand output")
	}
}

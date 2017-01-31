package statistics

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
)

// Example output on Linux.
//
// top - 09:41:59 up 3 days, 41 min,  2 users,  load average: 0.00, 0.01, 0.23
// Tasks: 198 total,   1 running, 197 sleeping,   0 stopped,   0 zombie
// %Cpu(s):  0.0 us,  0.4 sy,  0.7 ni, 98.0 id,  0.9 wa,  0.0 hi,  0.0 si,  0.0 st
// KiB Mem : 16162360 total, 14701408 free,   709896 used,   751056 buff/cache
// KiB Swap:  8191996 total,  8191996 free,        0 used. 15072028 avail Mem

// top - 10:50:30 up  1:49,  2 users,  load average: 0.00, 0.01, 0.05
var averageLoadRegex = regexp.MustCompile("(\\d+\\.\\d+)")

// Tasks: 198 total,   1 running, 197 sleeping,   0 stopped,   0 zombie
var tasksRegex = regexp.MustCompile("(\\d+)")

// %Cpu(s):  0.0 us,  0.0 sy,  0.0 ni,100.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
var cpuUsageRegex = regexp.MustCompile("(\\d+\\.\\d+)")

// KiB Mem : 16162360 total, 15411480 free,   215844 used,   535036 buff/cache
var memUsageRegex = regexp.MustCompile("\\d+")

var swapUsageRegex = regexp.MustCompile("\\d+")

// PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND
// 1 root      20   0  193768   6872   3952 S   0.0  0.0   0:01.28 systemd
var executorUsageRegex = regexp.MustCompile("(\\d+\\.\\d+)\\s+(\\d+\\.\\d+)\\s+(\\d+:\\d+\\.\\d+)")

type LinuxTop struct {
	command string
	args    []string

	// mutex protects stats
	mutex *sync.RWMutex
	stats *Stats
}

func NewLinuxTop() LinuxTop {
	lt := LinuxTop{
		"top",
		[]string{"-p", strconv.Itoa(os.Getppid()), "-b"},
		&sync.RWMutex{},
		&Stats{},
	}
	lt.stats.Executor.PID = os.Getppid()
	return lt
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
				panic("The statistics server has failed.")
			}
			lt.Dispatch(line, out.Text())
		}
	}
}

func (lt LinuxTop) Statistics() Stats {
	lt.mutex.RLock()
	defer lt.mutex.RUnlock()
	// Return a copy of the current stats.
	return *lt.stats
}

func (lt *LinuxTop) AverageLoad(line string) {
	matches := averageLoadRegex.FindAllStringSubmatch(line, -1)
	oneMin := matches[0][0]
	fiveMin := matches[1][0]
	fifteenMin := matches[2][0]
	if err := lt.parseAndAssignFloat64(oneMin, &lt.stats.CPU.Load1); err != nil {
		log.Fatalln("One minute average load failed to parse.")
	}
	if err := lt.parseAndAssignFloat64(fiveMin, &lt.stats.CPU.Load5); err != nil {
		log.Fatalln("Five minute average load failed to parse.")
	}
	if err := lt.parseAndAssignFloat64(fifteenMin, &lt.stats.CPU.Load15); err != nil {
		log.Fatalln("Fifteen minute average load failed to parse.")
	}
}

func (lt *LinuxTop) Tasks(line string) {
	// @TODO
}

// 0.0 us,  0.0 sy,100.0 ni,  0.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st]
func (lt *LinuxTop) CPUUsage(line string) {
	matches := cpuUsageRegex.FindAllStringSubmatch(line, -1)
	percentActive := matches[2][0]
	if err := lt.parseAndAssignFloat64(percentActive, &lt.stats.CPU.Active); err != nil {
		log.Fatalln("CPU Usage failed to parse.")
	}
}

// KiB Mem : 16162360 total, 15411480 free,   215844 used,   535036 buff/cache
func (lt *LinuxTop) MemUsage(line string) {
	matches := memUsageRegex.FindAllStringSubmatch(line, -1)
	total := matches[0][0]
	free := matches[1][0]
	used := matches[2][0]
	cache := matches[3][0]
	if err := lt.parseAndAssignInt64(total, &lt.stats.Mem.Total); err != nil {
		log.Fatalln("Total memory load failed to parse.")
	}
	if err := lt.parseAndAssignInt64(free, &lt.stats.Mem.Free); err != nil {
		log.Fatalln("Free memory load failed to parse.")
	}
	if err := lt.parseAndAssignInt64(used, &lt.stats.Mem.Used); err != nil {
		log.Fatalln("Used memory load failed to parse.")
	}
	if err := lt.parseAndAssignInt64(cache, &lt.stats.Mem.Cache); err != nil {
		log.Fatalln("Cache memory load failed to parse.")
	}
}

func (lt *LinuxTop) ExecutorUsage(line string) {
	// @TODO this is capturing sudo and not the executor
	matches := executorUsageRegex.FindAllStringSubmatch(line, -1)
	cpu := matches[0][1]
	mem := matches[0][2]
	uptime := matches[0][3]
	if err := lt.parseAndAssignFloat64(cpu, &lt.stats.Executor.CPU); err != nil {
		log.Fatalln("CPU Usage failed to parse.")
	}
	if err := lt.parseAndAssignFloat64(mem, &lt.stats.Executor.Mem); err != nil {
		log.Fatalln("CPU Usage failed to parse.")
	}
	lt.mutex.Lock()
	lt.stats.Executor.Uptime = uptime
	lt.mutex.Unlock()
}

func (lt *LinuxTop) Dispatch(line int, output string) {
	switch line {
	case 0:
		lt.AverageLoad(output)
	case 1:
		// lt.Tasks(output)
	case 2:
		lt.CPUUsage(output)
	case 3:
		lt.MemUsage(output)
	case 4:
		// lt.Swap(output)
	case 5:
		// Blank line
	case 6:
		// PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND
	case 7:
		// lt.ExecutorUsage(output)
	case 8:
		// Blank line
	default:
		panic("failed to understand output")
	}
}

func (lt *LinuxTop) parseAndAssignFloat64(source string, target *float64) error {
	value, err := strconv.ParseFloat(source, 64)
	if err != nil {
		return err
	}
	lt.mutex.Lock()
	defer lt.mutex.Unlock()
	*target = value
	return nil
}

func (lt *LinuxTop) parseAndAssignInt64(source string, target *int64) error {
	value, err := strconv.ParseInt(source, 10, 64)
	if err != nil {
		return err
	}
	lt.mutex.Lock()
	defer lt.mutex.Unlock()
	*target = value
	return nil
}

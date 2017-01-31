package statistics

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
)

// top - 10:50:30 up  1:49,  2 users,  load average: 0.00, 0.01, 0.05
var averageLoadRegex = regexp.MustCompile("(\\d+\\.\\d+)")

// %Cpu(s):  0.0 us,  0.0 sy,  0.0 ni,100.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
// the id column is "Idle". We can use that to get that to get current usage.
var cpuUsageRegex = regexp.MustCompile("(\\d+\\.\\d+)")

// KiB Mem : 16162360 total, 15411480 free,   215844 used,   535036 buff/cache
var memUsageRegex = regexp.MustCompile("\\d+")

// PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND
// 1 root      20   0  193768   6872   3952 S   0.0  0.0   0:01.28 systemd
// Capturing user, %CPU, %MEM, and Time
var executorUsageRegex = regexp.MustCompile("\\d+\\s+(\\D+)\\s+.*(\\d+\\.\\d+).*(\\d+\\.\\d+).*(\\d+:\\d+\\.\\d+).*")

type LinuxTop struct {
	// 9
	command string
	args    []string
	mutex   sync.RWMutex
	stats   Stats
}

func NewLinuxTop() LinuxTop {
	return LinuxTop{
		"top",
		[]string{"-p", strconv.Itoa(os.Getppid()), "-b"},
		sync.RWMutex{},
		NewStats(),
	}
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

func (lt LinuxTop) Statistics() *Stats {
	lt.mutex.RLock()
	defer lt.mutex.RUnlock()
	log.Println(lt.stats.CPU.Active)
	return &lt.stats
}

func (lt *LinuxTop) AverageLoad(line string) {
	matches := averageLoadRegex.FindAllStringSubmatch(line)
	oneMin := matches[1][0]
	fiveMin := matches[2][0]
	fifteenMin := matches[3][0]
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
	matches := memUsageRegex.FindAllStringSubmatch(line)
	total := matches[1][0]
	free := matches[2][0]
	used := matches[3][0]
	cache := matches[4][0]
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

func (lt *LinuxTop) Dispatch(line int, output string) {
	switch line {
	case 0:
		lt.AverageLoad(output)
	case 1:
		lt.Tasks(output)
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
		// lt.ThisProcess(output)
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

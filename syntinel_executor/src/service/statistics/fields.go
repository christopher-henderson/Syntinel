package statistics

// top - 09:41:59 up 3 days, 41 min,  2 users,  load average: 0.00, 0.01, 0.23
// Tasks: 198 total,   1 running, 197 sleeping,   0 stopped,   0 zombie
// %Cpu(s):  0.0 us,  0.4 sy,  0.7 ni, 98.0 id,  0.9 wa,  0.0 hi,  0.0 si,  0.0 st
// KiB Mem : 16162360 total, 14701408 free,   709896 used,   751056 buff/cache
// KiB Swap:  8191996 total,  8191996 free,        0 used. 15072028 avail Mem
//
// PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND
// 	1 root      20   0  193768   6872   3952 S   0.0  0.0   0:04.17 systemd

type CPUStats struct {
	Active float64 `json:"active"`
	Load1  float64 `json:"load5minute"`
	Load5  float64 `json:"load10minute"`
	Load15 float64 `json:"load15minute"`
}

type MemStats struct {
	Total int64 `json:"sys"`
	Free  int64 `json:"free"`
	Used  int64 `json:"used"`
	Cache int64 `json:"cache"`
}

type SwapStats struct {
	Total int64 `json:"total"`
	Free  int64 `json:"free"`
	Used  int64 `json:"used"`
}

type ExecutorStats struct {
	PID    int     `json:"PID"`
	Uptime string  `json:"uptime"`
	CPU    float64 `json:"CPUPercentage"`
	Mem    float64 `json:"MemPercentage"`
}

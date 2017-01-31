package statistics

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
	PID           int32   `json:"PID"`
	User          string  `json:"user"`
	Uptime        string  `json:"uptime"`
	CPUPercentage float64 `json:"CPUPercentage"`
	MemPercentage float64 `json:"MemPercentage"`
}

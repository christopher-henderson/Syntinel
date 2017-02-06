package statistics

type Stats struct {
	CPU      CPUStats      `json:"cpu"`
	Mem      MemStats      `json:"memory"`
	Swap     SwapStats     `json:"swap"`
	Executor ExecutorStats `json:"executor"`
}

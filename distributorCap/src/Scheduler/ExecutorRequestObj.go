package Scheduler

type ExecutorRequestObj struct {
	ID                   int    `json:"id"`
	Dockerfile           string `json:"dockerfile"`
	Script               string `json:"script"`
	EnvironmentVariables string `json:"environmentVairables"`
}

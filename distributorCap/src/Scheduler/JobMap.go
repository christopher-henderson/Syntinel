package Scheduler

import (
	"sync"
)

type JobMap struct {
	m map[int]Job
	sync.RWMutex
}

func NewJobMap() *JobMap {
	return &JobMap{make(map[int]Job), sync.RWMutex{}}
}

func (j *JobMap) Get(key int) Job {
	j.RLock()
	defer j.RUnlock()
	return j.m[key]
}

func (j *JobMap) Put(key int, jo Job) {
	j.Lock()
	j.m[key] = jo
	j.Unlock()
}

func (j *JobMap) Delete(key int) {
	j.Lock()
	defer j.Unlock()
	delete(j.m, key)
}

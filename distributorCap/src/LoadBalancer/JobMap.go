package main

import "sync"

type JobMap struct {
	m map[int]job
	sync.RWMutex
}

func NewJobMap() *JobMap {
	return &JobMap{make(map[int]job), sync.RWMutex{}}
}

func (j *JobMap) Get(key int) job {
	j.RLock()
	defer j.RUnlock()
	return j.m[key]
}

func (j *JobMap) Put(key int, jo job) {
	j.Lock()
	j.m[key] = jo
	j.Unlock()
}

func (j *JobMap) Delete(key int) {
	j.Lock()
	defer j.Unlock()
	delete(j.m, key)
}

package PAO

import (
	"log"
	"sync"
	"syntinel_executor/PAO/work"
)

type TestQueue struct {
	dockerID    int
	scriptID    int
	queue       chan *work.Work
	running     bool
	currentWork *work.Work
	mutex       sync.Mutex
}

func NewTestQueue(dockerID int, scriptID int) *TestQueue {
	t := &TestQueue{dockerID, scriptID, make(chan *work.Work), false, nil, sync.Mutex{}}
	go t.consume()
	return t
}

func (t *TestQueue) consume() {
	for test := range t.queue {
		t.execute(test)
	}
}

func (t *TestQueue) execute(test *work.Work) {
	t.mutex.Lock()
	t.currentWork = test
	t.running = true
	t.mutex.Unlock()
	defer t.teardown()
	test.Run()
}

func (t *TestQueue) teardown() {
	if err := recover(); err != nil {
		log.Println(err)
	}
	t.mutex.Lock()
	close(t.currentWork.Cancel)
	t.running = false
	t.currentWork = nil
	t.mutex.Unlock()
}

func (t *TestQueue) Run() {
	t.queue <- work.NewWork()
}

func (t *TestQueue) Kill() {
	t.mutex.Lock()
	if t.running {
		t.currentWork.Cancel <- 1
	}
	t.mutex.Unlock()
}

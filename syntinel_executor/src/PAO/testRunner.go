package PAO

import (
	"log"
	"strconv"
	"sync"
	"syntinel_executor/DAO"
	"syntinel_executor/DAO/database"
)

type TestRunner struct {
	testID      int
	queue       chan int
	currentWork *database.TestRunEntity
	mutex       sync.RWMutex
}

func NewTestRunner(testID int) *TestRunner {
	t := &TestRunner{testID, make(chan int), nil, sync.RWMutex{}}
	go t.consume()
	return t
}

func (t *TestRunner) Run(testRunID int) {
	log.Println("Got test ID: " + strconv.Itoa(testRunID))
	t.queue <- testRunID
}

func (t *TestRunner) Kill(id int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.currentWork != nil && t.currentWork.ID == id {
		t.currentWork.Cancel <- 1
	}
}

func (t *TestRunner) Query(id int) int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if t.currentWork != nil && t.currentWork.ID == id {
		return t.currentWork.Query()
	}
	return 0
}

func (t *TestRunner) consume() {
	for test := range t.queue {
		t.execute(test)
	}
}

func (t *TestRunner) execute(id int) {
	log.Println("Executing ID: " + strconv.Itoa(id))
	test, err := DAO.TestRun.Get(id)
	if err != nil {
		log.Println(err)
		return
	}
	t.mutex.Lock()
	t.currentWork = test
	t.mutex.Unlock()
	defer t.teardown()
	test.Run()
}

func (t *TestRunner) teardown() {
	if err := recover(); err != nil {
		// This MUST be recovered from or else the t.consumer goroutine WILL go down.
		log.Println(err)
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if err := database.DeleteTestRun(t.currentWork.ID); err != nil {
		log.Println(err)
	}
	close(t.currentWork.Cancel)
	t.currentWork = nil
}

package PAO

import (
	"log"
	"sync"
	"syntinel_executor/DAO"
)

type ThreadSafeMapIntTestRun struct {
	m map[int]*TestRun
	sync.RWMutex
}

func (m *ThreadSafeMapIntTestRun) GetTestRun(id int) (*TestRun, bool) {
	m.RLock()
	defer m.RUnlock()
	testRun, ok := m.m[id]
	return testRun, ok
}

func (m *ThreadSafeMapIntTestRun) SetTestRun(id int, testRun *TestRun) {
	m.Lock()
	defer m.Unlock()
	m.m[id] = testRun
}

func (m *ThreadSafeMapIntTestRun) DeleteTestRun(id int) {
	m.Lock()
	defer m.Unlock()
	delete(m.m, id)
}

type TestRunner struct {
	testID      int
	queue       chan int
	currentWork *TestRun
	mutex       sync.RWMutex
	testRunMap  ThreadSafeMapIntTestRun
}

func NewTestRunner(testID int) *TestRunner {
	t := &TestRunner{testID, make(chan int), nil, sync.RWMutex{}, ThreadSafeMapIntTestRun{make(map[int]*TestRun), sync.RWMutex{}}}
	go t.consume()
	return t
}

func (t *TestRunner) Run(testRunID int) {
	if _, ok := t.testRunMap.GetTestRun(testRunID); ok {
		log.Println("Received attempt to execute test that is already queued.")
		return
	}
	if test, ok := DAO.GetTest(t.testID); !ok {
		log.Println("Request for non-existent test run.")
	} else {
		t.testRunMap.SetTestRun(testRunID, NewTestRun(t.testID, testRunID, test.DockerPath, test.ScriptPath))
		t.queue <- testRunID
	}
}

func (t *TestRunner) Kill(id int) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if t.currentWork != nil && t.currentWork.ID == id {
		t.currentWork.Cancel <- 1
	} else {
		t.testRunMap.DeleteTestRun(id)
	}
}

func (t *TestRunner) Query(testRunID int) int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if test, ok := t.testRunMap.GetTestRun(testRunID); ok {
		return test.Query()
	}
	return NotFound
}

func (t *TestRunner) consume() {
	for test := range t.queue {
		t.execute(test)
	}
}

func (t *TestRunner) execute(id int) {
	test, ok := t.testRunMap.GetTestRun(id)
	if !ok {
		log.Println("Ticket for expired test received.")
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
	t.testRunMap.DeleteTestRun(t.currentWork.ID)
	close(t.currentWork.Cancel)
	t.currentWork = nil
}

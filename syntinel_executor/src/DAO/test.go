package DAO

import "sync"

type Test struct {
	ID         int
	DockerPath string
	ScriptPath string
}

var testTable = make(map[int]*Test)
var mutex = sync.Mutex{}

func NewTest(ID, dockerID, scriptID int) *Test {
	dockerPath := (&Docker{dockerID}).Path()
	scriptPath := (&Script{scriptID}).Path()
	return &Test{ID, dockerPath, scriptPath}
}

func GetTest(id int) (*Test, bool) {
	mutex.Lock()
	defer mutex.Unlock()
	test, ok := testTable[id]
	return test, ok
}

func PutTest(testID, dockerID, scriptID int) {
	mutex.Lock()
	defer mutex.Unlock()
	testTable[testID] = NewTest(testID, dockerID, scriptID)
}

func DeleteTest(id int) {
	mutex.Lock()
	defer mutex.Unlock()
	delete(testTable, id)
}

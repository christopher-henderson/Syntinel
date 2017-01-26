package DAO

import "sync"

type Test struct {
	ID         int
	DockerPath string
	ScriptPath string
}

type ThreadSafeMapIntTest struct {
	m map[int]*Test
	sync.Mutex
}

func (m *ThreadSafeMapIntTest) GetTest(id int) (*Test, bool) {
	m.Lock()
	defer m.Unlock()
	test, ok := m.m[id]
	return test, ok
}

func (m *ThreadSafeMapIntTest) SetTest(id int, test *Test) {
	m.Lock()
	defer m.Unlock()
	m.m[id] = test
}

func (m *ThreadSafeMapIntTest) DeleteTest(id int) {
	m.Lock()
	defer m.Unlock()
	delete(m.m, id)
}

var testTable = ThreadSafeMapIntTest{make(map[int]*Test), sync.Mutex{}}

func NewTest(ID, dockerID, scriptID int) *Test {
	dockerPath := (&Docker{dockerID}).Path()
	scriptPath := (&Script{scriptID}).Path()
	return &Test{ID, dockerPath, scriptPath}
}

func GetTest(id int) (*Test, bool) {
	return testTable.GetTest(id)
}

func PutTest(testID, dockerID, scriptID int) {
	testTable.SetTest(testID, NewTest(testID, dockerID, scriptID))
}

func DeleteTest(id int) {
	testTable.DeleteTest(id)
}

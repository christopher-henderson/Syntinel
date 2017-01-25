package DAO

import "syntinel_executor/DAO/database"

var TestRun = &testRun{}

func (t *testRun) Get(id int) (*database.TestRunEntity, error) {
	return database.GetTestRun(id)
}

func (t *testRun) Save(id int, test int, environmentVariables, dockerfile, script string) error {
	return database.InsertTestRun(id, test, environmentVariables, dockerfile, script)
}

func (t *testRun) Delete(id int) error {
	return database.DeleteTestRun(id)
}

type testRun struct{}

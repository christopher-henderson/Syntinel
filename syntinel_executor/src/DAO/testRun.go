package DAO

import (
	"syntinel_executor/DAO/database"
	"syntinel_executor/DAO/database/entities"
)

var TestRun = &testRun{}

func (t *testRun) Get(id int) (*entities.TestRunEntity, error) {
	return database.GetTestRun(id)
}

func (t *testRun) Save(testRun *entities.TestRunEntity) error {
	return database.InsertTestRun(testRun)
}

func (t *testRun) Delete(id int) error {
	return database.DeleteTestRun(id)
}

type testRun struct{}

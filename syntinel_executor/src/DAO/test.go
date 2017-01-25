package DAO

import (
	"syntinel_executor/DAO/database"
	"syntinel_executor/DAO/database/entities"
)

var Test = &test{}

func (t *test) Get(id int) (*entities.TestEntity, error) {
	return database.GetTest(id)
}

func (t *test) Save(id, dockerfile, script int) error {
	return database.InsertTest(id, dockerfile, script)
}

func (t *test) Delete(id int) error {
	return database.DeleteTest(id)
}

func (t *test) Update(id, dockerfile, script int) error {
	return database.UpdateTest(id, dockerfile, script)
}

type test struct{}

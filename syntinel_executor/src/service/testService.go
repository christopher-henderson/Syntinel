package service

import (
	"log"
	"syntinel_executor/DAO"
)

var ts = &testService{}

var GetTestService func() *testService = func() func() *testService {
	return func() *testService {
		return ts
	}
}()

type testService struct {
}

func (t *testService) Delete(id int) error {
	DAO.DeleteTest(id)
	return nil
}

func (t *testService) Register(id, dockerID, scriptID int) error {
	log.Printf("Registering test ID %v with Docker ID %v and script ID %v.\n", id, dockerID, scriptID)
	DAO.PutTest(id, dockerID, scriptID)
	return nil
}

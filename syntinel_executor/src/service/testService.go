package service

import (
	"log"
	"syntinel_executor/DAO"
)

var TestService = &testService{}

func (t *testService) Delete(id int) {
	log.Printf("Deleting test ID %v.\n", id)
	DAO.DeleteTest(id)
}

func (t *testService) Register(id, dockerID, scriptID int) {
	log.Printf("Registering test ID %v with Docker ID %v and script ID %v.\n", id, dockerID, scriptID)
	DAO.PutTest(id, dockerID, scriptID)
}

type testService struct{}

package service

import (
	"log"
	"syntinel_executor/DAO"
)

var TestService = &testService{}

func (t *testService) Delete(id int) error {
	log.Printf("Deleting test ID %v.\n", id)
	return DAO.Test.Delete(id)
}

func (t *testService) Register(id, dockerID, scriptID int) error {
	log.Printf("Registering test ID %v with Docker ID %v and script ID %v.\n", id, dockerID, scriptID)
	return DAO.Test.Save(id, dockerID, scriptID)
}

type testService struct{}

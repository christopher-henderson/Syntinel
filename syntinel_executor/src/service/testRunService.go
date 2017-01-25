package service

import (
	"log"
	"syntinel_executor/DAO"
	"syntinel_executor/DAO/database"
	"syntinel_executor/PAO"
)

var TestRunService = &testRunService{}

func (t *testRunService) Get(id int) (*database.TestRunEntity, error) {
	return DAO.TestRun.Get(id)
}

func (t *testRunService) Save(id int, test int) error {
	testEntity, err := DAO.Test.Get(test)
	if err != nil {
		log.Println("Error occur")
		return err
	}
	dockerfile, err := DAO.Docker.Get(testEntity.Dockerfile)
	if err != nil {
		log.Println("Error occur")
		return err
	}
	script, err := DAO.Script.Get(testEntity.Script)
	if err != nil {
		log.Println("Error occur")
		return err
	}
	if err := DAO.TestRun.Save(id, test, "", dockerfile.Content, script.Content); err != nil {
		log.Println("Error occur")
		return err
	}
	PAO.Run(test, id)
	return nil
}

func (t *testRunService) Delete(id int, test int) error {
	if err := DAO.TestRun.Delete(id); err != nil {
		return err
	}
	PAO.Kill(test, id)
	return nil
}

func (t *testRunService) Query(testID, testRunID int) int {
	return PAO.Query(testID, testRunID)
}

type testRunService struct{}

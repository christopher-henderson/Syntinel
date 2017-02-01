package service

import (
	"syntinel_executor/DAO"
	"syntinel_executor/DAO/database/entities"
	"syntinel_executor/PAO"
)

var TestRunService = &testRunService{}

func (t *testRunService) Get(id int) (*entities.TestRunEntity, error) {
	return DAO.TestRun.Get(id)
}

func (t *testRunService) Save(testRun *entities.TestRunEntity) error {
	err := DAO.TestRun.Save(testRun)
	if err != nil {
		return err
	}
	PAO.Run(testRun.TestID, testRun.ID)
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

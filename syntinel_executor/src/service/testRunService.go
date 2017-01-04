package service

import "syntinel_executor/PAO"

var trs = &testRunService{}

var GetTestRunService func() *testRunService = func() func() *testRunService {
	PAO.StartPAO()
	return func() *testRunService {
		return trs
	}
}()

type testRunService struct {
}

func (t *testRunService) Run(testID, testRunID int) error {
	PAO.Run(testID, testRunID)
	return nil
}

func (t *testRunService) Kill(testID, testRunID int) error {
	PAO.Kill(testID, testRunID)
	return nil
}

func (t *testRunService) Query(testID, testRunID int) (interface{}, error) {
	return testID, nil
}

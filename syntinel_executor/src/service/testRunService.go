package service

import "syntinel_executor/PAO"

var TestRunService = newTestRunService()

func newTestRunService() *testRunService {
	PAO.StartPAO()
	return &testRunService{}
}

type testRunService struct {
}

func (t *testRunService) Run(testID, testRunID int) {
	PAO.Run(testID, testRunID)
}

func (t *testRunService) Kill(testID, testRunID int) {
	PAO.Kill(testID, testRunID)
}

func (t *testRunService) Query(testID, testRunID int) int {
	return PAO.Query(testID, testRunID)
}

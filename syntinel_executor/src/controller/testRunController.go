package controller

import (
	"net/http"
	"strconv"
	"syntinel_executor/PAO"
	"syntinel_executor/service"
)

func KillTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusAccepted
	defer WriteJsonResponse(w, payload, status)
	query := r.URL.Query()
	testID, err := strconv.Atoi(query.Get("testID"))
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Test ID must be an integer."
		return
	}
	testRunID, err := strconv.Atoi(query.Get("testRunID"))
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Test Run ID must be an integer."
		return
	}
	payload.Error = service.GetTestRunService().Kill(testID, testRunID)
}

func RunTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusAccepted
	defer WriteJsonResponse(w, payload, status)
	query := r.URL.Query()
	testID, err := strconv.Atoi(query.Get("testID"))
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Test ID must be an integer."
		return
	}
	testRunID, err := strconv.Atoi(query.Get("testRunID"))
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Test Run ID must be an integer."
		return
	}
	payload.Error = service.GetTestRunService().Run(testID, testRunID)
}

func QueryTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusOK
	defer WriteJsonResponse(w, payload, status)
	query := r.URL.Query()
	testID, err := strconv.Atoi(query.Get("testID"))
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Test ID must be an integer."
		return
	}
	testRunID, err := strconv.Atoi(query.Get("testRunID"))
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Test Run ID must be an integer."
		return
	}
	state := service.GetTestRunService().Query(testID, testRunID)
	payload.Data = PAO.TestStateToString(state)
}

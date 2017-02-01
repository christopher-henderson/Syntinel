package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"syntinel_executor/DAO/database/entities"
	"syntinel_executor/service"
)

func KillTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Running test.")
	payload := &Payload{nil, nil}
	status := http.StatusAccepted
	defer WriteJsonResponse(w, payload, status)
	query := r.URL.Query()
	testID, err := strconv.Atoi(query.Get("testID"))
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Bad test ID."
		return
	}
	testRunID, err := strconv.Atoi(query.Get("testRunID"))
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Bad test run ID."
		return
	}
	payload.Error = service.TestRunService.Delete(testID, testRunID)
}

func RunTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Running test.")
	payload := &Payload{nil, nil}
	status := http.StatusAccepted
	defer WriteJsonResponse(w, payload, status)
	testRun := &entities.TestRunEntity{}
	if err := json.NewDecoder(r.Body).Decode(testRun); err != nil {
		status = http.StatusBadRequest
		payload.Error = err
		return
	}
	payload.Error = service.TestRunService.Save(testRun)
	log.Println(payload.Error)
}

func QueryTest(w http.ResponseWriter, r *http.Request) {
	log.Println("Running test.")
	payload := &Payload{nil, nil}
	status := http.StatusOK
	defer WriteJsonResponse(w, payload, status)
	query := r.URL.Query()
	testID, err := strconv.Atoi(query.Get("testID"))
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Bad test ID."
		return
	}
	testRunID, err := strconv.Atoi(query.Get("testRunID"))
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Bad test run ID."
		return
	}
	state := service.TestRunService.Query(testID, testRunID)
	payload.Data = entities.TestStateToString(state)
}

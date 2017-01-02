package controller

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"syntinel_executor/service"
	"syntinel_executor/utils"
)

func KillTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusAccepted
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	payload.Error = service.GetTestService().Kill(id)
}

func RunTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusAccepted
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	payload.Error = service.GetTestService().Run(id)
}

func QueryTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusOK
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	payload.Data, payload.Error = service.GetTestService().Query(id)
}

func RegisterTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusCreated
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		status = http.StatusServiceUnavailable
		return
	}
	err = service.GetTestService().Register(id, body)
	payload.Data = string(body)
	payload.Error = err
	if err != nil {
		status = http.StatusServiceUnavailable
	}
}

func DeleteTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusNoContent
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	payload.Error = service.GetTestService().Delete(id)
}

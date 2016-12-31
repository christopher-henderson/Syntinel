package controller

import (
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
	payload.Error = service.TestService.Kill(id)
}

func RunTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusAccepted
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	payload.Error = service.TestService.Run(id)
}

func QueryTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusOK
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	payload.Data, payload.Error = service.TestService.Query(id)
}

func RegisterTest(w http.ResponseWriter, r *http.Request) {
	// @TODO
}

func DeleteTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusNoContent
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	payload.Data, payload.Error = service.TestService.Query(id)
}

func UpdateTest(w http.ResponseWriter, r *http.Request) {
	// @TODO
}

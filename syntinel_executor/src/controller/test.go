package controller

import (
	"net/http"

	"github.com/gorilla/mux"

	"syntinel_executor/service"
	"syntinel_executor/utils"
)

func KillTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	defer WriteJsonResponse(w, payload, http.StatusAccepted)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	payload.Error = service.TestService.Kill(id)
}

func RunTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	defer WriteJsonResponse(w, payload, http.StatusAccepted)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	payload.Error = service.TestService.Run(id)
}

func QueryTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	defer WriteJsonResponse(w, payload, http.StatusOK)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	payload.Data, payload.Error = service.TestService.Query(id)
}

func RegisterTest(w http.ResponseWriter, r *http.Request) {
	// @TODO
}

func DeleteTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	defer WriteJsonResponse(w, payload, http.StatusNoContent)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	payload.Data, payload.Error = service.TestService.Query(id)
}

func UpdateTest(w http.ResponseWriter, r *http.Request) {
	// @TODO
}

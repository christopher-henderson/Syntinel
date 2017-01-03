package controller

import (
	"net/http"
	"strconv"

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
	query := r.URL.Query()
	dockerIDString := query.Get("dockerID")
	scriptIDString := query.Get("scriptID")
	dockerID, err := strconv.Atoi(dockerIDString)
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Docker ID must be an integer."
		return
	}
	scriptID, err := strconv.Atoi(scriptIDString)
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Script ID must be an integer."
		return
	}
	service.GetTestService().Register(id, dockerID, scriptID)
}

func DeleteTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusNoContent
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	payload.Error = service.GetTestService().Delete(id)
}

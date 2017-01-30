package controller

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"syntinel_executor/service"
)

func RegisterTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusCreated
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id, err := strconv.Atoi(variables["id"])
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Bad Test ID."
		return
	}
	query := r.URL.Query()
	dockerIDString := query.Get("dockerID")
	scriptIDString := query.Get("scriptID")
	dockerID, err := strconv.Atoi(dockerIDString)
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Bad Docker ID."
		return
	}
	scriptID, err := strconv.Atoi(scriptIDString)
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Bad Script ID."
		return
	}
	log.Println("Registinerg")
	payload.Error = service.TestService.Register(id, dockerID, scriptID)
	log.Println(payload.Error)
}

func DeleteTest(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusNoContent
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id, err := strconv.Atoi(variables["id"])
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Bad Test ID."
		return
	}
	payload.Error = service.TestService.Delete(id)
}

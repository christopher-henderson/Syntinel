package controller

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"syntinel_executor/service"

	"github.com/gorilla/mux"
)

func RegisterDocker(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusCreated
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id, err := strconv.Atoi(variables["id"])
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Bad Docker ID."
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	service.DockerService.Register(id, body)
}

func DeleteDocker(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusNoContent
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id, err := strconv.Atoi(variables["id"])
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Bad Docker ID."
		return
	}
	service.DockerService.Delete(id)
}

func UpdateDocker(w http.ResponseWriter, r *http.Request) {
	RegisterDocker(w, r)
}

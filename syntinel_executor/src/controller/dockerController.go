package controller

import (
	"math"
	"net/http"
	"strconv"
	"syntinel_executor/service"

	"github.com/gorilla/mux"
)

func RegisterDocker(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusCreated
	defer WriteJsonResponse(w, payload, status)
	r.ParseMultipartForm(int64(math.Pow(10, 9)))
	file, _, _ := r.FormFile("docker")
	variables := mux.Vars(r)
	id, err := strconv.Atoi(variables["id"])
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Bad Docker ID."
		return
	}
	go service.DockerService.Register(id, file)
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

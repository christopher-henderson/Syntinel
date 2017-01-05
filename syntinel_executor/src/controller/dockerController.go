package controller

import (
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
	r.ParseMultipartForm(0)
	data, header, _ := r.FormFile("docker")
	defer data.Close()
	f, _ := header.Open()
	defer f.Close()
	service.DockerService.Register(id, f)
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

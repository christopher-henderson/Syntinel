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
	// Zero in http.Request.ParseMultipartForm means "do not save to memory,
	// and dump the whole thing to a temporary file on disk". Which is convenient
	// because we want to get this out of memory and onto the file system
	// as fast as possible. Once in the temp file we can link it the original.
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

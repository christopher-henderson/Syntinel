package controller

import (
	"math"
	"net/http"
	"syntinel_executor/service"
	"syntinel_executor/utils"

	"github.com/gorilla/mux"
)

func RegisterDocker(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusCreated
	defer WriteJsonResponse(w, payload, status)
	r.ParseMultipartForm(int64(math.Pow(10, 9)))
	file, _, _ := r.FormFile("docker")
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	service.GetDockerService().Register(id, file)
}

func DeleteDocker(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusNoContent
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	service.GetDockerService().Delete(id)
}

func UpdateDocker(w http.ResponseWriter, r *http.Request) {
	RegisterDocker(w, r)
}

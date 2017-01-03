package controller

import (
	"io/ioutil"
	"net/http"
	"syntinel_executor/service"
	"syntinel_executor/utils"

	"github.com/gorilla/mux"
)

func RegisterScript(w http.ResponseWriter, r *http.Request) {
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
	err = service.GetScriptService().Register(id, body)
	payload.Data = string(body)
	payload.Error = err
	if err != nil {
		status = http.StatusServiceUnavailable
	}
}

func DeleteScript(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusNoContent
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id := utils.AtoI(variables["id"])
	payload.Error = service.GetScriptService().Delete(id)
}

func UpdateScript(w http.ResponseWriter, r *http.Request) {
	RegisterScript(w, r)
}

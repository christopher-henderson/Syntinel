package controller

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"syntinel_executor/service"

	"github.com/gorilla/mux"
)

func RegisterScript(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusCreated
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id, err := strconv.Atoi(variables["id"])
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Bad Script ID."
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	service.ScriptService.Register(id, body)
}

func DeleteScript(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusNoContent
	defer WriteJsonResponse(w, payload, status)
	variables := mux.Vars(r)
	id, err := strconv.Atoi(variables["id"])
	if err != nil {
		status = http.StatusBadRequest
		payload.Data = "Bad Script ID."
		return
	}
	service.ScriptService.Delete(id)
}

func UpdateScript(w http.ResponseWriter, r *http.Request) {
	RegisterScript(w, r)
}

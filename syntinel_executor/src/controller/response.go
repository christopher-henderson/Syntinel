package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter, payload *Payload, status int) {
	w.Header().Set("Content-Type", "application/json")
	if err := recover(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.(error).Error())
		return
	}
	serializedPayload, err := json.MarshalIndent(payload, " ", "    ")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}
	w.WriteHeader(status)
	io.WriteString(w, string(append(serializedPayload, '\n')))
}

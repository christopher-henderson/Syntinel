package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter, payload *Payload, status int) {
	w.Header().Set("Content-Type", "application/json")
	if err := recover(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(""))
		return
	}
	serializedPayload, err := json.MarshalIndent(payload, " ", "    ")
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(""))
		return
	}
	w.WriteHeader(status)
	w.Write(append(serializedPayload, '\n'))
}

package controller

import (
	"encoding/json"
	"log"
)

type Payload struct {
	Error error       `json:"error"`
	Data  interface{} `json:"data"`
}

func NewPayload(requestError error, data interface{}) []byte {
	payload, err := json.MarshalIndent(Payload{requestError, data}, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	return payload
}

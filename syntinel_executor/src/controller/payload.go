package controller

type Payload struct {
	Error error       `json:"error"`
	Data  interface{} `json:"data"`
}

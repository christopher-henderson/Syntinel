package controller

import (
	"net/http"
	"syntinel_executor/service/statistics"
)

func Statistics(w http.ResponseWriter, r *http.Request) {
	payload := &Payload{nil, nil}
	status := http.StatusOK
	defer WriteJsonResponse(w, payload, status)
	payload.Data = statistics.Statistics()
}

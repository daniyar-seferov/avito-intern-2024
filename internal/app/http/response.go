package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Reason string `json:"reason"`
}

func GetErrorResponse(w http.ResponseWriter, handlerName string, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	data, jsonErr := json.Marshal(&ErrorResponse{fmt.Sprintf("%s: %s", handlerName, err)})
	if jsonErr != nil {
		fmt.Println(jsonErr)
	}
	_, _ = w.Write(data)
}

func GetSuccessResponseWithBody(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

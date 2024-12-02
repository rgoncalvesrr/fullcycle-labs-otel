package api

import (
	"encoding/json"
	"net/http"
)

type ResultError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func WriteJsonResult(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(w).Encode(&ResultError{
		StatusCode: statusCode,
		Message:    message,
	})
	if err != nil {
		return
	}
}

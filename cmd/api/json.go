package main

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ErrorJSON(w http.ResponseWriter, err error, status int) error {
	type envelope struct {
		message string
	}

	return WriteJSON(w, status, envelope{message: err.Error()})
}

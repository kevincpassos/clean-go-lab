package platformhttp

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ErrorMapper func(err error) (status int, message string)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, err error, mapper ErrorMapper) {
	status, message := mapper(err)
	WriteJSON(w, status, ErrorResponse{Error: message})
}

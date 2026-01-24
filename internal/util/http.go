package util

import (
	"encoding/json/v2"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.MarshalWrite(w, data)
}

func WriteActivityJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/activity+json")
	w.WriteHeader(status)
	_ = json.MarshalWrite(w, data)
}

func WriteWebFingerJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/jrd+json")
	w.WriteHeader(status)
	_ = json.MarshalWrite(w, data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	WriteJSON(w, status, errorResponse{Error: message})
}

func ReadJSON(r *http.Request, dst any) error {
	return json.UnmarshalRead(r.Body, dst)
}

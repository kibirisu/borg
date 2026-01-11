package util

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"net/http"
	"time"
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

func DeliverToEndpoint(endpoint string, payload any) {
	if endpoint == "" {
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return
		}
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			endpoint,
			bytes.NewBuffer(jsonData),
		)
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		var client http.Client
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		_ = resp.Body.Close()
	}()
}

package util

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/kibirisu/borg/internal/service"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteActivityJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/activity+json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func WriteWebFingerJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/jrd+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	WriteJSON(w, status, errorResponse{Error: message})
}

func ReadJSON(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
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
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(jsonData))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
	}()
}

func DeliverToFollowers(app service.AppService,
	w http.ResponseWriter, r *http.Request, userID int,
	build func(recipientURI string) any,
) {
	followers, err := app.GetAccountFollowers(r.Context(), userID);
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
	for _, follower := range followers {
		payload := build(follower.Uri)
		DeliverToEndpoint(follower.InboxUri, payload)
	}
}


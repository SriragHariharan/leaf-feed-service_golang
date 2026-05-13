package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sriraghariharan/feed-service-go/internal/models"
)

type feedResponse struct {
	Data       []models.Feed `json:"data"`
	NextCursor string        `json:"next_cursor,omitempty"`
}

type errorResponse struct {
	Error errorBody `json:"error"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{
		Error: errorBody{
			Code:    code,
			Message: message,
		},
	})
}

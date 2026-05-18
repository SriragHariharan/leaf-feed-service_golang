package handler

import (
	"net/http"

	"github.com/sriraghariharan/feed-service-go/internal/httputil"
	"github.com/sriraghariharan/feed-service-go/internal/models"
)

type feedResponse struct {
	Data       []models.Feed `json:"data"`
	NextCursor string        `json:"next_cursor,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	httputil.WriteJSON(w, status, payload)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	httputil.WriteError(w, status, code, message)
}

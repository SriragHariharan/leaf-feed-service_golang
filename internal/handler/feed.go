package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sriraghariharan/feed-service-go/internal/middleware"
	"github.com/sriraghariharan/feed-service-go/internal/service"
)

const feedRequestTimeout = 5 * time.Second

type FeedHandler struct {
	service service.IService
}

func NewFeedHandler(s service.IService) *FeedHandler {
	return &FeedHandler{service: s}
}

// GetFeed serves GET /feed and returns cursor-paginated feed items.
func (h *FeedHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeError(w, http.StatusInternalServerError, "service_unavailable", "service dependency is not initialized")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "authenticated user is required")
		return
	}

	cursor := strings.TrimSpace(r.URL.Query().Get("cursor"))
	ctx, cancel := context.WithTimeout(r.Context(), feedRequestTimeout)
	defer cancel()

	feeds, nextCursor, err := h.service.GetFeed(ctx, userID, cursor)
	if err != nil {
		status := http.StatusInternalServerError
		code := "internal_error"
		message := "failed to fetch feed"
		if strings.Contains(strings.ToLower(err.Error()), "decode cursor") || errors.Is(err, context.DeadlineExceeded) {
			status = http.StatusBadRequest
			code = "invalid_cursor"
			message = "cursor is invalid"
		}
		writeError(w, status, code, message)
		return
	}

	writeJSON(w, http.StatusOK, feedResponse{
		Data:       feeds,
		NextCursor: nextCursor,
	})
}

// GetTimeline serves GET /timeline/{user_id} and returns cursor-paginated timeline items.
func (h *FeedHandler) GetTimeline(w http.ResponseWriter, r *http.Request) {
	if h.service == nil {
		writeError(w, http.StatusInternalServerError, "service_unavailable", "service dependency is not initialized")
		return
	}

	vars := mux.Vars(r)
	userID := strings.TrimSpace(vars["user_id"])
	if userID == "" {
		writeError(w, http.StatusBadRequest, "invalid_request", "user_id is required")
		return
	}

	if userID == "self" {
		authUserID, ok := middleware.UserIDFromContext(r.Context())
		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized", "authenticated user is required")
			return
		}
		userID = authUserID
	}
	cursor := strings.TrimSpace(r.URL.Query().Get("cursor"))
	ctx, cancel := context.WithTimeout(r.Context(), feedRequestTimeout)
	defer cancel()

	viewerUserID, _ := middleware.UserIDFromContext(r.Context())
	feeds, nextCursor, err := h.service.GetTimeline(ctx, userID, viewerUserID, cursor)
	if err != nil {
		status := http.StatusInternalServerError
		code := "internal_error"
		message := "failed to fetch timeline"
		if strings.Contains(strings.ToLower(err.Error()), "decode cursor") || errors.Is(err, context.DeadlineExceeded) {
			status = http.StatusBadRequest
			code = "invalid_cursor"
			message = "cursor is invalid"
		}
		writeError(w, status, code, message)
		return
	}

	writeJSON(w, http.StatusOK, feedResponse{
		Data:       feeds,
		NextCursor: nextCursor,
	})
}

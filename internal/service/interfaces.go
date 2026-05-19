package service

import (
	"context"

	"github.com/sriraghariharan/feed-service-go/internal/models"
)

type IService interface {
	// get the feed of logged in user
	GetFeed(ctx context.Context, userId string, cursor string) ([]models.Feed, string, error)

	// get the timeline of visited profile
	GetTimeline(ctx context.Context, userId string, viewerUserId string, cursor string) ([]models.Feed, string, error)
}
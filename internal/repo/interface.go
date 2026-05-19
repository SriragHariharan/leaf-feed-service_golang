package repo

import (
	"context"

	"github.com/sriraghariharan/feed-service-go/internal/models"
)

type IRepository interface {

	// get the feed of logged in user
	GetFeed(ctx context.Context, userId string, cursor string) ([]models.Feed, string, error)

	// get the timeline of visited profile
	GetTimeline(ctx context.Context, userId string, viewerUserId string, cursor string) ([]models.Feed, string, error)

	// update interaction flags on a viewer's feed row
	UpdateFeedInteraction(ctx context.Context, actorUserID, postID string, isLiked, isCommented *bool) error
}

type IUserRepository interface {
	UpsertUser(ctx context.Context, userID, username, profilePic string) error
}

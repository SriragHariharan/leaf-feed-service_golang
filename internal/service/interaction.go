package service

import (
	"context"
	"errors"
	"strings"

	"github.com/sriraghariharan/feed-service-go/internal/kafka/events"
	"github.com/sriraghariharan/feed-service-go/internal/repo"
)

type InteractionService struct {
	repo repo.IRepository
}

func NewInteractionService(r repo.IRepository) *InteractionService {
	return &InteractionService{repo: r}
}

func (s *InteractionService) ApplyInteractionEvent(ctx context.Context, event events.InteractionEvent) error {
	if s.repo == nil {
		return errors.New("repo dependency is nil")
	}

	actorUserID := strings.TrimSpace(event.ActorUserId)
	postID := strings.TrimSpace(event.PostId)
	eventType := strings.TrimSpace(event.EventType)

	if actorUserID == "" || postID == "" || eventType == "" {
		return errors.New("invalid interaction event")
	}

	var isLiked *bool
	var isCommented *bool

	switch eventType {
	case "post.liked":
		value := true
		isLiked = &value
	case "post.unliked":
		value := false
		isLiked = &value
	case "post.commented":
		value := true
		isCommented = &value
	case "post.uncommented":
		value := false
		isCommented = &value
	default:
		return nil
	}

	return s.repo.UpdateFeedInteraction(ctx, actorUserID, postID, isLiked, isCommented)
}

package service

import (
	"context"
	"errors"

	"github.com/sriraghariharan/feed-service-go/internal/models"
	"github.com/sriraghariharan/feed-service-go/internal/repo"
)

type Service struct {
	repo repo.IRepository
}

func NewService(r repo.IRepository) *Service {
	return &Service{repo: r}
}

// get the feed of logged in user
func (s *Service) GetFeed(ctx context.Context, userId string, cursor string) ([]models.Feed, string, error) {
	if s.repo == nil {
		return nil, "", errors.New("repo dependency is nil")
	}

	feeds, nextCursor, err := s.repo.GetFeed(ctx, userId, cursor)
	if err != nil {
		return nil, "", err
	}
	return feeds, nextCursor, nil
}

// get the timeline of visited profile
func (s *Service) GetTimeline(ctx context.Context, userId string, cursor string) ([]models.Feed, string, error) {
	if s.repo == nil {
		return nil, "", errors.New("repo dependency is nil")
	}
	return s.repo.GetTimeline(ctx, userId, cursor)
}

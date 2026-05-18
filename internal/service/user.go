package service

import (
	"context"
	"errors"

	"github.com/sriraghariharan/feed-service-go/internal/repo"
)

type UserService struct {
	repo repo.IUserRepository
}

func NewUserService(r repo.IUserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) SyncUserFromEvent(ctx context.Context, userID, username string, profilePicture *string) error {
	if s.repo == nil {
		return errors.New("repo dependency is nil")
	}

	profilePic := ""
	if profilePicture != nil {
		profilePic = *profilePicture
	}

	return s.repo.UpsertUser(ctx, userID, username, profilePic)
}

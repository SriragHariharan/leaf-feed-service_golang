package service

import "context"

type IUserService interface {
	SyncUserFromEvent(ctx context.Context, userID, username string, profilePicture *string) error
}

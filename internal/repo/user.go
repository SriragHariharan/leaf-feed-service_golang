package repo

import (
	"context"
	"errors"
	"time"

	"github.com/sriraghariharan/feed-service-go/internal/models"
	"gorm.io/gorm/clause"
)

func (r *Repo) UpsertUser(ctx context.Context, userID, username, profilePic string) error {
	if r.db == nil {
		return errors.New("repo db dependency is nil")
	}

	now := time.Now().UTC()
	user := models.User{
		UserID:     userID,
		Username:   username,
		ProfilePic: profilePic,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"username":    username,
			"profile_pic": profilePic,
			"updated_at":  now,
		}),
	}).Create(&user).Error
}

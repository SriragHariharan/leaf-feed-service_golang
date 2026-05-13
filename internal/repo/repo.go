package repo

import (
	"context"
	"errors"

	"github.com/sriraghariharan/feed-service-go/internal/models"
	"gorm.io/gorm"
)

const (
	FeedsPerCursor  = 20
	TimelinePerPage = 20
)

type Repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// GetFeed returns latest feeds for a user using cursor-based pagination.
// Pass empty cursor for the first page.
func (r *Repo) GetFeed(ctx context.Context, userId string, cursor string) ([]models.Feed, string, error) {
	if r.db == nil {
		return nil, "", errors.New("repo db dependency is nil")
	}

	cursorData, err := decodeFeedCursor(cursor)
	if err != nil {
		return nil, "", errors.New("failed to decode cursor: " + err.Error())
	}

	var feeds []models.Feed
	query := r.db.WithContext(ctx).
		Model(&models.Feed{}).
		Preload("Author").
		Preload("Post").
		Where("user_id = ?", userId).
		Order("created_at DESC, feed_id DESC").
		Limit(FeedsPerCursor + 1)

	if cursorData != nil {
		query = query.Where(
			"(created_at < ?) OR (created_at = ? AND feed_id < ?)",
			cursorData.CreatedAt, cursorData.CreatedAt, cursorData.FeedID,
		)
	}

	err = query.Find(&feeds).Error
	if err != nil {
		return nil, "", errors.New("failed to get feed: " + err.Error())
	}

	nextCursor := ""
	if len(feeds) > FeedsPerCursor {
		feeds = feeds[:FeedsPerCursor]
		lastFeed := feeds[len(feeds)-1]

		nextCursor, err = encodeFeedCursor(feedCursor{
			CreatedAt: lastFeed.CreatedAt,
			FeedID:    lastFeed.FeedID,
		})
		if err != nil {
			return nil, "", errors.New("failed to encode next cursor: " + err.Error())
		}
	}

	return feeds, nextCursor, nil
}

// GetTimeline returns timeline feeds with cursor based pagination.
func (r *Repo) GetTimeline(ctx context.Context, userId string, cursor string) ([]models.Feed, string, error) {
	if r.db == nil {
		return nil, "", errors.New("repo db dependency is nil")
	}

	cursorData, err := decodeFeedCursor(cursor)
	if err != nil {
		return nil, "", errors.New("failed to decode cursor: " + err.Error())
	}

	query := r.db.WithContext(ctx).
		Model(&models.Feed{}).
		Preload("Author").
		Preload("Post").
		Where("author_id = ?", userId).
		Order("created_at DESC, feed_id DESC").
		Limit(TimelinePerPage + 1)

	if cursorData != nil {
		query = query.Where(
			"(created_at < ?) OR (created_at = ? AND feed_id < ?)",
			cursorData.CreatedAt, cursorData.CreatedAt, cursorData.FeedID,
		)
	}
	var feeds []models.Feed
	err = query.Find(&feeds).Error
	if err != nil {
		return nil, "", errors.New("failed to get timeline: " + err.Error())
	}

	nextCursor := ""
	if len(feeds) > TimelinePerPage {
		feeds = feeds[:TimelinePerPage]
		lastFeed := feeds[len(feeds)-1]

		nextCursor, err = encodeFeedCursor(feedCursor{
			CreatedAt: lastFeed.CreatedAt,
			FeedID:    lastFeed.FeedID,
		})
		if err != nil {
			return nil, "", errors.New("failed to encode next cursor: " + err.Error())
		}
	}

	return feeds, nextCursor, nil
}

package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	fmt.Println("feeds: ", feeds)
	fmt.Println("length of feeds: ", len(feeds))
	return feeds, nextCursor, nil
}

// GetTimeline returns posts authored by userId (feed_posts.owner_id) with cursor pagination.
func (r *Repo) GetTimeline(ctx context.Context, userId string, viewerUserId string, cursor string) ([]models.Feed, string, error) {
	if r.db == nil {
		return nil, "", errors.New("repo db dependency is nil")
	}

	cursorData, err := decodePostCursor(cursor)
	if err != nil {
		return nil, "", errors.New("failed to decode cursor: " + err.Error())
	}

	query := r.db.WithContext(ctx).
		Model(&models.Post{}).
		InnerJoins("Owner").
		Where("feed_posts.owner_id = ?", userId).
		Order("feed_posts.created_at DESC, feed_posts.post_id DESC").
		Limit(TimelinePerPage + 1)

	if cursorData != nil {
		query = query.Where(
			"(feed_posts.created_at < ?) OR (feed_posts.created_at = ? AND feed_posts.post_id < ?)",
			cursorData.CreatedAt, cursorData.CreatedAt, cursorData.PostID,
		)
	}

	var posts []models.Post
	if err := query.Find(&posts).Error; err != nil {
		return nil, "", errors.New("failed to get timeline: " + err.Error())
	}

	nextCursor := ""
	if len(posts) > TimelinePerPage {
		posts = posts[:TimelinePerPage]
		lastPost := posts[len(posts)-1]

		nextCursor, err = encodePostCursor(postCursor{
			CreatedAt: lastPost.CreatedAt,
			PostID:    lastPost.PostID,
		})
		if err != nil {
			return nil, "", errors.New("failed to encode next cursor: " + err.Error())
		}
	}

	feeds := make([]models.Feed, 0, len(posts))
	for i := range posts {
		post := posts[i]
		feed := models.Feed{
			FeedID:    post.PostID,
			AuthorID:  post.OwnerID,
			PostID:    post.PostID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			Post:      &post,
		}
		if post.Owner != nil {
			owner := *post.Owner
			feed.Author = &owner
		}
		feeds = append(feeds, feed)
	}

	if viewerUserId != "" && len(feeds) > 0 {
		postIDs := make([]string, 0, len(feeds))
		for i := range feeds {
			postIDs = append(postIDs, feeds[i].PostID)
		}

		var interactionRows []models.Feed
		if err := r.db.WithContext(ctx).
			Model(&models.Feed{}).
			Where("user_id = ? AND post_id IN ?", viewerUserId, postIDs).
			Find(&interactionRows).Error; err != nil {
			return nil, "", errors.New("failed to load timeline interaction flags: " + err.Error())
		}

		interactionByPostID := make(map[string]models.Feed, len(interactionRows))
		for i := range interactionRows {
			interactionByPostID[interactionRows[i].PostID] = interactionRows[i]
		}

		for i := range feeds {
			if row, ok := interactionByPostID[feeds[i].PostID]; ok {
				feeds[i].IsLiked = row.IsLiked
				feeds[i].IsCommented = row.IsCommented
			}
		}
	}

	return feeds, nextCursor, nil
}

func (r *Repo) UpdateFeedInteraction(ctx context.Context, actorUserID, postID string, isLiked, isCommented *bool) error {
	if r.db == nil {
		return errors.New("repo db dependency is nil")
	}

	updates := map[string]interface{}{}
	if isLiked != nil {
		updates["is_liked"] = *isLiked
	}
	if isCommented != nil {
		updates["is_commented"] = *isCommented
	}
	if len(updates) == 0 {
		return nil
	}

	result := r.db.WithContext(ctx).
		Model(&models.Feed{}).
		Where("user_id = ? AND post_id = ?", actorUserID, postID).
		Updates(updates)
	if result.Error != nil {
		return errors.New("failed to update feed interaction: " + result.Error.Error())
	}

	if result.RowsAffected > 0 {
		return nil
	}

	var post models.Post
	if err := r.db.WithContext(ctx).
		Model(&models.Post{}).
		Where("post_id = ?", postID).
		First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return errors.New("failed to lookup post for feed interaction upsert: " + err.Error())
	}

	isLikedVal := false
	isCommentedVal := false
	if isLiked != nil {
		isLikedVal = *isLiked
	}
	if isCommented != nil {
		isCommentedVal = *isCommented
	}

	now := time.Now().UTC()
	newFeed := models.Feed{
		FeedID:      uuid.NewString(),
		UserID:      actorUserID,
		AuthorID:    post.OwnerID,
		PostID:      postID,
		IsLiked:     isLikedVal,
		IsCommented: isCommentedVal,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := r.db.WithContext(ctx).Create(&newFeed).Error; err != nil {
		return errors.New("failed to create feed interaction row: " + err.Error())
	}

	return nil
}

package seed

import (
	"context"
	"errors"
	"log"
	"math/rand/v2"
	"time"

	"github.com/google/uuid"
	"github.com/sriraghariharan/feed-service-go/internal/models"
	"gorm.io/gorm"
)

const (
	feedCount     = 500_000
	feedBatchSize = 5_000
	feedLogEvery  = 50_000
)

// SeedFeeds inserts exactly feedCount feed rows using CreateInBatches(feedBatchSize).
// Random (viewer, post) pairs may repeat; do not assume UNIQUE(user_id, post_id).
func SeedFeeds(ctx context.Context, db *gorm.DB, log *log.Logger, r *rand.Rand, userIDs []string, posts []postRef) error {
	if len(userIDs) == 0 {
		return errors.New("seed: SeedFeeds requires non-empty userIDs")
	}
	if len(posts) == 0 {
		return errors.New("seed: SeedFeeds requires non-empty posts")
	}

	now := time.Now().UTC()
	batch := make([]models.Feed, 0, feedBatchSize)
	committed := 0

	for range feedCount {
		p := posts[r.IntN(len(posts))]
		viewerID := randomUsers(userIDs, r)
		ts := randomTime(r, now)

		batch = append(batch, models.Feed{
			FeedID:      uuid.NewString(),
			UserID:      viewerID,
			AuthorID:    p.OwnerID,
			PostID:      p.PostID,
			IsLiked:     randomBool(r, 0.30),
			IsCommented: randomBool(r, 0.10),
			CreatedAt:   ts,
			UpdatedAt:   ts,
		})

		if len(batch) == feedBatchSize {
			if err := db.WithContext(ctx).CreateInBatches(batch, feedBatchSize).Error; err != nil {
				return err
			}
			committed += feedBatchSize
			if committed%feedLogEvery == 0 {
				log.Printf("Seeded %d feeds...", committed)
			}
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		if err := db.WithContext(ctx).CreateInBatches(batch, feedBatchSize).Error; err != nil {
			return err
		}
		committed += len(batch)
		if committed%feedLogEvery == 0 {
			log.Printf("Seeded %d feeds...", committed)
		}
	}
	return nil
}

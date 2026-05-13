package seed

import (
	"context"
	"log"
	"math/rand/v2"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/sriraghariharan/feed-service-go/internal/models"
	"gorm.io/gorm"
)

const (
	postCount     = 10_000
	postBatchSize = 1_000
)

// postRef is a minimal projection kept in memory for feed generation (post id + owner).
type postRef struct {
	PostID  string
	OwnerID string
}

// SeedPosts inserts exactly postCount posts in batches of postBatchSize via CreateInBatches.
func SeedPosts(ctx context.Context, db *gorm.DB, log *log.Logger, r *rand.Rand, userIDs []string) ([]postRef, error) {
	if len(userIDs) == 0 {
		panic("seed: SeedPosts requires non-empty userIDs")
	}
	now := time.Now().UTC()
	refs := make([]postRef, 0, postCount)
	batch := make([]models.Post, 0, postBatchSize)

	for range postCount {
		postID := uuid.NewString()
		ownerID := randomUsers(userIDs, r)
		ts := randomTime(r, now)
		// Paragraph-length body exercises API payloads and preloads.
		content := faker.Paragraph()
		if len(content) < 80 {
			content = faker.Paragraph() + " " + faker.Sentence()
		}

		batch = append(batch, models.Post{
			PostID:    postID,
			MediaURL:  generateMediaURL(r),
			Content:   content,
			OwnerID:   ownerID,
			CreatedAt: ts,
			UpdatedAt: ts,
		})
		refs = append(refs, postRef{PostID: postID, OwnerID: ownerID})

		if len(batch) == postBatchSize {
			if err := db.WithContext(ctx).CreateInBatches(batch, postBatchSize).Error; err != nil {
				return nil, err
			}
			log.Printf("Seeded %d posts...", len(refs))
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		if err := db.WithContext(ctx).CreateInBatches(batch, postBatchSize).Error; err != nil {
			return nil, err
		}
		log.Printf("Seeded %d posts...", len(refs))
	}
	return refs, nil
}

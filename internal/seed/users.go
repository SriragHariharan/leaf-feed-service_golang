package seed

import (
	"context"
	"log"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/google/uuid"
	"github.com/sriraghariharan/feed-service-go/internal/models"
	"gorm.io/gorm"
)

const (
	userCount     = 500
	userBatchSize = 500
)

// SeedUsers inserts exactly userCount users using a single CreateInBatches call (batch size 500).
func SeedUsers(ctx context.Context, db *gorm.DB, log *log.Logger, r *rand.Rand) ([]string, error) {
	now := time.Now().UTC()
	users := make([]models.User, 0, userCount)
	ids := make([]string, 0, userCount)

	for range userCount {
		id := uuid.NewString()
		ids = append(ids, id)
		ts := randomTime(r, now)
		// Realistic handle: faker username + numeric suffix keeps collisions negligible.
		username := faker.Username() + "_" + strconv.Itoa(r.IntN(1_000_000))
		profilePic := fmtProfilePicURL(id, r)
		users = append(users, models.User{
			UserID:     id,
			Username:   username,
			ProfilePic: profilePic,
			CreatedAt:  ts,
			UpdatedAt:  ts,
		})
	}

	if err := db.WithContext(ctx).CreateInBatches(users, userBatchSize).Error; err != nil {
		return nil, err
	}
	log.Printf("seed: inserted %d users (batch size %d)", userCount, userBatchSize)
	return ids, nil
}

func fmtProfilePicURL(userID string, r *rand.Rand) string {
	// Pravatar-style deterministic avatar URL per user; size varies slightly for realism.
	size := 120 + r.IntN(80)
	return "https://i.pravatar.cc/" + strconv.Itoa(size) + "?u=" + userID
}

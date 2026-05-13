package seed

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// Options controls seed behavior.
type Options struct {
	// Truncate when true: TRUNCATE feed tables before insert (destructive).
	Truncate bool
}

// Run migrates are expected to be done by the caller; this only seeds data.
// Order: users → posts → feeds. Uses a pooled session with FullSaveAssociations disabled.
func Run(ctx context.Context, db *gorm.DB, log *log.Logger, opts Options) error {
	sess := db.Session(&gorm.Session{FullSaveAssociations: false, PrepareStmt: true})

	r, err := newRandFromEnv()
	if err != nil {
		return err
	}

	if opts.Truncate {
		if err := truncateTables(ctx, sess); err != nil {
			return fmt.Errorf("truncate: %w", err)
		}
		log.Printf("seed: truncated feed_feeds, feed_posts, feed_users")
	}

	userIDs, err := SeedUsers(ctx, sess, log, r)
	if err != nil {
		return fmt.Errorf("users: %w", err)
	}

	posts, err := SeedPosts(ctx, sess, log, r, userIDs)
	if err != nil {
		return fmt.Errorf("posts: %w", err)
	}

	if err := SeedFeeds(ctx, sess, log, r, userIDs, posts); err != nil {
		return fmt.Errorf("feeds: %w", err)
	}

	log.Printf("seed: completed (%d users, %d posts, %d feeds)", userCount, len(posts), feedCount)
	return nil
}

func truncateTables(ctx context.Context, db *gorm.DB) error {
	// Child tables first; CASCADE still recommended for FK graphs.
	q := `TRUNCATE TABLE feed_feeds, feed_posts, feed_users RESTART IDENTITY CASCADE`
	return db.WithContext(ctx).Exec(q).Error
}

// newRandFromEnv returns PCG-backed *rand.Rand. If SEED_RNG_SEED is set (decimal uint64),
// it is used for reproducible runs; otherwise crypto/rand seeds PCG state.
func newRandFromEnv() (*rand.Rand, error) {
	if v := strings.TrimSpace(os.Getenv("SEED_RNG_SEED")); v != "" {
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("SEED_RNG_SEED: %w", err)
		}
		return rand.New(rand.NewPCG(n, n^0x9e3779b97f4a7c15)), nil
	}
	var b [16]byte
	if _, err := cryptorand.Read(b[:]); err != nil {
		return nil, fmt.Errorf("crypto/rand: %w", err)
	}
	s0 := binary.LittleEndian.Uint64(b[0:8])
	s1 := binary.LittleEndian.Uint64(b[8:16])
	return rand.New(rand.NewPCG(s0, s1)), nil
}

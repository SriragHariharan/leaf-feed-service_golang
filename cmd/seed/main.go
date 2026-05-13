package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sriraghariharan/feed-service-go/internal/db"
	"github.com/sriraghariharan/feed-service-go/internal/models"
	"github.com/sriraghariharan/feed-service-go/internal/seed"
)

// Example:
//
//	DATABASE_URL="postgres://user:pass@localhost:5432/feed?sslmode=disable" go run ./cmd/seed
//
// Optional:
//
//	SEED_TRUNCATE=1   — TRUNCATE feed_feeds, feed_posts, feed_users before seeding (destructive).
//	SEED_RNG_SEED=42  — fixed PCG seed for reproducible data.
func main() {
	loadDotEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.LUTC)

	gdb, err := db.Connect(ctx)
	if err != nil {
		logger.Fatalf("database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Printf("database close: %v", err)
		}
	}()

	if err := gdb.WithContext(ctx).AutoMigrate(&models.User{}, &models.Post{}, &models.Feed{}); err != nil {
		logger.Fatalf("migrate: %v", err)
	}
	logger.Printf("migrate: applied models.User, models.Post, models.Feed")

	opts := seed.Options{Truncate: truncateFromEnv()}
	if err := seed.Run(ctx, gdb, logger, opts); err != nil {
		logger.Fatalf("seed: %v", err)
	}
	logger.Printf("seed: finished successfully")
}

func truncateFromEnv() bool {
	v := strings.TrimSpace(os.Getenv("SEED_TRUNCATE"))
	switch strings.ToLower(v) {
	case "1", "true", "yes", "y":
		return true
	default:
		return false
	}
}

// loadDotEnv loads the first existing .env from cwd or parent dirs (e.g. module root when
// running `go run .` from cmd/seed).
func loadDotEnv() {
	candidates := []string{
		".env",
		filepath.Join("..", ".env"),
		filepath.Join("..", "..", ".env"),
		filepath.Join("..", "..", "..", ".env"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err != nil {
			continue
		}
		if err := godotenv.Load(p); err != nil {
			log.Printf("env: load %s: %v", p, err)
			return
		}
		return
	}
}

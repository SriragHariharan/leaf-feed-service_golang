package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sriraghariharan/feed-service-go/internal/db"
	"github.com/sriraghariharan/feed-service-go/internal/handler"
	"github.com/sriraghariharan/feed-service-go/internal/kafka"
	"github.com/sriraghariharan/feed-service-go/internal/kafka/consumers"
	"github.com/sriraghariharan/feed-service-go/internal/middleware"
	"github.com/sriraghariharan/feed-service-go/internal/models"
	"github.com/sriraghariharan/feed-service-go/internal/repo"
	"github.com/sriraghariharan/feed-service-go/internal/service"
)

func main() {
	loadDotEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	_, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("database close: %v", err)
		}
	}()

	if err := db.DB.WithContext(ctx).AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	if err := kafka.Connect(); err != nil {
		log.Fatalf("kafka: %v", err)
	}

	database := db.DB
	feedRepo := repo.NewRepo(database)
	feedService := service.NewService(feedRepo)
	feedHandler := handler.NewFeedHandler(feedService)

	userService := service.NewUserService(feedRepo)
	interactionService := service.NewInteractionService(feedRepo)
	consumers.Run(ctx, userService, interactionService)

	if os.Getenv("ACCESS_TOKEN_SECRET") == "" {
		log.Fatal("ACCESS_TOKEN_SECRET is required")
	}

	fmt.Println("Hello, Welcome to the Feed Service!")

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/test", testHandler).Methods("GET")

	protected := router.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.ValidateAccessToken)
	protected.HandleFunc("/feed", feedHandler.GetFeed).Methods("GET")
	protected.HandleFunc("/timeline/{user_id}", feedHandler.GetTimeline).Methods("GET")
	port := os.Getenv("PORT")
	if port == "" {
		port = "2004"
	}
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Feed service running!"})
}

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

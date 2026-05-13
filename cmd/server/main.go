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
)

func main() {
	loadDotEnv()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if _, err := db.Connect(ctx); err != nil {
		log.Fatalf("database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("database close: %v", err)
		}
	}()

	fmt.Println("Hello, Welcome to the Feed Service!")

	//gorilla mux router
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":4040", router))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello, Welcome to the Feed Service!"})
}

// loadDotEnv loads the first existing .env from cwd or parent dirs (e.g. module root when
// running `go run .` from cmd/server). Go does not read .env files without an explicit load.
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
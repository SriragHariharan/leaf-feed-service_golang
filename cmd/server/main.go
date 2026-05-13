package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
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
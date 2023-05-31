package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	store, err := NewSqliteStore("database.db")
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}

	server := NewServer(store)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"message": "pong"
		}`))
	})

	http.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		todoLists, err := server.store.GetAll()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to fetch todo lists"))
			return
		}

		jsonBytes, err := json.Marshal(todoLists)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to marshal todo lists"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonBytes)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

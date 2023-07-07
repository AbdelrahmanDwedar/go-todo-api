package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	router, err := NewRouter()
	if err != nil {
		log.Fatalf("Failed to initialize router: %v", err)
	}

	store, err := NewSqliteStore("database.db")
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}

	server := NewServer(store)

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"message": "pong"
		}`))
	})

	router.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
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

	router.HandleFunc("/todo/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var todoItem TodoItem
		err := json.NewDecoder(r.Body).Decode(&todoItem)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid request payload"))
			return
		}

		err = server.store.Post(todoItem)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to create todo item"))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte("Todo item created"))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

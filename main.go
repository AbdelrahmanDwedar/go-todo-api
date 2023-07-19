package main

import (
	"log"
	"net/http"
)

func main() {
	defer log.Fatal("Server started....")

	router, err := NewRouter()
	if err != nil {
		log.Fatalf("Failed to initialize router: %v", err)
	}

	store, err := NewSqliteStore("database.db")
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}

	server := NewServer(store)

	router.HandleFunc("/ping", MakeHandleFunc(server.HandlePing))
	router.HandleFunc("/todo", MakeHandleFunc(server.HandleTodoList))
	router.HandleFunc("/todo/new", MakeHandleFunc(server.HandleNewItem))
	router.HandleFunc("/todo/lists/{id}", MakeHandleFunc(server.HandleGetList))
	router.HandleFunc("/todo/lists/{id}/new", MakeHandleFunc(server.HandleNewItemInList))
	router.HandleFunc("/todo/lists/new", MakeHandleFunc(server.HandleNewList))

	log.Println("Server Started at port: 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

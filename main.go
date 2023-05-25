package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{
			"message": "pong"
		}`))
	})

	http.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		/* TODO
		* List all the needed todo items in todoList
		* Add json responding with a list of the todo list
		 */
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Pong!\n")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/", apiRoute)

	log.Print("Listening on 8080...\n")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}

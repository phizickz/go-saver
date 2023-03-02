package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/joho/godotenv"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	path := "/health"

	if r.URL.Path != path  {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}


	fmt.Fprintf(w, "Healthy!")
}


func main() {
	http.HandleFunc("/health", healthHandler)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
	}
	fmt.Printf("Starting server at port 8080.\n")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal(err)
	}
}
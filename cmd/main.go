package main

import (
	"log"
	"net/http"
	"os"
	// ...existing imports...
)

func main() {
	// ...existing API router and middleware setup...

	// Serve frontend static files from the frontend folder
	fs := http.FileServer(http.Dir("frontend"))
	http.Handle("/", fs)

	// Use PORT from env, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

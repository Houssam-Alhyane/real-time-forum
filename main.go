package main

import (
	"log"
	"net/http"

	"zone/backend/database"
	"zone/backend/handlers"
)

func main() {

	// Init DB
	if err := database.Init(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Serve static files (CSS + JS)
	fs := http.FileServer(http.Dir("./frontend/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", handlers.Forum)
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", handlers.Logout)

	// (later)
	// http.HandleFunc("/posts", handlers.Posts)
	// http.HandleFunc("/comments", handlers.Comments)
	// http.HandleFunc("/ws", handlers.WebSocketHandler)

	// Start server
	port := ":8080"
	log.Println("Server running on http://localhost" + port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
package main

import (
	"log"
	"net/http"

	"zone/backend/database"
	"zone/backend/handlers"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Static files
	fs := http.FileServer(http.Dir("./frontend/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Auth routes
	// Add this to main.go:
	http.HandleFunc("/register", handlers.Register)
	http.HandleFunc("/login", handlers.Login)
	http.HandleFunc("/logout", handlers.Logout)
http.HandleFunc("/api/me", handlers.Me) 
	// API routes
	http.HandleFunc("/api/posts", handlers.GetPostsAPI)
	http.HandleFunc("/api/categories", handlers.GetCategoriesAPI)
	http.HandleFunc("/api/posts/create", handlers.CreatePostAPI)

	// Catch-all: serve index.html for SPA routes, 404 for unknown /api/* paths
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Any unknown /api/* path → JSON 404
		if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"endpoint not found"}`))
			return
		}
		// Everything else → SPA (index.html handles routing + error display)
		http.ServeFile(w, r, "./frontend/index.html")
	})

	port := ":8080"
	log.Println("Server running on http://localhost" + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
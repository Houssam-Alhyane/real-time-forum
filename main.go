package main

import (
	"log"
	"net/http"

	"zone/backend/database"
	"zone/backend/routing"
)

func main() {
	if err := database.Init(); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Static files
	fs := http.FileServer(http.Dir("./frontend/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

		routing.RegisterRout()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
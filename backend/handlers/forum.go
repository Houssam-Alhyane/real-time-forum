package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

// Forum handles the root URL "/" and serves the core single-page application (SPA) shell.
func Forum(w http.ResponseWriter, r *http.Request) {
	// Strictly allow only the exact root path
	if r.URL.Path != "/" {
		HandleError(w, http.StatusNotFound, "Page not found")
		return
	}

	// Strictly allow only GET requests to load the frontend application
	if r.Method != http.MethodGet {
		HandleError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse the single HTML file that drives the whole frontend SPA
	t, err := template.ParseFiles("./frontend/index.html")
	if err != nil {
		log.Printf("Error parsing SPA main template: %v", err)
		HandleError(w, http.StatusInternalServerError, "Server error loading template")
		return
	}

	// Render into a buffer FIRST. If execution fails halfway through,
	// nothing has been written to the client yet, so we can still
	// send a clean 500 instead of a half-rendered 200 page.
	var buf bytes.Buffer
	if err := t.Execute(&buf, nil); err != nil {
		log.Printf("Error executing SPA template: %v", err)
		HandleError(w, http.StatusInternalServerError, "Server error rendering page")
		return
	}

	// Only now, once we know rendering succeeded, send headers + body.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}
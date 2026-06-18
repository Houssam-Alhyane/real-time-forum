package handlers

import (
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, status int, file string, data any) {
	w.WriteHeader(status)

	tmpl, err := template.ParseFiles("./frontend/" + file)
	if err != nil {
		http.Error(w, "Critical template error: "+err.Error(), 500)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Template execution error: "+err.Error(), 500)
	}
}
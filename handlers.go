package main

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*
var templatesFS embed.FS

// تحميل الـ Template من النظام المدمج
var tmpl = template.Must(template.ParseFS(templatesFS, "templates/index.html"))

type PageData struct {
	Text   string
	Banner string
	Result string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 Page Not Found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl.Execute(w, PageData{})
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	bannerName := r.FormValue("banner")

	result, statusCode, err := ProcessAsciiArt(text, bannerName)
	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	data := PageData{
		Text:   text,
		Banner: bannerName,
		Result: result,
	}

	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

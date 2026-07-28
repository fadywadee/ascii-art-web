package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

// تجهيز وتخزين القالب مرة واحدة فقط عند بداية تشغيل البرنامج
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// داخل الدالة: نستخدم tmpl الجاهز في الذاكرة مباشرة دون قراءة الهارد ديسك!
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// 1. تعريف الـ Struct لنقل البيانات إلى قالب الـ HTML
type PageData struct {
	Text   string
	Banner string
	Result string
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	// الـ Mux يضمن تلقائياً أن الطلب POST وأن المسار هو "/ascii-art"

	// Read form values
	text := r.FormValue("text")
	bannerName := r.FormValue("banner")

	// Validate text
	if text == "" {
		http.Error(w, "Text is required", http.StatusBadRequest)
		return
	}

	// Choose banner file
	var bannerFile string

	switch bannerName {
	case "standard":
		bannerFile = "standard.txt"

	case "shadow":
		bannerFile = "shadow.txt"

	case "thinkertoy":
		bannerFile = "thinkertoy.txt"

	default:
		http.Error(w, "Invalid banner", http.StatusBadRequest)
		return
	}

	// Load banner
	banner, err := loadBanner(bannerFile)
	if err != nil {
		http.Error(w, "Banner not found", http.StatusNotFound)
		return
	}

	// Validate input characters
	valid, badChar := validateInput(text, banner)
	if !valid {
		http.Error(w, "Unsupported character: "+string(badChar), http.StatusBadRequest)
		return
	}

	// Render ASCII Art
	result := render(banner, text)

	// 2. تجهيز البيانات وتغليفها داخل PageData
	data := PageData{
		Text:   text,
		Banner: bannerName,
		Result: result,
	}

	// 3. استبدال w.Write بـ tmpl.Execute لضخ البيانات داخل الـ HTML
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", homeHandler)
	mux.HandleFunc("POST /ascii-art", asciiArtHandler)

	// 💡 قراءة المنفذ من السيرفر السحابي، أو استخدام 8080 كافتراضي محلياً
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server started on http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}

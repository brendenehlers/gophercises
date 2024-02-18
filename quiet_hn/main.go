package main

import (
	"gophercises/quiet_hn/hn"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static", fs))

	http.HandleFunc("/", serveTemplate)

	log.Println("Listening on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("templates", "layout.html")
	tmpl := template.Must(template.ParseFiles(lp))

	items, err := hn.TopItems(30)
	check(err)

	err = tmpl.ExecuteTemplate(w, "layout", items)
	check(err)
}

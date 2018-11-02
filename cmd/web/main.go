package main

import (
	"flag"
	"log"
	"net/http"

	"snippetbox.org/pkg/models"
)

func main() {
	addr := flag.String("addr", ":4040", "HTTP network address")
	htmlDir := flag.String("html-dir", "./ui/html", "Path to HTML templates")
	staticDir := flag.String("static-dir", "./ui/static", "Path to static assets")

	flag.Parse()

	app := &App{
		Database:  &models.Database{},
		HTMLDir:   *htmlDir,
		StaticDir: *staticDir,
	}

	// log.Println("Starting server on :4040")
	log.Printf("Starting server on %s", *addr)
	// err := http.ListenAndServe(*addr, mux)
	err := http.ListenAndServe(*addr, app.Routes())
	log.Fatal(err)
}

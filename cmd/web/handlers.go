package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// Home func
func (app *App) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		// w.WriteHeader(404)
		// w.Write([]byte("Not Found"))
		app.NotFound(w)
		return
	}
	// w.Write([]byte("Hello from Snippetbox"))

	// The file paths that you pass to template.ParseFiles() must be relative
	// to your current working directory or absolute paths.
	// The code below is relative to the root of the project repository (i.e. snippetbox.org).
	// files := []string{
	// 	"./ui/html/base.html",
	// 	"./ui/html/home.page.html",
	// }

	app.RenderHTML(w, "home.page.html")
}

// ShowSnippet func
func (app *App) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// w.Write([]byte("Display a specific snippet..."))
	// fmt.Fprintf(w, "Display a specific snippet (ID %d)...", id)
	snippet, err := app.Database.GetSnippet(id)

	if err != nil {
		app.ServerError(w, err)
		return
	}
	if snippet == nil {
		app.NotFound(w)
		return
	}

	fmt.Fprint(w, snippet)
}

// NewSnippet func
func (app *App) NewSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the new snippet form..."))
}

package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// Home func
func (app *App) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.NotFound(w)
		return
	}

	app.RenderHTML(w, "home.page.html")
}

// ShowSnippet func
func (app *App) ShowSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

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

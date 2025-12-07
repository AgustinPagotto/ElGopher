package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/{$}", app.home)
	mux.HandleFunc("/about", app.about)
	mux.HandleFunc("/articles", app.viewArticles)
	mux.HandleFunc("/projects", app.viewProjects)
	mux.HandleFunc("/article/create", articleCreate)
	mux.HandleFunc("/article/view/{id}", articleView)
	return alice.New(app.recoverPanic, app.logRequest, commonHeaders).Then(mux)
}

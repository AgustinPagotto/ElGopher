package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /about", app.about)
	mux.HandleFunc("GET /articles", app.viewArticles)
	mux.HandleFunc("GET /projects", app.viewProjects)
	mux.HandleFunc("GET /articles/create", app.articleCreate)
	mux.HandleFunc("POST /articles/create", app.articleCreatePost)
	mux.HandleFunc("GET /articles/view/{id}", articleView)
	return alice.New(app.recoverPanic, app.logRequest, commonHeaders, timeoutMiddleware).Then(mux)
}

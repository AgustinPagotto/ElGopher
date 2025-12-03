package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", home)
	mux.HandleFunc("/about", app.about)
	mux.HandleFunc("/article/create", articleCreate)
	mux.HandleFunc("/article/view/{id}", articleView)
	return mux
}

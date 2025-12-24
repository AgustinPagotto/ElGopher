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
	mux.HandleFunc("GET /article/view/{slug}", app.articleView)
	mux.HandleFunc("GET /article/create", app.articleCreate)
	mux.HandleFunc("POST /article/create", app.articleCreatePost)
	mux.HandleFunc("GET /projects", app.viewProjects)
	mux.HandleFunc("GET /user/signup", app.userSignup)
	mux.HandleFunc("POST /user/signup", app.userSignupPost)
	mux.HandleFunc("GET /user/login", app.userLogin)
	mux.HandleFunc("POST /user/login", app.userLoginPost)
	return alice.New(app.recoverPanic, app.logRequest, commonHeaders, timeoutMiddleware).Then(mux)
}

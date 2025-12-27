package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	dynamic := alice.New(app.sessionManager.LoadAndSave, app.authenticated, app.preferences)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))
	mux.Handle("GET /articles", dynamic.ThenFunc(app.viewArticles))
	mux.Handle("GET /article/view/{slug}", dynamic.ThenFunc(app.articleView))
	mux.Handle("GET /projects", dynamic.ThenFunc(app.viewProjects))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("GET /pref/lng", dynamic.ThenFunc(app.setLanguage))
	mux.Handle("GET /pref/thm", dynamic.ThenFunc(app.setTheme))
	protected := dynamic.Append(app.requireAuthentiation)
	mux.Handle("GET /article/create", protected.ThenFunc(app.articleCreate))
	mux.Handle("POST /article/create", protected.ThenFunc(app.articleCreatePost))
	standars := alice.New(app.recoverPanic, app.logRequest, commonHeaders, timeoutMiddleware)
	return standars.Then(mux)
}

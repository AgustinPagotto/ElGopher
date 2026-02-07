package main

import (
	"net/http"

	"github.com/AgustinPagotto/ElGopher/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))
	dynamic := alice.New(app.sessionManager.LoadAndSave, app.authenticated, app.preferences, app.registerEvents, noSurf)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.HandleFunc("GET /ping", ping)
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))
	mux.Handle("GET /articles", dynamic.ThenFunc(app.viewArticles))
	mux.Handle("GET /article/view/{slug}", dynamic.ThenFunc(app.articleView))
	mux.Handle("GET /projects", dynamic.ThenFunc(app.viewProjects))
	//	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	//	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("POST /user/logout", dynamic.ThenFunc(app.logoutPost))
	mux.Handle("GET /pref/lng", dynamic.ThenFunc(app.setLanguage))
	mux.Handle("GET /pref/thm", dynamic.ThenFunc(app.setTheme))
	protected := dynamic.Append(app.requireAuthentiation)
	mux.Handle("GET /article/edit/{slug}", dynamic.ThenFunc(app.articleEdit))
	mux.Handle("PATCH /article/{id}", dynamic.ThenFunc(app.articlePatch))
	mux.Handle("GET /article/create", protected.ThenFunc(app.articleCreate))
	mux.Handle("POST /article/create", protected.ThenFunc(app.articleCreatePost))
	standars := alice.New(app.recoverPanic, app.logRequest, commonHeaders, timeoutMiddleware)
	return standars.Then(mux)
}

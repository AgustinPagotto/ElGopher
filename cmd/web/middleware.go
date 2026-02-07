package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/justinas/nosurf"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(
			"Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' https://cdn.jsdelivr.net 'unsafe-inline'; "+
				"style-src 'self' https://cdn.jsdelivr.net https://fonts.googleapis.com 'unsafe-inline'; "+
				"font-src 'self' https://fonts.gstatic.com",
		)
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("Server", "Go")
		if IsProd() {
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
		next.ServeHTTP(w, r)
	})
}

func timeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(
			r.Context(),
			5*time.Second,
		)
		defer cancel()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (a *application) requireAuthentiation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.isAuthenticated(r) {
			a.sessionManager.Put(r.Context(), "redirectPathAfterLogin", r.URL.Path)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func (a *application) authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := a.sessionManager.GetInt(r.Context(), "authenticatedUserID")
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}
		exists, err := a.users.Exists(r.Context(), id)
		if err != nil {
			a.serverError(w, r, err)
			return
		}
		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}

func (a *application) preferences(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		theme := a.sessionManager.GetBool(r.Context(), "isLightTheme")
		lang := a.sessionManager.GetBool(r.Context(), "isSpanish")
		ctx := context.WithValue(r.Context(), isLightThemeContextKey, theme)
		ctx = context.WithValue(ctx, isSpanishContextKey, lang)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	prod := IsProd()
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   prod,
	})
	csrfHandler.SetIsTLSFunc(func(r *http.Request) bool { return r.TLS != nil })
	return csrfHandler
}

func (a *application) registerEvents(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		theme := a.sessionManager.GetBool(r.Context(), "isLightTheme")
		lang := a.sessionManager.GetBool(r.Context(), "isSpanish")
		path := r.URL.Path
		const articlePrefix = "/article/view/"
		var articleID *int
		if slug, found := strings.CutPrefix(path, articlePrefix); found {
			if slug != "" {
				if article, err := a.articles.GetWithSlug(r.Context(), slug); err == nil {
					articleID = &article.ID
				}
			}
		}
		go func(articleID *int, path string, lang, theme bool) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			err := a.events.Insert(ctx, articleID, path, lang, theme)
			if err != nil {
				a.logger.Info("couldn't append the log", "path", path)
			}
		}(articleID, path, lang, theme)
		next.ServeHTTP(w, r)
	})
}

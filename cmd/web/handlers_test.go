package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/AgustinPagotto/ElGopher/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	status, _, body := ts.get(t, "/ping")
	assert.Equal(t, status, http.StatusOK)
	assert.Equal(t, string(body), "OK")
}

func TestProjects(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	status, _, body := ts.get(t, "/projects")
	assert.Equal(t, status, http.StatusOK)
	assert.StringContains(t, string(body), "Go Webcrawler")
}

func TestHome(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	status, _, body := ts.get(t, "/")
	assert.Equal(t, status, http.StatusOK)
	assert.StringContains(t, string(body), "Who am I?")
}

func TestAbout(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	status, _, body := ts.get(t, "/about")
	assert.Equal(t, status, http.StatusOK)
	assert.StringContains(t, string(body), "Thanks!")
}

func TestSetLanguage(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	status, _, _ := ts.get(t, "/pref/lng")
	assert.Equal(t, status, http.StatusOK)
	status, _, body := ts.get(t, "/")
	assert.StringContains(t, string(body), "Bienvenido")
}

func TestSetLightMode(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	status, _, _ := ts.get(t, "/pref/thm")
	assert.Equal(t, status, http.StatusOK)
	status, _, body := ts.get(t, "/")
	assert.StringContains(t, string(body), "light")
}

func TestArticleView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Latest Article",
			urlPath:  "/article/view/latest",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond",
		},
		{
			name:     "Existent Article",
			urlPath:  "/article/view/an-old-silent-pond",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond",
		},
		{
			name:     "Non Existent Slug",
			urlPath:  "/article/view/asdf-asdf",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Number Slug",
			urlPath:  "/article/view/123-123",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty Slug",
			urlPath:  "/article/view/",
			wantCode: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)
			assert.Equal(t, tt.wantCode, code)
			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}

func TestArticleCreate(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()
	t.Run("Unauthorized", func(t *testing.T) {
		status, header, _ := ts.get(t, "/article/create")
		assert.Equal(t, status, http.StatusSeeOther)
		assert.Equal(t, header.Get("Location"), "/")
	})
	t.Run("Authenticated", func(t *testing.T) {
		_, _, body := ts.get(t, "/user/login")
		token := extractCSRFToken(t, body)
		form := url.Values{}
		form.Add("email", "alice@example.com")
		form.Add("password", "pa$$word")
		form.Add("csrf_token", token)
		status, _, body := ts.postForm(t, "/user/login", form)
		status, _, body = ts.get(t, "/article/create")
		assert.Equal(t, status, http.StatusOK)
		assert.StringContains(t, body, `hx-post="/article/create"`)
	})
}

func TestLogin(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	_, _, body := ts.get(t, "/user/login")
	token := extractCSRFToken(t, body)
	form := url.Values{}
	form.Add("email", "alice@example.com")
	form.Add("password", "pa$$word")
	form.Add("csrf_token", token)
	status, _, _ := ts.postForm(t, "/user/login", form)
	assert.Equal(t, status, http.StatusSeeOther)
}

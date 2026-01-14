package main

import (
	"net/http"
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

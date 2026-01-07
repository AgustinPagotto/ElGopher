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

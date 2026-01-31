package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AgustinPagotto/ElGopher/internal/assert"
)

func TestCommonHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("OK")); err != nil {
			t.Fatal(err)
		}
	})
	commonHeaders(next).ServeHTTP(rr, r)
	rs := rr.Result()
	expectedValue := "default-src 'self'; " +
		"script-src 'self' https://cdn.jsdelivr.net 'unsafe-inline'; " +
		"style-src 'self' https://cdn.jsdelivr.net https://fonts.googleapis.com 'unsafe-inline'; " +
		"font-src 'self' https://fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)
	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)
	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)
	expectedValue = "Go"
	assert.Equal(t, rs.Header.Get("Server"), expectedValue)
	assert.Equal(t, rs.StatusCode, http.StatusOK)
	defer func() {
		if err := rs.Body.Close(); err != nil {
			t.Fatal(err)
		}
	}()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}

func TestTimeoutMiddleware(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deadline, ok := r.Context().Deadline()
		assert.Equal(t, ok, true)
		assert.Equal(t, time.Until(deadline) > 0, true)
	})
	timeoutMiddleware(next).ServeHTTP(rr, r)
}

func TestRecoverPanic(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("ops")
	})
	app := newTestApplication(t)
	app.recoverPanic(next).ServeHTTP(rr, r)
	res := rr.Result()
	if res.Header.Get("Connection") != "close" {
		t.Fatalf("expected Connection: close header")
	}
}

func TestRequireAuthentication(t *testing.T) {
	app := newTestApplication(t)
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/admin", nil)
	if err != nil {
		t.Fatal(err)
	}
	ctx, err := app.sessionManager.Load(r.Context(), "")
	if err != nil {
		t.Fatal(err)
	}
	req := r.WithContext(ctx)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})
	app.requireAuthentiation(next).ServeHTTP(rr, req)
	res := rr.Result()

	if res.StatusCode != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d", res.StatusCode)
	}

	if res.Header.Get("Location") != "/" {
		t.Fatalf("expected redirect to /")
	}
}

package main

import (
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
func articleCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("articleCreate"))
}
func articleView(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
	}
	w.Write([]byte("articleView"))
}

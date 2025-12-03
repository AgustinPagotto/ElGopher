package main

import (
	"html/template"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}
func (a *application) about(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./ui/html/pages/about.html")
	if err != nil {
		a.logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		a.logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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

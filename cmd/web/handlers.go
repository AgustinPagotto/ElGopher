package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/AgustinPagotto/ElGopher/internal/validator"
)

type articleCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Publish             bool   `form:"publish"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "home.html", app.newTemplateData(r))
}
func (app *application) about(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/about.html",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func (app *application) articleCreate(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "createArticle.html", app.newTemplateData(r))
}
func (app *application) articleCreatePost(w http.ResponseWriter, r *http.Request) {
	var form articleCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w)
		return
	}
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	if !form.Valid() {
		return
	}
	id, err := app.articles.Insert(form.Title, form.Content, form.Publish)
	if err != nil {
		app.serverError(w, r, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/articles/view/%d", id), http.StatusSeeOther)
}

func articleView(w http.ResponseWriter, r *http.Request) {
	_, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
	}
	w.Write([]byte("articleView"))
}

func (app *application) viewArticles(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "articles.html", app.newTemplateData(r))
}

func (app *application) viewProjects(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "projects.html", app.newTemplateData(r))
}

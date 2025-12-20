package main

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/AgustinPagotto/ElGopher/internal/models"
	"github.com/AgustinPagotto/ElGopher/internal/validator"
)

type articleCreateForm struct {
	Title               string `form:"title"`
	Body                string `form:"body"`
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
	hxTrigger := r.Header.Get("HX-Trigger")
	app.logger.Info(hxTrigger)
	switch hxTrigger {
	case "title":
		title := r.URL.Query().Get("title")
		var errMsg string
		if !validator.NotBlank(title) {
			w.WriteHeader(http.StatusOK)
			errMsg = "Title cannot be blank"
		}
		app.renderHtmxPartial(w, r, "field_error", errMsg)
	case "body":
		body := r.URL.Query().Get("body")
		var errMsg string
		if !validator.NotBlank(body) {
			w.WriteHeader(http.StatusOK)
			errMsg = "Body cannot be blank"
		}
		app.renderHtmxPartial(w, r, "field_error", errMsg)
	default:
		app.render(w, r, http.StatusOK, "createArticle.html", app.newTemplateData(r))
	}
}
func (app *application) articleCreatePost(w http.ResponseWriter, r *http.Request) {
	var form articleCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w)
		return
	}
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Body), "body", "This field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.renderHtmxPartial(w, r, "form_errors", data)
		return
	}
	_, err = app.articles.Insert(r.Context(), form.Title, form.Body, form.Publish)
	if err != nil {
		app.serverError(w, r, err)
	}
	w.Header().Set("HX-Redirect", "/articles")
	http.Redirect(w, r, "/articles", http.StatusSeeOther)
}

func (app *application) articleView(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	if slug == "" {
		http.NotFound(w, r)
		return
	}
	data := app.newTemplateData(r)
	article, err := app.articles.GetWithSlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	app.logger.Info("here is the issue", "article", article)
	data.Article = article
	app.render(w, r, http.StatusOK, "article.html", data)
}

func (app *application) viewArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := app.articles.GetLastFive(r.Context())
	if err != nil {
		app.serverError(w, r, err)
	}
	data := app.newTemplateData(r)
	data.Articles = articles
	app.render(w, r, http.StatusOK, "articles.html", data)
}

func (app *application) viewProjects(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "projects.html", app.newTemplateData(r))
}

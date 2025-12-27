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

type signUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type logInForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
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
	htmlBody, err := app.MarkToHTML(article.Body)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	data.Article = article
	data.ArticleBody = template.HTML(htmlBody)
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

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	hxTrigger := r.Header.Get("HX-Trigger")
	app.logger.Info(hxTrigger)
	switch hxTrigger {
	case "name":
		name := r.URL.Query().Get("name")
		var errMsg string
		if !validator.NotBlank(name) {
			w.WriteHeader(http.StatusOK)
			errMsg = "Name cannot be blank"
		}
		app.renderHtmxPartial(w, r, "field_error", errMsg)
	case "email":
		email := r.URL.Query().Get("email")
		var errMsg string
		if !validator.EmailValidator(email) {
			w.WriteHeader(http.StatusOK)
			errMsg = "Email contains errors"
		}
		app.renderHtmxPartial(w, r, "field_error", errMsg)
	case "password":
		password := r.URL.Query().Get("password")
		var errMsg string
		if !validator.PasswordValidator(password) {
			w.WriteHeader(http.StatusOK)
			errMsg = "Password cannot be blank or less than 8 characters"
		}
		app.renderHtmxPartial(w, r, "field_error", errMsg)
	default:
		app.render(w, r, http.StatusOK, "signup.html", app.newTemplateData(r))
	}
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form signUpForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w)
		return
	}
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.renderHtmxPartial(w, r, "form_errors", data)
		return
	}
	err = app.users.Insert(r.Context(), form.Name, form.Email, form.Password)
	if err != nil {
		app.serverError(w, r, err)
	}
	w.Header().Set("HX-Redirect", "/user/login")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	hxTrigger := r.Header.Get("HX-Trigger")
	app.logger.Info(hxTrigger)
	switch hxTrigger {
	case "email":
		email := r.URL.Query().Get("email")
		var errMsg string
		if !validator.EmailValidator(email) {
			w.WriteHeader(http.StatusOK)
			errMsg = "Email contain errors"
		}
		app.renderHtmxPartial(w, r, "field_error", errMsg)
	case "password":
		password := r.URL.Query().Get("password")
		var errMsg string
		if !validator.PasswordValidator(password) {
			w.WriteHeader(http.StatusOK)
			errMsg = "Password cannot be blank"
		}
		app.renderHtmxPartial(w, r, "field_error", errMsg)
	default:
		app.render(w, r, http.StatusOK, "login.html", app.newTemplateData(r))
	}
}
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form logInForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w)
		return
	}
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.renderHtmxPartial(w, r, "form_errors", data)
		return
	}
	id, err := app.users.Authenticate(r.Context(), form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or Password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.renderHtmxPartial(w, r, "form_errors", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	path := app.sessionManager.PopString(r.Context(), "redirectPathAfterLogin")
	if path != "" {
		w.Header().Set("HX-Redirect", path)
		http.Redirect(w, r, path, http.StatusSeeOther)
	}
	w.Header().Set("HX-Redirect", "/article/create")
	http.Redirect(w, r, "/article/create", http.StatusSeeOther)
}

func (app *application) setLanguage(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
	}
	lang := app.sessionManager.PopBool(r.Context(), "isSpanish")
	app.sessionManager.Put(r.Context(), "isSpanish", !lang)
	w.Header().Set("HX-Redirect", r.Referer())
	w.WriteHeader(http.StatusOK)
}

func (app *application) setTheme(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
	}
	theme := app.sessionManager.PopBool(r.Context(), "isLightTheme")
	app.sessionManager.Put(r.Context(), "isLightTheme", !theme)
	w.Header().Set("HX-Redirect", r.Referer())
	w.WriteHeader(http.StatusOK)
}

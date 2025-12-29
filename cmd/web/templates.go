package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/AgustinPagotto/ElGopher/internal/i18n"
	"github.com/AgustinPagotto/ElGopher/internal/models"
)

type templateData struct {
	Articles        []models.Article
	Article         models.Article
	ArticleBody     template.HTML
	Errors          []string
	Form            any
	IsAuthenticated bool
	IsSpanish       bool
	IsLightTheme    bool
	Translator      i18n.Translator
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02 Jan 2006")
}

func getTranslation(t templateData, key string) string {
	return t.Translator.T(key)
}

var functions = template.FuncMap{
	"humanDate":      humanDate,
	"getTranslation": getTranslation,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			"./ui/html/partials/field_error.html",
			"./ui/html/partials/form_error.html",
			"./ui/html/partials/article_preview.html",
			page,
		}
		ts, err := template.New(name).Funcs(functions).ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

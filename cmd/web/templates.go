package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/AgustinPagotto/ElGopher/internal/i18n"
	"github.com/AgustinPagotto/ElGopher/internal/models"
	"github.com/AgustinPagotto/ElGopher/ui"
)

type templateData struct {
	Articles        []models.Article
	Article         models.Article
	ArticleBody     template.HTML
	TopArticles     []models.ArticleTop
	DailyViews      []models.DailyViews
	TotalViews      int
	Errors          []string
	Form            any
	IsAuthenticated bool
	IsSpanish       bool
	IsLightTheme    bool
	Translator      i18n.Translator
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format("02/01/2006")
}

func getTranslation(t templateData, key string) string {
	return t.Translator.T(key)
}

func addBreakLines(text string) template.HTML {
	return template.HTML(strings.ReplaceAll(
		template.HTMLEscapeString(text),
		"\n",
		"<br>",
	))
}

var functions = template.FuncMap{
	"humanDate":      humanDate,
	"getTranslation": getTranslation,
	"addBreakLines":  addBreakLines,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

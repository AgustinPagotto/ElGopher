package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"runtime/debug"

	"github.com/AgustinPagotto/ElGopher/internal/i18n"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var method = r.Method
	var uri = r.URL.RequestURI()
	var trace = string(debug.Stack())

	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) renderHtmxPartial(w http.ResponseWriter, r *http.Request, partial string, data any) {
	partialPath := fmt.Sprintf("./ui/html/partials/%s.html", partial)
	tmpl, err := template.ParseFiles(partialPath)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	err = tmpl.ExecuteTemplate(w, partial, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		Form:            map[string]string{},
		IsAuthenticated: app.isAuthenticated(r),
		IsSpanish:       app.isSpanish(r),
		IsLightTheme:    app.isLightTheme(r),
		Translator:      app.getTranslator(r),
		CSRFToken:       nosurf.Token(r),
	}
}

func (app *application) MarkToHTML(markdown string) (string, error) {
	var buf bytes.Buffer
	if err := app.markdownParser.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (a *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (a *application) isSpanish(r *http.Request) bool {
	isSpanish, ok := r.Context().Value(isSpanishContextKey).(bool)
	if !ok {
		return false
	}
	return isSpanish
}

func (a *application) isLightTheme(r *http.Request) bool {
	isLightTheme, ok := r.Context().Value(isLightThemeContextKey).(bool)
	if !ok {
		return false
	}
	return isLightTheme
}

func (a *application) getTranslator(r *http.Request) i18n.Translator {
	isSpanish, ok := r.Context().Value(isSpanishContextKey).(bool)
	if !ok || !isSpanish {
		return i18n.Translator{Messages: i18n.EN}
	}
	return i18n.Translator{Messages: i18n.ES}
}

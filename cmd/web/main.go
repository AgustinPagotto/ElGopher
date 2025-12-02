package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"text/template"
)

type application struct {
	logger        *slog.Logger
	templateCache map[string]*template.Template
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	log.Print("Starting server at :4000")
	err := http.ListenAndServe(":4000", mux)
	if err != nil {
		fmt.Print(err)
	}
}

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	app := &application{
		logger: logger,
	}
	app.logger.Info("Starting server at :4000")
	err := http.ListenAndServe(":4000", app.routes())
	if err != nil {
		fmt.Print(err)
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AgustinPagotto/ElGopher/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	logger        *slog.Logger
	templateCache map[string]*template.Template
	articles      models.ArticleModelInterface
}

func main() {
	dsn := flag.String("dsn", "postgres://postgres:admin@localhost:5432/elgopher", "PostgreSQL data source name")
	flag.Parse()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	var ctx = context.Background()
	db, err := openDB(*dsn, ctx)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("connected to the database")
	defer db.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	app := &application{
		logger:        logger,
		templateCache: templateCache,
		articles:      &models.ArticleModel{DB: db, Ctx: ctx},
	}
	srv := &http.Server{
		Addr:         ":4000",
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	app.logger.Info("Starting server at :4000")
	err = srv.ListenAndServe()
	if err != nil {
		fmt.Print(err)
	}
}

func openDB(dsn string, ctx context.Context) (*pgxpool.Pool, error) {
	var err error
	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return pool, nil
}

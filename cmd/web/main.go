package main

import (
	"context"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/AgustinPagotto/ElGopher/internal/models"
	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

type application struct {
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	articles       models.ArticleModelInterface
	users          models.UserModelInterface
	events         models.EventModelInterface
	markdownParser goldmark.Markdown
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	port := os.Getenv("PORT")
	if dsn == "" {
		dsn = "postgres://postgres:admin@localhost:5432/elgopher"
	}
	if port == "" {
		port = "4000"
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))
	var ctx = context.Background()
	pool, err := openDB(dsn, ctx)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("connected to the database")
	defer pool.Close()
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	formDecoder := form.NewDecoder()
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.New(pool)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	app := &application{
		logger:         logger,
		templateCache:  templateCache,
		articles:       &models.ArticleModel{POOL: pool},
		users:          &models.UserModel{POOL: pool},
		events:         &models.EventModel{POOL: pool},
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		markdownParser: goldmark.New(
			goldmark.WithExtensions(
				highlighting.NewHighlighting(
					highlighting.WithStyle("solarized-dark"),
				),
			),
		),
	}
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	app.logger.Info("Starting server at :", "port", port)
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

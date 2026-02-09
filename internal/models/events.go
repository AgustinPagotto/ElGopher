package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Event struct {
	ID           int64
	ArticleID    *int
	Page         string
	IsSpanish    bool
	IsLightTheme bool
	ViewedAt     time.Time
}

type DailyViews struct {
	Day              time.Time
	Views            int
	SpanishAmount    int
	LightThemeAmount int
}

type ArticleTop struct {
	ArticleSlug string
	Views       int
}

type EventModelInterface interface {
	Insert(ctx context.Context, articleId *int, page string, isSpanish, isLightTheme bool) error
	TotalViews(ctx context.Context) (int, error)
	ViewsPerDay(ctx context.Context) ([]DailyViews, error)
	TopArticles(ctx context.Context) ([]ArticleTop, error)
}

type EventModel struct {
	POOL *pgxpool.Pool
}

func (em *EventModel) Insert(ctx context.Context, articleId *int, page string, isSpanish, isLightTheme bool) error {
	sqlQuery := `
	INSERT INTO events (article_id, page, is_spanish, is_light_theme, viewed_at) 
	VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
	`
	_, err := em.POOL.Exec(ctx, sqlQuery, articleId, page, isSpanish, isLightTheme)
	return err
}

func (em *EventModel) TotalViews(ctx context.Context) (int, error) {
	var total int
	sqlQuery := `SELECT COUNT(*) FROM events;`
	err := em.POOL.QueryRow(ctx, sqlQuery).Scan(&total)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNoRecord
		} else {
			return 0, err
		}
	}
	return total, nil
}

func (em *EventModel) ViewsPerDay(ctx context.Context) ([]DailyViews, error) {
	sqlQuery := `SELECT DATE(viewed_at) AS day, COUNT(*), COUNT(*) FILTER (WHERE is_spanish = TRUE), COUNT(*) FILTER (WHERE is_light_theme = TRUE) FROM events GROUP BY day ORDER BY day DESC;`
	rows, err := em.POOL.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []DailyViews
	for rows.Next() {
		var d DailyViews
		if err := rows.Scan(&d.Day, &d.Views, &d.SpanishAmount, &d.LightThemeAmount); err != nil {
			return nil, err
		}
		result = append(result, d)
	}
	return result, nil
}

func (em *EventModel) TopArticles(ctx context.Context) ([]ArticleTop, error) {
	sqlQuery := `SELECT a.slug, COUNT(*) AS views FROM events e JOIN articles a ON a.id = e.article_id GROUP BY a.slug ORDER BY views DESC;`
	rows, err := em.POOL.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []ArticleTop
	for rows.Next() {
		var p ArticleTop
		if err := rows.Scan(&p.ArticleSlug, &p.Views); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

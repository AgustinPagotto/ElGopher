package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Article struct {
	ID          int
	Title       string
	Body        string
	Slug        string
	Excerpt     string
	IsPublished bool
	Created     time.Time
	UpdatedAt   time.Time
}

type ArticleModelInterface interface {
	Insert(ctx context.Context, title, body string, publish bool) (int, error)
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (Article, error)
	GetWithSlug(ctx context.Context, slug string) (Article, error)
	GetArticles(ctx context.Context) ([]Article, error)
	GetLatest(ctx context.Context) (Article, error)
}

type ArticleModel struct {
	POOL *pgxpool.Pool
}

func (am *ArticleModel) Insert(ctx context.Context, title, body string, publish bool) (int, error) {
	slug := slugifyTitle(title)
	excerpt := generateExcerpt(body)
	sqlQuery := `
	INSERT INTO articles (title, body, slug, excerpt, is_published, created, updated_at) 
	VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) 
	RETURNING id;
	`
	var id int
	err := am.POOL.QueryRow(ctx, sqlQuery, title, body, slug, excerpt, publish).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (am *ArticleModel) Get(ctx context.Context, id int) (Article, error) {
	var article Article
	sqlQuery := `SELECT id, title, body, slug, excerpt, is_published, created, updated_at FROM articles WHERE id = $1;`
	err := am.POOL.QueryRow(ctx, sqlQuery, id).Scan(&article.ID, &article.Title, &article.Body, &article.Slug, &article.Excerpt, &article.IsPublished, &article.Created)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Article{}, ErrNoRecord
		} else {
			return Article{}, err
		}
	}
	return article, nil
}

func (am *ArticleModel) GetWithSlug(ctx context.Context, slug string) (Article, error) {
	var article Article
	sqlQuery := `SELECT id, title, body, slug, excerpt, is_published, created, updated_at FROM articles WHERE slug = $1;`
	err := am.POOL.QueryRow(ctx, sqlQuery, slug).Scan(&article.ID, &article.Title, &article.Body, &article.Slug, &article.Excerpt, &article.IsPublished, &article.Created, &article.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Article{}, ErrNoRecord
		} else {
			return Article{}, err
		}
	}
	return article, nil
}

func (am *ArticleModel) GetArticles(ctx context.Context) ([]Article, error) {
	sqlQuery := `SELECT id, title, slug, excerpt, updated_at FROM articles ORDER BY created DESC;`
	rows, err := am.POOL.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var articles []Article
	for rows.Next() {
		var a Article
		err = rows.Scan(&a.ID, &a.Title, &a.Slug, &a.Excerpt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return articles, nil
}

func (am *ArticleModel) GetLatest(ctx context.Context) (Article, error) {
	var article Article
	sqlQuery := `SELECT id, title, body, slug, excerpt, is_published, created, updated_at FROM articles ORDER BY created DESC LIMIT 1;`
	err := am.POOL.QueryRow(ctx, sqlQuery).Scan(&article.ID, &article.Title, &article.Body, &article.Slug, &article.Excerpt, &article.IsPublished, &article.Created, &article.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Article{}, ErrNoRecord
		} else {
			return Article{}, err
		}
	}
	return article, nil
}

func (am *ArticleModel) Delete(ctx context.Context, id int) error {
	sqlQuery := `DELETE FROM articles WHERE id = ?;`
	_, err := am.POOL.Exec(context.Background(), sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}

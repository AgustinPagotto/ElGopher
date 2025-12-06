package models

import (
	"database/sql"
	"errors"
	"time"

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
	Insert(title, body string) error
	Delete(id int) error
	Get(id int) (Article, error)
	GetLastFive() ([]Article, error)
}

type ArticleModel struct {
	DB *pgxpool.Pool
}

func (am *ArticleModel) Insert(title, body string) error {
	//Create 2 helper functions, one that generates the SLUG, another the Excerpt
	sqlQuery := `INSERT INTO articles (title, body, slug, excerpt, created, updated) VALUES (?, ?, UTC_TIMESTAMP());`
	_, err := am.DB.Exec(sqlQuery, title, body)
	if err != nil {
		return err
	}
	return nil
}

func (am *ArticleModel) Get(id int) (Article, error) {
	var article Article
	sqlQuery := `SELECT id, title, body, slug, excerpt, is_published, created, updated_at FROM articles WHERE id = ?;`
	err := am.DB.QueryRow(sqlQuery, id).Scan(&article.ID, &article.Title, &article.Body, &article.Slug, &article.Excerpt, &article.IsPublished, &article.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Article{}, ErrNoRecord
		} else {
			return Article{}, err
		}
	}
	return article, nil
}

func (am *ArticleModel) GetLastFive() ([]Article, error) {
	sqlQuery := `SELECT id, title, slug, excerpt, created, updated_at FROM articles WHERE is_published != false ORDER BY created DESC LIMIT 5;`
	rows, err := am.DB.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var articles []Article
	for rows.Next() {
		var a Article
		err = rows.Scan(&a.ID, &a.Title, &a.Body, &a.Created)
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

func (am *ArticleModel) Delete(id int) error {
	sqlQuery := `DELETE FROM articles WHERE id = ?;`
	_, err := am.DB.Exec(sqlQuery, id)
	if err != nil {
		return err
	}
	return nil
}

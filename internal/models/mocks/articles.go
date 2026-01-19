package mocks

import (
	"context"
	"time"

	"github.com/AgustinPagotto/ElGopher/internal/models"
)

var mockArticle = models.Article{
	ID:    1,
	Title: "An old silent pond",
	Body: `
		# An old silent pond
		This is a short test article written in **Markdown**.
	`,
	Slug:        "an-old-silent-pond",
	Excerpt:     "An old silent pond...",
	IsPublished: false,
	Created:     time.Now(),
	UpdatedAt:   time.Now(),
}

type ArticleModel struct{}

func (am *ArticleModel) Insert(ctx context.Context, title, body string, publish bool) (int, error) {
	return 2, nil
}

func (am *ArticleModel) Delete(ctx context.Context, id int) error {
	return nil
}
func (am *ArticleModel) Get(ctx context.Context, id int) (models.Article, error) {
	switch id {
	case 1:
		return mockArticle, nil
	default:
		return models.Article{}, models.ErrNoRecord
	}
}
func (am *ArticleModel) Update(ctx context.Context, title, body string, is_published bool, id int) error {
	return nil
}
func (am *ArticleModel) GetWithSlug(ctx context.Context, slug string) (models.Article, error) {
	switch slug {
	case "an-old-silent-pond":
		return mockArticle, nil
	default:
		return models.Article{}, models.ErrNoRecord
	}
}
func (am *ArticleModel) GetArticles(ctx context.Context) ([]models.Article, error) {
	return []models.Article{}, nil
}
func (am *ArticleModel) GetLatest(ctx context.Context) (models.Article, error) {
	return mockArticle, nil
}

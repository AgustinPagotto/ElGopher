package mocks

import (
	"context"

	"github.com/AgustinPagotto/ElGopher/internal/models"
)

type EventModel struct{}

func (m *EventModel) Insert(ctx context.Context, articleId *int, page string, isSpanish, isLightTheme bool) error {
	return nil
}

func (m *EventModel) TotalViews(ctx context.Context) (int, error) {
	return 0, nil
}

func (m *EventModel) ViewsPerDay(ctx context.Context) ([]models.DailyViews, error) {
	return nil, nil
}

func (m *EventModel) TopArticles(ctx context.Context) ([]models.ArticleTop, error) {
	return nil, nil
}


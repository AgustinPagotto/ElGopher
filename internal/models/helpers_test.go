package models

import (
	"testing"

	"github.com/AgustinPagotto/ElGopher/internal/assert"
)

func TestSlugifyTitle(t *testing.T) {
	tests := []struct {
		name           string
		title          string
		wantResultSlug string
	}{
		{
			name:           "Normal title",
			title:          "Hello to go",
			wantResultSlug: "hello-to-go",
		},
		{
			name:           "Accents title",
			title:          "Qué es go?",
			wantResultSlug: "que-es-go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug := slugifyTitle(tt.title)
			assert.Equal(t, slug, tt.wantResultSlug)
		})
	}
}

func TestGenerateExcerpt(t *testing.T) {
	tests := []struct {
		name              string
		body              string
		wantResultExcerpt string
	}{
		{
			name:              "30+ words body",
			body:              "Lorem ipsum dolor sit amet consectetur adipiscing elit quisque faucibus ex sapien vitae pellentesque sem placerat in id cursus mi pretium tellus duis convallis tempus leo eu aenean sed diam urna tempor pulvinar vivamus fringilla lacus nec metus bibendum egestas iaculis massa nisl malesuada lacinia integer nunc posuere ut hendrerit.",
			wantResultExcerpt: "Lorem ipsum dolor sit amet consectetur adipiscing elit quisque faucibus ex sapien vitae pellentesque sem placerat in id cursus mi pretium tellus duis convallis tempus leo eu aenean sed diam",
		},
		{
			name:              "30 words body",
			body:              "Lorem ipsum dolor sit amet consectetur adipiscing elit quisque faucibus ex sapien vitae pellentesque sem placerat in id cursus mi pretium tellus duis convallis tempus leo eu aenean sed diam.",
			wantResultExcerpt: "Lorem ipsum dolor sit amet consectetur adipiscing elit quisque faucibus ex sapien vitae pellentesque sem placerat in id cursus mi pretium tellus duis convallis tempus leo eu aenean sed diam.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug := generateExcerpt(tt.body)
			assert.Equal(t, slug, tt.wantResultExcerpt)
		})
	}
}

package main

import (
	"testing"
	"time"

	"github.com/AgustinPagotto/ElGopher/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 3, 12, 15, 0, 0, 0, time.UTC),
			want: "12/03/2024",
		},
		{
			name: "empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 3, 12, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "12/03/2024",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)
			assert.Equal(t, hd, tt.want)
		})
	}
}

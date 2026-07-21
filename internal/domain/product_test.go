package domain_test

import (
	"errors"
	"testing"
	"time"

	"github.com/nicolas2601/go-graphql-products-api/internal/domain"
)

func TestNewProduct(t *testing.T) {
	createdAt := time.Date(2026, time.July, 20, 12, 0, 0, 0, time.UTC)

	t.Run("valid product is created", func(t *testing.T) {
		p, err := domain.NewProduct("id-1", "Keyboard", 49.90, 10, createdAt)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := domain.Product{ID: "id-1", Name: "Keyboard", Price: 49.90, Stock: 10, CreatedAt: createdAt}
		if p != want {
			t.Fatalf("got %+v, want %+v", p, want)
		}
	})

	t.Run("invalid inputs are rejected", func(t *testing.T) {
		cases := []struct {
			name        string
			productName string
			price       float64
			want        error
		}{
			{"empty name", "", 10, domain.ErrInvalidName},
			{"whitespace name", "   ", 10, domain.ErrInvalidName},
			{"zero price", "Mouse", 0, domain.ErrInvalidPrice},
			{"negative price", "Mouse", -1, domain.ErrInvalidPrice},
		}
		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				_, err := domain.NewProduct("id", c.productName, c.price, 1, createdAt)
				if !errors.Is(err, c.want) {
					t.Fatalf("got error %v, want %v", err, c.want)
				}
			})
		}
	})
}

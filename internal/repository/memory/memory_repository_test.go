package memory_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nicolas2601/go-graphql-products-api/internal/domain"
	"github.com/nicolas2601/go-graphql-products-api/internal/repository/memory"
)

func sampleProduct(id string) domain.Product {
	return domain.Product{
		ID:        id,
		Name:      "Product " + id,
		Price:     10,
		Stock:     1,
		CreatedAt: time.Date(2026, time.July, 20, 12, 0, 0, 0, time.UTC),
	}
}

func TestCreateAndGetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("creates and retrieves a product", func(t *testing.T) {
		repo := memory.New()
		want := sampleProduct("a")

		if err := repo.Create(ctx, want); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got, err := repo.GetByID(ctx, "a")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != want {
			t.Fatalf("got %+v, want %+v", got, want)
		}
	})

	t.Run("duplicate id returns ErrProductAlreadyExists", func(t *testing.T) {
		repo := memory.New()
		p := sampleProduct("a")
		_ = repo.Create(ctx, p)

		if err := repo.Create(ctx, p); !errors.Is(err, domain.ErrProductAlreadyExists) {
			t.Fatalf("got %v, want ErrProductAlreadyExists", err)
		}
	})

	t.Run("missing id returns ErrProductNotFound", func(t *testing.T) {
		repo := memory.New()

		if _, err := repo.GetByID(ctx, "missing"); !errors.Is(err, domain.ErrProductNotFound) {
			t.Fatalf("got %v, want ErrProductNotFound", err)
		}
	})
}

func TestList(t *testing.T) {
	ctx := context.Background()

	t.Run("empty repository returns empty slice", func(t *testing.T) {
		repo := memory.New()

		got, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("got %d products, want 0", len(got))
		}
	})

	t.Run("returns all products ordered by id", func(t *testing.T) {
		repo := memory.New()
		_ = repo.Create(ctx, sampleProduct("c"))
		_ = repo.Create(ctx, sampleProduct("a"))
		_ = repo.Create(ctx, sampleProduct("b"))

		got, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 3 || got[0].ID != "a" || got[1].ID != "b" || got[2].ID != "c" {
			t.Fatalf("expected order a,b,c, got %v", []string{got[0].ID, got[1].ID, got[2].ID})
		}
	})
}

package memory_test

import (
	"context"
	"errors"
	"strconv"
	"sync"
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

func TestUpdate(t *testing.T) {
	ctx := context.Background()

	t.Run("updates an existing product", func(t *testing.T) {
		repo := memory.New()
		_ = repo.Create(ctx, sampleProduct("a"))

		updated := sampleProduct("a")
		updated.Name = "Renamed"
		if err := repo.Update(ctx, updated); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got, _ := repo.GetByID(ctx, "a")
		if got.Name != "Renamed" {
			t.Fatalf("update not persisted: %+v", got)
		}
	})

	t.Run("missing product returns ErrProductNotFound", func(t *testing.T) {
		repo := memory.New()
		if err := repo.Update(ctx, sampleProduct("missing")); !errors.Is(err, domain.ErrProductNotFound) {
			t.Fatalf("got %v, want ErrProductNotFound", err)
		}
	})

	// CreatedAt es inmutable: aunque el producto entrante traiga otro CreatedAt, Update debe
	// preservar el original, igual que el repo PostgreSQL (que no toca la columna created_at).
	t.Run("preserves original CreatedAt", func(t *testing.T) {
		repo := memory.New()
		original := sampleProduct("a")
		_ = repo.Create(ctx, original)

		tampered := sampleProduct("a")
		tampered.CreatedAt = original.CreatedAt.Add(48 * time.Hour)
		if err := repo.Update(ctx, tampered); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got, _ := repo.GetByID(ctx, "a")
		if !got.CreatedAt.Equal(original.CreatedAt) {
			t.Fatalf("CreatedAt changed: got %v, want %v", got.CreatedAt, original.CreatedAt)
		}
	})
}

func TestDelete(t *testing.T) {
	ctx := context.Background()

	t.Run("deletes an existing product", func(t *testing.T) {
		repo := memory.New()
		_ = repo.Create(ctx, sampleProduct("a"))

		if err := repo.Delete(ctx, "a"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, err := repo.GetByID(ctx, "a"); !errors.Is(err, domain.ErrProductNotFound) {
			t.Fatalf("product still present after delete")
		}
	})

	t.Run("missing product returns ErrProductNotFound", func(t *testing.T) {
		repo := memory.New()
		if err := repo.Delete(ctx, "missing"); !errors.Is(err, domain.ErrProductNotFound) {
			t.Fatalf("got %v, want ErrProductNotFound", err)
		}
	})
}

func TestConcurrentAccess(t *testing.T) {
	repo := memory.New()
	ctx := context.Background()
	const workers = 50

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(n int) {
			defer wg.Done()
			id := strconv.Itoa(n)
			_ = repo.Create(ctx, sampleProduct(id))
			_, _ = repo.GetByID(ctx, id)
			_, _ = repo.List(ctx)
		}(i)
	}
	wg.Wait()

	got, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != workers {
		t.Fatalf("got %d products, want %d", len(got), workers)
	}
}

func TestConcurrentCreateSameID(t *testing.T) {
	repo := memory.New()
	ctx := context.Background()
	const workers = 50

	var (
		wg        sync.WaitGroup
		mu        sync.Mutex
		created   int
		duplicate int
	)
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			err := repo.Create(ctx, sampleProduct("same"))
			mu.Lock()
			switch {
			case err == nil:
				created++
			case errors.Is(err, domain.ErrProductAlreadyExists):
				duplicate++
			}
			mu.Unlock()
		}()
	}
	wg.Wait()

	if created != 1 || duplicate != workers-1 {
		t.Fatalf("got created=%d duplicate=%d, want created=1 duplicate=%d", created, duplicate, workers-1)
	}
}

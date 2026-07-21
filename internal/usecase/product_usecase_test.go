package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nicolas2601/go-graphql-products-api/internal/domain"
	"github.com/nicolas2601/go-graphql-products-api/internal/usecase"
)

// fakeRepo es un doble de prueba en memoria que implementa domain.ProductRepository.
type fakeRepo struct {
	items      map[string]domain.Product
	failCreate error
	failUpdate error
}

func newFakeRepo() *fakeRepo { return &fakeRepo{items: make(map[string]domain.Product)} }

func (f *fakeRepo) Create(_ context.Context, p domain.Product) error {
	if f.failCreate != nil {
		return f.failCreate
	}
	f.items[p.ID] = p
	return nil
}

func (f *fakeRepo) GetByID(_ context.Context, id string) (domain.Product, error) {
	p, ok := f.items[id]
	if !ok {
		return domain.Product{}, domain.ErrProductNotFound
	}
	return p, nil
}

func (f *fakeRepo) List(_ context.Context) ([]domain.Product, error) {
	out := make([]domain.Product, 0, len(f.items))
	for _, p := range f.items {
		out = append(out, p)
	}
	return out, nil
}

func (f *fakeRepo) Update(_ context.Context, p domain.Product) error {
	if _, ok := f.items[p.ID]; !ok {
		return domain.ErrProductNotFound
	}
	if f.failUpdate != nil {
		return f.failUpdate
	}
	f.items[p.ID] = p
	return nil
}

func (f *fakeRepo) Delete(_ context.Context, id string) error {
	if _, ok := f.items[id]; !ok {
		return domain.ErrProductNotFound
	}
	delete(f.items, id)
	return nil
}

const fixedID = "prod-1"

var fixedTime = time.Date(2026, time.July, 20, 12, 0, 0, 0, time.UTC)

func newUseCase(repo domain.ProductRepository) *usecase.ProductUseCase {
	return usecase.NewProductUseCase(repo,
		func() string { return fixedID },
		func() time.Time { return fixedTime },
	)
}

func TestCreate(t *testing.T) {
	t.Run("valid product is created and persisted", func(t *testing.T) {
		repo := newFakeRepo()
		uc := newUseCase(repo)

		p, err := uc.Create(context.Background(), "Keyboard", 49.90, 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.ID != fixedID || !p.CreatedAt.Equal(fixedTime) {
			t.Fatalf("id/createdAt not injected: %+v", p)
		}
		if _, ok := repo.items[fixedID]; !ok {
			t.Fatalf("product was not persisted")
		}
	})

	t.Run("invalid product returns domain error and is not persisted", func(t *testing.T) {
		repo := newFakeRepo()
		uc := newUseCase(repo)

		_, err := uc.Create(context.Background(), "", 49.90, 10)
		if !errors.Is(err, domain.ErrInvalidName) {
			t.Fatalf("got %v, want ErrInvalidName", err)
		}
		if len(repo.items) != 0 {
			t.Fatalf("nothing should have been persisted")
		}
	})
}

func TestGet(t *testing.T) {
	repo := newFakeRepo()
	repo.items[fixedID] = domain.Product{ID: fixedID, Name: "Mouse", Price: 10, Stock: 5, CreatedAt: fixedTime}
	uc := newUseCase(repo)

	t.Run("existing product is returned", func(t *testing.T) {
		p, err := uc.Get(context.Background(), fixedID)
		if err != nil || p.Name != "Mouse" {
			t.Fatalf("got %+v, %v", p, err)
		}
	})

	t.Run("missing product returns ErrProductNotFound", func(t *testing.T) {
		_, err := uc.Get(context.Background(), "does-not-exist")
		if !errors.Is(err, domain.ErrProductNotFound) {
			t.Fatalf("got %v, want ErrProductNotFound", err)
		}
	})
}

func TestList(t *testing.T) {
	repo := newFakeRepo()
	repo.items["a"] = domain.Product{ID: "a"}
	repo.items["b"] = domain.Product{ID: "b"}
	uc := newUseCase(repo)

	got, err := uc.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("got %d products, want 2", len(got))
	}
}

func TestUpdate(t *testing.T) {
	base := domain.Product{ID: fixedID, Name: "Mouse", Price: 10, Stock: 5, CreatedAt: fixedTime}

	t.Run("updates only the provided fields and persists", func(t *testing.T) {
		repo := newFakeRepo()
		repo.items[fixedID] = base
		uc := newUseCase(repo)

		newName := "Gaming Mouse"
		p, err := uc.Update(context.Background(), fixedID, &newName, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Name != "Gaming Mouse" || repo.items[fixedID].Name != "Gaming Mouse" {
			t.Fatalf("name not updated: %+v", repo.items[fixedID])
		}
		if p.Price != 10 {
			t.Fatalf("price should be unchanged, got %v", p.Price)
		}
	})

	t.Run("invalid update returns domain error and does not persist", func(t *testing.T) {
		repo := newFakeRepo()
		repo.items[fixedID] = base
		uc := newUseCase(repo)

		invalidPrice := 0.0
		_, err := uc.Update(context.Background(), fixedID, nil, &invalidPrice)
		if !errors.Is(err, domain.ErrInvalidPrice) {
			t.Fatalf("got %v, want ErrInvalidPrice", err)
		}
		if repo.items[fixedID].Price != 10 {
			t.Fatalf("price should remain unchanged")
		}
	})

	t.Run("missing product returns ErrProductNotFound", func(t *testing.T) {
		repo := newFakeRepo()
		uc := newUseCase(repo)

		name := "X"
		_, err := uc.Update(context.Background(), "missing", &name, nil)
		if !errors.Is(err, domain.ErrProductNotFound) {
			t.Fatalf("got %v, want ErrProductNotFound", err)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("existing product is deleted", func(t *testing.T) {
		repo := newFakeRepo()
		repo.items[fixedID] = domain.Product{ID: fixedID}
		uc := newUseCase(repo)

		if err := uc.Delete(context.Background(), fixedID); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := repo.items[fixedID]; ok {
			t.Fatalf("product was not deleted")
		}
	})

	t.Run("missing product returns ErrProductNotFound", func(t *testing.T) {
		repo := newFakeRepo()
		uc := newUseCase(repo)

		if err := uc.Delete(context.Background(), "missing"); !errors.Is(err, domain.ErrProductNotFound) {
			t.Fatalf("got %v, want ErrProductNotFound", err)
		}
	})
}

func TestRepositoryErrorPropagation(t *testing.T) {
	sentinel := errors.New("repository failure")

	t.Run("create propagates repository error", func(t *testing.T) {
		repo := newFakeRepo()
		repo.failCreate = sentinel
		uc := newUseCase(repo)

		if _, err := uc.Create(context.Background(), "Mouse", 10, 1); !errors.Is(err, sentinel) {
			t.Fatalf("got %v, want sentinel error", err)
		}
	})

	t.Run("update propagates repository error", func(t *testing.T) {
		repo := newFakeRepo()
		repo.items[fixedID] = domain.Product{ID: fixedID, Name: "Mouse", Price: 10, CreatedAt: fixedTime}
		repo.failUpdate = sentinel
		uc := newUseCase(repo)

		name := "New"
		if _, err := uc.Update(context.Background(), fixedID, &name, nil); !errors.Is(err, sentinel) {
			t.Fatalf("got %v, want sentinel error", err)
		}
	})
}

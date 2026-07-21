package postgres_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/nicolas2601/go-graphql-products-api/internal/domain"
	"github.com/nicolas2601/go-graphql-products-api/internal/repository/postgres"
)

func newTestRepo(t *testing.T) (*postgres.Repository, *pgxpool.Pool) {
	t.Helper()
	testcontainers.SkipIfProviderIsNotHealthy(t)
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx, "postgres:16-alpine",
		tcpostgres.WithDatabase("products"),
		tcpostgres.WithUsername("products"),
		tcpostgres.WithPassword("products"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("5432/tcp").WithStartupTimeout(90*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	testcontainers.CleanupContainer(t, container)

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := postgres.Migrate(ctx, pool); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return postgres.New(pool), pool
}

func sampleProduct(id string) domain.Product {
	return domain.Product{
		ID:        id,
		Name:      "Product " + id,
		Price:     10,
		Stock:     1,
		CreatedAt: time.Date(2026, time.July, 21, 12, 0, 0, 0, time.UTC),
	}
}

func TestPostgresRepository(t *testing.T) {
	repo, pool := newTestRepo(t)
	ctx := context.Background()
	reset := func() {
		if _, err := pool.Exec(ctx, "TRUNCATE products"); err != nil {
			t.Fatalf("truncate: %v", err)
		}
	}

	t.Run("create and get", func(t *testing.T) {
		reset()
		want := sampleProduct("11111111-1111-1111-1111-111111111111")
		if err := repo.Create(ctx, want); err != nil {
			t.Fatalf("create: %v", err)
		}
		got, err := repo.GetByID(ctx, want.ID)
		if err != nil {
			t.Fatalf("get: %v", err)
		}
		if got.ID != want.ID || got.Name != want.Name || got.Price != want.Price ||
			got.Stock != want.Stock || !got.CreatedAt.Equal(want.CreatedAt) {
			t.Fatalf("got %+v, want %+v", got, want)
		}
	})

	t.Run("duplicate id returns ErrProductAlreadyExists", func(t *testing.T) {
		reset()
		p := sampleProduct("22222222-2222-2222-2222-222222222222")
		if err := repo.Create(ctx, p); err != nil {
			t.Fatalf("first create: %v", err)
		}
		if err := repo.Create(ctx, p); !errors.Is(err, domain.ErrProductAlreadyExists) {
			t.Fatalf("got %v, want ErrProductAlreadyExists", err)
		}
	})

	t.Run("get missing returns ErrProductNotFound", func(t *testing.T) {
		reset()
		if _, err := repo.GetByID(ctx, "33333333-3333-3333-3333-333333333333"); !errors.Is(err, domain.ErrProductNotFound) {
			t.Fatalf("got %v, want ErrProductNotFound", err)
		}
	})

	t.Run("list ordered by id", func(t *testing.T) {
		reset()
		_ = repo.Create(ctx, sampleProduct("cccccccc-cccc-cccc-cccc-cccccccccccc"))
		_ = repo.Create(ctx, sampleProduct("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))
		_ = repo.Create(ctx, sampleProduct("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"))
		list, err := repo.List(ctx)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(list) != 3 || list[0].ID >= list[1].ID || list[1].ID >= list[2].ID {
			t.Fatalf("expected list ordered by id, got %v", list)
		}
	})

	t.Run("update existing and missing", func(t *testing.T) {
		reset()
		p := sampleProduct("44444444-4444-4444-4444-444444444444")
		_ = repo.Create(ctx, p)
		p.Name = "Renamed"
		p.Price = 99
		if err := repo.Update(ctx, p); err != nil {
			t.Fatalf("update: %v", err)
		}
		got, _ := repo.GetByID(ctx, p.ID)
		if got.Name != "Renamed" || got.Price != 99 {
			t.Fatalf("update not applied: %+v", got)
		}
		if err := repo.Update(ctx, sampleProduct("55555555-5555-5555-5555-555555555555")); !errors.Is(err, domain.ErrProductNotFound) {
			t.Fatalf("got %v, want ErrProductNotFound", err)
		}
	})

	t.Run("delete existing and missing", func(t *testing.T) {
		reset()
		p := sampleProduct("66666666-6666-6666-6666-666666666666")
		_ = repo.Create(ctx, p)
		if err := repo.Delete(ctx, p.ID); err != nil {
			t.Fatalf("delete: %v", err)
		}
		if _, err := repo.GetByID(ctx, p.ID); !errors.Is(err, domain.ErrProductNotFound) {
			t.Fatalf("product still present after delete")
		}
		if err := repo.Delete(ctx, "77777777-7777-7777-7777-777777777777"); !errors.Is(err, domain.ErrProductNotFound) {
			t.Fatalf("got %v, want ErrProductNotFound", err)
		}
	})
}

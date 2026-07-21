package main

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/nicolas2601/go-graphql-products-api/internal/config"
)

func TestBuildRepository(t *testing.T) {
	t.Run("memory driver returns a repository", func(t *testing.T) {
		repo, cleanup, err := buildRepository(config.Config{RepoDriver: "memory"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo == nil {
			t.Fatal("expected a non-nil repository")
		}
		cleanup()
	})

	t.Run("postgres driver requires DATABASE_URL", func(t *testing.T) {
		if _, _, err := buildRepository(config.Config{RepoDriver: "postgres"}); err == nil {
			t.Fatal("expected an error when DATABASE_URL is empty")
		}
	})

	t.Run("unknown driver returns an error", func(t *testing.T) {
		if _, _, err := buildRepository(config.Config{RepoDriver: "mongo"}); err == nil {
			t.Fatal("expected an error for an unknown driver")
		}
	})
}

func TestServeGracefulShutdown(t *testing.T) {
	srv := &http.Server{
		Addr:              "127.0.0.1:0",
		Handler:           http.NewServeMux(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() { result <- serve(ctx, srv) }()

	time.Sleep(100 * time.Millisecond) // dar tiempo a que el servidor empiece a escuchar
	cancel()

	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("graceful shutdown should return nil, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("serve did not return after context cancellation")
	}
}

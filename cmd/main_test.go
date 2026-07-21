package main

import (
	"testing"

	"github.com/nicolas2601/go-graphql-products-api/internal/config"
)

func TestBuildRepository(t *testing.T) {
	t.Run("memory driver returns a repository", func(t *testing.T) {
		repo, err := buildRepository(config.Config{RepoDriver: "memory"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if repo == nil {
			t.Fatal("expected a non-nil repository")
		}
	})

	t.Run("postgres driver is not implemented yet", func(t *testing.T) {
		if _, err := buildRepository(config.Config{RepoDriver: "postgres"}); err == nil {
			t.Fatal("expected an error for the postgres driver")
		}
	})

	t.Run("unknown driver returns an error", func(t *testing.T) {
		if _, err := buildRepository(config.Config{RepoDriver: "mongo"}); err == nil {
			t.Fatal("expected an error for an unknown driver")
		}
	})
}

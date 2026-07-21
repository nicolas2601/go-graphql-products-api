package graphql

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nicolas2601/go-graphql-products-api/internal/domain"
)

func TestToGraphQLError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want string
	}{
		{"not found", domain.ErrProductNotFound, "PRODUCT_NOT_FOUND"},
		{"invalid name", domain.ErrInvalidName, "INVALID_NAME"},
		{"invalid price", domain.ErrInvalidPrice, "INVALID_PRICE"},
		{"invalid stock", domain.ErrInvalidStock, "INVALID_STOCK"},
		{"already exists", domain.ErrProductAlreadyExists, "PRODUCT_ALREADY_EXISTS"},
		{"wrapped domain error", fmt.Errorf("create product: %w", domain.ErrInvalidPrice), "INVALID_PRICE"},
		{"unknown error", errors.New("boom"), "INTERNAL_ERROR"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gqlErr := toGraphQLError(c.err)
			if got := gqlErr.Extensions["code"]; got != c.want {
				t.Fatalf("got code %v, want %s", got, c.want)
			}
			if gqlErr.Message == "" {
				t.Fatal("expected a non-empty message")
			}
		})
	}
}

package graphql_test

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	gql "github.com/nicolas2601/go-graphql-products-api/internal/delivery/graphql"
	"github.com/nicolas2601/go-graphql-products-api/internal/delivery/graphql/generated"
	"github.com/nicolas2601/go-graphql-products-api/internal/repository/memory"
	"github.com/nicolas2601/go-graphql-products-api/internal/usecase"
)

func newTestClient() *client.Client {
	repo := memory.New()
	seq := 0
	newID := func() string { seq++; return "id-" + strconv.Itoa(seq) }
	now := func() time.Time { return time.Date(2026, time.July, 20, 12, 0, 0, 0, time.UTC) }
	uc := usecase.NewProductUseCase(repo, newID, now)
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: gql.NewResolver(uc)}))
	return client.New(srv)
}

func TestGraphQLOperations(t *testing.T) {
	t.Run("create then query products and by id", func(t *testing.T) {
		c := newTestClient()

		var created struct {
			CreateProduct struct {
				ID, Name, CreatedAt string
				Price               float64
				Stock               *int
			}
		}
		c.MustPost(`mutation { createProduct(input:{name:"Keyboard", price:49.9, stock:10}) { id name price stock createdAt } }`, &created)
		if created.CreateProduct.ID == "" || created.CreateProduct.Name != "Keyboard" || created.CreateProduct.CreatedAt == "" {
			t.Fatalf("unexpected create result: %+v", created.CreateProduct)
		}

		var list struct {
			Products []struct{ ID, Name string }
		}
		c.MustPost(`{ products { id name } }`, &list)
		if len(list.Products) != 1 {
			t.Fatalf("want 1 product, got %d", len(list.Products))
		}

		id := created.CreateProduct.ID
		var single struct {
			Product struct{ ID, Name string }
		}
		c.MustPost(`{ product(id:"`+id+`") { id name } }`, &single)
		if single.Product.ID != id {
			t.Fatalf("unexpected product: %+v", single.Product)
		}
	})

	t.Run("update and delete", func(t *testing.T) {
		c := newTestClient()

		var created struct{ CreateProduct struct{ ID string } }
		c.MustPost(`mutation { createProduct(input:{name:"Mouse", price:10}) { id } }`, &created)
		id := created.CreateProduct.ID

		var updated struct {
			UpdateProduct struct {
				Name  string
				Price float64
			}
		}
		c.MustPost(`mutation { updateProduct(id:"`+id+`", input:{name:"Gaming Mouse"}) { name price } }`, &updated)
		if updated.UpdateProduct.Name != "Gaming Mouse" || updated.UpdateProduct.Price != 10 {
			t.Fatalf("update failed: %+v", updated.UpdateProduct)
		}

		var deleted struct{ DeleteProduct bool }
		c.MustPost(`mutation { deleteProduct(id:"`+id+`") }`, &deleted)
		if !deleted.DeleteProduct {
			t.Fatal("delete returned false")
		}
	})

	t.Run("errors surface to the client", func(t *testing.T) {
		c := newTestClient()

		var resp struct{ CreateProduct struct{ ID string } }
		if err := c.Post(`mutation { createProduct(input:{name:"", price:10}) { id } }`, &resp); err == nil {
			t.Fatal("expected error for empty name")
		}

		var resp2 struct{ Product *struct{ ID string } }
		if err := c.Post(`{ product(id:"missing") { id } }`, &resp2); err == nil {
			t.Fatal("expected error for missing product")
		}
	})

	t.Run("not-found surfaces PRODUCT_NOT_FOUND code end-to-end", func(t *testing.T) {
		c := newTestClient()
		queries := []string{
			`{ product(id:"missing") { id } }`,
			`mutation { updateProduct(id:"missing", input:{name:"X"}) { id } }`,
			`mutation { deleteProduct(id:"missing") }`,
		}
		for _, q := range queries {
			err := c.Post(q, &struct{}{})
			if err == nil {
				t.Fatalf("query %q: expected a GraphQL error", q)
			}
			// El error del cliente contiene el JSON de la respuesta, incluidas las extensions.
			if !strings.Contains(err.Error(), "PRODUCT_NOT_FOUND") {
				t.Fatalf("query %q: expected code PRODUCT_NOT_FOUND in response, got %v", q, err)
			}
		}
	})
}

package server_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/nicolas2601/go-graphql-products-api/internal/config"
	"github.com/nicolas2601/go-graphql-products-api/internal/repository/memory"
	"github.com/nicolas2601/go-graphql-products-api/internal/server"
	"github.com/nicolas2601/go-graphql-products-api/internal/usecase"
)

func newTestHandler() http.Handler {
	repo := memory.New()
	uc := usecase.NewProductUseCase(repo,
		func() string { return "id-1" },
		func() time.Time { return time.Date(2026, time.July, 21, 0, 0, 0, 0, time.UTC) },
	)
	return server.NewHandler(config.Config{AppEnv: "development", GraphQLPlayground: true}, uc)
}

func TestHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestHandler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "ok") {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestGraphQLEndpoint(t *testing.T) {
	body := `{"query":"mutation { createProduct(input:{name:\"Keyboard\", price:49.9}) { id name } }"}`
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	newTestHandler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want 200: %s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "Keyboard") {
		t.Fatalf("expected the created product in the response, got: %s", rec.Body.String())
	}
}

package server_test

import (
	"encoding/json"
	"fmt"
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

func newHandler(cfg config.Config) http.Handler {
	repo := memory.New()
	uc := usecase.NewProductUseCase(repo,
		func() string { return "id-1" },
		func() time.Time { return time.Date(2026, time.July, 21, 0, 0, 0, 0, time.UTC) },
	)
	return server.NewHandler(cfg, uc)
}

func newTestHandler() http.Handler {
	return newHandler(config.Config{AppEnv: "development", GraphQLPlayground: true})
}

// graphQLRequest arma un POST /query a partir de una operacion GraphQL, escapando el body como JSON.
func graphQLRequest(t *testing.T, query string) *http.Request {
	t.Helper()
	body, err := json.Marshal(map[string]string{"query": query})
	if err != nil {
		t.Fatalf("marshal query: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	return req
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

// En produccion (fail-safe) el playground no debe estar montado: GET / da 404.
func TestPlaygroundDisabledInProduction(t *testing.T) {
	rec := httptest.NewRecorder()
	newHandler(config.Config{AppEnv: "production", GraphQLPlayground: true}).
		ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got status %d, want 404 (playground must be off in production)", rec.Code)
	}
}

// La introspection solo se habilita en desarrollo; en produccion una query de introspection falla.
func TestIntrospectionGatedByEnvironment(t *testing.T) {
	const introspection = "{ __schema { queryType { name } } }"

	t.Run("enabled in development", func(t *testing.T) {
		rec := httptest.NewRecorder()
		newTestHandler().ServeHTTP(rec, graphQLRequest(t, introspection))
		if !strings.Contains(rec.Body.String(), "queryType") {
			t.Fatalf("expected introspection to work in development, got: %s", rec.Body.String())
		}
	})

	t.Run("disabled in production", func(t *testing.T) {
		rec := httptest.NewRecorder()
		newHandler(config.Config{AppEnv: "production"}).ServeHTTP(rec, graphQLRequest(t, introspection))
		body := rec.Body.String()
		if !strings.Contains(body, "errors") || strings.Contains(body, "queryType") {
			t.Fatalf("expected introspection to be rejected in production, got: %s", body)
		}
	})
}

// El limite de complejidad (anti-DoS) rechaza operaciones que exceden el umbral configurado.
func TestComplexityLimitRejectsExpensiveQuery(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("query {")
	// Cada alias "products { id }" suma complejidad; 150 aliases superan holgadamente el limite de 200.
	for i := 0; i < 150; i++ {
		fmt.Fprintf(&sb, " a%d: products { id }", i)
	}
	sb.WriteString(" }")

	rec := httptest.NewRecorder()
	newTestHandler().ServeHTTP(rec, graphQLRequest(t, sb.String()))

	if !strings.Contains(strings.ToLower(rec.Body.String()), "complexity") {
		t.Fatalf("expected a complexity-limit error, got: %s", rec.Body.String())
	}
}

// http.MaxBytesHandler acota el tamano del body: un payload por encima de 1 MiB se rechaza.
func TestRequestBodyLimitRejectsHugePayload(t *testing.T) {
	huge := `{"query":"` + strings.Repeat("#", (1<<20)+1024) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(huge))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	newTestHandler().ServeHTTP(rec, req)

	body := strings.ToLower(rec.Body.String())
	if rec.Code == http.StatusOK && !strings.Contains(body, "error") {
		t.Fatalf("expected an oversized body to be rejected, got status %d: %s", rec.Code, rec.Body.String())
	}
}

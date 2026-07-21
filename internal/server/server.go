// Package server arma el handler HTTP de la aplicacion: endpoint GraphQL, health check y,
// en desarrollo, el playground de GraphQL.
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/nicolas2601/go-graphql-products-api/internal/config"
	graphqldelivery "github.com/nicolas2601/go-graphql-products-api/internal/delivery/graphql"
	"github.com/nicolas2601/go-graphql-products-api/internal/delivery/graphql/generated"
	"github.com/nicolas2601/go-graphql-products-api/internal/usecase"
)

const (
	graphQLPath = "/query"
	// maxRequestBytes limita el tamano del body para evitar payloads abusivos.
	maxRequestBytes = 1 << 20 // 1 MiB
	// queryComplexityLimit acota la complejidad de una operacion GraphQL (anti-DoS).
	queryComplexityLimit = 200
)

// NewHandler arma el handler HTTP cableando el caso de uso a los resolvers de GraphQL.
// El playground y la introspection quedan habilitados solo en modo desarrollo (fail-safe).
func NewHandler(cfg config.Config, products *usecase.ProductUseCase) http.Handler {
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: graphqldelivery.NewResolver(products)})

	gql := handler.New(schema)
	gql.AddTransport(transport.POST{})
	gql.Use(extension.FixedComplexityLimit(queryComplexityLimit))
	if cfg.IsDevelopment() {
		gql.Use(extension.Introspection{})
	}

	mux := http.NewServeMux()
	mux.Handle(graphQLPath, http.MaxBytesHandler(gql, maxRequestBytes))
	mux.HandleFunc("/healthz", handleHealth)
	if cfg.IsDevelopment() && cfg.GraphQLPlayground {
		mux.Handle("/", playground.Handler("Products API", graphQLPath))
	}
	return mux
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		slog.Error("failed to write health response", "error", err)
	}
}

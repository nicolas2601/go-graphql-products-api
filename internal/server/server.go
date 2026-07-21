// Package server arma el handler HTTP de la aplicacion: endpoint GraphQL, health check y,
// en desarrollo, el playground de GraphQL.
package server

import (
	"encoding/json"
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

const graphQLPath = "/query"

// NewHandler arma el handler HTTP cableando el caso de uso a los resolvers de GraphQL.
// El playground y la introspection quedan habilitados solo fuera de production.
func NewHandler(cfg config.Config, products *usecase.ProductUseCase) http.Handler {
	schema := generated.NewExecutableSchema(generated.Config{Resolvers: graphqldelivery.NewResolver(products)})

	gql := handler.New(schema)
	gql.AddTransport(transport.POST{})
	if cfg.AppEnv != "production" {
		gql.Use(extension.Introspection{})
	}

	mux := http.NewServeMux()
	mux.Handle(graphQLPath, gql)
	mux.HandleFunc("/healthz", handleHealth)
	if cfg.GraphQLPlayground && cfg.AppEnv != "production" {
		mux.Handle("/", playground.Handler("Products API", graphQLPath))
	}
	return mux
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

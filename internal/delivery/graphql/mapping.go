package graphql

import (
	"time"

	"github.com/nicolas2601/go-graphql-products-api/internal/delivery/graphql/model"
	"github.com/nicolas2601/go-graphql-products-api/internal/domain"
)

// toModel traduce una entidad de dominio al modelo de GraphQL. createdAt se expone en formato
// RFC 3339; el dominio no conoce el formato de presentacion.
func toModel(product domain.Product) *model.Product {
	stock := product.Stock
	return &model.Product{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     &stock,
		CreatedAt: product.CreatedAt.Format(time.RFC3339),
	}
}

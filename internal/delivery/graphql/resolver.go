package graphql

import (
	"context"

	"github.com/nicolas2601/go-graphql-products-api/internal/delivery/graphql/generated"
	"github.com/nicolas2601/go-graphql-products-api/internal/delivery/graphql/model"
	"github.com/nicolas2601/go-graphql-products-api/internal/usecase"
)

// Resolver es la raiz de los resolvers de GraphQL. No contiene logica de negocio: delega en
// el caso de uso y traduce entre GraphQL y el dominio.
type Resolver struct {
	products *usecase.ProductUseCase
}

// NewResolver construye el resolver raiz con el caso de uso inyectado.
func NewResolver(products *usecase.ProductUseCase) *Resolver {
	return &Resolver{products: products}
}

// CreateProduct is the resolver for the createProduct field.
func (r *mutationResolver) CreateProduct(ctx context.Context, input model.CreateProductInput) (*model.Product, error) {
	stock := 0
	if input.Stock != nil {
		stock = *input.Stock
	}
	product, err := r.products.Create(ctx, input.Name, input.Price, stock)
	if err != nil {
		return nil, toGraphQLError(err)
	}
	return toModel(product), nil
}

// UpdateProduct is the resolver for the updateProduct field.
func (r *mutationResolver) UpdateProduct(ctx context.Context, id string, input model.UpdateProductInput) (*model.Product, error) {
	product, err := r.products.Update(ctx, id, input.Name, input.Price)
	if err != nil {
		return nil, toGraphQLError(err)
	}
	return toModel(product), nil
}

// DeleteProduct is the resolver for the deleteProduct field.
func (r *mutationResolver) DeleteProduct(ctx context.Context, id string) (bool, error) {
	if err := r.products.Delete(ctx, id); err != nil {
		return false, toGraphQLError(err)
	}
	return true, nil
}

// Products is the resolver for the products field.
func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	products, err := r.products.List(ctx)
	if err != nil {
		return nil, toGraphQLError(err)
	}
	result := make([]*model.Product, 0, len(products))
	for _, product := range products {
		result = append(result, toModel(product))
	}
	return result, nil
}

// Product is the resolver for the product field.
func (r *queryResolver) Product(ctx context.Context, id string) (*model.Product, error) {
	product, err := r.products.Get(ctx, id)
	if err != nil {
		return nil, toGraphQLError(err)
	}
	return toModel(product), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type (
	mutationResolver struct{ *Resolver }
	queryResolver    struct{ *Resolver }
)

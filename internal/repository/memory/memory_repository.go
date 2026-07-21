// Package memory provee una implementacion en memoria de domain.ProductRepository,
// segura para uso concurrente.
package memory

import (
	"context"
	"sort"
	"sync"

	"github.com/nicolas2601/go-graphql-products-api/internal/domain"
)

// Repository almacena productos en memoria, protegidos por un RWMutex.
type Repository struct {
	mu    sync.RWMutex
	items map[string]domain.Product
}

// New crea un repositorio en memoria vacio.
func New() *Repository {
	return &Repository{items: make(map[string]domain.Product)}
}

// Create inserta un producto nuevo. Devuelve ErrProductAlreadyExists si el id ya existe.
func (r *Repository) Create(_ context.Context, product domain.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.items[product.ID]; exists {
		return domain.ErrProductAlreadyExists
	}
	r.items[product.ID] = product
	return nil
}

// GetByID devuelve el producto con el id dado, o ErrProductNotFound si no existe.
func (r *Repository) GetByID(_ context.Context, id string) (domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	product, ok := r.items[id]
	if !ok {
		return domain.Product{}, domain.ErrProductNotFound
	}
	return product, nil
}

// List devuelve todos los productos ordenados por id.
func (r *Repository) List(_ context.Context) ([]domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	products := make([]domain.Product, 0, len(r.items))
	for _, product := range r.items {
		products = append(products, product)
	}
	sort.Slice(products, func(i, j int) bool {
		return products[i].ID < products[j].ID
	})
	return products, nil
}

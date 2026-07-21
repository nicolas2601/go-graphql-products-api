package domain

import "context"

// ProductRepository define el contrato de persistencia de productos. Vive en la capa
// de dominio: los casos de uso dependen de esta interfaz, no de una implementacion
// concreta. Las implementaciones (memoria, PostgreSQL) viven en la capa de repositorio.
type ProductRepository interface {
	// Create persiste un producto nuevo. Debe devolver ErrProductAlreadyExists si ya existe
	// un producto con el mismo id.
	Create(ctx context.Context, product Product) error
	// GetByID devuelve el producto con el id dado, o ErrProductNotFound si no existe.
	GetByID(ctx context.Context, id string) (Product, error)
	// List devuelve todos los productos.
	List(ctx context.Context) ([]Product, error)
	// Update reemplaza un producto existente, o devuelve ErrProductNotFound si no existe.
	Update(ctx context.Context, product Product) error
	// Delete elimina el producto con el id dado, o devuelve ErrProductNotFound si no existe.
	Delete(ctx context.Context, id string) error
}

// Package usecase contiene la logica de aplicacion. Depende de las interfaces del
// dominio, nunca de una implementacion concreta del repositorio ni de la capa de delivery.
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/nicolas2601/go-graphql-products-api/internal/domain"
)

// ProductUseCase orquesta las operaciones de negocio sobre productos.
type ProductUseCase struct {
	repo  domain.ProductRepository
	newID func() string
	now   func() time.Time
}

// NewProductUseCase construye el caso de uso con sus dependencias inyectadas.
// newID y now se inyectan para hacer las operaciones deterministas en los tests.
func NewProductUseCase(repo domain.ProductRepository, newID func() string, now func() time.Time) *ProductUseCase {
	if repo == nil || newID == nil || now == nil {
		// Fail-fast: dependencias faltantes son un error de programacion en el wiring,
		// no una condicion de ejecucion. Mejor romper en el arranque que con un nil
		// pointer lejano en la primera operacion.
		panic("usecase: repo, newID and now must not be nil")
	}
	return &ProductUseCase{repo: repo, newID: newID, now: now}
}

// Create valida y persiste un producto nuevo, asignando id y createdAt.
func (uc *ProductUseCase) Create(ctx context.Context, name string, price float64, stock int) (domain.Product, error) {
	product, err := domain.NewProduct(uc.newID(), name, price, stock, uc.now())
	if err != nil {
		return domain.Product{}, err
	}
	if err := uc.repo.Create(ctx, product); err != nil {
		return domain.Product{}, fmt.Errorf("create product: %w", err)
	}
	return product, nil
}

// Get devuelve un producto por id, o ErrProductNotFound si no existe.
func (uc *ProductUseCase) Get(ctx context.Context, id string) (domain.Product, error) {
	return uc.repo.GetByID(ctx, id)
}

// List devuelve todos los productos.
func (uc *ProductUseCase) List(ctx context.Context) ([]domain.Product, error) {
	return uc.repo.List(ctx)
}

// Update aplica los campos enviados (nombre y/o precio), valida el resultado y persiste.
// Los punteros nil representan campos no enviados, que se dejan sin cambios. Segun el
// enunciado, solo se actualizan nombre y precio; el stock no es modificable por esta via.
func (uc *ProductUseCase) Update(ctx context.Context, id string, name *string, price *float64) (domain.Product, error) {
	product, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}
	if name != nil {
		product.Name = *name
	}
	if price != nil {
		product.Price = *price
	}
	if err := product.Validate(); err != nil {
		return domain.Product{}, err
	}
	if err := uc.repo.Update(ctx, product); err != nil {
		return domain.Product{}, fmt.Errorf("update product: %w", err)
	}
	return product, nil
}

// Delete elimina un producto por id, o devuelve ErrProductNotFound si no existe.
func (uc *ProductUseCase) Delete(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}

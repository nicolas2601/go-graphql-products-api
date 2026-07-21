// Package domain contiene las entidades del negocio, sus invariantes y las
// interfaces de repositorio. No depende de ningun detalle de infraestructura.
package domain

import (
	"strings"
	"time"
)

// Product es la entidad central del sistema.
type Product struct {
	ID        string
	Name      string
	Price     float64
	Stock     int
	CreatedAt time.Time
}

// NewProduct construye un producto valido. El id y createdAt los provee el caso de
// uso; el constructor solo valida los invariantes de negocio.
func NewProduct(id, name string, price float64, stock int, createdAt time.Time) (Product, error) {
	p := Product{ID: id, Name: name, Price: price, Stock: stock, CreatedAt: createdAt}
	if err := p.Validate(); err != nil {
		return Product{}, err
	}
	return p, nil
}

// Validate verifica los invariantes de negocio del producto.
func (p Product) Validate() error {
	if strings.TrimSpace(p.Name) == "" {
		return ErrInvalidName
	}
	if p.Price <= 0 {
		return ErrInvalidPrice
	}
	return nil
}

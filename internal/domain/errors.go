package domain

import "errors"

// Errores de dominio tipados. Las capas externas los mapean a sus propias
// representaciones (por ejemplo, extensiones de error en GraphQL).
var (
	// ErrProductNotFound indica que no existe un producto con el identificador dado.
	ErrProductNotFound = errors.New("product not found")
	// ErrInvalidName indica que el nombre del producto esta vacio.
	ErrInvalidName = errors.New("product name must not be empty")
	// ErrInvalidPrice indica que el precio del producto no es mayor a cero.
	ErrInvalidPrice = errors.New("product price must be greater than zero")
	// ErrInvalidStock indica que el stock del producto es negativo.
	ErrInvalidStock = errors.New("product stock must not be negative")
	// ErrProductAlreadyExists indica que ya existe un producto con el identificador dado.
	ErrProductAlreadyExists = errors.New("product already exists")
)

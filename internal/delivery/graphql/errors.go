package graphql

import (
	"errors"

	"github.com/nicolas2601/go-graphql-products-api/internal/domain"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Codigos de error expuestos en las extensions de GraphQL.
const (
	codeProductNotFound      = "PRODUCT_NOT_FOUND"
	codeProductAlreadyExists = "PRODUCT_ALREADY_EXISTS"
	codeInvalidName          = "INVALID_NAME"
	codeInvalidPrice         = "INVALID_PRICE"
	codeInvalidStock         = "INVALID_STOCK"
	codeInternal             = "INTERNAL_ERROR"
)

// toGraphQLError traduce un error de dominio tipado a un error de GraphQL con extensions.code.
// Usa errors.Is, por lo que reconoce errores envueltos por las capas superiores.
func toGraphQLError(err error) *gqlerror.Error {
	return &gqlerror.Error{
		Message: err.Error(),
		Extensions: map[string]interface{}{
			"code": codeFor(err),
		},
	}
}

func codeFor(err error) string {
	switch {
	case errors.Is(err, domain.ErrProductNotFound):
		return codeProductNotFound
	case errors.Is(err, domain.ErrProductAlreadyExists):
		return codeProductAlreadyExists
	case errors.Is(err, domain.ErrInvalidName):
		return codeInvalidName
	case errors.Is(err, domain.ErrInvalidPrice):
		return codeInvalidPrice
	case errors.Is(err, domain.ErrInvalidStock):
		return codeInvalidStock
	default:
		return codeInternal
	}
}

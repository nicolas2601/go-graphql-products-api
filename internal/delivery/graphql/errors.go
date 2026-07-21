package graphql

import (
	"errors"
	"log/slog"

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
	code := codeFor(err)
	message := err.Error()
	if code == codeInternal {
		// No se filtran detalles internos al cliente: se loguea el error real y se devuelve un
		// mensaje generico. Los errores de dominio conocidos si exponen su mensaje de negocio.
		slog.Error("unhandled resolver error", "error", err)
		message = "internal server error"
	}
	return &gqlerror.Error{
		Message:    message,
		Extensions: map[string]interface{}{"code": code},
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

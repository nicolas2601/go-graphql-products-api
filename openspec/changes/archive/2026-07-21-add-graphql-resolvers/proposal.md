## Why

Los resolvers de GraphQL generados por gqlgen son stubs que hacen `panic("not implemented")`,
por lo que la API no puede responder ninguna query ni mutation. Esta capa conecta el esquema
con los casos de uso ya implementados para que el servidor funcione de punta a punta.

## What Changes

- Implementar los cinco resolvers (`products`, `product`, `createProduct`, `updateProduct`,
  `deleteProduct`) delegando en `ProductUseCase`, sin logica de negocio en el resolver.
- Mapear los errores de dominio tipados a errores de GraphQL con `extensions.code` y `message`.
- Mapear la entidad `domain.Product` al modelo de GraphQL (`createdAt` en formato RFC 3339).
- Inyectar el caso de uso en el `Resolver` mediante un constructor.

## Capabilities

### New Capabilities
- `graphql-api`: comportamiento observable de las operaciones GraphQL de productos (queries,
  mutations y mapeo de errores de dominio a codigos de GraphQL).

### Modified Capabilities

## Impact

- Codigo: `internal/delivery/graphql` (resolver.go + helpers de mapeo de errores y modelos).
- Depende de `internal/usecase` y `internal/domain`. No cambia el esquema GraphQL.

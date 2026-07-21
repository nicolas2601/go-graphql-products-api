## 1. Implementation

- [x] 1.1 Inyectar `ProductUseCase` en el `Resolver` via `NewResolver`
- [x] 1.2 Implementar los resolvers de query (`products`, `product`)
- [x] 1.3 Implementar los resolvers de mutation (`createProduct`, `updateProduct`, `deleteProduct`)
- [x] 1.4 Agregar el helper de mapeo de errores de dominio a GraphQL (`extensions.code`)
- [x] 1.5 Agregar el helper de mapeo de `domain.Product` al modelo de GraphQL

## 2. Verification

- [x] 2.1 Test de integracion ejecutando las operaciones GraphQL contra el esquema
- [x] 2.2 `go test ./... -race`, `go build ./...` y `go vet ./...` en verde

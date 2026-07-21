## 1. Implementation

- [ ] 1.1 Inyectar `ProductUseCase` en el `Resolver` via `NewResolver`
- [ ] 1.2 Implementar los resolvers de query (`products`, `product`)
- [ ] 1.3 Implementar los resolvers de mutation (`createProduct`, `updateProduct`, `deleteProduct`)
- [ ] 1.4 Agregar el helper de mapeo de errores de dominio a GraphQL (`extensions.code`)
- [ ] 1.5 Agregar el helper de mapeo de `domain.Product` al modelo de GraphQL

## 2. Verification

- [ ] 2.1 Test de integracion ejecutando las operaciones GraphQL contra el esquema
- [ ] 2.2 `go test ./... -race`, `go build ./...` y `go vet ./...` en verde

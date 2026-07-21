## 1. Implementation

- [ ] 1.1 `postgres.Repository` (pgxpool) con las 5 operaciones + mapeo de errores
- [ ] 1.2 Migracion embebida (tabla `products`) + `Migrate(ctx, pool)` idempotente
- [ ] 1.3 Cablear `buildRepository`: conectar `DATABASE_URL`, migrar y devolver el repo PostgreSQL

## 2. Verification

- [ ] 2.1 Tests de integracion con testcontainers (persistir/recuperar, duplicado, not found, update, delete, orden)
- [ ] 2.2 `go test ./... -race`, `go build ./...` y `go vet ./...` en verde

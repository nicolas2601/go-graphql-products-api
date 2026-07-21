## 1. Implementation

- [x] 1.1 `postgres.Repository` (pgxpool) con las 5 operaciones + mapeo de errores
- [x] 1.2 Migracion embebida (tabla `products`) + `Migrate(ctx, pool)` idempotente
- [x] 1.3 Cablear `buildRepository`: conectar `DATABASE_URL`, migrar y devolver el repo PostgreSQL

## 2. Verification

- [x] 2.1 Tests de integracion con testcontainers (persistir/recuperar, duplicado, not found, update, delete, orden)
- [x] 2.2 `go test ./... -race`, `go build ./...` y `go vet ./...` en verde
- [x] 2.3 Verificacion end-to-end: servidor con REPO_DRIVER=postgres persiste en Postgres real (conteo directo confirmado)

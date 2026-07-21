## Why

El sistema ya soporta `REPO_DRIVER=postgres` pero hoy devuelve "not implemented yet". Esta
feature agrega la implementacion PostgreSQL del repositorio de productos, intercambiable por
inyeccion de dependencias, cumpliendo el deseable del enunciado.

## What Changes

- `internal/repository/postgres`: implementacion de `domain.ProductRepository` sobre `pgx`
  (pgxpool).
- Esquema de la tabla `products` aplicado por una migracion embebida.
- `buildRepository` conecta a `DATABASE_URL`, aplica la migracion y devuelve el repositorio
  PostgreSQL cuando `REPO_DRIVER=postgres`.
- Tests de integracion con testcontainers-go (Postgres efimero en Docker).

## Capabilities

### New Capabilities

### Modified Capabilities
- `product-storage`: agrega el requerimiento de persistencia en PostgreSQL.

## Impact

- Nuevas dependencias: `jackc/pgx/v5` (runtime) y `testcontainers-go` (solo tests). Se cablea en
  `cmd/main.go`. Requiere Docker corriendo para los tests de integracion.

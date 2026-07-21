## Context

El contrato del repositorio (`product-storage`) ya existe con la implementacion en memoria.
Falta la implementacion PostgreSQL para el modo `REPO_DRIVER=postgres`.

## Decisions

### pgxpool
Se usa `jackc/pgx/v5/pgxpool` por rendimiento y connection pooling. `postgres.New(pool)`
implementa `domain.ProductRepository`.

### Migraciones
El esquema (tabla `products`: id UUID PK, name TEXT, price DOUBLE PRECISION, stock INTEGER,
created_at TIMESTAMPTZ) se embebe con `go:embed` y se aplica con una funcion `Migrate(ctx, pool)`
idempotente (`CREATE TABLE IF NOT EXISTS`). Se usa DOUBLE PRECISION para mapear limpio al
`float64` del dominio (heredado del schema GraphQL); para dinero exacto lo correcto seria
NUMERIC, queda como deuda documentada. Simple y suficiente para una sola tabla; en la Prueba 2
se usaria goose/migrate versionado.

### Mapeo de errores
Violacion de la primary key de id -> `ErrProductAlreadyExists`; `pgx.ErrNoRows` ->
`ErrProductNotFound`. Se mantiene el mismo contrato que el repositorio en memoria.

### Tests
testcontainers-go levanta un Postgres efimero, aplica la migracion y corre la misma bateria de
comportamiento que el repositorio en memoria (persistir/recuperar, duplicado, not found, update,
delete, listado ordenado). Si Docker no esta disponible, el test hace `t.Skip` con un mensaje
claro en vez de fallar.

## Testing

Integracion con testcontainers (requiere Docker). El resto de la suite del proyecto no se ve
afectada.

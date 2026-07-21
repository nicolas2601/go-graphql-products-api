## Why

El proyecto compila y tiene tests, pero NO es ejecutable: no hay composition root ni servidor
HTTP. Esta feature agrega el binario que arranca la API, cablea las dependencias y expone el
endpoint GraphQL, el health check y el playground.

## What Changes

- `cmd/main.go` como composition root: carga la configuracion, arma el repositorio segun
  `REPO_DRIVER`, los casos de uso y los resolvers, y levanta el servidor HTTP.
- Seleccion de repositorio por `REPO_DRIVER` (memory disponible; postgres queda para una feature
  posterior y por ahora devuelve un error explicito).
- Endpoint GraphQL en `/query`, health check en `/healthz`, y playground gateado por entorno.
- Apagado ordenado (graceful shutdown) ante SIGINT/SIGTERM. Logging estructurado con `slog`.

## Capabilities

### New Capabilities
- `server-runtime`: arranque del servidor, inyeccion de dependencias, health check y apagado ordenado.

### Modified Capabilities

## Impact

- Nuevo `cmd/main.go` e `internal/server`. Reusa config, usecase, repository/memory y
  delivery/graphql ya existentes. Depende de `google/uuid` (ya en el modulo) para generar ids.

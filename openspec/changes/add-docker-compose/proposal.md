## Why

El proyecto es ejecutable pero no esta empaquetado para desplegarse ni para correrse de forma
reproducible con su base de datos. Esta feature agrega la imagen Docker y un docker-compose que
levanta la API conectada a PostgreSQL.

## What Changes

- Dockerfile multi-stage: build en `golang:alpine`, runtime minimo en `alpine` con usuario
  no-root y healthcheck sobre `/healthz`.
- `docker-compose.yml`: servicio `api` (`REPO_DRIVER=postgres`) + servicio `db`
  (`postgres:16-alpine`) con healthcheck y volumen persistente; el `api` espera a que la `db`
  este healthy.
- `.dockerignore` para un contexto de build chico.

## Capabilities

### New Capabilities
- `deployment`: ejecucion del sistema como contenedor y via docker-compose con PostgreSQL.

### Modified Capabilities

## Impact

- Nuevos `Dockerfile`, `docker-compose.yml`, `.dockerignore`. No cambia codigo de aplicacion.

## Context

El binario ya arranca y soporta `REPO_DRIVER=memory|postgres`. Falta empaquetarlo y orquestarlo
con su base de datos para un despliegue/demo reproducible.

## Decisions

### Dockerfile multi-stage
Etapa build: `golang:1.26-alpine` compila un binario estatico (`CGO_ENABLED=0`, `-trimpath`).
Etapa runtime: `alpine:3.20` con usuario no-root, `ca-certificates` y `HEALTHCHECK` via `wget`
sobre `/healthz`. Se elige alpine (no distroless) para poder incluir el healthcheck; distroless
seria una alternativa mas minima pero sin shell ni herramientas para el healthcheck.

### docker-compose
El servicio `api` corre en `REPO_DRIVER=postgres` apuntando al servicio `db` por el DNS de red de
compose (host `db`). `depends_on` con `condition: service_healthy` asegura que la api arranque
despues de que Postgres este listo. El puerto de la db se expone solo en `127.0.0.1`. Volumen
nombrado para persistencia. `APP_ENV=development` para exponer el playground en la demo.

### Dos modos de ejecucion
Con Postgres: `docker compose up`. Sin Postgres (en memoria): correr el binario con
`REPO_DRIVER=memory`. Ambos se documentan en el README (feature posterior).

## Testing

Verificacion manual: `docker build` de la imagen; `docker compose up`; la API responde `/healthz`
y una operacion GraphQL persiste en el Postgres del compose; `docker compose down` limpia.

## Limitaciones conocidas (deuda documentada)

- La migracion embebida (`CREATE TABLE IF NOT EXISTS`) no es 100% segura con varias replicas del
  api arrancando en paralelo (race en el catalogo de Postgres). Con 1 replica no aplica; la
  solucion futura seria un advisory lock o un migrator dedicado.
- `/healthz` es un liveness check estatico; no verifica la conectividad a la base. Un `/readyz`
  con `pool.Ping` seria la evolucion natural.
- Endurecimiento runtime adicional (`read_only`, `cap_drop`, `no-new-privileges`, limites de
  recursos) y escaneo de vulnerabilidades en CI quedan fuera del alcance de esta feature.

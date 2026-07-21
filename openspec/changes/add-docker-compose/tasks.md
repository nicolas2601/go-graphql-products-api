## 1. Implementation

- [x] 1.1 Dockerfile multi-stage (build `golang:alpine`; runtime `alpine` no-root + healthcheck)
- [x] 1.2 `docker-compose.yml` (api `REPO_DRIVER=postgres` + db postgres healthy + volumen; db interna, api port configurable)
- [x] 1.3 `.dockerignore`

## 2. Verification

- [x] 2.1 `docker build` de la imagen sin errores
- [x] 2.2 `docker compose up`: la db queda healthy, la api arranca despues, `/healthz` responde y `createProduct` persiste en el Postgres del compose (conteo directo confirmado)
- [x] 2.3 `go build`/`vet`/`test` siguen en verde (sin cambios de codigo)

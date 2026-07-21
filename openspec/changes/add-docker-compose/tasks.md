## 1. Implementation

- [ ] 1.1 Dockerfile multi-stage (build `golang:alpine`; runtime `alpine` no-root + healthcheck)
- [ ] 1.2 `docker-compose.yml` (api `REPO_DRIVER=postgres` + db postgres healthy + volumen)
- [ ] 1.3 `.dockerignore`

## 2. Verification

- [ ] 2.1 `docker build` de la imagen sin errores
- [ ] 2.2 `docker compose up`: `/healthz` responde y una operacion GraphQL persiste en Postgres
- [ ] 2.3 `go build`/`vet`/`test` siguen en verde (sin cambios de codigo)

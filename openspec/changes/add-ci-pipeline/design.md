## Context

El proyecto tiene una suite completa (unit + integracion) pero sin automatizacion. Falta un
pipeline que la corra en cada cambio y bloquee regresiones.

## Decisions

### Jobs separados
`lint`, `test` y `docker` en jobs paralelos para feedback claro y rapido. Triggers en `push` y
`pull_request` a `main` y `develop`.

### golangci-lint via go install
Se instala `golangci-lint v2.12.2` con `go install` (compilado contra el Go del runner tomado de
`go.mod`) en vez de la action, para evitar incompatibilidades de version entre la action y
golangci-lint v2.

### Tests de Postgres en CI
`go test ./... -race` corre tambien los tests de integracion con testcontainers; los runners
`ubuntu-latest` traen Docker, asi que funcionan sin configuracion extra.

### go mod tidy check
Se verifica que `go.mod`/`go.sum` esten ordenados (`go mod tidy` + `git diff --exit-code`).

### Escaneo de vulnerabilidades informativo
`govulncheck` corre como paso informativo (`continue-on-error`) para no bloquear el pipeline por
vulnerabilidades en dependencias transitivas; se puede endurecer a bloqueante mas adelante.

## Testing

El workflow debe correr VERDE en GitHub Actions (se confirma con `gh run`/`gh pr checks`), no solo
existir.

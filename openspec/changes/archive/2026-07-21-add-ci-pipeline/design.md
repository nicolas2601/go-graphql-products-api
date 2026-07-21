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

### Escaneo de vulnerabilidades bloqueante
`govulncheck` (pineado `@v1.6.0`) corre en su propio job y ES bloqueante. Las vulnerabilidades
detectadas eran de la biblioteca estandar de Go 1.26.3 (GO-2026-5037/5039/5856), no de
dependencias; se resolvieron bumpeando el toolchain a `go1.26.5` (`toolchain go1.26.5` en
`go.mod`), con lo que el scan queda limpio y el gate tiene sentido.

### Endurecimiento del pipeline
- Actions pineadas por SHA (checkout/setup-go/upload-artifact v7, sobre Node 24) para
  reproducibilidad y superficie de supply chain.
- `concurrency` con `cancel-in-progress` para no acumular runs viejos.
- `timeout-minutes` por job para no quemar runners ante un cuelgue.
- Cobertura reportada (`go tool cover -func`) y subida como artifact.

## Testing

El workflow debe correr VERDE en GitHub Actions (se confirma con `gh run`/`gh pr checks`), no solo
existir. Verificado: run 29852440901 con lint, test, vuln y docker en success.

## Why

No hay verificacion automatica: cada cambio depende de correr los checks a mano. Esta feature
agrega un pipeline de CI que corre lint, vet, build, tests (con race) y el build de la imagen
Docker en cada push y pull request a `main` y `develop`.

## What Changes

- `.github/workflows/ci.yml` con jobs: `lint` (golangci-lint), `test` (vet + build + go mod tidy
  check + `go test -race` + govulncheck informativo) y `docker` (docker build).
- `.golangci.yml` con el set estandar de linters, excluyendo el codigo generado.

## Capabilities

### New Capabilities
- `continuous-integration`: verificacion automatica de calidad en cada cambio.

### Modified Capabilities

## Impact

- Nuevos `.github/workflows/ci.yml` y `.golangci.yml`. Los tests de Postgres corren via
  testcontainers sobre el Docker del runner. No cambia codigo de aplicacion (salvo un fix de formato).

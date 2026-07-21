## 1. Implementation

- [x] 1.1 `internal/server.NewHandler` con `/query`, `/healthz` y playground gateado por entorno
- [x] 1.2 `buildRepository` segun `REPO_DRIVER` (memory; postgres -> error explicito)
- [x] 1.3 `cmd/main.go`: composition root + `run()` con apagado ordenado y logging `slog`

## 2. Verification

- [x] 2.1 Tests: `httptest` de `/healthz` y `/query`; `buildRepository` (memory/postgres/desconocido)
- [x] 2.2 `go test ./... -race`, `go build ./...` y `go vet ./...` en verde
- [x] 2.3 Verificacion manual: binario levantado; `/healthz` responde 200 y `createProduct`/`products` funcionan

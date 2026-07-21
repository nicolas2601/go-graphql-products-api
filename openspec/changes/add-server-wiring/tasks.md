## 1. Implementation

- [ ] 1.1 `internal/server.NewHandler` con `/query`, `/healthz` y playground gateado por entorno
- [ ] 1.2 `buildRepository` segun `REPO_DRIVER` (memory; postgres -> error explicito)
- [ ] 1.3 `cmd/main.go`: composition root + `run()` con apagado ordenado y logging `slog`

## 2. Verification

- [ ] 2.1 Tests: `httptest` de `/healthz` y `/query`; `buildRepository` (memory/postgres/desconocido)
- [ ] 2.2 `go test ./... -race`, `go build ./...` y `go vet ./...` en verde
- [ ] 2.3 Verificacion manual: levantar el binario y consultar `/healthz` y una operacion GraphQL

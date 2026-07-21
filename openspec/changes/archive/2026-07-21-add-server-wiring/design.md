## Context

Todas las capas internas existen y estan probadas. Falta el composition root que las conecte y
un servidor HTTP que las exponga.

## Decisions

### Composition root en cmd/main.go
`main()` delega en `run() error` para ser verificable. `run()` carga la config, construye el
repositorio, el caso de uso (con `uuid.NewString` y `time.Now` inyectados), el handler HTTP y
sirve con apagado ordenado.

### Handler HTTP en internal/server
`NewHandler(cfg, uc)` arma un `http.ServeMux` con: `/query` (GraphQL via `handler.New` +
`transport.POST`, no el `NewDefaultServer` deprecado), `/healthz` (JSON 200) y, solo fuera de
production y si esta habilitado, el playground e introspection. Es testeable con `httptest` sin
levantar un puerto real.

### Seleccion de repositorio
`buildRepository(cfg)` mapea `REPO_DRIVER`: `memory` -> repositorio en memoria; `postgres` ->
error "no implementado aun" (hasta la feature de PostgreSQL); desconocido -> error. Testeable.

### Apagado ordenado
`signal.NotifyContext` + `srv.Shutdown` con timeout. Un canal captura el error de
`ListenAndServe` para no tragarse un fallo de arranque (por ejemplo, puerto ocupado).

## Testing

- `internal/server`: `httptest` sobre `/healthz` y `/query`.
- `cmd`: `buildRepository` (memory ok, postgres error, desconocido error).
- Verificacion manual: levantar el binario y consultar `/healthz` y una operacion GraphQL.

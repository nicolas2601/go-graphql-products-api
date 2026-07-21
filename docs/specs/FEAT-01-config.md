# FEAT-01 — Configuracion y entorno

## Objetivo

Centralizar la configuracion de la aplicacion en un unico paquete que lee variables de
entorno con valores por defecto sensatos, de forma testeable y sin acoplar el resto del
codigo a `os.Getenv`.

## Alcance

Incluye:
- Paquete `internal/config` con una funcion pura de carga.
- Archivo `.env.example` versionado y documentado.
- Confirmacion de que `.env` queda ignorado por git.

No incluye:
- Uso de la configuracion en el servidor (ocurre en FEAT-06, wiring).

## Variables de entorno

| Variable | Default | Uso |
|---|---|---|
| `PORT` | `8080` | Puerto del servidor HTTP |
| `APP_ENV` | `development` | Modo de ejecucion (gatea playground e introspection) |
| `GRAPHQL_PLAYGROUND` | `true` | Habilita el playground de GraphQL |
| `LOG_LEVEL` | `info` | Nivel de logging estructurado |
| `REPO_DRIVER` | `memory` | Selecciona el repositorio (`memory` o `postgres`) por inyeccion de dependencias |
| `DATABASE_URL` | vacio | Cadena de conexion, solo para el driver `postgres` |

## Criterios de aceptacion (Given-When-Then)

- Dado un entorno sin variables definidas, cuando se carga la configuracion, entonces cada
  campo toma su valor por defecto.
- Dado un entorno con variables definidas, cuando se carga la configuracion, entonces cada
  campo toma el valor del entorno.
- Dado `GRAPHQL_PLAYGROUND` con un valor no booleano, cuando se carga la configuracion,
  entonces el campo toma su valor por defecto en lugar de fallar.

## Diseno

- La funcion de carga recibe una funcion de busqueda (`func(string) string`) inyectada, de
  modo que sea pura y testeable sin tocar el entorno real del proceso.
- `Config` es una estructura de valores simple, sin comportamiento.

## Definition of Done

- Tests de la carga de configuracion en verde (defaults, overrides y booleano invalido).
- `go build ./...` y `go vet ./...` limpios.
- `.env.example` versionado; `.env` ignorado por git.
- Commits atomicos en `feature/config`, integrados a `develop` por pull request.

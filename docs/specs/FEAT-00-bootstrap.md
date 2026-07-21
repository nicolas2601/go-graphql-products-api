# FEAT-00 — Bootstrap del proyecto

## Objetivo

Dejar el andamiaje del proyecto compilando y con la generacion de codigo de gqlgen
funcionando a partir del esquema GraphQL existente. Esta feature no implementa logica de
negocio; establece la estructura, las dependencias y las herramientas de desarrollo.

## Alcance

Incluye:
- Inicializacion del modulo Go.
- Configuracion de gqlgen y generacion del codigo del servidor a partir del esquema.
- Estructura de carpetas segun Clean Architecture.
- Makefile con objetivos de desarrollo.

No incluye:
- Entidades de dominio, casos de uso, repositorios ni resolvers con logica (features siguientes).

## Criterios de aceptacion (Given-When-Then)

- Dado el repositorio recien inicializado, cuando se ejecuta `go build ./...`, entonces
  compila sin errores.
- Dado el esquema en `graph/schema.graphqls`, cuando se ejecuta la generacion de gqlgen,
  entonces se produce el codigo ejecutable y los modelos sin errores.
- Dado el Makefile, cuando se ejecuta `make generate`, entonces regenera el codigo de gqlgen
  de forma reproducible.
- Dado `go vet ./...`, cuando se ejecuta, entonces no reporta problemas.

## Diseno

- Modulo: `github.com/nicolas2601/go-graphql-products-api`.
- gqlgen genera el codigo ejecutable y los modelos de GraphQL en la capa de delivery, de modo
  que el dominio quede independiente de la libreria de GraphQL.
- Los modelos de GraphQL (generados) se mantienen separados de las entidades de dominio; el
  mapeo entre ambos vivira en los resolvers (FEAT-05).

## Tareas

1. Inicializar el modulo Go.
2. Agregar la dependencia de gqlgen y el archivo de herramientas (`tools.go`).
3. Crear `gqlgen.yml` apuntando al esquema y a las rutas de salida.
4. Generar el codigo con gqlgen.
5. Agregar el Makefile con objetivos `run`, `test`, `lint`, `generate` y `docker`.

## Definition of Done

- `go build ./...` compila.
- `go vet ./...` limpio.
- La generacion de gqlgen corre sin errores y el resultado esta versionado.
- Commits atomicos en la rama `feature/bootstrap`, integrados a `develop` por pull request.

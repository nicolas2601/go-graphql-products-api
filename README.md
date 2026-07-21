# API de Productos - Go + GraphQL

API de gestion de productos construida con Go, GraphQL (gqlgen) y Clean Architecture, con dos
implementaciones de repositorio intercambiables por inyeccion de dependencias: en memoria y
PostgreSQL. Es mi solucion a la Prueba 1 de la evaluacion tecnica.

## Stack

- Go 1.26
- GraphQL con [gqlgen](https://gqlgen.com/)
- PostgreSQL (via `pgx`) - repositorio opcional, intercambiable
- Docker + docker-compose
- Tests de integracion con [testcontainers-go](https://testcontainers.com/)
- CI con GitHub Actions

## Prerequisitos

- Go >= 1.26
- Docker (para `docker compose` y para los tests de integracion de PostgreSQL con testcontainers)
- `golangci-lint` v2 (solo para correr `make lint` localmente)

## Arquitectura (Clean Architecture)

Las dependencias apuntan hacia adentro: el dominio no conoce ningun detalle de infraestructura,
y los casos de uso dependen de la interfaz del repositorio (definida en el dominio), no de una
implementacion concreta.

```
cmd/main.go                     composition root: carga config, cablea dependencias, sirve
        |
        v
internal/server                 handler HTTP: GraphQL (/query), health (/healthz), playground
        |
        v
internal/delivery/graphql       resolvers (SIN logica) + mapeo de errores a extensions.code
        |
        v
internal/usecase                casos de uso (Create/Get/List/Update/Delete)
        |
        v
internal/domain                 entidad Product, invariantes, errores tipados, interfaz del repo
        ^
        |
internal/repository/{memory,postgres}   implementaciones intercambiables de la interfaz
```

## Configuracion

La aplicacion lee su configuracion directamente del entorno; **no carga un archivo `.env`**. El
`.env` sirve para la interpolacion de `docker compose` (`API_PORT`, `POSTGRES_PASSWORD`); para
ejecucion local directa, las variables se pasan inline o se exportan en el shell. `.env.example`
documenta los valores disponibles y el `.env` real no se versiona.

| Variable | Default | Uso |
|---|---|---|
| `PORT` | `8080` | Puerto del servidor HTTP |
| `APP_ENV` | `production` | Modo de ejecucion; `development` habilita playground e introspection |
| `GRAPHQL_PLAYGROUND` | `true` | Habilita el playground (solo tiene efecto en `development`) |
| `LOG_LEVEL` | `info` | Nivel de logging estructurado (`debug`/`info`/`warn`/`error`) |
| `REPO_DRIVER` | `memory` | Repositorio a usar: `memory` o `postgres` |
| `DATABASE_URL` | vacio | Cadena de conexion, requerida solo si `REPO_DRIVER=postgres` |

El default de `APP_ENV` es `production` a proposito (fail-safe): el playground y la introspection
solo se exponen si se setea explicitamente `APP_ENV=development`.

## Como ejecutarlo

Hay dos modos, ambos soportados.

### Modo en memoria (sin base de datos)

```bash
REPO_DRIVER=memory APP_ENV=development go run ./cmd
# make run tambien arranca en memoria, pero en modo production (sin playground);
# para habilitar el playground: APP_ENV=development make run
```

### Modo PostgreSQL (con docker-compose)

Levanta la API conectada a un PostgreSQL, todo en contenedores:

```bash
docker compose up --build
```

`docker compose` espera a que PostgreSQL este healthy antes de arrancar la API, que corre en
`REPO_DRIVER=postgres` y aplica su migracion al conectarse. El puerto de la API es configurable
con `API_PORT` (default 8080). El compose usa `APP_ENV=development` para dejar el playground
disponible en la demo; para un despliegue tipo produccion se cambia a `APP_ENV=production`.

### Endpoints

Una vez arriba (cualquiera de los dos modos):

- **GraphQL**: `POST http://localhost:8080/query`
- **Playground** (solo en `development`): `http://localhost:8080/`
- **Health check**: `GET http://localhost:8080/healthz`

## API GraphQL

**Queries**
- `products`: lista todos los productos.
- `product(id: ID!)`: obtiene un producto por id.

**Mutations**
- `createProduct(input: CreateProductInput!)`: crea un producto.
- `updateProduct(id: ID!, input: UpdateProductInput!)`: actualiza nombre o precio.
- `deleteProduct(id: ID!)`: elimina un producto.

Ejemplo de mutation (crear):

```graphql
mutation {
  createProduct(input: { name: "Teclado", price: 49.90, stock: 10 }) {
    id
    name
    price
    createdAt
  }
}
```

Ejemplo de query (listar):

```graphql
query {
  products {
    id
    name
    price
  }
}
```

Los errores de dominio se mapean a errores de GraphQL con un codigo en `extensions.code`
(por ejemplo `INVALID_PRICE`, `PRODUCT_NOT_FOUND`); los errores internos no filtran detalles.

## Tests

```bash
make test    # go test ./... -race -covermode=atomic -coverprofile=coverage.out
```

La cobertura de las capas de logica (dominio, casos de uso, repositorio en memoria) es 100%. Los
tests de integracion del repositorio PostgreSQL levantan una base efimera con testcontainers, por
lo que requieren Docker corriendo; si Docker no esta disponible, esos tests se saltan solos.

## Desarrollo

El servidor GraphQL usa codigo generado por gqlgen. Tras modificar `graph/schema.graphqls` hay que
regenerarlo:

```bash
make generate   # regenera el codigo de gqlgen desde el esquema
make lint       # golangci-lint
make build      # compila el binario en bin/
```

## Integracion continua

El pipeline de GitHub Actions (`.github/workflows/ci.yml`) corre en cada push y pull request a
`main` y `develop`: lint (`golangci-lint`), `go build`, verificacion de `go mod tidy`, tests con
race detector (incluye los de PostgreSQL via testcontainers), `govulncheck` (bloqueante) y build de
la imagen Docker.

## Estructura del proyecto

```
cmd/main.go                     composition root
internal/
  config/                       configuracion desde el entorno
  domain/                       entidad, invariantes, errores, interfaz del repo
  usecase/                      casos de uso
  repository/memory/            repositorio en memoria (RWMutex)
  repository/postgres/          repositorio PostgreSQL (pgx) + migracion embebida
  delivery/graphql/             resolvers + mapeo de errores/modelos
  server/                       handler HTTP
graph/schema.graphqls           esquema GraphQL
openspec/specs/                 especificaciones por capacidad (framework OpenSpec, formato BDD)
Dockerfile, docker-compose.yml  empaquetado y orquestacion
.github/workflows/ci.yml        pipeline de CI
```

## Decisiones de diseno y deuda conocida

- **`price` como `float64` / `DOUBLE PRECISION`**: heredado del tipo `Float` del esquema GraphQL.
  Para un sistema de dinero real lo correcto seria `NUMERIC`/decimal en centavos, para evitar
  errores de redondeo; se documenta como deuda consciente.
- **Sin paginacion en `products`**: fuera del alcance del enunciado; a escala se resolveria con
  cursor-based pagination.
- **Migracion embebida** (`CREATE TABLE IF NOT EXISTS`): suficiente para una tabla, pero no es
  100% segura con varias replicas de la API arrancando en paralelo. Evolucion: advisory lock o un
  migrator versionado (goose/migrate).
- **`/healthz` es un liveness check** estatico; no verifica la conectividad a la base. Un
  `/readyz` con `pool.Ping` seria el siguiente paso.
- **Aislamiento del dominio**: las validaciones viven en el dominio y se fuerzan tambien en el
  caso de uso antes de persistir; los resolvers no contienen logica de negocio.
- **Paridad de contrato entre implementaciones**: ambos repositorios (memoria y PostgreSQL)
  respetan el mismo contrato de forma verificada por tests: `CreatedAt` es inmutable en `Update`
  y un id con formato invalido (no-UUID contra la columna `UUID` de Postgres) se traduce a
  `ErrProductNotFound`, no a un error interno, igual que en memoria. Ademas del linter estandar
  se corre `gosec` para analisis estatico de seguridad.

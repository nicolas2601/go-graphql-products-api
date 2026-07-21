# FEAT-02 — Capa de dominio

## Objetivo

Modelar la entidad `Product` con sus invariantes, definir la interfaz del repositorio en la
capa de dominio (no en infraestructura) y declarar los errores de dominio tipados. Esta capa
no depende de gqlgen, de la base de datos ni de ningun detalle externo.

## Alcance

Incluye:
- Entidad `Product` y constructor `NewProduct` que valida sus invariantes.
- Interfaz `ProductRepository` en el paquete de dominio.
- Errores de dominio tipados.

No incluye:
- Casos de uso (FEAT-03), implementacion del repositorio (FEAT-04) ni resolvers (FEAT-05).

## Invariantes de negocio

- El nombre no puede estar vacio ni contener solo espacios.
- El precio debe ser mayor a cero.

## Criterios de aceptacion (Given-When-Then)

- Dados datos validos, cuando se construye un producto, entonces devuelve el producto sin error.
- Dado un nombre vacio o de solo espacios, cuando se construye un producto, entonces devuelve
  `ErrInvalidName`.
- Dado un precio menor o igual a cero, cuando se construye un producto, entonces devuelve
  `ErrInvalidPrice`.

## Diseno

- `Product` usa `time.Time` para `CreatedAt`; el mapeo a `String` de GraphQL vive en la capa de
  delivery (FEAT-05), no en el dominio.
- El `id` y el `createdAt` los provee el caso de uso (generador de id y reloj inyectados); el
  constructor solo valida los invariantes de negocio y ensambla una entidad valida.
- `ProductRepository` usa `context.Context` para soportar cancelacion y la futura implementacion
  en PostgreSQL. `GetByID`, `Update` y `Delete` devuelven `ErrProductNotFound` cuando el producto
  no existe.

## Definition of Done

- Tests del constructor y sus invariantes en verde (valido, nombre invalido, precio invalido).
- `go build ./...` y `go vet ./...` limpios.
- Revision del diff con el subagente go-reviewer antes del merge.
- Commits atomicos en `feature/domain-product`, integrados a `develop` por pull request.

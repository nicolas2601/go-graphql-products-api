## Context

El dominio, los casos de uso y el repositorio en memoria ya existen y estan probados al 100%.
Falta conectar el esquema GraphQL (resolvers generados por gqlgen, hoy stubs) con esa logica.

## Decisions

### Inyeccion del caso de uso en el Resolver
El `Resolver` recibe `*usecase.ProductUseCase` por un constructor `NewResolver`. Los resolvers
solo traducen entre GraphQL y el caso de uso; no contienen logica de negocio (requisito del
enunciado).

### Mapeo de errores de dominio a GraphQL
Un helper `toGraphQLError` traduce cada error de dominio tipado a un `*gqlerror.Error` con
`extensions.code` y `message`. Los errores desconocidos se mapean a `INTERNAL_ERROR`.
- Alternativa considerada: devolver `null` en `product(id)` cuando no existe (idiomatico en
  GraphQL). Se descarto para exponer un codigo de error consistente en todas las operaciones,
  que es lo que la prueba evalua ("errores tipados mapeados a GraphQL con code y message").

### Mapeo de modelos
`domain.Product` (con `time.Time`) se mapea al modelo generado, con `createdAt` como `string`
en formato RFC 3339. El dominio no conoce el formato de presentacion.

## Testing

Test de integracion que ejecuta las operaciones GraphQL reales contra el esquema ejecutable,
cableado a un repositorio en memoria, usando el cliente de pruebas de gqlgen. Esto valida el
funcionamiento correcto de las operaciones sin testear los resolvers por dentro.

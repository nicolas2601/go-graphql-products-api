# API de Productos · Go + GraphQL

Esta es mi solución a la **Prueba 1** de la evaluación técnica: una API de gestión de
productos para una tienda. El objetivo no es que tenga muchas features, sino que el código
quede bien organizado y las capas queden claramente separadas.

## Stack

- **Go**
- **GraphQL** con [gqlgen](https://gqlgen.com/)
- Almacenamiento **en memoria** (con un segundo repositorio en PostgreSQL intercambiable
  por inyección de dependencias)

## Arquitectura

Sigo **Clean Architecture** con las capas separadas, y las interfaces del repositorio viven
en la capa de dominio (no en infraestructura). Los resolvers no contienen lógica de negocio.

```
cmd/                        arranque de la aplicación
internal/
  domain/                   entidad Product + interfaces del repositorio
  usecase/                  lógica de negocio (casos de uso)
  repository/               implementación en memoria del repositorio
  delivery/graphql/         resolvers de GraphQL
graph/schema.graphqls       esquema GraphQL
```

## Cómo correrlo

> 🚧 Las instrucciones de arranque se completan al terminar la implementación.

## Estado

En construcción — el historial de commits refleja el progreso incremental.

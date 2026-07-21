# FEAT-04 — Repositorio en memoria

## Objetivo (explore)

Implementar `domain.ProductRepository` con almacenamiento en memoria, seguro para uso
concurrente. Es la implementacion por defecto del sistema (`REPO_DRIVER=memory`) y el doble
real que valida el contrato del dominio de punta a punta.

## Alcance

Incluye:
- Tipo `Repository` en `internal/repository/memory` con las cinco operaciones.
- Seguridad ante concurrencia.
- Tests, incluyendo un test de concurrencia bajo el race detector.

No incluye:
- Repositorio PostgreSQL (feature posterior) ni resolvers ni wiring.

## Decisiones de diseno (propose)

1. **Concurrencia con `sync.RWMutex`.**
   - Alternativa A (elegida): `sync.RWMutex`. Lecturas concurrentes sin bloqueo mutuo,
     escrituras exclusivas. Simple y suficiente para un mapa.
   - Alternativa B: `sync.Map`. Pensada para claves dispares con muchas escrituras; pierde el
     tipado estatico y complica el listado ordenado. Rechazada por sobredimensionada.
   - Alternativa C: sin sincronizacion. Rechazada: data race bajo acceso concurrente.

2. **`Create` detecta duplicados en vez de sobrescribir.** Create significa insertar un
   producto nuevo; sobrescribir un id existente seria un upsert silencioso. Ante un id ya
   presente devuelve `ErrProductAlreadyExists`.

3. **`List` devuelve un orden determinista (por id).** La iteracion de un mapa en Go no es
   determinista; ordenar por id da una salida estable, mas facil de testear y de consumir.

4. **Nombre del paquete `memory`, tipo `Repository`.** Evita el "stutter"
   `memory.MemoryRepository`; el constructor es `memory.New()`.

## Criterios de aceptacion (Given-When-Then)

- Dado un producto nuevo, cuando se crea, entonces queda persistido y recuperable por id.
- Dado un id ya existente, cuando se crea, entonces devuelve `ErrProductAlreadyExists`.
- Dado un id inexistente, cuando se consulta, actualiza o elimina, entonces devuelve
  `ErrProductNotFound`.
- Dados varios productos, cuando se listan, entonces se devuelven todos en orden estable por id.
- Dado acceso concurrente de lectura y escritura, cuando se opera, entonces no hay data races.

## Diseno

- `Repository` guarda `map[string]domain.Product` protegido por `sync.RWMutex`.
- `domain.Product` es un tipo por valor sin punteros internos; guardar y devolver copias evita
  aliasing sin necesidad de clonar manualmente.
- Verificacion en tiempo de compilacion de que `*Repository` satisface `domain.ProductRepository`.

## Tareas

1. Definir `Repository` y `New`.
2. Implementar Create (con deteccion de duplicados) y GetByID.
3. Implementar List (orden estable).
4. Implementar Update y Delete.
5. Agregar test de concurrencia bajo `-race`.

## Verificacion

- `go test ./... -race -cover` en verde, cubriendo caminos felices y de error.
- `go build ./...` y `go vet ./...` limpios.
- Revision del diff con el subagente go-reviewer antes del merge.

## Definition of Done

- Contrato del dominio satisfecho (assertion de interfaz en compilacion).
- Verificacion en verde y revision aplicada.
- Commits atomicos en `feature/repo-memory`, integrados a `develop` por pull request
  (la rama se conserva).
- Registro de cierre (archive) al mergear.

# FEAT-03 — Casos de uso

## Objetivo (explore)

Implementar la logica de aplicacion de productos (`ProductUseCase`) sobre la interfaz de
repositorio del dominio, con las operaciones Create, Get, List, Update y Delete. Aqui viven
los tests unitarios exigidos por el enunciado (sobre casos de uso, no sobre resolvers).

## Alcance

Incluye:
- `ProductUseCase` con sus cinco operaciones.
- Validacion de invariantes antes de persistir (Create via `NewProduct`, Update via `Validate`).
- Tests unitarios con un repositorio falso (test double).

No incluye:
- Implementacion real del repositorio (FEAT-04) ni resolvers (FEAT-05).

## Decisiones de diseno (propose)

1. **Generacion de id y timestamp: inyectados como funciones.** El caso de uso recibe
   `newID func() string` y `now func() time.Time` por constructor.
   - Alternativa A (elegida): funciones inyectadas. Simple, idiomatica y determinista en tests.
   - Alternativa B: interfaces `IDGenerator` y `Clock`. Mas formal pero mas ceremonia para dos
     dependencias triviales.
   - Alternativa C: generar id/timestamp dentro del dominio. Rechazada: acopla el dominio a
     `uuid` y al reloj del sistema, y rompe la testeabilidad determinista.

2. **Update con campos opcionales via punteros.** `Update(ctx, id, name *string, price *float64)`.
   El enunciado dice "actualiza nombre o precio"; los punteros distinguen "no enviado" de
   "enviado vacio", igual que los inputs opcionales de GraphQL. Se valida con `Validate()`
   despues de aplicar los cambios, para no confiar en el llamador (el struct es mutable).

3. **Los errores de no-existencia se propagan del repositorio.** El caso de uso no inventa
   `ErrProductNotFound`; lo devuelve tal cual lo entrega la interfaz del repositorio.

## Criterios de aceptacion (Given-When-Then)

- Dado un input valido, cuando se crea un producto, entonces se persiste con id y createdAt
  provistos por las funciones inyectadas y se devuelve el producto.
- Dado un input invalido (nombre vacio, precio <= 0, stock < 0), cuando se crea un producto,
  entonces devuelve el error de dominio correspondiente y no persiste nada.
- Dado un id inexistente, cuando se consulta, actualiza o elimina, entonces devuelve
  `ErrProductNotFound`.
- Dado un producto existente y un cambio de nombre o precio, cuando se actualiza, entonces se
  aplican solo los campos enviados y se valida el resultado antes de persistir.
- Dado un cambio que viola un invariante, cuando se actualiza, entonces devuelve el error de
  dominio y no persiste.

## Diseno

- `ProductUseCase` depende de `domain.ProductRepository` (interfaz), no de una implementacion.
- Constructor: `NewProductUseCase(repo, newID, now)`.
- Los tests usan un repositorio falso en memoria que implementa la interfaz del dominio.

## Tareas

1. Definir `ProductUseCase` y su constructor con dependencias inyectadas.
2. Implementar Create (valida via `NewProduct`, persiste).
3. Implementar Get y List (delegan al repositorio).
4. Implementar Update (aplica campos opcionales, valida, persiste).
5. Implementar Delete (delega al repositorio).

## Verificacion

- `go test ./... -race -cover` en verde, cubriendo los cinco casos de uso y sus caminos de error.
- `go build ./...` y `go vet ./...` limpios.
- Revision del diff con el subagente go-reviewer antes del merge.

## Definition of Done

- Al menos los tres tests unitarios exigidos sobre casos de uso, superados con holgura.
- Verificacion en verde y revision aplicada.
- Commits atomicos en `feature/usecase-product`, integrados a `develop` por pull request
  (la rama se conserva).
- Registro de cierre (archive) en el reporte del merge.

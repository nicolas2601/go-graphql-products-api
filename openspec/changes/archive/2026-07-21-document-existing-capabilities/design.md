## Context

El sistema ya implementa las capas de Clean Architecture (dominio, casos de uso, repositorio,
delivery) con cobertura completa. Este change solo documenta esas capacidades en OpenSpec.

## Decisions

### Documentacion retroactiva como ADDED
Las capacidades preexistentes se documentan como `ADDED Requirements` porque en OpenSpec eso
representa la creacion de la especificacion (la primera vez que se registra la capacidad), aunque
el codigo ya exista. La correspondencia codigo-especificacion se respalda con la suite de tests.

### Separacion de capacidades
Se separan cuatro capacidades por responsabilidad: dominio (invariantes), aplicacion
(orquestacion), almacenamiento (persistencia) y configuracion (entorno). La capacidad de API
GraphQL ya fue documentada y archivada por el change `add-graphql-resolvers`.

## Testing

Sin cambios de codigo: la validacion es que la suite existente (`go test ./... -race`) siga en
verde y que cada requirement documentado corresponda a comportamiento ya probado.

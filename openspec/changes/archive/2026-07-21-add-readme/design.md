## Context

El proyecto esta completo (dominio, casos de uso, dos repositorios, servidor, Docker, CI) pero el
README no refleja ese estado. Es lo primero que lee un evaluador.

## Decisions

### Documentar ambos modos de ejecucion
Se documentan explicitamente los dos modos: en memoria (`REPO_DRIVER=memory`, sin base) y con
PostgreSQL via `docker compose up`. Es un requisito practico para que la API sea usable sin
depender de una base.

### Diagrama de arquitectura en ASCII
Se incluye un diagrama simple de las capas y la direccion de las dependencias (hacia el dominio),
para comunicar la aplicacion de Clean Architecture de un vistazo.

### Decisiones y deuda visibles
Se documentan las decisiones con tradeoffs (float64 para precio, sin paginacion, migracion
embebida, liveness vs readiness). En una prueba tecnica, nombrar la deuda consciente suma mas que
esconderla.

## Testing

Verificacion: los comandos del README funcionan (modo memoria y modo compose ya verificados
end-to-end en features anteriores); la estructura descrita coincide con el arbol real.

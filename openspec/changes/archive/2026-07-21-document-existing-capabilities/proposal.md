## Why

Las capacidades internas del sistema (dominio, aplicacion, almacenamiento y configuracion) se
construyeron antes de adoptar OpenSpec y quedaron documentadas en `docs/specs/`. Este change las
formaliza como especificaciones de OpenSpec para que todo el proyecto quede bajo el mismo
framework y OpenSpec sea la fuente de verdad.

## What Changes

- Documentar como especificaciones OpenSpec las capacidades ya implementadas y probadas: dominio
  de producto, aplicacion (casos de uso), almacenamiento (repositorio) y configuracion.
- No cambia comportamiento ni codigo: es documentacion retroactiva de lo ya construido.

## Capabilities

### New Capabilities
- `product-domain`: entidad Producto, invariantes de negocio y errores tipados.
- `product-application`: casos de uso de producto (crear, consultar, listar, actualizar, eliminar).
- `product-storage`: contrato del repositorio de productos y su implementacion en memoria.
- `configuration`: configuracion de la aplicacion desde el entorno con valores por defecto.

### Modified Capabilities

## Impact

- Documentacion en `openspec/specs/`. No toca codigo de aplicacion.

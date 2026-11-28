# Modos de falla

## Sesiones simuladas

| Falla                                           | Prevención                                                | Verificación                                       |
| ----------------------------------------------- | --------------------------------------------------------- | -------------------------------------------------- |
| Inicio con credenciales incorrectas             | Se responde `401`.                                        | El E2E envía una clave incorrecta.                 |
| Consulta sin sesión                             | Se responde `401` si falta la cookie o no coincide.       | El E2E consulta antes de iniciar sesión.           |
| La cookie sigue activa al cerrar                | La respuesta expira la cookie con el mismo nombre y ruta. | El E2E consulta después del cierre y espera `401`. |
| El inicio simulado queda expuesto en producción | Las rutas se registran solo en `dev` y `test`.            | La prueba de arranque en `prod` espera `404`.      |

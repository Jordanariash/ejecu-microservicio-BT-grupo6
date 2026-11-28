# E2E: salud y sesiones del servicio

- GET /healthz con la base de datos arriba: 200 application/json {"estado":"ok"}
- GET /readyz con la base de datos arriba: 200 application/json {"estado":"ok"}
- Consulta sin sesión: 401
- Inicio con clave incorrecta: 401
- Inicio de sesión: 200
- Consulta con sesión activa: 200
- Cierre de sesión: 204
- Consulta después del cierre: 401
- GET /readyz con la base de datos detenida: 503 application/problem+json {"title":"Service Unavailable","status":503,"detail":"La base de datos no responde"}
- GET /healthz con la base de datos detenida: 200 application/json {"estado":"ok"}
- Apagado con SIGTERM: código de salida 0, registra la señal de término

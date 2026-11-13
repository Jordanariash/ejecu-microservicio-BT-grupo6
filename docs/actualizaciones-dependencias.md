# Actualizaciones de dependencias

Versiones que cambiamos respecto de las que fija el DAS. Ayudantes indicaron en la reunión del 13 de noviembre
que el DAS quedó desactualizado y que hay que subir las versiones.

## Go: 1.26.5 → 1.26.11

- Dónde: `toolchain` en `apps/api/go.mod` y la imagen `golang` de `apps/api/Dockerfile`.
- Motivo: govulncheck encontró 12 vulnerabilidades de la biblioteca estándar que usa el servicio
  (net/http, crypto/tls, crypto/x509, net/textproto y encoding/asn1), corregidas hasta la 1.26.11.
- Imagen: `golang:1.26.11-alpine3.23@sha256:db99415e2625ce0ee03eed09edf7f32ae8c58dc399f28cc602a87fe7238119e4`

## golang.org/x/text: v0.29.0 → v0.41.0

- Dependencia indirecta: la usan pgx y Huma.
- Motivo: corrige GO-2026-6629 y GO-2026-5970.

## golang.org/x/sync: v0.17.0 → v0.22.0

- Dependencia indirecta: la exige `golang.org/x/text` v0.41.0.

#!/bin/sh
# Compila la API con cobertura, corre el e2e contra PostgreSQL real y deja coverage.out para SonarQube.
set -e
cd "$(dirname "$0")/.."

rm -rf .cobertura
mkdir -p .cobertura/datos

go build -cover -o .cobertura/servicio ./cmd/api
SERVICIO_BIN="$PWD/.cobertura/servicio" GOCOVERDIR="$PWD/.cobertura/datos" go test ./e2e/ -count=1 -v "$@"

go tool covdata textfmt -i .cobertura/datos -o ../../coverage.out

#!/bin/sh
set -eu

temporal=$(mktemp -d)
trap 'rm -rf "$temporal"' EXIT

go -C apps/api run ./cmd/api openapi > "$temporal/openapi.yaml"
openapi-typescript "$temporal/openapi.yaml" -o "$temporal/tipos-api.d.ts"

mv "$temporal/openapi.yaml" contracts/openapi.yaml
mv "$temporal/tipos-api.d.ts" apps/web/src/tipos-api.d.ts

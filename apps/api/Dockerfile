# syntax=docker/dockerfile:1.7
FROM --platform=$BUILDPLATFORM golang:1.26.11-alpine3.23@sha256:db99415e2625ce0ee03eed09edf7f32ae8c58dc399f28cc602a87fe7238119e4 AS build
WORKDIR /src
COPY apps/api/go.mod apps/api/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY apps/api/ ./
ARG TARGETOS TARGETARCH VERSION=dev
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -buildvcs=false -trimpath \
    -ldflags="-s -w -X main.version=$VERSION" \
    -o /out/base-tecnica-grupo6 ./cmd/api

FROM gcr.io/distroless/static-debian13:nonroot@sha256:f7f8f729987ad0fdf6b05eeeae94b26e6a0f613bdf46feea7fc40f7bd72953e6
WORKDIR /app
COPY --from=build --chown=nonroot:nonroot /out/base-tecnica-grupo6 /app/base-tecnica-grupo6
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/base-tecnica-grupo6"]
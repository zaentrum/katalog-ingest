# Base-image registry prefix. Empty default = public Docker Hub; a private
# deploy mirror passes --build-arg BASE=registry.example/library/ .
ARG BASE=
# syntax=docker/dockerfile:1.7
FROM ${BASE}golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/katalog-ingest ./cmd/server

FROM ${BASE}distroless/static-debian12:nonroot
WORKDIR /
COPY --from=build /out/katalog-ingest /katalog-ingest
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/katalog-ingest"]

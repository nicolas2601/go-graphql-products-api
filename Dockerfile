# syntax=docker/dockerfile:1

# Etapa de build: compila un binario estatico.
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/products-api ./cmd

# Etapa de runtime: imagen minima, usuario no-root, con healthcheck sobre /healthz.
FROM alpine:3.20 AS production
RUN apk add --no-cache ca-certificates wget \
    && adduser -D -u 1001 appuser
COPY --from=build /out/products-api /usr/local/bin/products-api
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
    CMD wget -qO- http://localhost:8080/healthz || exit 1
ENTRYPOINT ["products-api"]

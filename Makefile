.PHONY: help run build test lint generate tidy docker

APP_NAME := products-api
PKG := ./...

help:
	@echo "Targets disponibles:"
	@echo "  run       Ejecuta el servidor"
	@echo "  build     Compila el binario en bin/"
	@echo "  test      Corre los tests con race detector y cobertura"
	@echo "  lint      Ejecuta golangci-lint"
	@echo "  generate  Regenera el codigo de gqlgen desde el esquema"
	@echo "  tidy      Ordena las dependencias del modulo"
	@echo "  docker    Construye la imagen de la aplicacion"

run:
	go run ./cmd

build:
	go build -o bin/$(APP_NAME) ./cmd

test:
	go test $(PKG) -race -covermode=atomic -coverprofile=coverage.out

lint:
	golangci-lint run

generate:
	go run github.com/99designs/gqlgen generate

tidy:
	go mod tidy

docker:
	docker build -t $(APP_NAME):latest .

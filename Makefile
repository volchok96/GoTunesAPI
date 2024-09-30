.PHONY: run swag-generate all

all: swag-generate run

# Запуск приложения

run:
	go run cmd/main.go

# Доработка Swagger докуметации

swag-generate:
	cd cmd && swag init -g ../cmd/main.go -d ../config,../models,../controllers,../database,../repository -o ../docs


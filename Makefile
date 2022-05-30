NAME=split-specs-v2
ROOT=github.com/Shelex/${NAME}
GO111MODULE=on
SHELL=/bin/bash

.PHONY: start
start:
	go run main.go

.PHONY: build
build:
	rm -r bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/split-specs
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" scripts/migrate/migrate.go
	mv migrate bin

.PHONY: open
open-web:
	open http://localhost:8080

.PHONY: prof
prof:
	go tool pprof -web http://localhost:6060/debug/pprof/heap

.PHONY: doc
doc:
	swag init

.PHONY: migration
migration:
	go run $(shell pwd)/scripts/migrate ${version} ${direction}

.PHONY: lint
lint: 
	golangci-lint run

.PHONY: deps
deps:
	go mod tidy
	go mod download

.PHONY: web-dev
web-dev: 
	cd web && npm start

.PHONY: web-build
web-build: 
	cd web && npm run build

.PHONY: clear
clear: 
	rm -r web/build && rm -r split-specs-v2 && rm -r web.tar.gz

.PHONY: pack
pack:
	cd web && npm install && npm run build:tailwind && REACT_APP_API_PROTOCOL=http REACT_APP_API_HOST=localhost:3000 npm run build && \
	cd ../ && tar -czf web.tar.gz ./web/build && \
	CGO_ENABLED=0 GOOS=linux go build \
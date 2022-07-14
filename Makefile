.PHONY: build
build:
	go build -v  -o "./artifacts/bin" ./cmd/urlapi

.PHONY: test_api
test_api:
	go test -race -v -timeout 40s ./internal/app/server

.PHONY: test_store
	go test -race -v -timeout 40s ./internal/app/store

.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -v  -o "./artifacts/bin" ./cmd/urlapi

.PHONY: test
test:
	go test -race -v -cover -timeout 40s ./internal/app/server
	go test -race -v -cover -timeout 40s ./internal/app/store
	go test -race -v -cover -timeout 40s ./internal/app/algorithm

.DEFAULT_GOAL := build

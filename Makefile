.PHONY: build
build:
	go build -v  -o "./artifacts/bin" ./cmd/urlapi

.PHONY: test
test:
	go test -race -v -timeout 40s ./internal/app/urlapi

.DEFAULT_GOAL := build
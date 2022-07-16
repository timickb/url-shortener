.PHONY: build
build:
	go build -v  -o "./artifacts/bin" ./cmd/urlapi

.PHONY: test
test:
	go test -race -v -cover -coverprofile=coverage_server.out -timeout 40s ./internal/app/server
	go test -race -v -cover -coverprofile=coverage_store.out -timeout 40s ./internal/app/store
	go test -race -v -cover -coverprofile=coverage_algorithm.out -timeout 40s ./internal/app/algorithm

.PHONY: test_report
test_report:
	make test
	go tool cover -html=coverage_server.out
	go tool cover -html=coverage_store.out
	go tool cover -html=coverage_algorithm.out

.DEFAULT_GOAL := build

.PHONY: all tidy lint

all: tidy lint

tidy:
	gofmt -w .
	go mod tidy -v

lint:
	@go vet ./...
	@golangci-lint run
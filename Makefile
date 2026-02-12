.PHONY: build test lint clean

build:
	go build -o bin/bookstore ./cmd/bookstore

test:
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

clean:
	rm -rf bin/ coverage.out

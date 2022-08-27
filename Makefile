.PHONY: build format lint vet test

build:
	go build -o build/systembolaget cmd/*.go

format:
	go fmt ./...

lint:
	staticcheck ./...

vet:
	go vet ./...

test:
	go test -v ./...

.PHONY: build format lint vet test package

build:
	go build -o build/systembolaget cmd/*.go

format:
	go fmt ./...

lint:
	go run honnef.co/go/tools/cmd/staticcheck@2023.1.7 ./...

vet:
	go vet ./...

test:
	go test -v ./...

package:
	GOOS=linux GOARCH=arm64 go build -o build/linux_arm64 cmd/*.go
	tar -czf build/linux_arm64.tgz build/linux_arm64

	GOOS=linux GOARCH=amd64 go build -o build/linux_amd64 cmd/*.go
	tar -czf build/linux_amd64.tgz build/linux_amd64

	GOOS=darwin GOARCH=arm64 go build -o build/darwin_arm64 cmd/*.go
	zip build/darwin_arm64.zip build/darwin_arm64

	GOOS=darwin GOARCH=amd64 go build -o build/darwin_amd64 cmd/*.go
	zip build/darwin_amd64.zip build/darwin_amd64

	GOOS=windows GOARCH=amd64 go build -o build/windows_amd64 cmd/*.go
	zip build/windows_amd64.zip build/windows_amd64

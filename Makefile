.PHONY: all
all: format test build

.PHONY: format
format:
	go fmt ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	go build -o mysql-process-exporter.exe .

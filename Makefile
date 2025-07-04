# Go project automation

.PHONY: lint format test

lint:
	golangci-lint run

format:
	gofmt -w .

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.txt ./... && go tool cover -func=coverage.txt

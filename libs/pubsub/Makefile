PWD = $(shell pwd)
test:
	go test -race ./...

lint-docker:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.47-alpine golangci-lint run

lint-docker-verbose:
	docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.47-alpine golangci-lint run -v
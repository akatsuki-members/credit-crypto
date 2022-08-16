test:
	go test -race ./...

lint-docker:
	docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.47-alpine golangci-lint run -v --enable-all
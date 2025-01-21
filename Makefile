.PHONY: run test build

run:
	@echo "Running Curlify application..."
	go run 

build:
	@echo "Building Curlify application..."
	go build -o ./bin/curlify

test:
	@echo "Running tests..."
	go test ./... -v


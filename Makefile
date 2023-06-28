build:
	@echo "Building..."
	@go build -o "tmp/main" "./src"
	@echo "Build complete"

# Path: Makefile
run: build
	@echo "Running..."
	@./tmp/main

# Path: Makefile
air: build
	@echo "Running in debug..."
	@air ./cmd/main.go -b 0.0.0.0

# Path: Makefile
test:
	@go test -v ./...
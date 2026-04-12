.PHONY: build run clean test help docs

# Default target
help:
	@echo "AWDoc - Automatic Web Documentation Generator"
	@echo ""
	@echo "Available targets:"
	@echo "  make build       - Build the project"
	@echo "  make run         - Run the analyzer on examples"
	@echo "  make run-pkg     - Run the analyzer on pkg directory"
	@echo "  make test        - Run tests"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make docs        - Generate documentation"

# Build the project
build:
	@echo "🔨 Building AWDoc..."
	go build -v -o awdoc main.go
	@echo "✓ Build complete"

# Run analyzer on examples
run: build
	@echo "🚀 Running analyzer on examples..."
	./awdoc -source ./examples -lang go -output-dir output
	@echo "✓ Documentation generated: output/analysis.md"

# Run analyzer on pkg directory
run-pkg: build
	@echo "🚀 Running analyzer on pkg directory..."
	./awdoc -source ./pkg -lang go -output-dir output
	@echo "✓ Documentation generated: docs/pkg-docs.md"

# Generate documentation for the entire project
docs: build
	@echo "📚 Generating full project documentation..."
	./awdoc -source . -lang go -output-dir output
	@echo "✓ Full documentation generated: output/analysis.md"

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f awdoc
	rm -rf docs/
	@echo "✓ Clean complete"

# Run linter
lint:
	@echo "🔍 Linting..."
	golangci-lint run ./...

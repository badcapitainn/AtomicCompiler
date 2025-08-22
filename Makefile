.PHONY: build test run clean examples

# Build the compiler
build:
	go build -o simplelang cmd/compiler/main.go

# Run tests
test:
	go test ./tests/...

# Run a specific example
run: build
	@echo "Available examples:"
	@echo "  make run-hello     - Hello World program"
	@echo "  make run-vars      - Variables and operations"
	@echo "  make run-loops     - Loops demonstration"
	@echo "  make run-if        - Conditional statements"
	@echo "  make run-funcs     - Functions demonstration"

# Run hello world example
run-hello: build
	./simplelang examples/hello.sl

# Run variables example
run-vars: build
	./simplelang examples/variables.sl

# Run loops example
run-loops: build
	./simplelang examples/loops.sl

# Run if-else example
run-if: build
	./simplelang examples/if_else.sl

# Run functions example
run-funcs: build
	./simplelang examples/functions.sl

# Clean build artifacts
clean:
	rm -f simplelang
	go clean

# Install dependencies
deps:
	go mod tidy

# Show help
help:
	@echo "SimpleLang Compiler - Available targets:"
	@echo "  build      - Build the compiler"
	@echo "  test       - Run tests"
	@echo "  run        - Show available examples"
	@echo "  run-hello  - Run Hello World example"
	@echo "  run-vars   - Run variables example"
	@echo "  run-loops  - Run loops example"
	@echo "  run-if     - Run if-else example"
	@echo "  run-funcs  - Run functions example"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Install dependencies"
	@echo "  help       - Show this help message"

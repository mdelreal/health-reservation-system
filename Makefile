.PHONY: all setup build run test clean protoc-check

# Default target
all: setup build run

# Install dependencies
setup:
	@echo "Setting up the project..."
	go mod tidy

# Build the server
build:
	@echo "Building the server..."
	go build -o health-reservation-server ./cmd/server

# Run the server
run:
	@echo "Running the server..."
	./health-reservation-server

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Generate code from .proto files
generate:
	@echo "Generating Twirp and Go code from .proto files..."
	protoc --twirp_out=. --go_out=. api/reservation.proto

# Check for protoc installation
protoc-check:
	@which protoc > /dev/null || (echo "protoc not installed. Please install it." && exit 1)

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -f health-reservation-server

.PHONY: proto proto-clean test build

# Proto generation
proto:
	@echo "Generating protobuf code..."
	@mkdir -p internal/mnemosyne/pb
	protoc \
		--go_out=internal/mnemosyne/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=internal/mnemosyne/pb \
		--go-grpc_opt=paths=source_relative \
		--proto_path=proto \
		proto/mnemosyne/v1/*.proto
	@echo "Protobuf code generated in internal/mnemosyne/pb/"

proto-clean:
	@echo "Cleaning generated protobuf code..."
	@rm -rf internal/mnemosyne/pb
	@echo "Cleaned!"

# Testing
test:
	go test ./... -v

test-short:
	go test ./... -short

# Build
build:
	go build -o pedantic_raven .

run:
	go run main.go

# Install proto tools (run once)
install-proto-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Installed protoc-gen-go and protoc-gen-go-grpc"
	@echo "Make sure protoc is installed on your system:"
	@echo "  macOS: brew install protobuf"
	@echo "  Linux: apt-get install protobuf-compiler"

help:
	@echo "Available targets:"
	@echo "  proto              - Generate Go code from protobuf files"
	@echo "  proto-clean        - Remove generated protobuf code"
	@echo "  test               - Run all tests"
	@echo "  test-short         - Run short tests only"
	@echo "  build              - Build the binary"
	@echo "  run                - Run the application"
	@echo "  install-proto-tools - Install protoc Go plugins"

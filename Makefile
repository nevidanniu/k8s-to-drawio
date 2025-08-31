.PHONY: build clean test run install

# Build the application
build:
	go build -o k8s-to-drawio main.go

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o k8s-to-drawio-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -o k8s-to-drawio-windows-amd64.exe main.go
	GOOS=darwin GOARCH=amd64 go build -o k8s-to-drawio-darwin-amd64 main.go

# Clean build artifacts
clean:
	rm -f k8s-to-drawio*
	go clean

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install the application
install:
	go install

# Run with example
run-example:
	go run main.go convert -i ./examples/simple-app -o ./output/example.drawio

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Generate documentation
docs:
	go doc -all > docs/API.md
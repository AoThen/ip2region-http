.PHONY: build run test clean docker docker-build docker-run help

# Build binary
build:
	go build -o ip2region-http .

# Run server locally
run:
	./ip2region-http

# Test API
test:
	@echo "Testing API..."
	@curl -s http://localhost:8080/health | jq .
	@echo ""
	@curl -s "http://localhost:8080/search?ip=1.2.3.4" | jq .

# Clean build artifacts
clean:
	rm -f ip2region-http
	rm -f xdb/*.xdb
	rm -rf data/*.xdb

# Build Docker image
docker-build:
	docker build -t ip2region-http .

# Run Docker container
docker-run:
	docker run -p 8080:8080 ip2region-http

# Stop Docker container
docker-stop:
	docker stop ip2region-http || true
	docker rm ip2region-http || true

# Docker Compose
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Show help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  run           - Run the server"
	@echo "  test          - Test the API"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  docker-up     - Start with Docker Compose"
	@echo "  docker-down   - Stop Docker Compose"
	@echo "  help          - Show this help message"

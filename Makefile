.PHONY: build run test clean swagger

BINARY_NAME=robot-scheduler
CONFIG_PATH=config/config.yaml

build:
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME) cmd/main.go

run: build
	@echo "Running..."
	@./bin/$(BINARY_NAME) $(CONFIG_PATH)

dev:
	@echo "Running in development mode..."
	@go run cmd/main.go config/config_dev.yaml

test:
	@echo "Testing..."
	@go test ./... -v

clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf logs/
	@rm -rf data/
	@rm -rf docs/swagger.*

swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/main.go -o docs --parseDependency --parseInternal 
	@echo "Swagger docs generated in docs/"

swagger-serve: swagger
	@echo "Starting Swagger UI server..."
	@docker run -p 8081:8080 -e SWAGGER_JSON=/docs/swagger.json -v $(PWD)/docs:/docs swaggerapi/swagger-ui

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(BINARY_NAME):latest

migrate:
	@echo "Running migrations..."
	@go run scripts/migrate.go

help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run in development mode"
	@echo "  test          - Run tests"
	@echo "  clean         - Clean build artifacts"
	@echo "  swagger       - Generate Swagger documentation"
	@echo "  swagger-serve - Serve Swagger UI in Docker"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  migrate       - Run database migrations"
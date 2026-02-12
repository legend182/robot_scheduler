.PHONY: build run test clean swagger mocks test-unit test-integration test-coverage

BINARY_NAME=robot_scheduler
CONFIG_PATH=configs/config.yaml

build:
	@echo "Building..."
	@go build -o bin/$(BINARY_NAME) cmd/main.go

run: build
	@echo "Running..."
	@./bin/$(BINARY_NAME) $(CONFIG_PATH)

dev:
	@echo "Running in development mode..."
	@go run cmd/main.go config/config_dev.yaml

# Generate mocks for all DAO interfaces
mocks:
	@echo "Generating mocks..."
	@go install go.uber.org/mock/mockgen@latest
	@mkdir -p internal/testutil/mocks
	@mockgen -source=internal/dao/interfaces/user_dao.go -destination=internal/testutil/mocks/mock_user_dao.go -package=mocks
	@mockgen -source=internal/dao/interfaces/device.go -destination=internal/testutil/mocks/mock_device_dao.go -package=mocks
	@mockgen -source=internal/dao/interfaces/task.go -destination=internal/testutil/mocks/mock_task_dao.go -package=mocks
	@mockgen -source=internal/dao/interfaces/semantic.go -destination=internal/testutil/mocks/mock_semantic_dao.go -package=mocks
	@mockgen -source=internal/dao/interfaces/user_operation.go -destination=internal/testutil/mocks/mock_user_operation_dao.go -package=mocks
	@mockgen -source=internal/dao/interfaces/pcd_dao.go -destination=internal/testutil/mocks/mock_pcd_dao.go -package=mocks
	@echo "Mocks generated successfully"

# Run all tests
test:
	@echo "Running all tests..."
	@go test ./... -v

# Run unit tests (fast)
test-unit:
	@echo "Running unit tests..."
	@go test ./internal/service/... ./internal/utils/... -v -short

# Run integration tests (slower)
test-integration:
	@echo "Running integration tests..."
	@go test ./internal/dao/... -v

# Generate test coverage report
test-coverage:
	@echo "Generating test coverage report..."
	@go test ./... -coverprofile=coverage.out -covermode=atomic
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out | grep total
	@echo "Coverage report generated: coverage.html"

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
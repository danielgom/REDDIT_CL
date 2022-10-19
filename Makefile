.PHONY: genM test lint image

# Generates mocks Repositories and Services mocks for testing
genM:
	@mockgen -source=pkg/services/services.go -destination=pkg/routes/mock_services/mocks.go -package mock_services
	@mockgen -source=pkg/services/services.go -destination=pkg/services/mock_services/mocks.go -package mock_services
	@mockgen -source=pkg/db/repositories.go -destination=pkg/services/mock_repositories/mocks.go -package mock_repositories
	@mockgen -source=pkg/config/datasource.go -destination=pkg/db/mock_db/mocks.go -package mock_db

# Test with coverage (dev/local environment)
testL: genM
	@go test ./... --cover -v

# Test with coverage (CI)
testCI:
	@go test ./... --cover -v

# Checks code with golangci-lint linters
lint:
	@golangci-lint run
	@hadolint Dockerfile

# Run the api
run:
	@go run ./cmd/api

# Create Docker image
image:
	@docker build -t reddit-clone .
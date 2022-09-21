.PHONY: genM test lint

# Generates mocks Repositories and Services mocks for testing
genM:
	@mockgen -source=pkg/services/services.go -destination=pkg/routes/mock_services/mocks.go -package mock_services
	@mockgen -source=pkg/services/services.go -destination=pkg/services/mock_services/mocks.go -package mock_services
	@mockgen -source=pkg/db/repositories.go -destination=pkg/services/mock_repositories/mocks.go -package mock_repositories

# Test with coverage (dev/local environment)
testL: genM
	@go test ./... --cover

# Test with coverage (CI)
testCI:
	@go test ./... --cover

# Checks code with golangci-lint linters
lint:
	@golangci-lint run
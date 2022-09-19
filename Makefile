.PHONY: genM test

# Generates mocks Repositories and Services mocks for testing
genM:
	@mockgen -source=pkg/services/services.go -destination=pkg/routes/mock_services/mocks.go -package mock_services
	@mockgen -source=pkg/db/repositories.go -destination=pkg/services/mock_repositories/mocks.go -package mock_repositories

# Test with coverage
test: genM
	@go test ./... --cover
	
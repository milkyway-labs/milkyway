# Milkyway Development Guidelines

## Build Commands
- `make build` - Build milkywayd binary
- `make install` - Install milkywayd binary
- `make test` - Run all unit tests
- `make test-unit` - Run unit tests
- `make test-cover` - Run tests with coverage
- `make test-race` - Run tests with race detection
- `go test ./x/module/... -v -run TestSpecificTest` - Run single test

## Lint & Format
- `make lint` - Run linter checks
- `make lint-fix` - Run linter and fix issues
- `make format` - Format code with gofmt, goimports

## Code Style Guidelines
- Follow standard Go conventions
- Imports ordered: stdlib, 3rd party, local (github.com/milkyway-labs/milkyway)
- Use proper error wrapping: `fmt.Errorf("failed to do X: %w", err)`
- Module-specific code in `x/` directory
- Tests should be thorough with descriptive names
- Prefer table-driven tests with meaningful test cases
- Use Cosmos SDK conventions for keeper implementations
- Follow protobuf naming conventions for RPC services

## Git Workflow
- Descriptive commit messages with module prefix (e.g., "x/restaking: fix unbonding")
- PRs should include tests for new functionality
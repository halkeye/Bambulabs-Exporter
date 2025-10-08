# Testing Guide

This document describes the testing setup and coverage for the BambuLabs Exporter project.

## Test Structure

The project uses Go's built-in testing framework with the following structure:

```
internal/exporter/
├── exporter.go          # Main exporter implementation
└── exporter_test.go     # Test suite for the exporter
```

## Test Coverage

Current test coverage: **72.5%** of statements

### Test Categories

1. **Configuration Tests**
   - Environment variable loading
   - Configuration validation
   - Default values handling

2. **HTTP Endpoint Tests**
   - Home endpoint (`/`)
   - Health check endpoint (`/healthz`)
   - Metrics endpoint (`/metrics`)

3. **MQTT Message Handling Tests**
   - Valid JSON message processing
   - Invalid JSON error handling
   - Wrong command filtering
   - Data structure validation

4. **Prometheus Metrics Tests**
   - Metric registration
   - Metric value setting
   - Label handling for multi-dimensional metrics

## Running Tests

### Basic Test Execution
```bash
go test -v ./...
```

### Test with Coverage
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Test with Race Detection
```bash
go test -v -race ./...
```

### Using Makefile
```bash
make test              # Run all tests
make test-coverage     # Run tests with coverage report
make test-race         # Run tests with race detection
make test-verbose      # Run tests with verbose output
```

## Test Data

Sample MQTT messages are stored in the `testdata/` directory:
- `sample_mqtt_message.json` - Complete sample message with AMS data

## Mock Objects

The test suite includes mock implementations for:
- MQTT Client
- MQTT Message
- MQTT Token

These mocks allow testing without requiring actual MQTT connections.

## Test Environment

Tests use isolated Prometheus registries to avoid metric registration conflicts between test runs.

## CI/CD Integration

The project includes GitHub Actions workflow (`.github/workflows/test.yml`) that runs:
- Unit tests
- Race detection tests
- Linting
- Build verification

## Test Best Practices

1. **Isolation**: Each test uses its own Prometheus registry
2. **Cleanup**: Environment variables are cleaned up after each test
3. **Mocking**: External dependencies are mocked for reliable testing
4. **Coverage**: Tests cover both happy path and error scenarios
5. **Documentation**: Test names clearly describe what is being tested

## Adding New Tests

When adding new functionality:

1. Add corresponding test cases in `exporter_test.go`
2. Ensure tests cover both success and failure scenarios
3. Use descriptive test names
4. Clean up any resources created during tests
5. Update this documentation if adding new test categories

## Test Dependencies

The test suite requires:
- Go 1.25+
- Prometheus client libraries
- MQTT client libraries
- Standard Go testing packages

All dependencies are managed through `go.mod`.
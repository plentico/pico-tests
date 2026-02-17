# Pico Tests

End-to-end tests for the [Pico](https://github.com/plentico/pico) template engine using go-playwright.

## Structure

```
pico-tests/
├── site/           # Test site fixtures
│   ├── views/      # Template files
│   ├── static/     # Static assets
│   └── props.json  # Test props
├── e2e/            # End-to-end tests
│   └── pico_test.go
├── public/         # Generated output (gitignored)
├── go.mod
└── README.md
```

## Prerequisites

- Go 1.21+
- Playwright browsers (auto-installed on first run)

## Running Tests

```bash
# Install dependencies
go mod tidy

# Run all tests
go test ./e2e/... -v

# Run a specific test
go test ./e2e/... -v -run TestPageLoads
```

## What's Tested

- **TestPageLoads** - Verifies the page loads and has correct title
- **TestReactiveCounter** - Tests Pattr reactivity (clicking buttons updates state)
- **TestPRootData** - Verifies p-root-data script contains JSON props

## Adding New Tests

1. Add new test functions to `e2e/pico_test.go`
2. Use `newPage(t)` to create a fresh browser page
3. Tests automatically build the site and start a server

## Local Development

The `go.mod` uses a `replace` directive to reference the local pico repo:

```go
replace github.com/plentico/pico => ../pico
```

For CI/production, update to use a specific version.

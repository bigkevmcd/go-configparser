# go-configparser Development Instructions

**ALWAYS follow these instructions first and only fallback to additional search and context gathering if the information here is incomplete or found to be in error.**

go-configparser is a Go library that implements a Python ConfigParser-compatible configuration file parser with support for interpolation, custom options, and all major ConfigParser features.

## Quick Start & Validation

Bootstrap and validate the repository with these commands:
- `go version` -- ensure Go 1.21+ is installed
- `go build -v ./...` -- builds in ~7 seconds, NEVER CANCEL
- `go test -v ./...` -- runs in ~5 seconds with 93.1% test coverage, NEVER CANCEL
- `go test -cover ./...` -- shows test coverage statistics

## Essential Development Commands

### Building
- `go build -v ./...` -- build all packages (7 seconds)
- `go build` -- build current package only

### Testing  
- `go test -v ./...` -- run all tests with verbose output (5 seconds)
- `go test -cover ./...` -- run tests with coverage report
- `go test ./chainmap` -- run tests for chainmap subpackage only
- `go test -check.f TestSpecificFunction` -- run specific test (uses check framework)

### Code Quality
- `go vet ./...` -- static analysis, runs instantly
- `gofmt -l .` -- check code formatting, runs instantly  
- `go mod tidy` -- clean up module dependencies
- `~/go/bin/golangci-lint run` -- comprehensive linting (install first with: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`, takes ~2 minutes to install)

### CRITICAL: Always run before committing
- `go build -v ./...` -- ensure code builds
- `go test -v ./...` -- ensure tests pass
- `go vet ./...` -- static analysis
- `gofmt -l .` -- verify formatting (should return no output)

## Validation Scenarios

ALWAYS test functionality after making changes by running this validation scenario:

1. Create a test program to validate core functionality:
```go
package main

import (
    "fmt"
    "log"
    "github.com/bigkevmcd/go-configparser"
)

func main() {
    // Test parsing the example configuration file
    p, err := configparser.NewConfigParserFromFile("testdata/example.cfg")
    if err != nil {
        log.Fatal(err)
    }

    // Test basic functionality
    sections := p.Sections()
    fmt.Printf("Sections: %v\n", sections)

    // Test getting a value
    value, err := p.Get("follower", "FrobTimeout")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("FrobTimeout: %s\n", value)

    // Test interpolation
    interpolated, err := p.GetInterpolated("follower", "builder_command")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Interpolated: %s\n", interpolated)
    
    fmt.Println("✓ Working correctly!")
}
```

2. Save this as `/tmp/test_validation.go` and run: `go run /tmp/test_validation.go`
3. Expected output should show sections, values, and interpolated results
4. ALWAYS verify this works after making changes to core parsing logic

## Codebase Navigation

### Key Files & Their Purpose
- `configparser.go` -- Core ConfigParser struct and parsing logic
- `methods.go` -- Main API methods (Get, Set, Sections, etc.)
- `interpolation.go` -- Value interpolation handling (%(var)s syntax)
- `options.go` -- Configuration options and customization
- `section.go` -- Section management and operations
- `chainmap/chainmap.go` -- ChainMap implementation for interpolation lookups

### Test Files
- `configparser_test.go` -- Core parsing tests
- `methods_test.go` -- API method tests  
- `interpolation_test.go` -- Interpolation tests
- `options_test.go` -- Options and customization tests
- `chainmap/chainmap_test.go` -- ChainMap tests

### Example & Test Data
- `testdata/example.cfg` -- Example configuration file with interpolation
- Uses standard INI format with sections, key=value pairs, and %(var)s interpolation

## Common Development Tasks

### Adding New Features
1. ALWAYS add tests first in the appropriate `*_test.go` file
2. Follow the existing test patterns using `gopkg.in/check.v1`
3. Implement the feature in the corresponding main file
4. Run validation: `go test -v ./... && go build -v ./...`
5. Test manually with validation scenario above

### Debugging Configuration Parsing
- Use `testdata/example.cfg` as a reference for valid syntax
- The library supports: sections `[name]`, key=value pairs, comments with `#` or `;`, interpolation with `%(var)s`
- Check `interpolation.go` for interpolation logic
- Check `options.go` for customization options

### Working with Interpolation
- Interpolation uses `%(variable_name)s` syntax
- Variables resolved from DEFAULT section and current section
- ChainMap in `chainmap/` package handles lookup chain
- Test interpolation changes with example config file

## Project Structure

```
.
├── .github/workflows/go.yml    # CI pipeline (Go 1.22, build + test)
├── README.md                   # Project documentation
├── go.mod                      # Go module definition (requires Go 1.21+)
├── configparser.go             # Main parser implementation
├── methods.go                  # Core API methods
├── interpolation.go            # Interpolation logic
├── options.go                  # Configuration options
├── section.go                  # Section management
├── chainmap/                   # ChainMap package for interpolation
│   ├── chainmap.go
│   └── chainmap_test.go
├── testdata/
│   └── example.cfg             # Example configuration file
└── *_test.go                   # Test files
```

## CI/CD Information

The GitHub Actions workflow (`.github/workflows/go.yml`):
- Runs on: Ubuntu latest
- Go version: 1.22 (but works with 1.21+)
- Steps: Checkout → Setup Go → Build → Test
- Build command: `go build -v ./...`
- Test command: `go test -v ./...`

ALWAYS ensure your changes pass both build and test steps locally before pushing.

## Dependencies

- Main dependency: `gopkg.in/check.v1` for testing
- No external runtime dependencies
- Standard library only for core functionality
- Uses Go modules for dependency management

## Performance Notes

- Build time: ~7 seconds (very fast)
- Test time: ~5 seconds (very fast)  
- golangci-lint install: ~2 minutes (one-time setup)
- golangci-lint run: ~few seconds

## Troubleshooting

### Build Issues
- Ensure Go 1.21+ is installed: `go version`
- Clean module cache: `go clean -modcache`
- Verify dependencies: `go mod download`

### Test Issues  
- Run specific test: `go test -run TestName`
- Check test coverage: `go test -cover ./...`
- Validate with example: `go run /tmp/test_validation.go`

### Linting Issues
- Minor linting issues exist in test files (unchecked errors)
- Main code passes all linters
- Install linter first: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

Remember: This is a library, not an application. Focus on API correctness, test coverage, and compatibility with Python ConfigParser behavior.
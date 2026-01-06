# Changelog

All notable changes to this project will be documented in this file.

## [1.1.0] - 2026-01-06

### 💥 Breaking Changes
- **Generics Refactoring**: `Client` methods now return strongly-typed structs (e.g., `*SerpResponse`, `*UniversalResponse`) instead of `map[string]any`.
- **API Signatures**: Updated `SerpSearch`, `UniversalScrape`, `ListTasks`, etc. to return specific types.

### Added
- **Models**: Added comprehensive struct definitions in `models.go` for all API responses.
- **Resource Management**: Added `CloseIdleConnections` method to `Client` to prevent connection leaks.
- **User-Agent**: Standardized User-Agent format.

### Fixed
- **Linting**: Fixed `bodyclose` and `errcheck` issues; upgraded `golangci-lint` configuration.

## [1.0.1] - 2026-01-05

### Added
- **CloseIdleConnections**: Added `CloseIdleConnections` to `Client` to prevent connection leaks.
- **User-Agent**: Standardized User-Agent format to `thordata-go-sdk/{version} go/{ver} ({os}/{arch})`.

### Fixed
- **Linting**: Fixed `bodyclose` issues in examples and updated `.golangci.yml` configuration.
- **CI**: Upgraded `golangci-lint-action` to fix compatibility issues.
- **Stability**: Improved error handling in API response parsing.

# Changelog

All notable changes to this project will be documented in this file.

## [1.0.1] - 2026-01-05

### Added
- **CloseIdleConnections**: Added `CloseIdleConnections` to `Client` to prevent connection leaks.
- **User-Agent**: Standardized User-Agent format to `thordata-go-sdk/{version} go/{ver} ({os}/{arch})`.

### Fixed
- **Linting**: Fixed `bodyclose` issues in examples and updated `.golangci.yml` configuration.
- **CI**: Upgraded `golangci-lint-action` to fix compatibility issues.
- **Stability**: Improved error handling in API response parsing.

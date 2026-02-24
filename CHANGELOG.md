# Changelog

All notable changes to this project will be documented in this file.

## [1.2.0] - 2026-02-24

### Added
- **Browser Connection URL**: Added `GetBrowserConnectionURL()` method for connecting to Thordata Scraping Browser (use with Playwright Go or CDP clients)
- **Complete Four Scraping Methods**: Now supports all four major scraping methods:
  - ✅ SERP Search (`SerpSearch`)
  - ✅ Web Unlocker (`UniversalScrape`)
  - ✅ Browser Scraper (`GetBrowserConnectionURL`)
  - ✅ Web Scraper (`CreateScraperTask`, `RunTask` - tool-based specialized scraping)
- **Browser Connection Example**: Added `examples/browser_connection` demonstrating browser URL generation

### Changed
- **Documentation**: Streamlined README with clear role positioning (Infrastructure / High-Performance Networking)
- **Examples**: Removed redundant examples (`serp_google_news`, `verify`, `verify_run_task`), kept core examples
- **Configuration**: Updated `.env.example` to use unified upstream proxy port (7898)
- **Integration Tests**: Improved error handling for network interference (EOF handling in non-strict mode)
- **Error Handling**: Improved error handling in `execute` function to properly check `io.ReadAll` errors

### Added
- **Capability Matrix**: Added P0/P1/P2 capability matrix in README aligned with multi-language SDK strategy
- **Examples README**: Enhanced examples documentation with clear requirements and usage
- **Upstream Proxy Support**: Added support for `THORDATA_UPSTREAM_PROXY` (e.g., Clash on 127.0.0.1:7898) for proxy chaining
- **Environment Validation**: Added `ValidateEnv()` function to check required environment variables with helpful error messages
- **Comprehensive Example**: Added `examples/comprehensive` demonstrating all major SDK features in one place

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

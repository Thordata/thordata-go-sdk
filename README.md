# thordata-go-sdk

Official Go SDK for Thordata APIs:
- Proxy Network (via standard HTTP proxy configuration)
- SERP API
- Universal / Web Unlocker API
- Web Scraper API (task-based)
- Location API

## Installation

```bash
go get github.com/Thordata/thordata-go-sdk
```

## Environment Variables

```env
THORDATA_SCRAPER_TOKEN=...
THORDATA_PUBLIC_TOKEN=...
THORDATA_PUBLIC_KEY=...

THORDATA_SCRAPERAPI_BASE_URL=https://scraperapi.thordata.com
THORDATA_UNIVERSALAPI_BASE_URL=https://universalapi.thordata.com
THORDATA_WEB_SCRAPER_API_BASE_URL=https://api.thordata.com/api/web-scraper-api
THORDATA_LOCATIONS_BASE_URL=https://api.thordata.com/api/locations
```

## Quick Start

```go
package main

import (
  "context"
  "fmt"
  "os"
  "time"

  "github.com/Thordata/thordata-go-sdk/thordata"
)

func main() {
  cfg := thordata.Config{
    ScraperToken: os.Getenv("THORDATA_SCRAPER_TOKEN"),
    PublicToken:  os.Getenv("THORDATA_PUBLIC_TOKEN"),
    PublicKey:    os.Getenv("THORDATA_PUBLIC_KEY"),
    Timeout:      30 * time.Second,
  }

  client, err := thordata.NewClient(cfg)
  if err != nil {
    panic(err)
  }

  ctx := context.Background()
  data, err := client.SerpSearch(ctx, thordata.SerpOptions{
    Query:  "pizza",
    Engine: "google",
    Country: "us",
    SearchType: "news",
    OutputFormat: "json",
  })
  if err != nil {
    panic(err)
  }

  fmt.Printf("%T\n", data)
}
```

## Development

This repository includes a git submodule (`sdk-spec`) for cross-SDK parity checks.

```bash
git submodule update --init --recursive
go test ./...
```

### Windows note

If `go test` fails with `Access is denied` when executing test binaries, your antivirus may be blocking temporary executables.
Try disabling the antivirus for the repository, or set `GOTMPDIR` to a project-local directory.
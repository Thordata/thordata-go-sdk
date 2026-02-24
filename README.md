# Thordata Go SDK

<div align="center">

<img src="https://img.shields.io/badge/Thordata-AI%20Infrastructure-blue?style=for-the-badge" alt="Thordata Logo">

**The Official Go Client for Thordata APIs**

*Infrastructure • High-Performance Networking • Connection Pooling*

[![Go Reference](https://pkg.go.dev/badge/github.com/Thordata/thordata-go-sdk.svg)](https://pkg.go.dev/github.com/Thordata/thordata-go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/Thordata/thordata-go-sdk)](https://goreportcard.com/report/github.com/Thordata/thordata-go-sdk)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

</div>

---

## 📖 Introduction

The **Thordata Go SDK** is designed for **high-concurrency proxy networking and infrastructure scenarios**. It serves as a reliable, performant foundation for building proxy gateways, forwarding services, and high-throughput data collection platforms.

**Complete Four Scraping Methods Support:**
- ✅ **SERP Search** - Google, Bing, Yandex, DuckDuckGo
- ✅ **Web Unlocker** - Universal scrape with JS rendering
- ✅ **Browser Scraper** - Browser connection URL (use with Playwright Go or CDP)
- ✅ **Web Scraper** - Tool-based specialized scraping (100+ pre-built tools)

**Role & Positioning:**
- **Infrastructure / High-Performance Networking SDK** (per [multi-language strategy](https://github.com/Thordata/thordata-sdk-spec))
- Optimized for connection pooling, Keep-Alive, and efficient resource usage
- Simple, reliable API client wrappers for Thordata services

**Key Features:**
- **🚀 Connection Pooling:** Intelligent TCP connection reuse to proxy gateways, preventing port exhaustion
- **🛡️ Type-Safe:** Full struct definitions for SERP, Tasks, Locations, and Public APIs
- **✨ Idiomatic Go:** Context-aware requests, standardized error handling
- **🧩 Environment-Based Config:** Lazy validation, configure via `THORDATA_*` env vars

---

## 📦 Installation

```bash
go get github.com/Thordata/thordata-go-sdk
```

---

## ⚙️ Configuration

Set environment variables (copy `.env.example` to `.env` and fill in values):

```bash
# Required: Scraper token (for SERP, Universal, Tasks)
export THORDATA_SCRAPER_TOKEN="your_scraper_token"

# Optional: Public token/key (for Task status/download, Locations, Public API)
export THORDATA_PUBLIC_TOKEN="your_public_token"
export THORDATA_PUBLIC_KEY="your_public_key"

# Proxy credentials (Residential, Datacenter, Mobile, ISP)
export THORDATA_RESIDENTIAL_USERNAME="your_username"
export THORDATA_RESIDENTIAL_PASSWORD="your_password"
export THORDATA_PROXY_HOST="vpnXXXX.pr.thordata.net"
export THORDATA_PROXY_PORT=9999
export THORDATA_PROXY_PROTOCOL="https"

# Optional: Upstream proxy (e.g., Clash on 127.0.0.1:7898)
# Supports proxy chaining: client -> upstream proxy -> Thordata proxy -> target
export THORDATA_UPSTREAM_PROXY="socks5://127.0.0.1:7898"
```

See [`.env.example`](.env.example) for the complete configuration template.

---

## 🚀 Quick Start

### 1. Initialize Client

```go
package main

import (
    "context"
    "os"
    "github.com/Thordata/thordata-go-sdk/thordata"
)

func main() {
    client, err := thordata.NewClient(thordata.Config{
        ScraperToken: os.Getenv("THORDATA_SCRAPER_TOKEN"),
        PublicToken:  os.Getenv("THORDATA_PUBLIC_TOKEN"),
        PublicKey:    os.Getenv("THORDATA_PUBLIC_KEY"),
    })
    if err != nil {
        panic(err)
    }
    
    // Use client...
}
```

### 2. Proxy Network (Core Feature)

High-performance proxy requests with connection pooling:

```go
proxy := &thordata.ProxyConfig{
    Product:  thordata.ProxyResidential,
    Username: os.Getenv("THORDATA_RESIDENTIAL_USERNAME"),
    Password: os.Getenv("THORDATA_RESIDENTIAL_PASSWORD"),
    Country:  "us",
    SessionID: "sess-01",
    SessionDuration: 10, // minutes
}

resp, err := client.ProxyGet(context.Background(), "https://httpbin.org/ip", proxy)
if err != nil {
    panic(err)
}
defer resp.Body.Close()

// Reuse the same client instance across goroutines for connection pooling
```

### 3. SERP Search

```go
results, err := client.SerpSearch(context.Background(), thordata.SerpOptions{
    Query:   "golang concurrency",
    Engine:  "google",
    Country:  "us",
    Num:      10,
})
if err != nil {
    panic(err)
}

for _, item := range results.Organic {
    fmt.Printf("%s - %s\n", item.Title, item.Link)
}
```

### 4. Universal Scrape (Web Unlocker)

```go
html, err := client.UniversalScrape(context.Background(), thordata.UniversalOptions{
    URL:          "https://example.com",
    JSRender:     true,
    OutputFormat: "html",
    Country:      "us",
    WaitFor:      ".content-loaded",
})
if err != nil {
    panic(err)
}

fmt.Println(html.HTML)
```

### 5. Web Scraper Tasks

```go
// Create task
taskID, err := client.CreateScraperTask(context.Background(), thordata.ScraperTaskOptions{
    FileName:   "task1",
    SpiderID:   "universal",
    SpiderName: "universal",
    Parameters: map[string]any{"url": "https://example.com"},
})
if err != nil {
    panic(err)
}

// Wait for completion
status, err := client.WaitForTask(context.Background(), taskID, 5*time.Second, 10*time.Minute)
if err != nil {
    panic(err)
}

// Get result
if status == "ready" {
    url, _ := client.GetTaskResult(context.Background(), taskID, "json")
    fmt.Println("Download:", url)
}

// Or use RunTask for complete workflow
url, err := client.RunTask(context.Background(), taskOpt, &thordata.RunTaskConfig{
    MaxWait: 10 * time.Minute,
})
```

### 6. Locations API

```go
countries, err := client.ListCountries(context.Background(), 1) // 1=Residential
states, err := client.ListStates(context.Background(), "us", 1)
cities, err := client.ListCities(context.Background(), "us", "ca", 1)
asns, err := client.ListASNs(context.Background(), "us", 1)
```

### 7. Browser Connection URL

Get WebSocket URL for connecting to Thordata Scraping Browser (use with Playwright Go or CDP):

```go
wsURL, err := client.GetBrowserConnectionURL(
    os.Getenv("THORDATA_BROWSER_USERNAME"),
    os.Getenv("THORDATA_BROWSER_PASSWORD"),
)
if err != nil {
    panic(err)
}

// Use with Playwright Go or Chrome DevTools Protocol
fmt.Println("Connect to:", wsURL)
```

### 8. Public API (Management)

```go
// Usage statistics
stats, err := client.GetUsageStatistics(ctx, "2024-01-01", "2024-01-31")

// Proxy users
users, err := client.ListProxyUsers(ctx, 1)

// Whitelist IP
_, err = client.AddWhitelistIP(ctx, "1.2.3.4", 1)

// Proxy servers
servers, err := client.ListProxyServers(ctx, 1)
```

---

## 🎯 Capability Matrix

**P0 (Core):**
- ✅ `proxyNetwork` - High-performance proxy with connection pooling (HTTPS, SOCKS5h)
- ✅ `publicApi` - Usage stats, proxy users, whitelist, proxy servers
- ✅ `locations` - Countries, states, cities, ASNs
- ✅ `unlimited` - Supported via proxy network
- ✅ `extract` - Supported via proxy network

**P1 (Service Clients):**
- ✅ `serp` - Google, Bing, Yandex, DuckDuckGo
- ✅ `webunlocker` - Universal scrape with JS rendering
- ✅ `webScraperTasksLifecycle` - Create, wait, status, download, RunTask

**P2 (Optional):**
- ✅ `browserConnection` - Browser connection URL generation (use with Playwright Go or CDP clients)

For detailed capability mapping, see [thordata-sdk-spec](https://github.com/Thordata/thordata-sdk-spec).

---

## 🧪 Testing

**Unit tests (no network):**
```bash
go test ./...
```

**Integration tests (requires `.env` with real credentials):**
```bash
go test -tags=integration ./thordata
```

**Run examples:**
```bash
# See examples/README.md for details
go run ./examples/proxy_basic
go run ./examples/serp_basic
go run ./examples/universal_basic
go run ./examples/tasks_basic
go run ./examples/locations_basic
go run ./examples/comprehensive  # All features in one demo
```

**Validate environment:**
```go
import "github.com/Thordata/thordata-go-sdk/thordata"

if err := thordata.ValidateEnv(); err != nil {
    if thordata.IsWarning(err) {
        // Non-fatal warning (e.g., missing proxy credentials)
        fmt.Printf("Warning: %v\n", err)
    } else {
        // Fatal error (missing required tokens)
        panic(err)
    }
}
```

---

## 📚 Examples

See [`examples/`](examples/) directory for complete working examples:
- `proxy_basic` - Proxy network usage
- `serp_basic` - SERP search
- `universal_basic` - Universal scrape
- `tasks_basic` - Web Scraper tasks lifecycle
- `locations_basic` - Locations API

---

## 🔗 Related

- **Python SDK** (Flagship / Full Coverage): [thordata-python-sdk](https://github.com/Thordata/thordata-python-sdk)
- **Node.js SDK** (Product Integration): [thordata-js-sdk](https://github.com/Thordata/thordata-js-sdk)
- **Java SDK** (Enterprise Integration): [thordata-java-sdk](https://github.com/Thordata/thordata-java-sdk)
- **SDK Specification**: [thordata-sdk-spec](https://github.com/Thordata/thordata-sdk-spec)

---

## 📄 License

MIT License. See [LICENSE](LICENSE) for details.

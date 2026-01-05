# Thordata Go SDK

<div align="center">

**Official Go Client for Thordata APIs**

*Proxy Network • SERP API • Web Unlocker • Web Scraper API*

[![Go Reference](https://pkg.go.dev/badge/github.com/Thordata/thordata-go-sdk.svg)](https://pkg.go.dev/github.com/Thordata/thordata-go-sdk)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

</div>

---

## 📦 Installation

```bash
go get github.com/Thordata/thordata-go-sdk
```

## 🔐 Configuration

```bash
export THORDATA_SCRAPER_TOKEN="your_token"
export THORDATA_PUBLIC_TOKEN="public_token"
export THORDATA_PUBLIC_KEY="public_key"
```

## 🚀 Quick Start

```go
package main

import (
    "context"
    "fmt"
    "os"
    "github.com/Thordata/thordata-go-sdk/thordata"
)

func main() {
    client, _ := thordata.NewClient(thordata.Config{
        ScraperToken: os.Getenv("THORDATA_SCRAPER_TOKEN"),
    })

    // SERP Search
    res, _ := client.SerpSearch(context.Background(), thordata.SerpOptions{
        Query:  "golang",
        Engine: "google",
    })
    fmt.Println(res)
}
```

## 📚 Core Features

### 🌐 Proxy Network

```go
// Use default proxy from env
resp, err := client.ProxyGet(ctx, "https://httpbin.org/ip", nil)

// Or custom configuration
proxy := &thordata.ProxyConfig{
    Username: "user",
    Password: "pass",
    Product:  thordata.ProxyResidential,
    Country:  "us",
    City:     "new_york",
}
resp, err := client.ProxyGet(ctx, "https://httpbin.org/ip", proxy)
```

### 🔍 SERP API

```go
opts := thordata.SerpOptions{
    Query:        "pizza",
    Engine:       "google_maps",
    Country:      "us",
    OutputFormat: "json",
}
result, err := client.SerpSearch(ctx, opts)
```

### 🔓 Universal Scraping API

```go
html, err := client.UniversalScrape(ctx, thordata.UniversalOptions{
    URL:          "https://example.com",
    JSRender:     true,
    WaitFor:      ".content",
    OutputFormat: "html",
})
```

### 🕷️ Web Scraper API (Tasks)

```go
// 1. Create Task
taskId, _ := client.CreateScraperTask(ctx, thordata.ScraperTaskOptions{
    FileName:   "task1",
    SpiderID:   "universal",
    SpiderName: "universal",
    Parameters: map[string]any{"url": "https://example.com"},
})

// 2. Wait
status, _ := client.WaitForTask(ctx, taskId, 5*time.Second, 10*time.Minute)

// 3. Download
if status == "ready" {
    url, _ := client.GetTaskResult(ctx, taskId, "json")
    fmt.Println(url)
}
```

### 📹 Video/Audio Tasks

```go
taskId, _ := client.CreateVideoTask(ctx, thordata.VideoTaskOptions{
    FileName:   "video_{{VideoID}}",
    SpiderID:   "youtube_video_by-url",
    SpiderName: "youtube.com",
    Parameters: map[string]any{"url": "..."},
    CommonSettings: thordata.CommonSettings{
        Resolution: "1080p",
    },
})
```

### 📊 Account Management

```go
// Usage Stats
stats, _ := client.GetUsageStatistics(ctx, "2024-01-01", "2024-01-31")

// Whitelist IP
client.AddWhitelistIP(ctx, "1.2.3.4", 1) // 1=Residential
```

## ⚙️ Advanced Usage

### Connection Management

For high-throughput applications, ensure you close idle connections if you create many clients.

```go
client.CloseIdleConnections()
```

## 📄 License

MIT License
# Thordata Go SDK

<div align="center">

<img src="https://img.shields.io/badge/Thordata-AI%20Infrastructure-blue?style=for-the-badge" alt="Thordata Logo">

**The Official Go Client for Thordata APIs**

*High Performance • Concurrency Safe • Connection Pooling*

[![Go Reference](https://pkg.go.dev/badge/github.com/Thordata/thordata-go-sdk.svg)](https://pkg.go.dev/github.com/Thordata/thordata-go-sdk)
[![Go Report Card](https://goreportcard.com/badge/github.com/Thordata/thordata-go-sdk)](https://goreportcard.com/report/github.com/Thordata/thordata-go-sdk)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

</div>

---

## 📖 Introduction

A high-performance Go client designed for scale. It features a custom `Transport` implementation to ensure **connection pooling (Keep-Alive)** works correctly with Thordata's dynamic proxy network, preventing port exhaustion in high-concurrency scenarios.

**Key Features:**
*   **🚀 Optimized Transport:** Intelligent reuse of TCP connections to proxy gateways.
*   **🛡️ Type-Safe:** Full struct definitions for SERP, Tasks, and Usage APIs.
*   **✨ Idiomatic:** Context-aware requests, standardized error handling.
*   **🧩 Lazy Validation:** Initialize the client once, configure credentials via env vars as needed.

---

## 📦 Installation

```bash
go get github.com/Thordata/thordata-go-sdk
```

---

## 🚀 Quick Start

### 1. Initialize Client

The client automatically reads `THORDATA_*` environment variables.

```go
package main

import (
    "context"
    "fmt"
    "os"
    "github.com/Thordata/thordata-go-sdk/thordata"
)

func main() {
    // Zero-config init (reads from env)
    client, err := thordata.NewClient(thordata.Config{
        ScraperToken: os.Getenv("THORDATA_SCRAPER_TOKEN"),
    })
    if err != nil {
        panic(err)
    }
    
    // ... use client
}
```

### 2. SERP Search

```go
res, err := client.SerpSearch(context.Background(), thordata.SerpOptions{
    Query:  "golang concurrency",
    Engine: "google",
    Num:    10,
})

if err != nil {
    panic(err)
}

fmt.Printf("Found %d results\n", len(res.Organic))
```

### 3. High-Performance Proxy Request

```go
// Create proxy config (Residential, US, Sticky Session)
proxy := &thordata.ProxyConfig{
    Username:        os.Getenv("THORDATA_RESIDENTIAL_USERNAME"),
    Password:        os.Getenv("THORDATA_RESIDENTIAL_PASSWORD"),
    Product:         thordata.ProxyResidential,
    Country:         "us",
    SessionID:       "sess-go-01",
    SessionDuration: 10,
}

// Perform Request (Uses internal connection pool)
resp, err := client.ProxyGet(context.Background(), "https://httpbin.org/ip", proxy)
if err != nil {
    panic(err)
}
defer resp.Body.Close()

fmt.Println("Status:", resp.StatusCode)
```

---

## ⚙️ Advanced Features

### Connection Pooling
The `Client` struct maintains an internal map of transports. You should share a single `Client` instance across goroutines to maximize performance.

```go
// Use this single instance for all your routines
client, _ := thordata.NewClient(...)

for i := 0; i < 100; i++ {
    go func() {
        // These calls will reuse underlying TCP connections to the proxy
        client.ProxyGet(ctx, url, proxy)
    }()
}
```

### Web Scraper Tasks

```go
// Create Task
taskID, _ := client.CreateScraperTask(ctx, thordata.ScraperTaskOptions{
    FileName:   "task1",
    SpiderID:   "universal",
    SpiderName: "universal",
    Parameters: map[string]any{"url": "https://example.com"},
})

// Wait for completion
status, _ := client.WaitForTask(ctx, taskID, 5*time.Second, 10*time.Minute)

// Get Result
if status == "ready" {
    url, _ := client.GetTaskResult(ctx, taskID, "json")
    fmt.Println("Result URL:", url)
}
```

---

## 📄 License

MIT License.
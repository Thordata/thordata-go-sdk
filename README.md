# thordata-go-sdk

Official Go SDK for Thordata APIs.

## Installation

```bash
go get github.com/Thordata/thordata-go-sdk
```

## Configuration

```bash
export THORDATA_SCRAPER_TOKEN=your_token
export THORDATA_PUBLIC_TOKEN=your_public_token
export THORDATA_PUBLIC_KEY=your_public_key
```

## Quick Start

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
        PublicToken:  os.Getenv("THORDATA_PUBLIC_TOKEN"),
        PublicKey:    os.Getenv("THORDATA_PUBLIC_KEY"),
    })

    // SERP Search
    res, _ := client.SerpSearch(context.Background(), thordata.SerpOptions{
        Query:  "golang",
        Engine: "google",
    })
    fmt.Println(res)
}
```

## Features

### Web Scraper API

```go
// Video Task
taskId, _ := client.CreateVideoTask(ctx, thordata.VideoTaskOptions{
    FileName: "video",
    SpiderID: "youtube_video_by-url",
    SpiderName: "youtube.com",
    Parameters: map[string]any{"url": "..."},
    CommonSettings: thordata.CommonSettings{Resolution: "1080p"},
})

// Wait & Result
status, _ := client.WaitForTask(ctx, taskId, 5*time.Second, 10*time.Minute)
url, _ := client.GetTaskResult(ctx, taskId, "json")
```

### Account Management

```go
// Usage
stats, _ := client.GetUsageStatistics(ctx, "2024-01-01", "2024-01-31")

// Proxy Users
users, _ := client.ListProxyUsers(ctx, 1)

// Whitelist
client.AddWhitelistIP(ctx, "1.2.3.4", 1)
```

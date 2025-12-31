package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Thordata/thordata-go-sdk/examples/internal/env"
	"github.com/Thordata/thordata-go-sdk/thordata"
)

func main() {
	_ = env.LoadDotEnv(".env")

	// ScraperToken is not required for Proxy Network requests, but NewClient requires it.
	// You can pass a dummy token if you're only testing Proxy Network.
	scraper := os.Getenv("THORDATA_SCRAPER_TOKEN")
	if scraper == "" {
		scraper = "dummy"
	}

	client, err := thordata.NewClient(thordata.Config{
		ScraperToken: scraper,
		Timeout:      60 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	// Use default proxy from env (residential/datacenter/mobile/whitelist)
	resp, err := client.ProxyGet(context.Background(), "https://httpbin.org/ip", nil)
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	b, _ := io.ReadAll(resp.Body)
	fmt.Println("status:", resp.StatusCode)
	fmt.Println(string(b))
}

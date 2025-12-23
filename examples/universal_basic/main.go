package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Thordata/thordata-go-sdk/examples/internal/env"
	"github.com/Thordata/thordata-go-sdk/thordata"
)

func main() {
	_ = env.LoadDotEnv(".env")

	token := os.Getenv("THORDATA_SCRAPER_TOKEN")
	if token == "" {
		fmt.Println("Missing THORDATA_SCRAPER_TOKEN. Set env var and re-run.")
		os.Exit(1)
	}

	client, err := thordata.NewClient(thordata.Config{
		ScraperToken: token,
		Timeout:      60 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	out, err := client.UniversalScrape(context.Background(), thordata.UniversalOptions{
		URL:          "https://httpbin.org/html",
		JSRender:     false,
		OutputFormat: "html",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Output type:", fmt.Sprintf("%T", out))
	fmt.Println("Preview:", truncate(fmt.Sprintf("%v", out), 300))
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

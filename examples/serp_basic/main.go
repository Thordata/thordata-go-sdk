package main

import (
	"context"
	"encoding/json"
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

	out, err := client.SerpSearch(context.Background(), thordata.SerpOptions{
		Query:        "pizza",
		Engine:       "google",
		Country:      "us",
		OutputFormat: "json",
	})
	if err != nil {
		panic(err)
	}

	b, _ := json.MarshalIndent(out, "", "  ")
	fmt.Println(string(b))
}

package main

import (
	"context"
	"time"

	"github.com/Thordata/thordata-go-sdk/examples/internal/example"
	"github.com/Thordata/thordata-go-sdk/thordata"
)

func main() {
	example.LoadEnv()
	if example.SkipIfMissing("THORDATA_SCRAPER_TOKEN") {
		return
	}

	client, err := example.NewClient(60 * time.Second)
	if err != nil {
		panic(err)
	}

	out, err := client.SerpSearch(context.Background(), thordata.SerpOptions{
		Query:        "AI regulation",
		Engine:       "google_news",
		Country:      "us",
		OutputFormat: "json",
		Extra: map[string]string{
			"so": "1", // 0=relevance, 1=date (if supported)
		},
	})
	if err != nil {
		panic(err)
	}

	example.PrintJSON(out)
}

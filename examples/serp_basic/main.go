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
		Query:        "pizza",
		Engine:       "google",
		Country:      "us",
		OutputFormat: "json",
	})
	if err != nil {
		panic(err)
	}

	example.PrintJSON(out)
}

package main

import (
	"context"
	"time"

	"github.com/Thordata/thordata-go-sdk/examples/internal/example"
)

func main() {
	example.LoadEnv()
	if example.SkipIfMissing("THORDATA_SCRAPER_TOKEN", "THORDATA_PUBLIC_TOKEN", "THORDATA_PUBLIC_KEY") {
		return
	}

	client, err := example.NewClient(60 * time.Second)
	if err != nil {
		panic(err)
	}

	items, err := client.ListCountries(context.Background(), 1)
	if err != nil {
		panic(err)
	}

	example.PrintJSON(items)
}

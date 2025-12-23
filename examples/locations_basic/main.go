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

	scraper := os.Getenv("THORDATA_SCRAPER_TOKEN")
	pub := os.Getenv("THORDATA_PUBLIC_TOKEN")
	key := os.Getenv("THORDATA_PUBLIC_KEY")

	if scraper == "" || pub == "" || key == "" {
		fmt.Println("Missing THORDATA_SCRAPER_TOKEN / THORDATA_PUBLIC_TOKEN / THORDATA_PUBLIC_KEY.")
		os.Exit(1)
	}

	client, err := thordata.NewClient(thordata.Config{
		ScraperToken: scraper,
		PublicToken:  pub,
		PublicKey:    key,
		Timeout:      60 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	items, err := client.ListCountries(context.Background(), 1)
	if err != nil {
		panic(err)
	}

	b, _ := json.MarshalIndent(items, "", "  ")
	fmt.Println(string(b))
}

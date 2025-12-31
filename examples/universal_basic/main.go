package main

import (
	"context"
	"fmt"
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

	out, err := client.UniversalScrape(context.Background(), thordata.UniversalOptions{
		URL:          "https://httpbin.org/html",
		JSRender:     false,
		OutputFormat: "html",
	})
	if err != nil {
		panic(err)
	}

	s := fmt.Sprintf("%v", out)
	fmt.Println("Output type:", fmt.Sprintf("%T", out))
	fmt.Println("Preview:", example.Truncate(s, 300))
}

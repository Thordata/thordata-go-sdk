package main

import (
	"fmt"
	"os"

	"github.com/Thordata/thordata-go-sdk/examples/internal/example"
	"github.com/Thordata/thordata-go-sdk/thordata"
)

func main() {
	example.LoadEnv()

	if example.SkipIfMissing("THORDATA_BROWSER_USERNAME", "THORDATA_BROWSER_PASSWORD") {
		return
	}

	client, err := thordata.NewClient(thordata.Config{
		ScraperToken: "dummy", // Not needed for browser URL
	})
	if err != nil {
		panic(err)
	}

	// Get browser connection URL
	wsURL, err := client.GetBrowserConnectionURL(
		os.Getenv("THORDATA_BROWSER_USERNAME"),
		os.Getenv("THORDATA_BROWSER_PASSWORD"),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Browser WebSocket Connection URL:")
	fmt.Println(wsURL)
	fmt.Println()
	fmt.Println("Use this URL with:")
	fmt.Println("  - Playwright Go: playwright.Connect(wsURL)")
	fmt.Println("  - Chrome DevTools Protocol (CDP) clients")
	fmt.Println("  - Any WebSocket-based browser automation tool")
}

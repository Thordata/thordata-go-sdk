package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Thordata/thordata-go-sdk/examples/internal/example"
	"github.com/Thordata/thordata-go-sdk/thordata"
)

func main() {
	example.LoadEnv()

	// Validate environment
	if err := thordata.ValidateEnv(); err != nil {
		if !thordata.IsWarning(err) {
			fmt.Printf("❌ Configuration error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("⚠️  Warning: %v\n", err)
	}

	client, err := example.NewClient(60 * time.Second)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	fmt.Println("=== Thordata Go SDK Comprehensive Demo ===")
	fmt.Println()

	// 1. SERP Search
	fmt.Println("1. Testing SERP Search...")
	if example.Env("THORDATA_SCRAPER_TOKEN") != "" {
		serpRes, err := client.SerpSearch(ctx, thordata.SerpOptions{
			Query:   "golang programming",
			Engine:  "google",
			Country: "us",
			Num:     5,
		})
		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
		} else {
			fmt.Printf("   ✅ Found %d results\n", len(serpRes.Organic))
			if len(serpRes.Organic) > 0 {
				fmt.Printf("   First result: %s\n", serpRes.Organic[0].Title)
			}
		}
	} else {
		fmt.Println("   ⏭️  Skipped (THORDATA_SCRAPER_TOKEN not set)")
	}

	// 2. Universal Scrape
	fmt.Println("\n2. Testing Universal Scrape...")
	if example.Env("THORDATA_SCRAPER_TOKEN") != "" {
		uniRes, err := client.UniversalScrape(ctx, thordata.UniversalOptions{
			URL:          "https://httpbin.org/html",
			JSRender:     false,
			OutputFormat: "html",
		})
		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
		} else {
			preview := uniRes.HTML
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			fmt.Printf("   ✅ Scraped %d bytes\n", len(uniRes.HTML))
			fmt.Printf("   Preview: %s\n", preview)
		}
	} else {
		fmt.Println("   ⏭️  Skipped (THORDATA_SCRAPER_TOKEN not set)")
	}

	// 3. Proxy Network
	fmt.Println("\n3. Testing Proxy Network...")
	proxy := &thordata.ProxyConfig{
		Product:  thordata.ProxyResidential,
		Username: example.Env("THORDATA_RESIDENTIAL_USERNAME"),
		Password: example.Env("THORDATA_RESIDENTIAL_PASSWORD"),
		Country:  "us",
	}
	if proxy.Username != "" && proxy.Password != "" {
		resp, err := client.ProxyGet(ctx, "https://httpbin.org/ip", proxy)
		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
		} else {
			defer func() { _ = resp.Body.Close() }()
			fmt.Printf("   ✅ Proxy request successful (status: %d)\n", resp.StatusCode)
		}
	} else {
		fmt.Println("   ⏭️  Skipped (proxy credentials not set)")
	}

	// 4. Locations API
	fmt.Println("\n4. Testing Locations API...")
	if example.Env("THORDATA_PUBLIC_TOKEN") != "" && example.Env("THORDATA_PUBLIC_KEY") != "" {
		countries, err := client.ListCountries(ctx, 1)
		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
		} else {
			fmt.Printf("   ✅ Found %d countries\n", len(countries))
			if len(countries) > 0 {
				fmt.Printf("   Example: %s (%s)\n", countries[0].CountryName, countries[0].CountryCode)
			}
		}
	} else {
		fmt.Println("   ⏭️  Skipped (THORDATA_PUBLIC_TOKEN/KEY not set)")
	}

	// 5. Public API - Usage Statistics
	fmt.Println("\n5. Testing Usage Statistics...")
	if example.Env("THORDATA_PUBLIC_TOKEN") != "" && example.Env("THORDATA_PUBLIC_KEY") != "" {
		now := time.Now()
		weekAgo := now.Add(-7 * 24 * time.Hour)
		stats, err := client.GetUsageStatistics(ctx, weekAgo.Format("2006-01-02"), now.Format("2006-01-02"))
		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
		} else {
			balanceGB := stats.TrafficBalance / (1024 * 1024 * 1024)
			fmt.Printf("   ✅ Traffic balance: %.2f GB\n", balanceGB)
			fmt.Printf("   Total usage: %.2f GB\n", stats.TotalUsageTraffic/(1024*1024*1024))
		}
	} else {
		fmt.Println("   ⏭️  Skipped (THORDATA_PUBLIC_TOKEN/KEY not set)")
	}

	// 6. Browser Connection URL
	fmt.Println("\n6. Testing Browser Connection URL...")
	browserUser := example.Env("THORDATA_BROWSER_USERNAME")
	browserPass := example.Env("THORDATA_BROWSER_PASSWORD")
	if browserUser != "" && browserPass != "" {
		wsURL, err := client.GetBrowserConnectionURL(browserUser, browserPass)
		if err != nil {
			fmt.Printf("   ❌ Error: %v\n", err)
		} else {
			fmt.Printf("   ✅ Browser WebSocket URL generated\n")
			fmt.Printf("   URL: %s\n", wsURL[:50]+"...") // Show first 50 chars
		}
	} else {
		fmt.Println("   ⏭️  Skipped (THORDATA_BROWSER_USERNAME/PASSWORD not set)")
	}

	fmt.Println("\n=== Demo Complete ===")
}

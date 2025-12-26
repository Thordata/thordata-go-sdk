package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Thordata/thordata-go-sdk/examples/internal/env"
	"github.com/Thordata/thordata-go-sdk/thordata"
)

func main() {
	_ = env.LoadDotEnv(".env")

	scraper := os.Getenv("THORDATA_SCRAPER_TOKEN")
	pub := os.Getenv("THORDATA_PUBLIC_TOKEN")
	key := os.Getenv("THORDATA_PUBLIC_KEY")
	sign := os.Getenv("THORDATA_SIGN")
	apiKey := os.Getenv("THORDATA_API_KEY")

	if scraper == "" {
		fmt.Println("❌ Error: THORDATA_SCRAPER_TOKEN is required")
		os.Exit(1)
	}

	fmt.Println("========================================")
	fmt.Println("Thordata SDK - New Features Verification")
	fmt.Println("========================================")

	client, err := thordata.NewClient(thordata.Config{
		ScraperToken: scraper,
		PublicToken:  pub,
		PublicKey:    key,
		Sign:         sign,
		ApiKey:       apiKey,
		Timeout:      60 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// 1. Video Task
	fmt.Println("\n--- Testing: Video Task Creation ---")
	taskID, err := client.CreateVideoTask(ctx, thordata.VideoTaskOptions{
		FileName:   "test_{{VideoID}}",
		SpiderID:   "youtube_video_by-url",
		SpiderName: "youtube.com",
		Parameters: map[string]any{
			"url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		},
		CommonSettings: thordata.CommonSettings{
			Resolution:  "720p",
			IsSubtitles: "false",
		},
	})
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ Video task created: %s\n", taskID)
	}

	// 2. Usage Stats
	fmt.Println("\n--- Testing: Usage Statistics ---")
	now := time.Now()
	weekAgo := now.Add(-7 * 24 * time.Hour)
	stats, err := client.GetUsageStatistics(ctx, weekAgo.Format("2006-01-02"), now.Format("2006-01-02"))
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		bal := getFloat(stats, "traffic_balance")
		fmt.Printf("✅ Stats Retrieved:\n")
		fmt.Printf("   Balance: %.2f GB\n", bal/(1024*1024))
	}

	// 3. Proxy Users
	fmt.Println("\n--- Testing: Proxy Users ---")
	users, err := client.ListProxyUsers(ctx, 1) // 1=Residential
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		count := getInt(users, "user_count")
		fmt.Printf("✅ Users Retrieved: %d\n", count)
	}

	// 4. Proxy Servers (ISP)
	fmt.Println("\n--- Testing: Proxy Servers (ISP) ---")
	servers, err := client.ListProxyServers(ctx, 1) // 1=ISP
	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
	} else {
		fmt.Printf("✅ ISP Servers: %d\n", len(servers))
		if len(servers) > 0 {
			// ProxyServer struct is not exported directly from response, it's []any
			// Manual check
			s := servers[0].(map[string]any)
			fmt.Printf("   Server 1: %v:%v\n", s["ip"], s["port"])
		}
	}

	// 5. API NEW - Balance
	fmt.Println("\n--- Testing: API NEW - Balance ---")
	balNew, err := client.GetResidentialBalance(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Sign and ApiKey") {
			fmt.Println("⚠️ Skipped (requires Sign/ApiKey)")
		} else {
			fmt.Printf("❌ Failed: %v\n", err)
		}
	} else {
		b := getFloat(balNew, "balance")
		fmt.Printf("✅ Balance: %.2f GB\n", b/(1024*1024*1024))
	}

	// 6. API NEW - ISP Regions
	fmt.Println("\n--- Testing: API NEW - ISP Regions ---")
	regions, err := client.GetIspRegions(ctx)
	if err != nil {
		if strings.Contains(err.Error(), "Sign and ApiKey") {
			fmt.Println("⚠️ Skipped (requires Sign/ApiKey)")
		} else {
			fmt.Printf("❌ Failed: %v\n", err)
		}
	} else {
		fmt.Printf("✅ Regions: %d\n", len(regions))
	}
}

func getFloat(m map[string]any, key string) float64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}

func getInt(m map[string]any, key string) int {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return int(f)
		}
	}
	return 0
}

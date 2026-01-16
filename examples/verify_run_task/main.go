package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Thordata/thordata-go-sdk/examples/internal/example"
	"github.com/Thordata/thordata-go-sdk/thordata"
)

func main() {
	example.LoadEnv()

	// Check Env
	token := os.Getenv("THORDATA_SCRAPER_TOKEN")
	spiderId := os.Getenv("THORDATA_TASK_SPIDER_ID")
	spiderName := os.Getenv("THORDATA_TASK_SPIDER_NAME")
	paramsJson := os.Getenv("THORDATA_TASK_PARAMETERS_JSON")

	if token == "" || spiderId == "" {
		fmt.Println("❌ Missing env vars. Please check .env")
		os.Exit(1)
	}

	// Parse Params
	var params map[string]any
	// Dashboard sometimes gives array [{...}]
	if strings.HasPrefix(strings.TrimSpace(paramsJson), "[") {
		var arr []map[string]any
		json.Unmarshal([]byte(paramsJson), &arr)
		if len(arr) > 0 {
			params = arr[0]
		}
	} else {
		json.Unmarshal([]byte(paramsJson), &params)
	}

	// Init Client
	client, err := example.NewClient(60 * time.Second)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n--- Testing Go RunTask [%s] ---\n", spiderName)

	// Call RunTask
	url, err := client.RunTask(context.Background(), thordata.ScraperTaskOptions{
		FileName:      "go_test_" + fmt.Sprint(time.Now().Unix()),
		SpiderID:      spiderId,
		SpiderName:    spiderName,
		Parameters:    params,
		IncludeErrors: true,
	}, &thordata.RunTaskConfig{
		InitialPollInterval: 3 * time.Second,
		MaxWait:             10 * time.Minute,
	})

	if err != nil {
		fmt.Printf("❌ Failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Success! Download URL: %s\n", url)
}

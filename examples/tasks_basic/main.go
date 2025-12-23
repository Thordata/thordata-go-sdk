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

	spiderID := os.Getenv("THORDATA_SPIDER_ID")
	spiderName := os.Getenv("THORDATA_SPIDER_NAME")
	fileName := os.Getenv("THORDATA_TASK_FILE_NAME")
	paramsJSON := os.Getenv("THORDATA_TASK_PARAMETERS_JSON")

	if fileName == "" {
		fileName = "demo_task"
	}

	if scraper == "" || pub == "" || key == "" {
		fmt.Println("Missing THORDATA_SCRAPER_TOKEN / THORDATA_PUBLIC_TOKEN / THORDATA_PUBLIC_KEY.")
		os.Exit(1)
	}

	// Task creation requires spider_id/spider_name from dashboard.
	if spiderID == "" || spiderName == "" || paramsJSON == "" {
		fmt.Println("Skipping tasks example.")
		fmt.Println("Set THORDATA_SPIDER_ID, THORDATA_SPIDER_NAME, THORDATA_TASK_PARAMETERS_JSON to run.")
		return
	}

	var params map[string]any
	if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
		panic(err)
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

	taskID, err := client.CreateScraperTask(context.Background(), thordata.ScraperTaskOptions{
		FileName:      fileName,
		SpiderID:      spiderID,
		SpiderName:    spiderName,
		Parameters:    params,
		IncludeErrors: true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Task created:", taskID)

	status, err := client.WaitForTask(context.Background(), taskID, 5*time.Second, 10*time.Minute)
	if err != nil {
		panic(err)
	}
	fmt.Println("Final status:", status)

	if status == "ready" || status == "success" || status == "finished" {
		dl, err := client.GetTaskResult(context.Background(), taskID, "json")
		if err != nil {
			panic(err)
		}
		fmt.Println("Download:", dl)
	}
}

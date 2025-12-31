package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Thordata/thordata-go-sdk/examples/internal/example"
	"github.com/Thordata/thordata-go-sdk/thordata"
)

func main() {
	example.LoadEnv()

	// Basic auth requirements
	if example.SkipIfMissing("THORDATA_SCRAPER_TOKEN", "THORDATA_PUBLIC_TOKEN", "THORDATA_PUBLIC_KEY") {
		return
	}

	// Task creation requires spider_id/spider_name/parameters from dashboard.
	if example.SkipIfMissing("THORDATA_SPIDER_ID", "THORDATA_SPIDER_NAME", "THORDATA_TASK_PARAMETERS_JSON") {
		fmt.Println("This example requires a Web Scraper task template from your dashboard (API Builder).")
		return
	}

	fileName := example.Env("THORDATA_TASK_FILE_NAME")
	if fileName == "" {
		fileName = "demo_task"
	}

	var params map[string]any
	if err := json.Unmarshal([]byte(example.Env("THORDATA_TASK_PARAMETERS_JSON")), &params); err != nil {
		panic(err)
	}

	client, err := example.NewClient(60 * time.Second)
	if err != nil {
		panic(err)
	}

	taskID, err := client.CreateScraperTask(context.Background(), thordata.ScraperTaskOptions{
		FileName:      fileName,
		SpiderID:      example.Env("THORDATA_SPIDER_ID"),
		SpiderName:    example.Env("THORDATA_SPIDER_NAME"),
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

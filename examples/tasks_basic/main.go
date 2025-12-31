// examples/tasks_basic/main.go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Thordata/thordata-go-sdk/examples/internal/example"
	"github.com/Thordata/thordata-go-sdk/thordata"
)

func main() {
	example.LoadEnv()

	// Auth requirements for builder/tasks APIs
	if example.SkipIfMissing("THORDATA_SCRAPER_TOKEN", "THORDATA_PUBLIC_TOKEN", "THORDATA_PUBLIC_KEY") {
		return
	}

	// Prefer THORDATA_TASK_* (matches .env.example & dashboard), keep backward compatibility.
	spiderID := example.Env("THORDATA_TASK_SPIDER_ID")
	if spiderID == "" {
		spiderID = example.Env("THORDATA_SPIDER_ID")
	}

	spiderName := example.Env("THORDATA_TASK_SPIDER_NAME")
	if spiderName == "" {
		spiderName = example.Env("THORDATA_SPIDER_NAME")
	}

	fileName := example.Env("THORDATA_TASK_FILE_NAME")
	if fileName == "" {
		fileName = "demo_task"
	}

	paramsJSON := example.Env("THORDATA_TASK_PARAMETERS_JSON")

	// Task creation requires spider_id/spider_name/parameters from dashboard.
	if spiderID == "" || spiderName == "" || paramsJSON == "" {
		fmt.Println("Skipping tasks example.")
		fmt.Println("Set THORDATA_TASK_SPIDER_ID, THORDATA_TASK_SPIDER_NAME, THORDATA_TASK_PARAMETERS_JSON to run.")
		return
	}

	// Dashboard curl uses spider_parameters as a JSON array string: [{...}]
	// For convenience, accept either "{...}" or "[{...}]".
	paramsJSON = strings.TrimSpace(paramsJSON)
	if strings.HasPrefix(paramsJSON, "{") {
		paramsJSON = "[" + paramsJSON + "]"
	}

	var raw any
	if err := json.Unmarshal([]byte(paramsJSON), &raw); err != nil {
		panic(err)
	}

	// Our SDK options currently expect a single map (Parameters: map[string]any).
	// If user provided an array, use the first element.
	var params map[string]any
	switch v := raw.(type) {
	case []any:
		if len(v) == 0 {
			panic("THORDATA_TASK_PARAMETERS_JSON is an empty array")
		}
		m, ok := v[0].(map[string]any)
		if !ok {
			panic("THORDATA_TASK_PARAMETERS_JSON[0] must be an object")
		}
		params = m
	case map[string]any:
		params = v
	default:
		panic("THORDATA_TASK_PARAMETERS_JSON must be a JSON object or an array of objects")
	}

	client, err := example.NewClient(60 * time.Second)
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

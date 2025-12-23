package thordata

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

type ScraperTaskOptions struct {
	FileName        string
	SpiderID        string
	SpiderName      string
	Parameters      map[string]any
	UniversalParams map[string]any
	IncludeErrors   bool
}

func (c *Client) CreateScraperTask(ctx context.Context, opt ScraperTaskOptions) (string, error) {
	if opt.FileName == "" || opt.SpiderID == "" || opt.SpiderName == "" {
		return "", errors.New("fileName, spiderId, and spiderName are required")
	}
	if opt.Parameters == nil {
		return "", errors.New("parameters is required")
	}

	paramsArr := []map[string]any{opt.Parameters}
	paramsJSON, _ := json.Marshal(paramsArr)

	payload := map[string]string{
		"file_name":         opt.FileName,
		"spider_id":         opt.SpiderID,
		"spider_name":       opt.SpiderName,
		"spider_parameters": string(paramsJSON),
		"spider_errors":     boolToLower(opt.IncludeErrors),
	}

	if opt.UniversalParams != nil {
		u, _ := json.Marshal(opt.UniversalParams)
		payload["spider_universal"] = string(u)
	}

	body := ToFormBody(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.scraperBuilderURL, bytes.NewBufferString(body))
	if err != nil {
		return "", err
	}
	for k, v := range BuildAuthHeaders(c.cfg.ScraperToken) {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	res, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	parsed, _ := SafeParseJSON(raw)

	obj, ok := parsed.(map[string]any)
	if !ok {
		return "", errors.New("invalid response from Scraper Builder API")
	}

	if cv, ok2 := obj["code"]; ok2 {
		if f, ok3 := cv.(float64); ok3 && int(f) != 200 {
			return "", RaiseForCode("Task creation failed", obj, res.StatusCode)
		}
	}

	dataObj, _ := obj["data"].(map[string]any)
	taskID := ""
	if dataObj != nil {
		taskID = toString(dataObj["task_id"])
	}
	if strings.TrimSpace(taskID) == "" {
		return "", errors.New("task_id missing in response")
	}
	return taskID, nil
}

func (c *Client) GetTaskStatus(ctx context.Context, taskID string) (string, error) {
	if strings.TrimSpace(taskID) == "" {
		return "", errors.New("taskId is required")
	}
	if strings.TrimSpace(c.cfg.PublicToken) == "" || strings.TrimSpace(c.cfg.PublicKey) == "" {
		return "", errors.New("publicToken and publicKey are required for task status")
	}

	payload := map[string]string{"tasks_ids": taskID}
	body := ToFormBody(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.scraperStatusURL, bytes.NewBufferString(body))
	if err != nil {
		return "", err
	}
	for k, v := range BuildPublicHeaders(c.cfg.PublicToken, c.cfg.PublicKey) {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	res, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	parsed, _ := SafeParseJSON(raw)

	obj, ok := parsed.(map[string]any)
	if !ok {
		return "unknown", nil
	}

	if cv, ok2 := obj["code"]; ok2 {
		if f, ok3 := cv.(float64); ok3 && int(f) != 200 {
			return "", RaiseForCode("Task status failed", obj, res.StatusCode)
		}
	}

	items, _ := obj["data"].([]any)
	for _, it := range items {
		m, _ := it.(map[string]any)
		if m == nil {
			continue
		}
		if toString(m["task_id"]) == taskID {
			st := toString(m["status"])
			if st == "" {
				return "unknown", nil
			}
			return st, nil
		}
	}
	return "unknown", nil
}

func (c *Client) GetTaskResult(ctx context.Context, taskID string, fileType string) (string, error) {
	if strings.TrimSpace(taskID) == "" {
		return "", errors.New("taskId is required")
	}
	if strings.TrimSpace(c.cfg.PublicToken) == "" || strings.TrimSpace(c.cfg.PublicKey) == "" {
		return "", errors.New("publicToken and publicKey are required for task download")
	}
	if strings.TrimSpace(fileType) == "" {
		fileType = "json"
	}

	payload := map[string]string{
		"tasks_id": taskID,
		"type":     fileType,
	}
	body := ToFormBody(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.scraperDownloadURL, bytes.NewBufferString(body))
	if err != nil {
		return "", err
	}
	for k, v := range BuildPublicHeaders(c.cfg.PublicToken, c.cfg.PublicKey) {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	res, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	parsed, _ := SafeParseJSON(raw)

	obj, ok := parsed.(map[string]any)
	if !ok {
		return "", errors.New("invalid response from task download API")
	}

	if cv, ok2 := obj["code"]; ok2 {
		if f, ok3 := cv.(float64); ok3 && int(f) == 200 {
			dataObj, _ := obj["data"].(map[string]any)
			if dataObj != nil {
				dl := toString(dataObj["download"])
				if dl != "" {
					return dl, nil
				}
			}
		}
	}

	return "", RaiseForCode("Get task result failed", obj, res.StatusCode)
}

func (c *Client) WaitForTask(ctx context.Context, taskID string, pollInterval time.Duration, maxWait time.Duration) (string, error) {
	if pollInterval <= 0 {
		pollInterval = 5 * time.Second
	}
	if maxWait <= 0 {
		maxWait = 10 * time.Minute
	}

	deadline := time.Now().Add(maxWait)

	for time.Now().Before(deadline) {
		st, err := c.GetTaskStatus(ctx, taskID)
		if err != nil {
			return "", err
		}
		l := strings.ToLower(st)
		if l == "ready" || l == "success" || l == "finished" || l == "failed" || l == "error" || l == "cancelled" {
			return st, nil
		}
		time.Sleep(pollInterval)
	}
	return "", errors.New("task did not complete within maxWait")
}

func boolToLower(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

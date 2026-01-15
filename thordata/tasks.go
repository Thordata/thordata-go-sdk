package thordata

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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

// CreateScraperTask returns taskID string
func (c *Client) CreateScraperTask(ctx context.Context, opt ScraperTaskOptions) (string, error) {
	if c.cfg.ScraperToken == "" {
		return "", errors.New("scraperToken is required for Task Builder")
	}
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
	for k, v := range BuildBuilderHeaders(c.cfg.ScraperToken, c.cfg.PublicToken, c.cfg.PublicKey) {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	resp, err := execute[APIResponse[TaskCreateResponse]](c, req)
	if err != nil {
		return "", err
	}
	return resp.Data.TaskId, nil
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

	// Status API returns list of statuses in "data"
	resp, err := execute[APIResponse[[]TaskStatus]](c, req)
	if err != nil {
		return "unknown", err
	}

	for _, item := range resp.Data {
		if item.TaskId == taskID {
			return item.Status, nil
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

	resp, err := execute[APIResponse[TaskDownload]](c, req)
	if err != nil {
		return "", err
	}
	return resp.Data.Download, nil
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

func (c *Client) CreateVideoTask(ctx context.Context, opt VideoTaskOptions) (string, error) {
	if c.cfg.ScraperToken == "" {
		return "", errors.New("scraperToken is required for Task Builder")
	}
	if opt.FileName == "" || opt.SpiderID == "" || opt.SpiderName == "" {
		return "", errors.New("fileName, spiderId, and spiderName are required")
	}
	if opt.Parameters == nil {
		return "", errors.New("parameters is required")
	}

	paramsArr := []map[string]any{opt.Parameters}
	paramsJSON, _ := json.Marshal(paramsArr)

	settingsJSON, _ := json.Marshal(opt.CommonSettings)

	payload := map[string]string{
		"file_name":         opt.FileName,
		"spider_id":         opt.SpiderID,
		"spider_name":       opt.SpiderName,
		"spider_parameters": string(paramsJSON),
		"spider_errors":     boolToLower(opt.IncludeErrors),
		"common_settings":   string(settingsJSON),
	}

	body := ToFormBody(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.videoBuilderURL, bytes.NewBufferString(body))
	if err != nil {
		return "", err
	}

	for k, v := range BuildBuilderHeaders(c.cfg.ScraperToken, c.cfg.PublicToken, c.cfg.PublicKey) {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	resp, err := execute[APIResponse[TaskCreateResponse]](c, req)
	if err != nil {
		return "", err
	}
	return resp.Data.TaskId, nil
}

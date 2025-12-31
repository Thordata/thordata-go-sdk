// thordata/public.go
package thordata

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// --- Task List ---

func (c *Client) ListTasks(ctx context.Context, page, size int) (map[string]any, error) {
	if c.cfg.PublicToken == "" || c.cfg.PublicKey == "" {
		return nil, errors.New("publicToken and publicKey required")
	}

	payload := map[string]string{
		"page": strconv.Itoa(page),
		"size": strconv.Itoa(size),
	}
	body := ToFormBody(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.taskListURL, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	for k, v := range BuildPublicHeaders(c.cfg.PublicToken, c.cfg.PublicKey) {
		req.Header.Set(k, v)
	}

	return c.executeRequest(req)
}

// --- Helper for execution ---

func (c *Client) executeRequest(req *http.Request) (map[string]any, error) {
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	raw, _ := io.ReadAll(res.Body)
	parsed, _ := SafeParseJSON(raw)

	obj, ok := parsed.(map[string]any)
	if !ok {
		return nil, errors.New("invalid JSON response")
	}

	if cv, ok2 := obj["code"]; ok2 {
		if f, ok3 := cv.(float64); ok3 && int(f) != 200 {
			return nil, RaiseForCode("API error", obj, res.StatusCode)
		}
	}

	// Return data field if exists, otherwise return whole object
	if data, ok := obj["data"].(map[string]any); ok {
		return data, nil
	}
	return obj, nil
}

// --- Usage Statistics ---

func (c *Client) GetUsageStatistics(ctx context.Context, fromDate, toDate string) (map[string]any, error) {
	if c.cfg.PublicToken == "" || c.cfg.PublicKey == "" {
		return nil, errors.New("publicToken and publicKey required")
	}

	u, _ := url.Parse(c.usageStatsURL)
	q := u.Query()
	q.Set("token", c.cfg.PublicToken)
	q.Set("key", c.cfg.PublicKey)
	q.Set("from_date", fromDate)
	q.Set("to_date", toDate)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.executeRequest(req)
}

// --- Proxy Users ---

func (c *Client) ListProxyUsers(ctx context.Context, proxyType int) (map[string]any, error) {
	if c.cfg.PublicToken == "" || c.cfg.PublicKey == "" {
		return nil, errors.New("publicToken and publicKey required")
	}

	u, _ := url.Parse(c.proxyUsersURL + "/user-list")
	q := u.Query()
	q.Set("token", c.cfg.PublicToken)
	q.Set("key", c.cfg.PublicKey)
	q.Set("proxy_type", strconv.Itoa(proxyType))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c.executeRequest(req)
}

// --- Whitelist IP ---

func (c *Client) AddWhitelistIP(ctx context.Context, ip string, proxyType int) (map[string]any, error) {
	if c.cfg.PublicToken == "" || c.cfg.PublicKey == "" {
		return nil, errors.New("publicToken and publicKey required")
	}

	payload := map[string]string{
		"ip":         ip,
		"proxy_type": strconv.Itoa(proxyType),
		"status":     "true",
	}
	body := ToFormBody(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.whitelistURL+"/add-ip", bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	for k, v := range BuildPublicHeaders(c.cfg.PublicToken, c.cfg.PublicKey) {
		req.Header.Set(k, v)
	}
	return c.executeRequest(req)
}

// --- Proxy List ---

func (c *Client) ListProxyServers(ctx context.Context, proxyType int) ([]any, error) {
	if c.cfg.PublicToken == "" || c.cfg.PublicKey == "" {
		return nil, errors.New("publicToken and publicKey required")
	}

	u, _ := url.Parse(c.proxyListURL)
	q := u.Query()
	q.Set("token", c.cfg.PublicToken)
	q.Set("key", c.cfg.PublicKey)
	q.Set("proxy_type", strconv.Itoa(proxyType))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// executeRequest returns map, but list returns array in data or root?
	// Manual implementation for list handling
	req.Header.Set("User-Agent", c.cfg.UserAgent)
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	raw, _ := io.ReadAll(res.Body)
	parsed, _ := SafeParseJSON(raw)

	if obj, ok := parsed.(map[string]any); ok {
		if cv, ok2 := obj["code"]; ok2 {
			if f, ok3 := cv.(float64); ok3 && int(f) != 200 {
				return nil, RaiseForCode("API error", obj, res.StatusCode)
			}
		}
		if list, ok := obj["data"].([]any); ok {
			return list, nil
		}
		if list, ok := obj["list"].([]any); ok {
			return list, nil
		}
	}
	return []any{}, nil
}

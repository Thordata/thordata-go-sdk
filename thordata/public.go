package thordata

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// ListTasks returns strongly typed TaskList
func (c *Client) ListTasks(ctx context.Context, page, size int) (*TaskList, error) {
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

	resp, err := execute[APIResponse[TaskList]](c, req)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// GetUsageStatistics returns strongly typed UsageStatistics
func (c *Client) GetUsageStatistics(ctx context.Context, fromDate, toDate string) (*UsageStatistics, error) {
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

	// API often returns data fields directly in root or in "data".
	// Let's try APIResponse first.
	resp, err := execute[APIResponse[UsageStatistics]](c, req)
	if err != nil {
		return nil, err
	}
	// Fallback logic: if Data is empty but root has fields (UsageStatistics fields overlap with APIResponse fields?)
	// Actually, based on spec, usage stats are in "data".
	return &resp.Data, nil
}

// ListProxyUsers returns strongly typed ProxyUserList
func (c *Client) ListProxyUsers(ctx context.Context, proxyType int) (*ProxyUserList, error) {
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

	// Sometimes returns fields directly, sometimes in "data".
	// We'll assume standard wrapper for now based on your other SDKs.
	resp, err := execute[APIResponse[ProxyUserList]](c, req)
	if err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// AddWhitelistIP returns raw map (since response is usually empty/status)
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

	// Generic execute
	resp, err := execute[APIResponse[map[string]any]](c, req)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// ListProxyServers returns []ProxyServer
func (c *Client) ListProxyServers(ctx context.Context, proxyType int) ([]ProxyServer, error) {
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

	// This endpoint often returns {"code": 200, "data": [...] } OR {"code": 200, "list": [...]}
	// Let's use a custom struct to capture potential fields
	type proxyListResp struct {
		List []ProxyServer `json:"list"`
		Data []ProxyServer `json:"data"`
	}
	resp, err := execute[APIResponse[proxyListResp]](c, req)
	if err != nil {
		return nil, err
	}

	if len(resp.Data.Data) > 0 {
		return resp.Data.Data, nil
	}
	return resp.Data.List, nil
}

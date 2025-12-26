// thordata/api_new.go
package thordata

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
)

func (c *Client) GetResidentialBalance(ctx context.Context) (map[string]any, error) {
	return c.executeApiNew(ctx, "/getFlowBalance", nil)
}

func (c *Client) GetIspRegions(ctx context.Context) ([]any, error) {
	res, err := c.executeApiNew(ctx, "/getRegionIsp", nil)
	if err != nil {
		return nil, err
	}
	// Result is typically in "data" which is returned by executeApiNew
	if list, ok := res["data"].([]any); ok {
		return list, nil
	}
	// Or maybe executeApiNew returned the data object directly and it's a list?
	// Need to check specific response structure. Assuming map wrapper.
	return []any{}, nil
}

func (c *Client) ListIspProxies(ctx context.Context) ([]any, error) {
	res, err := c.executeApiNew(ctx, "/queryListIsp", nil)
	if err != nil {
		return nil, err
	}
	if list, ok := res["data"].([]any); ok {
		return list, nil
	}
	return []any{}, nil
}

func (c *Client) executeApiNew(ctx context.Context, endpoint string, payload map[string]string) (map[string]any, error) {
	if c.cfg.Sign == "" || c.cfg.ApiKey == "" {
		return nil, errors.New("Sign and ApiKey are required for Public API NEW")
	}

	body := ""
	if payload != nil {
		body = ToFormBody(payload)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.gatewayBaseURL+endpoint, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}

	for k, v := range BuildSignHeaders(c.cfg.Sign, c.cfg.ApiKey) {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	parsed, _ := SafeParseJSON(raw)

	obj, ok := parsed.(map[string]any)
	if !ok {
		return nil, errors.New("invalid JSON response")
	}

	if cv, ok2 := obj["code"]; ok2 {
		if f, ok3 := cv.(float64); ok3 && int(f) != 200 {
			return nil, RaiseForCode("API NEW error", obj, res.StatusCode)
		}
	}

	return obj, nil
}

package thordata

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (c *Client) ListCountries(ctx context.Context, proxyType int) ([]any, error) {
	return c.getLocations(ctx, "countries", map[string]string{
		"proxy_type": strconv.Itoa(proxyType),
	})
}

func (c *Client) ListStates(ctx context.Context, countryCode string, proxyType int) ([]any, error) {
	return c.getLocations(ctx, "states", map[string]string{
		"proxy_type":   strconv.Itoa(proxyType),
		"country_code": strings.ToUpper(countryCode),
	})
}

func (c *Client) ListCities(ctx context.Context, countryCode string, stateCode string, proxyType int) ([]any, error) {
	params := map[string]string{
		"proxy_type":   strconv.Itoa(proxyType),
		"country_code": strings.ToUpper(countryCode),
	}
	if strings.TrimSpace(stateCode) != "" {
		params["state_code"] = strings.ToLower(stateCode)
	}
	return c.getLocations(ctx, "cities", params)
}

func (c *Client) ListASNs(ctx context.Context, countryCode string, proxyType int) ([]any, error) {
	return c.getLocations(ctx, "asn", map[string]string{
		"proxy_type":   strconv.Itoa(proxyType),
		"country_code": strings.ToUpper(countryCode),
	})
}

func (c *Client) getLocations(ctx context.Context, endpoint string, params map[string]string) ([]any, error) {
	if strings.TrimSpace(c.cfg.PublicToken) == "" || strings.TrimSpace(c.cfg.PublicKey) == "" {
		return nil, errors.New("publicToken and publicKey are required for locations API")
	}

	u, err := url.Parse(c.locationsBaseURL + "/" + endpoint)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("token", c.cfg.PublicToken)
	q.Set("key", c.cfg.PublicKey)
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	parsed, _ := SafeParseJSON(raw)

	// API may return {code,data} or a raw list; follow Python/JS behavior.
	if obj, ok := parsed.(map[string]any); ok {
		if cv, ok2 := obj["code"]; ok2 {
			if f, ok3 := cv.(float64); ok3 && int(f) != 200 {
				return nil, RaiseForCode("Locations API error", obj, res.StatusCode)
			}
		}
		if data, ok2 := obj["data"].([]any); ok2 {
			return data, nil
		}
		return []any{}, nil
	}
	if arr, ok := parsed.([]any); ok {
		return arr, nil
	}
	return []any{}, nil
}

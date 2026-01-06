package thordata

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (c *Client) ListCountries(ctx context.Context, proxyType int) ([]Country, error) {
	return getLocations[Country](c, ctx, "countries", map[string]string{
		"proxy_type": strconv.Itoa(proxyType),
	})
}

func (c *Client) ListStates(ctx context.Context, countryCode string, proxyType int) ([]State, error) {
	return getLocations[State](c, ctx, "states", map[string]string{
		"proxy_type":   strconv.Itoa(proxyType),
		"country_code": strings.ToUpper(countryCode),
	})
}

func (c *Client) ListCities(ctx context.Context, countryCode string, stateCode string, proxyType int) ([]City, error) {
	params := map[string]string{
		"proxy_type":   strconv.Itoa(proxyType),
		"country_code": strings.ToUpper(countryCode),
	}
	if strings.TrimSpace(stateCode) != "" {
		params["state_code"] = strings.ToLower(stateCode)
	}
	return getLocations[City](c, ctx, "cities", params)
}

func (c *Client) ListASNs(ctx context.Context, countryCode string, proxyType int) ([]ASN, error) {
	return getLocations[ASN](c, ctx, "asn", map[string]string{
		"proxy_type":   strconv.Itoa(proxyType),
		"country_code": strings.ToUpper(countryCode),
	})
}

func getLocations[T any](c *Client, ctx context.Context, endpoint string, params map[string]string) ([]T, error) {
	if strings.TrimSpace(c.cfg.PublicToken) == "" || strings.TrimSpace(c.cfg.PublicKey) == "" {
		return nil, errors.New("publicToken and publicKey are required")
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

	// Wrapper: {"code": 200, "data": [...]} OR raw list [...]
	// Locations API is a bit inconsistent. Let's try wrapped first.

	// Strategy: Use execute with APIResponse[[]T]. If data is empty but no error, maybe it returned direct list?
	// But according to spec/tests, it usually wraps in "data".

	resp, err := execute[APIResponse[[]T]](c, req)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

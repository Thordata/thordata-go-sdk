package thordata

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type UniversalOptions struct {
	URL            string
	JSRender       bool
	OutputFormat   string // "html" or "png"
	Country        string
	BlockResources string
	CleanContent   string
	Wait           *int
	WaitFor        string
	Headers        []map[string]string
	Cookies        []map[string]string
	Extra          map[string]string
}

// UniversalScrape returns *UniversalResponse
func (c *Client) UniversalScrape(ctx context.Context, opt UniversalOptions) (*UniversalResponse, error) {
	if strings.TrimSpace(opt.URL) == "" {
		return nil, errors.New("url is required")
	}

	format := strings.ToLower(strings.TrimSpace(opt.OutputFormat))
	if format == "" {
		format = "html"
	}
	if format != "html" && format != "png" {
		return nil, errors.New(`invalid outputFormat; supported: "html", "png"`)
	}

	payload := map[string]string{
		"url":       opt.URL,
		"js_render": boolToStr(opt.JSRender),
		"type":      format,
	}

	if opt.Country != "" {
		payload["country"] = strings.ToLower(opt.Country)
	}
	if opt.BlockResources != "" {
		payload["block_resources"] = opt.BlockResources
	}
	if opt.CleanContent != "" {
		payload["clean_content"] = opt.CleanContent
	}
	if opt.Wait != nil {
		payload["wait"] = intToStr(*opt.Wait)
	}
	if opt.WaitFor != "" {
		payload["wait_for"] = opt.WaitFor
	}
	if len(opt.Headers) > 0 {
		b, _ := json.Marshal(opt.Headers)
		payload["headers"] = string(b)
	}
	if len(opt.Cookies) > 0 {
		b, _ := json.Marshal(opt.Cookies)
		payload["cookies"] = string(b)
	}
	for k, v := range opt.Extra {
		payload[k] = v
	}

	body := ToFormBody(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.universalURL, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	for k, v := range BuildAuthHeaders(c.cfg.ScraperToken) {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	// Universal API returns fields directly in root (code, html, png)
	res, err := execute[UniversalResponse](c, req)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func boolToStr(v bool) string {
	if v {
		return "True"
	}
	return "False"
}

func intToStr(v int) string {
	return strconv.Itoa(v)
}

package thordata

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
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

func (c *Client) UniversalScrape(ctx context.Context, opt UniversalOptions) (any, error) {
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

	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = res.Body.Close() }()

	raw, _ := io.ReadAll(res.Body)

	// Try JSON first
	parsed, _ := SafeParseJSON(raw)
	if obj, ok := parsed.(map[string]any); ok {
		// API code check
		if cv, ok2 := obj["code"]; ok2 {
			if f, ok3 := cv.(float64); ok3 && int(f) != 200 {
				return nil, RaiseForCode("Universal API error", obj, res.StatusCode)
			}
		}

		// Extract html
		if hv, ok2 := obj["html"]; ok2 {
			return toString(hv), nil
		}

		// Extract png (base64)
		if pv, ok2 := obj["png"]; ok2 {
			s := toString(pv)
			data, err := decodeBase64MaybeDataURI(s)
			if err != nil {
				return nil, err
			}
			return data, nil
		}
	}

	// If response is not JSON, return raw content based on requested format
	if format == "png" {
		return raw, nil
	}
	return string(raw), nil
}

func decodeBase64MaybeDataURI(s string) ([]byte, error) {
	if strings.TrimSpace(s) == "" {
		return nil, errors.New("empty png data")
	}
	if idx := strings.Index(s, ","); idx >= 0 {
		s = s[idx+1:]
	}
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, " ", "")
	b, err := base64.StdEncoding.DecodeString(padBase64(s))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func padBase64(s string) string {
	switch len(s) % 4 {
	case 2:
		return s + "=="
	case 3:
		return s + "="
	default:
		return s
	}
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

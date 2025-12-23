package thordata

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type SerpOptions struct {
	Query        string
	Engine       string
	Num          int
	Start        int
	Country      string
	Language     string
	SearchType   string
	Device       string
	RenderJS     *bool
	NoCache      *bool
	OutputFormat string // "json" or "html"
	Extra        map[string]string
}

var tbmMap = map[string]string{
	"images":   "isch",
	"shopping": "shop",
	"news":     "nws",
	"videos":   "vid",
	"isch":     "isch",
	"shop":     "shop",
	"nws":      "nws",
	"vid":      "vid",
}

func (c *Client) SerpSearch(ctx context.Context, opt SerpOptions) (any, error) {
	if strings.TrimSpace(opt.Query) == "" {
		return nil, errors.New("query is required")
	}
	engine := strings.ToLower(strings.TrimSpace(opt.Engine))
	if engine == "" {
		engine = "google"
	}
	out := strings.ToLower(strings.TrimSpace(opt.OutputFormat))
	if out == "" {
		out = "json"
	}

	payload := map[string]string{
		"engine": engine,
		"json":   "1",
	}

	if out == "html" {
		payload["json"] = "0"
	}

	if engine == "yandex" {
		payload["text"] = opt.Query
	} else {
		payload["q"] = opt.Query
	}

	if opt.Num > 0 {
		payload["num"] = itoa(opt.Num)
	}
	if opt.Start > 0 {
		payload["start"] = itoa(opt.Start)
	}
	if opt.Country != "" {
		payload["gl"] = strings.ToLower(opt.Country)
	}
	if opt.Language != "" {
		payload["hl"] = strings.ToLower(opt.Language)
	}
	if opt.SearchType != "" {
		st := strings.ToLower(opt.SearchType)
		if v, ok := tbmMap[st]; ok {
			payload["tbm"] = v
		} else {
			payload["tbm"] = st
		}
	}
	if opt.Device != "" {
		payload["device"] = strings.ToLower(opt.Device)
	}
	if opt.RenderJS != nil {
		if *opt.RenderJS {
			payload["render_js"] = "True"
		} else {
			payload["render_js"] = "False"
		}
	}
	if opt.NoCache != nil {
		if *opt.NoCache {
			payload["no_cache"] = "True"
		} else {
			payload["no_cache"] = "False"
		}
	}

	for k, v := range opt.Extra {
		payload[k] = v
	}

	body := ToFormBody(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", c.serpURL, bytes.NewBufferString(body))
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
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)

	if out == "html" {
		return map[string]any{"html": string(raw)}, nil
	}

	parsed, _ := SafeParseJSON(raw)
	if obj, ok := parsed.(map[string]any); ok {
		if codeVal, ok2 := obj["code"]; ok2 {
			if f, ok3 := codeVal.(float64); ok3 && int(f) != 200 {
				return nil, RaiseForCode("SERP API error", obj, res.StatusCode)
			}
		}
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		if obj, ok := parsed.(map[string]any); ok {
			return nil, RaiseForCode("SERP HTTP error", obj, res.StatusCode)
		}
		return nil, errors.New("SERP request failed")
	}

	return parsed, nil
}

func itoa(n int) string {
	return strconv.Itoa(n)
}

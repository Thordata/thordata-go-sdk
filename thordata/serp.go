package thordata

import (
	"bytes"
	"context"
	"errors"
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

// SerpSearch now returns *SerpResponse (strong type) for JSON.
func (c *Client) SerpSearch(ctx context.Context, opt SerpOptions) (*SerpResponse, error) {
	if strings.TrimSpace(opt.Query) == "" {
		return nil, errors.New("query is required")
	}

	engine := normalizeEngine(strings.ToLower(strings.TrimSpace(opt.Engine)))
	if engine == "" {
		engine = "google"
	}

	out := strings.ToLower(strings.TrimSpace(opt.OutputFormat))
	if out == "" {
		out = "json"
	}
	if out == "html" {
		return nil, errors.New("outputFormat=html is not supported in strongly-typed SerpSearch, use raw HTTP request if needed")
	}

	payload := map[string]string{
		"engine": engine,
		"json":   "1",
	}

	if engine == "yandex" {
		payload["text"] = opt.Query
	} else {
		payload["q"] = opt.Query
	}

	if opt.Num > 0 {
		payload["num"] = strconv.Itoa(opt.Num)
	}
	if opt.Start > 0 {
		payload["start"] = strconv.Itoa(opt.Start)
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.serpURL, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	for k, v := range BuildAuthHeaders(c.cfg.ScraperToken) {
		req.Header.Set(k, v)
	}
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	// JSON
	res, err := execute[SerpResponse](c, req)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func normalizeEngine(engine string) string {
	switch strings.ToLower(strings.TrimSpace(engine)) {
	case "google_search":
		return "google"
	case "bing_search":
		return "bing"
	case "yandex_search":
		return "yandex"
	case "duckduckgo_search":
		return "duckduckgo"
	default:
		return strings.ToLower(strings.TrimSpace(engine))
	}
}

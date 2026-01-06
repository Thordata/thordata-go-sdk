// thordata/client.go
package thordata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type BaseURLs struct {
	ScraperAPIBaseURL    string
	UniversalAPIBaseURL  string
	WebScraperAPIBaseURL string
	LocationsBaseURL     string
}

type Config struct {
	ScraperToken string
	PublicToken  string
	PublicKey    string

	Timeout   time.Duration
	UserAgent string

	BaseURLs *BaseURLs
}

type Client struct {
	cfg  Config
	http *http.Client
	base BaseURLs

	serpURL            string
	scraperBuilderURL  string
	universalURL       string
	scraperStatusURL   string
	scraperDownloadURL string
	locationsBaseURL   string
	videoBuilderURL    string

	// Public API URLs
	usageStatsURL      string
	proxyUsersURL      string
	whitelistURL       string
	proxyListURL       string
	proxyExpirationURL string
	taskListURL        string
}

func NewClient(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.ScraperToken) == "" {
		return nil, errors.New("scraperToken is required")
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = BuildUserAgent("1.0.1")
	}

	base := resolveBaseURLs(cfg.BaseURLs)
	c := &Client{
		cfg:  cfg,
		http: &http.Client{Timeout: cfg.Timeout},
		base: base,
	}

	c.serpURL = strings.TrimRight(base.ScraperAPIBaseURL, "/") + "/request"
	c.scraperBuilderURL = strings.TrimRight(base.ScraperAPIBaseURL, "/") + "/builder"
	c.videoBuilderURL = strings.TrimRight(base.ScraperAPIBaseURL, "/") + "/video_builder"
	c.universalURL = strings.TrimRight(base.UniversalAPIBaseURL, "/") + "/request"
	c.scraperStatusURL = strings.TrimRight(base.WebScraperAPIBaseURL, "/") + "/tasks-status"
	c.scraperDownloadURL = strings.TrimRight(base.WebScraperAPIBaseURL, "/") + "/tasks-download"
	c.locationsBaseURL = strings.TrimRight(base.LocationsBaseURL, "/")

	// Public API endpoints
	apiBase := strings.Replace(c.locationsBaseURL, "/locations", "", 1)
	c.usageStatsURL = apiBase + "/account/usage-statistics"
	c.proxyUsersURL = apiBase + "/proxy-users"
	c.whitelistURL = "https://api.thordata.com/api/whitelisted-ips"
	c.proxyListURL = "https://api.thordata.com/api/proxy/proxy-list"
	c.proxyExpirationURL = apiBase + "/proxy/expiration-time"
	c.taskListURL = strings.TrimRight(base.WebScraperAPIBaseURL, "/") + "/tasks-list"

	return c, nil
}

func resolveBaseURLs(override *BaseURLs) BaseURLs {
	def := BaseURLs{
		ScraperAPIBaseURL:    getenvDefault("THORDATA_SCRAPERAPI_BASE_URL", "https://scraperapi.thordata.com"),
		UniversalAPIBaseURL:  getenvDefault("THORDATA_UNIVERSALAPI_BASE_URL", "https://universalapi.thordata.com"),
		WebScraperAPIBaseURL: getenvDefault("THORDATA_WEB_SCRAPER_API_BASE_URL", "https://openapi.thordata.com/api/web-scraper-api"),
		LocationsBaseURL:     getenvDefault("THORDATA_LOCATIONS_BASE_URL", "https://openapi.thordata.com/api/locations"),
	}
	if override == nil {
		return def
	}
	if override.ScraperAPIBaseURL != "" {
		def.ScraperAPIBaseURL = override.ScraperAPIBaseURL
	}
	if override.UniversalAPIBaseURL != "" {
		def.UniversalAPIBaseURL = override.UniversalAPIBaseURL
	}
	if override.WebScraperAPIBaseURL != "" {
		def.WebScraperAPIBaseURL = override.WebScraperAPIBaseURL
	}
	if override.LocationsBaseURL != "" {
		def.LocationsBaseURL = override.LocationsBaseURL
	}
	return def
}

func getenvDefault(key, def string) string {
	v := os.Getenv(key)
	if strings.TrimSpace(v) == "" {
		return def
	}
	return v
}

// CloseIdleConnections closes any idle connections in the client's transport.
// This is useful if you are creating many ephemeral clients.
func (c *Client) CloseIdleConnections() {
	c.http.CloseIdleConnections()
}

// execute sends the request and unmarshals the response into T.
// It handles standard Thordata error wrapping ({"code": != 200}).
func execute[T any](c *Client, req *http.Request) (T, error) {
	var zero T
	req.Header.Set("User-Agent", c.cfg.UserAgent)

	res, err := c.http.Do(req)
	if err != nil {
		return zero, err
	}
	defer res.Body.Close() //nolint:errcheck

	raw, _ := io.ReadAll(res.Body)

	// 1. Check for API error wrapper first (without consuming strict T structure)
	var meta map[string]any
	if err := json.Unmarshal(raw, &meta); err == nil {
		if err := checkAPIError(meta, res.StatusCode); err != nil {
			return zero, err
		}
	} else {
		// Not JSON? If T is string/[]byte we might allow it, but Thordata APIs usually return JSON.
		// For HTML scraping, we handle it separately in UniversalScrape.
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return zero, fmt.Errorf("http error %d: %s", res.StatusCode, string(raw))
		}
	}

	// 2. Unmarshal into target type
	var data T
	if err := json.Unmarshal(raw, &data); err != nil {
		return zero, fmt.Errorf("failed to decode response: %w (body: %s)", err, string(raw))
	}

	return data, nil
}

func checkAPIError(payload map[string]any, httpStatus int) error {
	apiCode := 0
	if v, ok := payload["code"]; ok {
		if f, ok2 := v.(float64); ok2 {
			apiCode = int(f)
		}
	}

	if (apiCode != 0 && apiCode != 200) || (httpStatus < 200 || httpStatus >= 300) {
		return RaiseForCode("API Error", payload, httpStatus)
	}
	return nil
}

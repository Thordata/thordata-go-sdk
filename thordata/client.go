package thordata

import (
	"errors"
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
}

func NewClient(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.ScraperToken) == "" {
		return nil, errors.New("scraperToken is required")
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = BuildUserAgent("0.0.0")
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

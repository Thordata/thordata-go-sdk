package thordata

import (
	"errors"
	"net/url"
	"os"
	"strings"
)

// GetBrowserConnectionURL returns the WebSocket URL for connecting to Thordata Scraping Browser.
// This URL can be used with Playwright Go, Chrome DevTools Protocol (CDP), or other browser automation tools.
//
// Example usage with Playwright Go:
//
//	wsURL, _ := client.GetBrowserConnectionURL("username", "password")
//	browser, _ := playwright.Connect(wsURL)
//
// If username/password are not provided, they will be read from environment variables:
// - THORDATA_BROWSER_USERNAME
// - THORDATA_BROWSER_PASSWORD
func (c *Client) GetBrowserConnectionURL(username, password string) (string, error) {
	// Use provided credentials or fall back to environment variables
	user := strings.TrimSpace(username)
	if user == "" {
		user = strings.TrimSpace(os.Getenv("THORDATA_BROWSER_USERNAME"))
	}

	pwd := strings.TrimSpace(password)
	if pwd == "" {
		pwd = strings.TrimSpace(os.Getenv("THORDATA_BROWSER_PASSWORD"))
	}

	if user == "" || pwd == "" {
		return "", errors.New("browser credentials missing: set THORDATA_BROWSER_USERNAME and THORDATA_BROWSER_PASSWORD, or pass them as arguments")
	}

	// Add prefix if not present
	prefix := "td-customer-"
	if !strings.HasPrefix(user, prefix) {
		user = prefix + user
	}

	// URL encode credentials (QueryEscape is safe for URL userinfo)
	safeUser := url.QueryEscape(user)
	safePass := url.QueryEscape(pwd)

	// Return WebSocket URL for Thordata Scraping Browser
	return "wss://" + safeUser + ":" + safePass + "@ws-browser.thordata.com", nil
}

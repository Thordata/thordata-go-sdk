package thordata

import (
	"os"
	"strings"
	"testing"
)

func TestGetBrowserConnectionURL(t *testing.T) {
	client, _ := NewClient(Config{ScraperToken: "dummy"})

	// Test with explicit credentials
	url, err := client.GetBrowserConnectionURL("testuser", "testpass")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(url, "wss://") {
		t.Fatalf("expected wss:// URL, got %s", url)
	}
	if !strings.Contains(url, "td-customer-testuser") {
		t.Fatalf("expected username in URL, got %s", url)
	}

	// Test with env vars
	_ = os.Setenv("THORDATA_BROWSER_USERNAME", "envuser")
	_ = os.Setenv("THORDATA_BROWSER_PASSWORD", "envpass")
	defer func() {
		_ = os.Unsetenv("THORDATA_BROWSER_USERNAME")
		_ = os.Unsetenv("THORDATA_BROWSER_PASSWORD")
	}()

	url2, err := client.GetBrowserConnectionURL("", "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(url2, "td-customer-envuser") {
		t.Fatalf("expected env username in URL, got %s", url2)
	}

	// Test missing credentials
	_ = os.Unsetenv("THORDATA_BROWSER_USERNAME")
	_ = os.Unsetenv("THORDATA_BROWSER_PASSWORD")
	_, err = client.GetBrowserConnectionURL("", "")
	if err == nil {
		t.Fatal("expected error for missing credentials")
	}
}

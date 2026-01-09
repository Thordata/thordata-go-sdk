//go:build integration
// +build integration

package thordata

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// loadDotEnv is a minimal .env loader for integration tests
func loadDotEnv(path string) {
	b, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		s := strings.TrimSpace(line)
		if s == "" || strings.HasPrefix(s, "#") {
			continue
		}

		i := strings.Index(s, "=")
		if i <= 0 {
			continue
		}

		k := strings.TrimSpace(s[:i])
		v := strings.TrimSpace(s[i+1:])

		// Skip proxy env vars to avoid double-proxying
		switch strings.ToUpper(k) {
		case "HTTP_PROXY", "HTTPS_PROXY", "ALL_PROXY", "NO_PROXY":
			continue
		}

		// Strip quotes
		if len(v) >= 2 {
			if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '\'' && v[len(v)-1] == '\'') {
				v = v[1 : len(v)-1]
			}
		}

		if os.Getenv(k) == "" {
			_ = os.Setenv(k, v)
		}
	}
}

func envNonEmpty(k string) string {
	return strings.TrimSpace(os.Getenv(k))
}

func isStrict() bool {
	return strings.ToLower(envNonEmpty("THORDATA_INTEGRATION_STRICT")) == "true"
}

// looksLikeNetworkInterference detects common errors from proxy/TUN interference
func looksLikeNetworkInterference(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	bad := []string{
		"http: server gave HTTP response to HTTPS client",
		"unexpected EOF",
		"forcibly closed",
		"connection reset",
		"timeout",
		"context deadline exceeded",
		"TLS handshake timeout",
		"wrong version number",
		"packet length too long",
	}
	for _, k := range bad {
		if strings.Contains(s, k) {
			return true
		}
	}
	return false
}

func TestProxyProtocolsIntegration(t *testing.T) {
	// Load .env from various possible locations
	loadDotEnv(".env")
	loadDotEnv("../.env")
	loadDotEnv("../../.env")

	host := envNonEmpty("THORDATA_PROXY_HOST")
	user := envNonEmpty("THORDATA_RESIDENTIAL_USERNAME")
	pass := envNonEmpty("THORDATA_RESIDENTIAL_PASSWORD")

	if host == "" || user == "" || pass == "" {
		t.Skip("missing THORDATA_PROXY_HOST / THORDATA_RESIDENTIAL_USERNAME / THORDATA_RESIDENTIAL_PASSWORD")
	}

	// Check for upstream proxy
	upstream := envNonEmpty("THORDATA_UPSTREAM_PROXY")
	if upstream != "" {
		t.Logf("🔗 Upstream proxy detected: %s", upstream)
		t.Log("   Note: Go SDK proxy chaining requires custom implementation")
	}

	// Determine protocols to test
	protos := []string{"https", "socks5h"}
	if strings.ToLower(envNonEmpty("THORDATA_INTEGRATION_HTTP")) == "true" {
		protos = append([]string{"http"}, protos...)
	}

	c, err := NewClient(Config{
		ScraperToken: "dummy",
		Timeout:      120 * time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	target := "https://ipinfo.thordata.com"

	for _, proto := range protos {
		t.Run(proto, func(t *testing.T) {
			t.Logf("\n--- Testing protocol: %s ---", proto)

			pcfg := &ProxyConfig{
				Product:  ProxyResidential,
				Username: user,
				Password: pass,
				Host:     host,
				Port:     9999,
				Protocol: proto,
				Country:  "us",
			}

			var lastErr error
			var success bool

			for attempt := 1; attempt <= 3; attempt++ {
				t.Logf("  Attempt %d/3...", attempt)

				resp, err := c.ProxyGet(ctx, target, pcfg)
				if err == nil {
					defer resp.Body.Close()

					if resp.StatusCode == 200 {
						t.Logf("  ✓ %s passed!", proto)
						success = true
						lastErr = nil
						break
					}

					lastErr = fmt.Errorf("non-200 status: %d", resp.StatusCode)
				} else {
					lastErr = err
				}

				t.Logf("  ✗ Error: %v", lastErr)

				if attempt < 3 {
					time.Sleep(2 * time.Second)
				}
			}

			if !success && lastErr != nil {
				if !isStrict() && looksLikeNetworkInterference(lastErr) {
					t.Skipf("⚠️  Skipping due to network interference: %v", lastErr)
				}
				t.Fatalf("❌ Protocol %s failed: %v", proto, lastErr)
			}
		})
	}
}

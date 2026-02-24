package thordata

import (
	"fmt"
	"os"
	"strings"
)

// ValidateEnv checks if required environment variables are set and returns helpful error messages.
func ValidateEnv() error {
	var missing []string
	var warnings []string

	// Check for at least one authentication method
	scraperToken := strings.TrimSpace(os.Getenv("THORDATA_SCRAPER_TOKEN"))
	publicToken := strings.TrimSpace(os.Getenv("THORDATA_PUBLIC_TOKEN"))
	publicKey := strings.TrimSpace(os.Getenv("THORDATA_PUBLIC_KEY"))

	if scraperToken == "" && (publicToken == "" || publicKey == "") {
		missing = append(missing, "THORDATA_SCRAPER_TOKEN (or THORDATA_PUBLIC_TOKEN + THORDATA_PUBLIC_KEY)")
	}

	// Check for proxy credentials (at least one product)
	hasProxyCreds := false
	proxyProducts := []string{
		"RESIDENTIAL", "DATACENTER", "MOBILE", "ISP",
	}
	for _, product := range proxyProducts {
		user := os.Getenv("THORDATA_" + product + "_USERNAME")
		pwd := os.Getenv("THORDATA_" + product + "_PASSWORD")
		if user != "" && pwd != "" {
			hasProxyCreds = true
			break
		}
	}

	// Proxy credentials are optional (some APIs don't need them)
	// But warn if user might expect proxy functionality
	if !hasProxyCreds {
		whitelist := strings.ToLower(strings.TrimSpace(os.Getenv("THORDATA_PROXY_WHITELIST")))
		if whitelist != "1" && whitelist != "true" && whitelist != "yes" && whitelist != "y" {
			warnings = append(warnings, "No proxy credentials found. Proxy Network features will not work.")
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	if len(warnings) > 0 {
		// Return warnings as a special error type that can be checked
		return &EnvWarning{Message: strings.Join(warnings, " ")}
	}

	return nil
}

// EnvWarning represents a non-fatal environment configuration warning.
type EnvWarning struct {
	Message string
}

func (e *EnvWarning) Error() string {
	return e.Message
}

// IsWarning checks if an error is just a warning (non-fatal).
func IsWarning(err error) bool {
	_, ok := err.(*EnvWarning)
	return ok
}

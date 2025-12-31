package example

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Thordata/thordata-go-sdk/examples/internal/env"
	"github.com/Thordata/thordata-go-sdk/thordata"
)

func LoadEnv() {
	_ = env.LoadDotEnv(".env")
}

func Env(name string) string {
	return strings.TrimSpace(os.Getenv(name))
}

// SkipIfMissing prints a friendly message and returns true if any required env is missing.
// Caller should "return" from main() if it returns true.
func SkipIfMissing(required ...string) bool {
	var missing []string
	for _, k := range required {
		if Env(k) == "" {
			missing = append(missing, k)
		}
	}
	if len(missing) == 0 {
		return false
	}
	fmt.Println("Skipping example: missing env:", strings.Join(missing, ", "))
	fmt.Println("Tip: copy .env.example to .env and fill values, then re-run.")
	return true
}

// NewClient creates a client using env vars. Requires THORDATA_SCRAPER_TOKEN to be set by the example.
func NewClient(timeout time.Duration) (*thordata.Client, error) {
	return thordata.NewClient(thordata.Config{
		ScraperToken: Env("THORDATA_SCRAPER_TOKEN"),
		PublicToken:  Env("THORDATA_PUBLIC_TOKEN"),
		PublicKey:    Env("THORDATA_PUBLIC_KEY"),
		Timeout:      timeout,
	})
}

// NewClientAllowDummyScraper allows examples that don't need ScraperToken (e.g. Proxy Network).
func NewClientAllowDummyScraper(timeout time.Duration) (*thordata.Client, error) {
	scraper := Env("THORDATA_SCRAPER_TOKEN")
	if scraper == "" {
		scraper = "dummy"
	}
	return thordata.NewClient(thordata.Config{
		ScraperToken: scraper,
		PublicToken:  Env("THORDATA_PUBLIC_TOKEN"),
		PublicKey:    Env("THORDATA_PUBLIC_KEY"),
		Timeout:      timeout,
	})
}

func PrintJSON(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("%v\n", v)
		return
	}
	fmt.Println(string(b))
}

func Truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

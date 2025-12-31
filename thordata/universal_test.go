package thordata

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestUniversalScrapeHTML_Offline(t *testing.T) {
	var gotAuth string
	var gotBody url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/request" {
			http.NotFound(w, r)
			return
		}
		gotAuth = r.Header.Get("Authorization")
		b := make([]byte, r.ContentLength)
		_, _ = r.Body.Read(b)
		_ = r.Body.Close()

		vals, _ := url.ParseQuery(string(b))
		gotBody = vals

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"html":"<h1>Hello</h1>"}`))
	}))
	defer srv.Close()

	client, err := NewClient(Config{
		ScraperToken: "token",
		Timeout:      10 * time.Second,
		BaseURLs: &BaseURLs{
			ScraperAPIBaseURL:    srv.URL,
			UniversalAPIBaseURL:  srv.URL,
			WebScraperAPIBaseURL: srv.URL,
			LocationsBaseURL:     srv.URL,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	out, err := client.UniversalScrape(context.Background(), UniversalOptions{
		URL:          "https://example.com",
		JSRender:     true,
		OutputFormat: "html",
	})
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer token" {
		t.Fatalf("expected auth header, got %q", gotAuth)
	}
	if gotBody.Get("js_render") != "True" || gotBody.Get("type") != "html" {
		t.Fatalf("unexpected payload: %v", gotBody)
	}
	if s, ok := out.(string); !ok || !strings.Contains(s, "Hello") {
		t.Fatalf("unexpected output: %#v", out)
	}
}

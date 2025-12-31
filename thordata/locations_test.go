package thordata

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLocations_Offline(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/countries" {
			http.NotFound(w, r)
			return
		}
		if r.URL.Query().Get("token") != "pub" || r.URL.Query().Get("key") != "key" {
			t.Fatalf("missing token/key query params")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":200,"data":[{"country_code":"US","country_name":"United States"}]}`))
	}))
	defer srv.Close()

	client, err := NewClient(Config{
		ScraperToken: "token",
		PublicToken:  "pub",
		PublicKey:    "key",
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

	items, err := client.ListCountries(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("unexpected list size: %d", len(items))
	}
}

package thordata

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTasks_Offline(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/builder":
			if r.Header.Get("Authorization") == "" {
				t.Fatalf("missing Authorization header")
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"code":200,"data":{"task_id":"t1"}}`))
		case "/tasks-status":
			if r.Header.Get("token") != "pub" || r.Header.Get("key") != "key" {
				t.Fatalf("missing public headers")
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"code":200,"data":[{"task_id":"t1","status":"ready"}]}`))
		case "/tasks-download":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"code":200,"data":{"download":"https://example.com/file.json"}}`))
		default:
			http.NotFound(w, r)
		}
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

	taskID, err := client.CreateScraperTask(context.Background(), ScraperTaskOptions{
		FileName:      "f",
		SpiderID:      "s1",
		SpiderName:    "example.com",
		Parameters:    map[string]any{"url": "https://example.com"},
		IncludeErrors: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if taskID != "t1" {
		t.Fatalf("unexpected task id: %s", taskID)
	}

	status, err := client.GetTaskStatus(context.Background(), "t1")
	if err != nil {
		t.Fatal(err)
	}
	if status != "ready" {
		t.Fatalf("unexpected status: %s", status)
	}

	dl, err := client.GetTaskResult(context.Background(), "t1", "json")
	if err != nil {
		t.Fatal(err)
	}
	if dl == "" {
		t.Fatalf("expected download url")
	}
}

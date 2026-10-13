package scheduler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckDB(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
	}{
		{
			name:       "healthy database",
			statusCode: http.StatusOK,
			body:       `{"status":"healthy"}`,
		},
		{
			name:       "unhealthy database",
			statusCode: http.StatusOK,
			body:       `{"status":"unhealthy","error":"boom"}`,
		},
		{
			name:       "non-200 response",
			statusCode: http.StatusServiceUnavailable,
			body:       `{"status":"unhealthy"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var hits int
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				hits++
				if r.URL.Path != "/db-health" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer ts.Close()

			s := &Scheduler{
				baseURL:  ts.URL,
				interval: time.Minute,
				client:   ts.Client(),
			}

			s.checkDB(context.Background())

			if hits != 1 {
				t.Fatalf("expected 1 request, got %d", hits)
			}
		})
	}
}

func TestCheckDBConnectionError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	baseURL := ts.URL
	ts.Close()

	s := &Scheduler{
		baseURL:  baseURL,
		interval: time.Minute,
		client:   &http.Client{Timeout: time.Second},
	}

	s.checkDB(context.Background())
}

func TestIntervalFromEnv(t *testing.T) {
	t.Setenv("DB_HEALTH_INTERVAL", "30s")
	if got := intervalFromEnv(); got != 30*time.Second {
		t.Fatalf("expected 30s, got %s", got)
	}

	t.Setenv("DB_HEALTH_INTERVAL", "not-a-duration")
	if got := intervalFromEnv(); got != defaultInterval {
		t.Fatalf("expected default %s, got %s", defaultInterval, got)
	}

	t.Setenv("DB_HEALTH_INTERVAL", "")
	if got := intervalFromEnv(); got != defaultInterval {
		t.Fatalf("expected default %s, got %s", defaultInterval, got)
	}
}

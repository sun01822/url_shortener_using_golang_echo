package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/gommon/log"
)

const (
	defaultInterval = 10 * time.Minute
	defaultTimeout  = 5 * time.Second
)

// Scheduler periodically calls the database health API to verify that the
// database is live.
type Scheduler struct {
	baseURL  string
	interval time.Duration
	client   *http.Client
}

// New creates a Scheduler that pings the DB health endpoint of the API running
// at baseURL (e.g. "http://localhost:8080").
func New(baseURL string) *Scheduler {
	return &Scheduler{
		baseURL:  baseURL,
		interval: intervalFromEnv(),
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func intervalFromEnv() time.Duration {
	raw := os.Getenv("DB_HEALTH_INTERVAL")
	if raw == "" {
		return defaultInterval
	}

	interval, err := time.ParseDuration(raw)
	if err != nil || interval <= 0 {
		log.Warnf("invalid DB_HEALTH_INTERVAL %q, defaulting to %s", raw, defaultInterval)
		return defaultInterval
	}

	return interval
}

// Start runs an immediate health check and then repeats it every interval until
// ctx is cancelled.
func (s *Scheduler) Start(ctx context.Context) {
	log.Infof("database health scheduler started (interval %s)", s.interval)

	s.checkDB(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("database health scheduler stopped")
			return
		case <-ticker.C:
			s.checkDB(ctx)
		}
	}
}

func (s *Scheduler) checkDB(ctx context.Context) {
	url := fmt.Sprintf("%s/db-health", s.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("failed to build db health request: %v", err)
		return
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Errorf("database is not live: %v", err)
		return
	}
	defer resp.Body.Close()

	var health map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		log.Errorf("database health check returned an invalid response: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK || health["status"] != "healthy" {
		log.Errorf("database is not live: status=%d response=%v", resp.StatusCode, health)
		return
	}

	log.Infof("database is live: %v", health)
}

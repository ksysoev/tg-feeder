// Package crawler provides a client for interacting with the SomeAPI service.
package crawler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultTimeout = 5 * time.Second

// Config holds configuration for the APIClient.
type Config struct{}

// APIClient is a client for communicating with the SomeAPI service.
type APIClient struct {
	cli *http.Client
	cfg Config
}

// New creates a new APIClient with the provided configuration.
func New(cfg Config) *APIClient {
	return &APIClient{
		cfg: cfg,
		cli: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (a *APIClient) FetchPage(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.cli.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("health check failed with status code: %d", resp.StatusCode)
	}

	byte, err := io.ReadAll(resp.Body)

	return string(byte), nil
}

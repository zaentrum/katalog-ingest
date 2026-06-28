// Package api holds thin REST clients for katalog-api (reads) and
// katalog-manager-api (writes). The clients are deliberately minimal —
// they only know enough to support the scheduler's enqueue policy.
package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

// KatalogClient calls the read-only katalog-api endpoints.
type KatalogClient struct {
	base string
	http *http.Client
}

// NewKatalogClient returns a client pointing at base (e.g.
// http://katalog-api.stube.svc.cluster.local).
func NewKatalogClient(base string) *KatalogClient {
	return &KatalogClient{
		base: base,
		http: &http.Client{Timeout: 10 * time.Second},
	}
}

// PendingItems returns a slice of candidate items that need processing.
// Scaffold: returns an empty slice.
func (c *KatalogClient) PendingItems(ctx context.Context) ([]string, error) {
	if c == nil || c.base == "" {
		return nil, errors.New("katalog client unconfigured")
	}
	// TODO: GET /api/v1/items?filter=needs_processing -> decode -> return ids
	slog.Debug("katalog PendingItems stub")
	return []string{}, nil
}

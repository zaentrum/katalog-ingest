package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
)

// ManagerClient calls the write-side katalog-manager-api endpoints
// (specifically the /svc/v1/* service-to-service controller).
type ManagerClient struct {
	base string
	http *http.Client
}

// NewManagerClient returns a client pointing at base (e.g.
// http://katalog-manager-api.stube.svc.cluster.local).
func NewManagerClient(base string) *ManagerClient {
	return &ManagerClient{
		base: base,
		http: &http.Client{Timeout: 30 * time.Second},
	}
}

// MarkProcessing tells the manager an item has been enqueued for processing
// (used so the scheduler doesn't re-emit the same task on every tick).
// Scaffold: returns nil.
func (c *ManagerClient) MarkProcessing(ctx context.Context, itemID, task string) error {
	if c == nil || c.base == "" {
		return errors.New("manager client unconfigured")
	}
	// TODO: POST /svc/v1/items/{id}/mark-processing with task name body.
	slog.Debug("manager MarkProcessing stub", "item", itemID, "task", task)
	return nil
}

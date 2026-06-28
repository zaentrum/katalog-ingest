// Package scheduler decides what work to enqueue on the
// stube.processing.task.* Kafka topics on a fixed cadence.
package scheduler

import (
	"context"
	"log/slog"
	"time"

	"github.com/zaentrum/katalog-ingest/internal/api"
	"github.com/zaentrum/katalog-ingest/internal/kafka"
)

// Scheduler runs a periodic tick. Each tick it queries katalog-api for
// items that need processing, decides what task to fan out, and emits the
// corresponding Kafka message via katalog-manager-api or directly.
type Scheduler struct {
	Period  time.Duration
	Reader  *api.KatalogClient
	Writer  *api.ManagerClient
	Emitter *kafka.Producer
}

// New constructs a Scheduler. Any of the dependencies may be nil — the
// scheduler will log + no-op on tick for the missing piece. This makes
// scaffold-mode boot trivial.
func New(period time.Duration, r *api.KatalogClient, w *api.ManagerClient, p *kafka.Producer) *Scheduler {
	return &Scheduler{Period: period, Reader: r, Writer: w, Emitter: p}
}

// Run blocks until ctx is cancelled. Each tick calls Tick().
func (s *Scheduler) Run(ctx context.Context) {
	t := time.NewTicker(s.Period)
	defer t.Stop()
	slog.Info("scheduler started", "period", s.Period)
	for {
		select {
		case <-ctx.Done():
			slog.Info("scheduler stopped")
			return
		case <-t.C:
			s.Tick(ctx)
		}
	}
}

// Tick is one scheduling cycle. Currently a stub.
func (s *Scheduler) Tick(ctx context.Context) {
	// TODO: query Reader for the candidate set; decide policy; emit via
	// Emitter to topic stube.processing.task.*.
	slog.Debug("scheduler tick")
}

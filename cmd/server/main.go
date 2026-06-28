package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/zaentrum/katalog-ingest/internal/api"
	"github.com/zaentrum/katalog-ingest/internal/kafka"
	"github.com/zaentrum/katalog-ingest/internal/scheduler"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	addr := envOr("KATALOG_INGEST_ADDR", ":8080")
	kafkaBrokers := envOr("KAFKA_BROKERS", "platform-kafka-kafka-bootstrap.platform-event-streaming.svc:9093")
	katalogAPIURL := envOr("KATALOG_API_URL", "http://katalog-api.stube.svc.cluster.local")
	katalogMgrURL := envOr("KATALOG_MANAGER_API_URL", "http://katalog-manager-api.stube.svc.cluster.local")
	period := 30 * time.Second

	slog.Info("starting katalog-ingest",
		"addr", addr,
		"kafka_brokers", kafkaBrokers,
		"katalog_api_url", katalogAPIURL,
		"katalog_manager_api_url", katalogMgrURL,
	)

	producer, err := kafka.NewProducer(kafkaBrokers)
	if err != nil {
		slog.Warn("kafka producer init failed; scheduler will no-op", "err", err)
	}

	readClient := api.NewKatalogClient(katalogAPIURL)
	writeClient := api.NewManagerClient(katalogMgrURL)

	s := scheduler.New(period, readClient, writeClient, producer)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go s.Run(ctx)

	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	r.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		slog.Info("http listening", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen failed", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

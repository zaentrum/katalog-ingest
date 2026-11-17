# katalog-ingest

Go service that schedules catalog processing work for the zaentrum platform.
Drives transcoder / packager / analyzer / metadata-enricher via Kafka. Reads
candidate items from `katalog-api`, marks them in-flight via
`katalog-manager-api`, and emits `stube.processing.task.*` messages on a tick.

## Status

**Scaffold; processing logic lands in a later phase.**

The HTTP server, scheduler tick, and Kafka emitter are wired as stubs. The
service compiles and serves `/healthz` + `/metrics` so it can be deployed as a
placeholder. The reads/writes split keeps the scheduler out of the database:
it reads from `katalog-api` and writes to `katalog-manager-api` over HTTP, and
publishes the `stube.processing.task.*` topics for downstream workers.

## Layout

```
cmd/server/main.go                # chi router, /healthz, /metrics, scheduler wiring
internal/scheduler/scheduler.go   # fixed-cadence tick that emits processing tasks
internal/kafka/kafka.go           # Kafka producer/consumer wiring
internal/api/katalog.go           # read client for katalog-api
internal/api/manager.go           # write client for katalog-manager-api
k8s/                              # Deployment, Service, ServiceAccount, ServiceMonitor, GrafanaDashboard
Dockerfile
```

## Local development

```bash
go run ./cmd/server
curl -fsS http://localhost:8080/healthz
```

Configuration is env-driven:

| Env | Default |
| --- | --- |
| `KATALOG_INGEST_ADDR` | `:8080` |
| `KAFKA_BROKERS` | `kafka:9092` |
| `KATALOG_API_URL` | `http://katalog-api` |
| `KATALOG_MANAGER_API_URL` | `http://katalog-manager-api` |

The defaults are the names the platform's chart gives the broker and the two
services, and the chart sets all four. Run elsewhere, set them to yours.

## Build the container

```bash
docker build -t zaentrum/katalog-ingest .
```

Build and push the image to your own registry, then apply the `k8s/`
manifests and update the image reference for your environment.

## License

[MPL-2.0](LICENSE).

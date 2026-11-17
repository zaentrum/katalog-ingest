# katalog-ingest

katalog-ingest is the scaffold of the catalog's processing scheduler, a Go
service. It starts, serves `/healthz` and `/metrics`, and runs a tick every 30
seconds. The tick (`internal/scheduler`), its Kafka producer (`internal/kafka`)
and its clients of katalog-api and katalog-manager (`internal/api`) are stubs:
it reads nothing, writes nothing and publishes nothing. No other service calls
it; the platform's chart runs it with the media pipeline (`features.pipeline`).
How a change is proposed and what a pull request has to pass is in the
[guide to contributing](https://github.com/zaentrum/.github/blob/main/CONTRIBUTING.md);
what may be built is in [building from the specification](https://github.com/zaentrum/zaentrum/blob/main/AGENTS.md).

## Build

```sh
go build ./...
docker build -t katalog-ingest .
```

It needs Go 1.23 or newer (`go.mod`; its `toolchain` line names go1.24.1).
The image builds the binary in `golang:1.23-alpine` and runs it on
`gcr.io/distroless/static-debian12:nonroot` as a non-root user, on port 8080.
`--build-arg BASE=<registry>/library/` pulls the build image from a mirror; the
runtime image is the same whatever `BASE` says.

## Test

```sh
go vet ./...
go test ./...
```

The workflow `.github/workflows/build.yml` runs `go build ./...`,
`go vet ./...` and `go test ./...` on every pull request and on main. The one
test, in `cmd/server/main_test.go`, starts the real `main` on a free loopback
address with the broker and both catalog APIs pointed at a closed port, waits
for `GET /healthz` to answer 200, and stops it with SIGTERM.

## Run it

```sh
go run ./cmd/server
curl -fsS http://localhost:8080/healthz
```

`/healthz` answers `{"status":"ok"}`. The service reads `KATALOG_INGEST_ADDR`
(the listen address, `:8080`), `KAFKA_BROKERS`, `KATALOG_API_URL` and
`KATALOG_MANAGER_API_URL`. It logs them and connects to none of them, so it
needs no database, broker or other service; the defaults of the last three name
one installation's services and are not used. There is nothing in it to run
against the reference demo or against recorded answers.

## What the tests leave unchecked

- Everything but the start: the one test starts the server and checks its
  health route, and nothing more.
- The tick, the producer and the two clients do nothing, so no behaviour of
  the scheduler exists to check.

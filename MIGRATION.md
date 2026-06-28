# MIGRATION

`katalog-ingest` is a fresh Go service with no upstream lift. It replaces an
in-process scheduler that previously lived inside the catalog service.

## Why a fresh repo?

- The processing fan-out moves from HTTP claim-loop polling to Kafka. The
  scheduler is no longer business-logic glue inside the monolith; it is its
  own deployable so it can be scaled / paused independently.
- The read/write boundary lives at the HTTP layer. The scheduler reads from
  `katalog-api` and writes to `katalog-manager-api` so it never holds DB
  connections directly.

## Open items

- Confirm Kafka topic naming (`stube.processing.task.transcode` vs
  `.transcode.requested` etc).
- Decide leader election (single replica vs lease-based active/passive).
- Wire the consumer side for `stube.processing.result.*` so the scheduler
  can mark items done.

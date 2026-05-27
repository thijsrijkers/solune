# Solune - A NoSQL Database

Solune is a TCP-based key-value database written in Go.
It supports multiple named stores, JSON line commands (`get`, `set`, `delete`), process supervision, and disk-backed persistence in `db/*.solstr`.

Roadmap items are tracked in `doc/README.md`.

## Highlights

- TCP server on port `9000` by default
- One JSON command per line
- Multiple named stores
- Auto-increment key assignment when `set` is called without `key`
- Supervisor process that restarts worker on crash
- Disk-backed storage with in-memory shard cache

## How Solune works

At runtime, Solune starts three layers:

1. **Launcher (`main.go`)**
   - Builds `worker` and `supervisor` binaries.
   - Cleans up existing supervisor processes and port `9000`.
   - Starts the supervisor.

2. **Supervisor (`cmd/supervisor`)**
   - Starts the worker on the configured port.
   - Waits for worker exit and restarts it.
   - Handles shutdown signals and forwards them to worker.

3. **Worker (`cmd/worker`)**
   - Loads stores from `db/*.solstr`.
   - Starts TCP server and handles client commands.

## Persistence model

Solune is **disk-backed** (not memory-only).

- Each store is persisted in `db/<store>.solstr`.
- File format is line-based: `key,base64(value)`.
- On startup, stores are discovered from files in `db/`.
- `set`/`update` and `delete` currently rewrite store files via temp-file + rename.
- In-memory shards cache values for fast access during runtime.

## Command protocol

Send one JSON object per line over TCP.

### Allowed fields

- `instruction` (string, required)
- `store` (string, optional depending on instruction)
- `key` (string or number, optional depending on instruction)
- `data` (string, optional depending on instruction)

Unknown fields are rejected.

### Supported instructions

- `get`
- `set`
- `delete`

Legacy non-JSON command format is rejected.

## Command behavior

### `get`

- `{"instruction":"get"}`
  Lists all stores.
- `{"instruction":"get","store":"users"}`
  Returns all records in store.
- `{"instruction":"get","store":"users","key":1}`
  Returns one record by key.

### `set`

- `{"instruction":"set","store":"users"}`
  Creates store if it does not exist; if it already exists, returns error.
- `{"instruction":"set","store":"users","data":"{\"name\":\"John\"}"}`
  Auto-assigns next key and stores data.
- `{"instruction":"set","store":"users","key":1,"data":"{\"name\":\"Jane\"}"}`
  Upserts key `1` (create or update).

### `delete`

- `{"instruction":"delete","store":"users","key":1}`
  Deletes one key.
- `{"instruction":"delete","store":"users"}`
  Deletes whole store and its backing file.

## Response format

Responses are newline-delimited JSON objects.

- Success write/delete:
  ```json
  {"status":200}
  ```

- Read single:
  ```json
  {"key":1,"value":"{\"name\":\"Jane\"}"}
  ```

- Read many:
  - One JSON object per line.

- Error:
  ```json
  {"error":"..."}
  ```

- Empty result in some read paths can return:
  ```json
  {"status":404}
  ```

## Getting started

### Requirements

- Go 1.20+
- Docker (optional)

### Run locally (native)

```bash
go run main.go
```

This starts Solune on `127.0.0.1:9000`.

### Run with Docker

```bash
docker compose up --build
```

## Sending commands

You can use any TCP client (for example `nc`):

```bash
printf '{"instruction":"set","store":"user_data"}\n' | nc 127.0.0.1 9000
printf '{"instruction":"set","store":"user_data","data":"{\"name\":\"John Doe\",\"age\":30}"}\n' | nc 127.0.0.1 9000
printf '{"instruction":"get","store":"user_data"}\n' | nc 127.0.0.1 9000
```

There is also a helper script at `scripts/communication.py`.

Note: it ships with `command = ''`; set `command` to your JSON command before running.

## Testing

### Unit tests

```bash
go test ./test/unit/...
```

### Benchmarks

```bash
go test -bench=. -benchmem ./test/unit/store/...
```

### Integration tests

Start server first:

```bash
go run main.go
```

Then in another terminal:

```bash
go test ./test/integration/... -v -count=1
```

Run only unhappy-path integration tests:

```bash
go test ./test/integration/... -run TestIntegrationUnhappyPaths -v -count=1
```

Manual request walkthrough:

```bash
go run ./test/integration/steps.go
```

## Notes and current limitations

- No authentication/authorization.
- No WAL yet.
- No replication.
- Persistence is file-based with rewrite-on-update/delete semantics.
- API is TCP JSON-line protocol only (no HTTP/gRPC layer).


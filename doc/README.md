# Solune Roadmap

This document tracks planned and in-progress work for Solune.

For setup, protocol, architecture, and usage, see the root `README.md`.

## Current focus

### 1) Storage correctness and guarantees
- Define storage guarantees and durability expectations.
- Document key/value constraints and error semantics.
- Document invariants (key allocation, visibility, delete behavior).
- Add a short concurrency + restart behavior note.

### 2) Durability improvements
- Introduce WAL format with framing + checksums.
- Append `set` / `delete` mutations before acknowledgment.
- Add configurable sync policy.
- Add WAL replay with corrupt-tail handling.
- Add checkpointing to bound startup replay time.

### 3) Persistence engine evolution
- Replace rewrite-on-update/delete behavior with append-only segments.
- Maintain in-memory index to latest record location.
- Add tombstone handling for deletes.
- Add compaction to reclaim stale data.

### 4) Recovery hardening
- Add atomic key reservation API.
- Seed keys from authoritative recovered state.
- Detect and surface index/data corruption explicitly.
- Add startup consistency checks.

### 5) Test hardening
- Add crash-recovery boundary tests.
- Add WAL corruption/truncation tests.
- Add concurrency race/duplicate-key tests.
- Add durability mode test matrix.
- Add long-running stress + restart + compaction tests.

## Backlog

### Scalability
- Expand sharding strategy and balancing behavior.
- Explore replication and failover.

### Security
- Add authentication and authorization.
- Add encryption in transit and at rest.

### Tooling
- Add CLI/admin tooling for easier operations.
- Improve observability and operator diagnostics.

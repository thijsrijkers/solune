# 🗺️ Solune Project Roadmap

## ✅ Completed Milestones

### 1. Core Infrastructure
- [x] Go-based architecture implemented.
- [x] Modular design with clear separation (`data`, `store`, `tcp`).

### 2. Networking and Communication
- [x] Basic TCP networking module created.
- [x] Communication script (`communication.py`) for interfacing or testing.

### 3. Deployment and Environment
- [x] Dockerfile and docker-compose setup for containerized deployment.

---

## 🚧 Future Developments

### 1. Scalability and Performance
- [ ] Expand sharding module for distributed data storage.
- [ ] Expand sharding with data processing splitting.
- [ ] Add replication and fault tolerance features.

### 2. Data Durability
- [ ] Implement persistent storage to allow data recovery after downtime (currently in-memory only).
- [ ] Explore write-ahead logging or snapshotting mechanisms.

### 3. Data Import & Export
- [ ] Add support for exporting data (e.g., to JSON, CSV, or binary formats).
- [ ] Implement import functionality for bootstrapping or migration.
- [ ] Support selective export/import (e.g., per key range or shard).

### 4. Security and Access Control
- [ ] Implement authentication and authorization.
- [ ] Add encryption for data at rest and in transit.

### 5. User Interface and Tooling
- [ ] Create an administrative web dashboard.
- [ ] Develop CLI tools for interaction and management.

### 6. Documentation
- [ ] Expand documentation in the `doc` module.

---

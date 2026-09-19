# ADR-005 — Initial Server Topology and Service Boundaries

**Status:** ACCEPTED  
**Date:** 2026-09-19  
**Applies to:** WP-001, WP-006, WP-009 and initial deployment architecture

## Context

GRIDWORKS needs strong transactional authority and realtime/social capabilities but should not begin as an unnecessarily fragmented microservice estate.

## Decision

Launch engineering begins as a **modular service architecture with a small number of native Linux processes**.

Initial process set:

1. **gridworks-api**
   - HTTP API;
   - identity/session application layer;
   - player/company/world commands;
   - contracts;
   - market orchestration;
   - manager/recruitment;
   - business lifecycle;
   - ledger transaction orchestration.

2. **gridworks-worker**
   - scheduled/background work;
   - offline settlement batches;
   - notifications;
   - content publication tasks;
   - cleanup/maintenance;
   - economic analytics jobs.

3. **gridworks-realtime**
   - WebSocket sessions;
   - DM/Consortium/object-thread realtime fan-out;
   - notification delivery;
   - selected market/contract updates.

4. **gridworks-sim-validator**
   - Rust-based deterministic validation/challenge execution;
   - may begin embedded/library-driven where operationally simpler, but remains an explicit isolation boundary.

Infrastructure:

- PostgreSQL 18 — transactional authority;
- NATS/JetStream — event transport/durable streams where required;
- S3-compatible object storage — attachments, bundles, replays/assets;
- Valkey — **not deployed by default**; add only after a measured requirement;
- dedicated search cluster — **not deployed by default**; PostgreSQL first.

## Native runtime

All estate services run as pinned native Linux/systemd services.

No Docker/Podman/OCI runtime dependency is permitted for the target estate.

## Scaling

Scale vertically and by process replication before introducing additional service decomposition.

A domain becomes a separate service only when justified by at least one of:

- distinct scaling profile;
- security/isolation requirement;
- operational blast-radius requirement;
- ownership/team boundary;
- materially different availability requirement.

## Rationale

This keeps the first system understandable and operable while preserving clear logical domains and future split points.

## Consequences

- WP-001 creates native service scaffolding and systemd conventions.
- WP-009 begins with a modular Go API rather than dozens of independently deployed services.
- Market/social may split later only on evidence.

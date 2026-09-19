# ADR-003 — Schema, Serialization and Contract Versioning

**Status:** ACCEPTED  
**Date:** 2026-09-19  
**Applies to:** simulation boundary, APIs, events, persistence-adjacent DTOs, content manifests

## Context

GRIDWORKS spans Rust, Go, Godot, TypeScript and persistent data. Contract drift is a major architecture risk.

## Decision

Use **JSON as the canonical human-readable external/API representation**, backed by a repository-owned schema catalogue.

For the Rust↔Godot simulation boundary and durable replay/snapshot payloads, use a **compact binary representation generated from the same canonical schema definitions** where performance or payload size justifies it.

The canonical schema source lives under:

`packages/schemas/`

Schemas are versioned and code generation may produce:

- Rust types;
- Go types;
- TypeScript types;
- JSON Schema/OpenAPI artifacts where appropriate.

## Rules

Every externally meaningful payload must carry or be associated with:

- schema version;
- rules/content version where simulation-relevant;
- semantic IDs rather than translated prose.

Breaking schema changes require:

- new version;
- migration/compatibility policy;
- tests proving old supported clients fail safely or remain compatible.

## API

Mobile/public API starts with HTTPS + JSON.

OpenAPI should describe HTTP contracts.

## Events

NATS domain events use versioned envelopes containing at minimum:

- event ID;
- event type;
- event version;
- occurred-at timestamp;
- producer;
- correlation/causation IDs;
- related entity IDs;
- payload.

## Simulation snapshots/replays

Snapshots and replay evidence must be deterministic, checksummable and versioned.

Do not serialize raw in-memory pointers, engine objects or language-specific object graphs.

## Rationale

JSON keeps operational debugging straightforward while generated contracts reduce cross-language drift. Binary encoding is reserved for high-volume/local simulation boundaries rather than forced everywhere.

## Consequences

- WP-001 creates schema tooling skeleton.
- WP-002 defines simulation snapshot/command schemas.
- WP-009 generates/validates HTTP contracts.
- CI detects incompatible schema changes.

# ADR-001 — Godot ↔ Rust Integration Boundary

**Status:** ACCEPTED  
**Date:** 2026-09-19  
**Decision owners:** GRIDWORKS Design/Architecture  
**Applies to:** WP-002, WP-010, WP-011 and all later client simulation work

## Context

GRIDWORKS requires one canonical deterministic rules engine shared by the mobile client, server validator, challenge runner, test harness and balance tooling. The Godot client must not reimplement simulation logic in GDScript/C#.

## Decision

The canonical simulation engine is a Rust library exposed to Godot through a **narrow native extension boundary using Godot's GDExtension model**.

The integration shall expose coarse-grained, versioned operations rather than fine-grained object calls.

Examples:

- load rules/content version;
- load deterministic simulation snapshot;
- submit command batch;
- advance simulation by an authoritative elapsed interval;
- query summarized facility state;
- serialize snapshot/replay evidence;
- validate deterministic digest.

The Rust library shall:

- have no network access;
- have no PostgreSQL/SQLite access;
- have no dependency on device wall clock;
- accept authoritative timestamps/elapsed durations as input;
- use deterministic seeded randomness only;
- expose explicit schema/rules versions;
- remain testable without Godot.

Godot owns:

- presentation;
- input;
- animation;
- local UX state;
- visual interpolation;
- non-authoritative UI calculations.

## Serialization at the boundary

Boundary payloads use the schema contract defined by ADR-003. Avoid passing complex engine-owned pointers or mutable shared state across the FFI boundary.

## Rationale

GDExtension gives us:

- native performance;
- one Rust implementation;
- clean engine separation;
- reusable server-side simulation crate;
- lower drift risk than parallel GDScript/Go implementations.

## Rejected alternatives

### Reimplement simulation in GDScript
Rejected because client/server rule drift would be inevitable.

### C# as the canonical rules engine
Rejected because it weakens reuse with the accepted Rust server-validator/tooling direction.

### Full Rust Godot application
Rejected because Godot remains valuable as the presentation/editor layer; replacing too much of it increases implementation complexity without product benefit.

## Consequences

- WP-002 must create a Godot-independent Rust simulation crate.
- WP-010 must implement the extension bridge and deterministic integration tests.
- The FFI API must remain intentionally small.
- Engine upgrades must include a compatibility test for the extension ABI/API.

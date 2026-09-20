# GRIDWORKS — WP-002 Rust Simulation Kernel

**Status:** AUTHORIZED FOR ENGINEERING  
**Architecture owner:** Architecture  
**Implementation model:** GPT-5.6 Luna XHigh if the defined contracts converge cleanly; escalate only if genuine determinism/integration investigation is required  
**Review:** GPT-5.6 Sol-class Architecture review  
**Exact required parent:** `1648af5a28114268b0b5d1a027e64e57297829e0`

## 1. Purpose

Implement the first real canonical GRIDWORKS simulation kernel inside `crates/gridworks-sim`.

WP-002 must establish a deterministic, replayable, versioned simulation execution model suitable for later facility, failure, production, offline and challenge systems.

It must **not** implement the facility/component graph, failure model, economy, marketplace, property simulation or other later work packages.

## 2. Controlling references

Engineering must read and obey:

- `AGENTS.md`
- `AI_AGENT_COLLABORATION_PROTOCOL.md`
- `docs/design/GRIDWORKS_PRODUCT_PHILOSOPHY_AND_GAME_DESIGN.md`
- `docs/design/GRIDWORKS_DP2_CORE_SYSTEMS_SPECIFICATION.md`
- `docs/design/GRIDWORKS_DP3_ENGINEERING_ARCHITECTURE_AND_WORK_PACKAGES.md`
- `docs/decisions/ADR-001_GODOT_RUST_INTEGRATION.md`
- `docs/decisions/ADR-003_SCHEMA_AND_SERIALIZATION.md`
- `docs/decisions/ADR-004_CONTENT_RULES_DISTRIBUTION.md`
- `docs/decisions/ADR-008_WORLD_TIME_SEASONS_AND_TIMERS.md`
- WP-001 accepted foundation.

If an accepted architecture decision appears unimplementable, stop that decision path and raise the evidence in the WP-002 GitHub issue. Do not redesign architecture silently.

## 3. Exact ancestry

Engineering must branch from exactly:

`1648af5a28114268b0b5d1a027e64e57297829e0`

Before implementation, verify:
- `main` resolves to that SHA;
- working tree is clean;
- the branch is created from that exact commit.

If `main` has moved before branch creation, do **not** silently use the newer head. Report it in the issue and wait for Architecture direction.

## 4. Runtime / host constraints

WP-002 is repository-only.

ERIS prerequisites from WP-001 are accepted as sufficient unless Engineering proves a concrete missing build/test dependency.

No additional host package installation is authorized by default.

Still prohibited:
- Docker / Podman / Compose / OCI / Kubernetes;
- deployment;
- application systemd installation;
- database/service provisioning;
- nginx/DNS/firewall/WireGuard changes;
- external port exposure;
- shared-service changes.

If an additional native development prerequisite is genuinely required, post the requirement and evidence to the WP-002 issue before installing it.

## 5. Canonical simulation ownership

All simulation mathematics and state transition logic introduced by WP-002 must live in Rust under:

`crates/gridworks-sim/`

Go, Godot, TypeScript and database code must not duplicate simulation rules.

WP-002 may add only the minimal surrounding repository support required to test/serialize the Rust kernel.

## 6. Required kernel concepts

Implement a clean, versioned foundation for:

### 6.1 Simulation state

A canonical state container with at least:
- schema version;
- rules version;
- simulation/operational time;
- deterministic RNG state or seed lineage;
- kernel revision/state metadata;
- domain state payload boundary suitable for later typed expansion.

The payload may remain deliberately small in WP-002, but it must prove real deterministic state evolution.

### 6.2 Commands

A canonical command model with:
- command ID;
- command type;
- schema/rules version context as appropriate;
- deterministic payload;
- authoritative effective time supplied externally;
- idempotency/replay identity where relevant.

Do not use device/local wall clock inside command execution.

### 6.3 Events / results

Command execution must produce deterministic outcomes suitable for:
- accepted/rejected result;
- emitted domain event(s);
- resulting simulation state;
- deterministic digest/checksum evidence.

The kernel must not perform network/database side effects.

### 6.4 Time advancement

Support explicit deterministic simulation-time advancement.

The kernel must accept elapsed/effective time as input and calculate state transitions from that input.

It must not read:
- system clock;
- timezone;
- network time;
- filesystem timestamps.

This is the basis for future offline catch-up.

### 6.5 Seeded randomness

Where randomness is required for proof, use a deterministic seeded RNG.

Requirements:
- same initial state + same commands + same seed/rules version => same result;
- RNG behavior must not depend on platform entropy;
- RNG state/seed lineage must be serializable or reproducible;
- no `thread_rng()`, OS entropy or implicit random seed in authoritative execution.

The selected RNG algorithm/version must be explicit and pinned enough for replay compatibility.

## 7. Deterministic execution contract

Provide one canonical execution entry point conceptually equivalent to:

`execute(state, command) -> transition/result`

and one explicit time advancement entry point conceptually equivalent to:

`advance(state, authoritative_elapsed_time) -> transition/result`

Names and exact Rust types are Engineering-owned, but the semantics must remain clear and testable.

Execution must be:
- pure with respect to external I/O;
- deterministic;
- version-aware;
- replayable;
- reject invalid input rather than silently normalizing ambiguous state.

## 8. Snapshot / serialization contract

Implement versioned serialization/deserialization for kernel state sufficient to prove:

1. construct state;
2. serialize snapshot;
3. deserialize snapshot;
4. replay the same command sequence;
5. reach exactly the same final digest.

Use the accepted schema/serialization architecture.

Human-readable JSON is acceptable for the repository proof unless a compact format is already justified by an accepted ADR.

Do not introduce a second canonical schema source.

If schema definitions must evolve, update `packages/schemas` as the canonical contract source and keep generated/manual Rust representation aligned according to ADR-003.

## 9. Rules-version behavior

Rules version must be explicit.

The kernel must:
- reject missing/invalid versions;
- preserve rules version in state/snapshot/replay;
- detect incompatible command/state version combinations;
- avoid silently running a command under a different rules version.

WP-002 does not need multi-version migration logic beyond what is necessary to prove correct rejection/compatibility behavior.

## 10. Deterministic digest

Provide a deterministic state/replay digest used only as verification evidence.

Requirements:
- identical logical state produces identical digest;
- field ordering/serialization must not make digest nondeterministic;
- digest algorithm must be explicitly selected and tested;
- digest should be strong enough for replay evidence, not merely an incidental debug hash.

Prefer a stable cryptographic digest for replay evidence.

## 11. Error model

Introduce typed kernel errors for at least:
- malformed command;
- unsupported command type;
- version mismatch;
- invalid time movement;
- duplicate/replayed command where disallowed;
- invalid state;
- serialization/deserialization failure.

Do not panic on normal invalid external input.

## 12. Minimal proof domain

WP-002 needs a tiny synthetic/internal proof domain so state actually changes.

This proof domain must remain intentionally generic and must **not** implement WP-003 facility/component mechanics.

Example acceptable proof:
- a counter/register or generic quantity;
- deterministic increment/decrement command;
- deterministic seeded event selecting from a tiny fixed set;
- explicit time advancement updating a kernel tick/time field.

The proof domain exists solely to demonstrate the engine contract.

It must be unmistakably marked as kernel test/demo state, not game design.

## 13. Replay proof

Add an executable replay test/harness that proves:

- same initial snapshot;
- same ordered command stream;
- same seed;
- same authoritative elapsed times;
- same schema/rules version;

produces:
- identical events/results;
- identical final serialized state;
- identical final digest.

Also prove at least one negative:
- altered command order, seed, rules version or elapsed time produces a different/rejected outcome as appropriate.

## 14. Cross-platform determinism constraints

The kernel must avoid nondeterministic constructs that could diverge between supported targets.

Review and document:
- floating-point use;
- unordered map/set iteration;
- locale-sensitive formatting;
- platform-dependent integer width;
- filesystem ordering;
- concurrency/race-dependent outcomes.

For WP-002, prefer:
- fixed-width integers;
- explicit ordering;
- deterministic collections or sorted iteration;
- no authoritative floating-point math unless explicitly normalized and tested.

Do not over-engineer a generic deterministic math library yet; simply establish safe conventions and guardrails.

## 15. No I/O in core

The canonical kernel crate must not depend on:
- HTTP;
- NATS;
- PostgreSQL;
- SQLite;
- S3;
- filesystem persistence;
- system time;
- process environment;
- Godot.

Serialization of in-memory values is permitted.

Any future adapter boundary must live outside authoritative transition logic.

## 16. Godot / Go boundaries

WP-002 does **not** implement GDExtension integration; that remains WP-010.

WP-002 does **not** implement server simulation validation integration beyond any minimal compile/test fixture needed to preserve the boundary.

Do not add duplicate Go simulation logic.

## 17. Required tests

At minimum:

### Unit tests
- initial state validity;
- valid command transition;
- invalid command rejection;
- version mismatch rejection;
- backwards time rejection;
- duplicate/idempotency behavior;
- deterministic seeded behavior;
- stable digest.

### Replay tests
- snapshot round trip;
- command-stream replay;
- identical final state and digest;
- negative replay case.

### Property/invariant tests
Use property-based testing where it adds real value, especially for:
- serialization round trips;
- monotonic time;
- deterministic replay over generated valid command sequences.

Do not add property-testing complexity merely for decoration.

## 18. CI / validation

Extend existing WP-001 gates so CI verifies the WP-002 kernel.

Required:
- `cargo fmt --check`;
- `cargo test`;
- `cargo clippy ... -D warnings`;
- deterministic/replay tests;
- schema validation if contracts change;
- secret scan;
- forbidden-runtime scan.

All existing WP-001 gates must remain green.

## 19. Evidence

Create repository evidence for WP-002 sufficient to reproduce:

- toolchain/version;
- exact parent;
- test commands;
- replay proof;
- determinism assumptions;
- selected RNG algorithm/version;
- selected digest algorithm;
- any schema changes;
- any dependency additions;
- known limitations.

Keep evidence concise and technical.

## 20. Allowed paths

Primary:
- `crates/gridworks-sim/**`
- `tests/deterministic/**`
- `packages/schemas/**` where necessary
- root Rust workspace/lock/toolchain files where necessary
- CI/build files only as required for WP-002
- `docs/architecture/**`, `docs/evidence/**` for WP-002

Minimal supportive changes elsewhere are allowed only when directly necessary to compile/test the kernel and must be explained in the handback.

## 21. Out of scope

Do not implement:
- facility/system/component graph;
- bottleneck calculation;
- production recipes/inventory;
- causal failures;
- repairs/interventions;
- business/economic ledger;
- company gameplay;
- accounts/auth;
- market/trade;
- contracts;
- real estate;
- managers/skills;
- social/messaging;
- full offline catch-up policies;
- Godot GDExtension;
- live server validator integration;
- challenge sandbox;
- production deployment.

## 22. Done-when

WP-002 is complete only when Architecture can independently verify:

1. exact authorized ancestry;
2. canonical Rust-only transition logic;
3. explicit versioned state/command model;
4. deterministic external-time advancement;
5. deterministic seeded RNG;
6. versioned snapshot round-trip;
7. replay produces identical events/state/digest;
8. negative replay/version/time tests behave correctly;
9. no wall-clock/network/database/filesystem dependency in authoritative kernel code;
10. CI and all pre-existing gates pass;
11. no later-WP gameplay has leaked into the kernel;
12. complete GitHub handback exists.

## 23. Engineering handback

Post to the WP-002 GitHub issue and draft PR:

- branch;
- exact parent SHA;
- exact HEAD SHA;
- commit chain;
- files changed;
- dependency changes;
- commands/tests run;
- replay/determinism evidence;
- RNG/digest choices;
- schema changes;
- CI run links/results;
- deviations;
- unresolved risks;
- architecture conflicts;
- explicit statement that no host/deployment work occurred.

## 24. Stop condition

After posting the handback, Engineering stops.

Do not start WP-003 without separate Architecture/owner authorization.

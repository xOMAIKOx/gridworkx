# GRIDWORKS — WP-003 Facility / System / Component Dependency Graph

**Status:** AUTHORIZED FOR ENGINEERING  
**Architecture owner:** Architecture  
**Implementation model:** GPT-5.6 Luna XHigh  
**Review:** GPT-5.6 Sol-class Architecture review  
**Exact required parent:** `f2e925f3392399c36f203c47a2c70dfb3f94bbf5`

## 1. Purpose

Implement the canonical deterministic facility/system/component graph in the Rust simulation engine.

WP-003 establishes the structural model that later work packages will use for failures, diagnosis, production, maintenance and the opening aggregate-plant vertical slice.

It must implement:
- Facility → System → Component hierarchy;
- directed dependency relationships;
- dependency semantics;
- operational availability propagation;
- capacity/bottleneck calculation;
- deterministic graph validation and evaluation.

It must **not** implement causal failures, repair/intervention mechanics, recipes/inventory, economics, contracts or UI.

## 2. Controlling references

Engineering must read and obey:
- `AGENTS.md`
- `AI_AGENT_COLLABORATION_PROTOCOL.md`
- Product Philosophy
- DP2 Core Systems Specification
- DP3 Engineering Architecture
- ADR-001 Godot/Rust
- ADR-003 Schema/Serialization
- ADR-004 Content/Rules Distribution
- ADR-008 World Time
- accepted WP-002 deterministic kernel

If implementation conflicts with an accepted architecture decision, stop that path and raise evidence in the WP-003 issue.

## 3. Exact ancestry

Branch from exactly:

`f2e925f3392399c36f203c47a2c70dfb3f94bbf5`

Verify:
- exact parent;
- clean working tree;
- no unrelated changes.

If `main` moves later, do not silently rebase onto the newer head unless Architecture instructs it.

## 4. Runtime / host constraints

Repository-only.

No additional ERIS package installation is authorized by default.

Prohibited:
- Docker/Podman/Compose/OCI/Kubernetes;
- deployment;
- systemd installation;
- database/service provisioning;
- nginx/DNS/firewall/WireGuard changes;
- external port activation.

## 5. Canonical domain hierarchy

Implement explicit deterministic types for:

`World → Region → Site/Parcel → Facility → System → Component`

WP-003's simulation scope begins at Facility/System/Component. World/Region/Site identifiers may appear as ownership/context references without implementing full world simulation.

Required stable identities:
- facility ID;
- system ID;
- component ID.

IDs must be semantic/stable enough for snapshots, commands, diagnostics and later persistence references.

## 6. Facility

A Facility is a graph container representing one operational asset/business site.

At minimum it must support:
- identity;
- facility type/category reference;
- systems;
- dependency edges;
- current evaluated operational state;
- nominal capacity basis;
- deterministic validation.

Do not encode translated display text as authoritative state.

## 7. System

A System is a logical operational grouping within a Facility.

Examples for future content:
- feed handling;
- crushing;
- screening;
- utilities;
- water;
- electrical;
- storage.

WP-003 should implement the generic structural type, not quarry-specific game rules.

A System must:
- contain or reference components;
- expose deterministic evaluated state;
- contribute to facility availability/capacity.

## 8. Component

A Component is the smallest operational graph node for this work package.

At minimum:
- stable component ID;
- component type reference;
- nominal capacity;
- current availability/operability input;
- current capacity factor or effective-capacity input;
- enabled/disabled state as structurally required;
- deterministic metadata needed for graph evaluation.

WP-003 must not yet implement failure causes, wear curves, maintenance history or repair actions.

Those arrive in WP-004.

## 9. Dependency edge types

Implement the accepted edge classes:

### HARD
Downstream operation requires the upstream dependency to be operational.

If required upstream availability is zero/unavailable, the dependent path cannot operate.

### CAPACITY
Upstream capacity constrains downstream achievable throughput.

### QUALITY
Represents a dependency that can affect output quality/grade later.

WP-003 only needs the structural/evaluation hook; do not build a complete quality simulation.

### RELIABILITY
Represents a relationship that affects reliability later.

Structural support only in WP-003.

### COST
Represents a relationship that affects operating cost later.

Structural support only in WP-003.

### OPTIONAL
Represents a non-mandatory dependency/path.

Optional edges must not accidentally become HARD dependencies.

Edge semantics must be explicit in serialized contracts and tests.

## 10. Directed graph rules

The facility model is a directed graph.

Engineering must define and test:
- node existence validation;
- duplicate node/edge rejection;
- self-edge rejection unless Architecture explicitly approves a legitimate case;
- deterministic edge ordering;
- cycle handling.

Default architecture expectation: operational flow/dependency graphs for WP-003 should be acyclic unless a specific future system justifies feedback loops.

Therefore:
- detect dependency cycles;
- reject invalid cycles in the initial kernel;
- do not silently break/ignore cycles.

If Engineering finds a legitimate immediate need for cyclic graphs, stop and raise it to Architecture before changing this rule.

## 11. Operational state

Keep **condition** separate from **operability**.

WP-003 may model only the inputs required to evaluate whether a component/system/facility can currently operate.

Do not infer:
- good condition means operational;
- poor condition means failed.

The structural model must leave room for WP-004 to introduce component condition, symptoms and causal faults independently.

## 12. Capacity model

Implement deterministic nominal/effective capacity mechanics using safe deterministic arithmetic.

Core principle:

> Facility throughput emerges from bottlenecks; it is not a manually assigned final percentage.

At minimum support:
- component nominal capacity;
- local effective capacity factor;
- HARD dependency availability;
- CAPACITY dependency constraints;
- deterministic aggregation into system/facility achievable capacity.

Avoid floating-point authoritative math if practical.

Prefer fixed-point integer/basis-point style representation where percentages/factors are required.

Example:
- 10000 = 100%;
- 5000 = 50%;
- 250 = 2.5%.

Document the selected unit and overflow rules.

## 13. Bottleneck identification

The evaluator must be able to identify the limiting node/path contributing to achievable capacity.

Return deterministic bottleneck evidence sufficient for future diagnostics/UI.

At minimum:
- resulting facility/system effective capacity;
- one or more limiting component IDs;
- reason/class of limitation where determinable from WP-003 state.

Do not add advisor text or user-facing prose in the kernel.

## 14. Partial operation

The graph must support partial operation as a first-class outcome.

A facility may be:
- technically operational;
- severely constrained;
- producing at 2–5% or another low percentage.

Do not collapse all non-100% states into “failed.”

This is a critical product invariant.

## 15. Multiple paths / redundancy

Support deterministic evaluation of facilities with:
- parallel components;
- redundant paths;
- optional paths;
- alternate supply routes where structurally expressible.

A failure/unavailability on one path should not necessarily stop the entire facility if another valid path remains.

Keep the implementation bounded; do not create a generic industrial process solver beyond WP-003 needs.

## 16. Aggregate-plant proof fixture

Add a **test fixture only** representing the accepted future opening aggregate-processing facility:

`feed stockpile → feed conveyor/motor → crusher → screen → output conveyor → finished stockpile`

Purpose:
- prove hierarchy;
- prove HARD/CAPACITY dependencies;
- prove bottleneck calculation;
- prove partial-operation behavior;
- prove deterministic serialization/evaluation.

Do not implement the WP-011 tutorial, UI, economy, inventory settlement or WP-004 failure model.

The fixture may represent component availability/capacity inputs directly.

## 17. Required acceptance scenarios

Tests must include at least:

1. **Healthy line**
   - all required nodes available;
   - expected nominal/effective capacity;
   - deterministic bottleneck result.

2. **Hard dependency unavailable**
   - critical upstream required node unavailable;
   - downstream path becomes unavailable;
   - facility output reflects zero where appropriate.

3. **Partial capacity**
   - one component constrained to ~3%;
   - facility remains operational;
   - achievable capacity resolves to ~3%;
   - limiting component identified.

4. **Non-critical degraded/optional element**
   - a worn/degraded proxy input or optional node is constrained;
   - it does not incorrectly stop the line if not required by current graph semantics.

5. **Parallel/redundant path**
   - one branch unavailable;
   - valid alternate branch preserves bounded capacity.

6. **Cycle rejection**
   - invalid dependency cycle rejected deterministically.

7. **Duplicate/missing references**
   - rejected with typed errors.

8. **Serialization/replay**
   - facility graph survives snapshot round trip;
   - same graph/state produces same evaluated result/digest.

## 18. Graph evaluation determinism

Evaluation must not depend on:
- insertion order unless explicitly part of the contract;
- hash-map iteration order;
- locale;
- thread scheduling;
- system time;
- platform width;
- floating-point implementation variance.

Use:
- explicit stable ordering;
- fixed-width integer types;
- deterministic collections or sorting.

## 19. Commands and kernel integration

Extend WP-002 cleanly.

WP-003 may add commands needed to prove facility graph state transitions, for example:
- create/register proof facility;
- set component operational input;
- set component capacity factor;
- evaluate facility.

Exact command structure is Engineering-owned, but:
- commands remain versioned;
- replay/idempotency semantics remain intact;
- no out-of-band mutation of authoritative state;
- no direct database/network interaction.

Do not turn every internal calculation into a public command merely to satisfy tests.

## 20. Snapshot/schema integration

Update canonical schemas where required.

The snapshot must preserve all authoritative graph information required for deterministic replay.

Schema additions must:
- use canonical schema ownership under `packages/schemas`;
- remain strict;
- reject unknown/invalid graph state as appropriate;
- preserve stable IDs and edge types.

## 21. Error model

Add typed errors for at least:
- duplicate facility/system/component IDs;
- missing node references;
- invalid/self dependency;
- duplicate dependency;
- dependency cycle;
- invalid capacity/factor;
- arithmetic overflow;
- invalid graph state.

Normal bad input must not panic.

## 22. Performance posture

WP-003 does not need enormous-world optimization, but avoid obviously pathological algorithms in frequently evaluated paths.

Target fixture-scale and realistic facility-scale graphs cleanly.

Document complexity for:
- validation;
- cycle detection;
- evaluation.

Do not introduce a graph database or external graph engine.

## 23. No external I/O

Authoritative facility evaluation remains pure Rust.

No:
- PostgreSQL;
- SQLite;
- NATS;
- S3;
- HTTP;
- filesystem;
- wall clock;
- Godot dependency.

## 24. Tests / CI

Required:
- unit tests for graph primitives and typed errors;
- aggregate fixture tests;
- deterministic replay/snapshot tests;
- capacity/bottleneck tests;
- cycle/reference validation tests;
- existing WP-002 tests unchanged or intentionally extended;
- `cargo fmt --check`;
- `cargo test --workspace`;
- `cargo clippy --workspace --all-targets --all-features -- -D warnings`;
- existing repository validation/security/Godot gates remain green.

Property-based tests may be used where useful, but are not mandatory if deterministic fixture coverage is stronger.

## 25. Architecture constraint carried from WP-002

Do not expand `executed_commands` into an unbounded general production history.

WP-003 may preserve the current WP-002 proof behavior, but must not make facility graph size or history retention depend on every command ever executed.

If WP-003 requires a change to idempotency retention semantics, raise it to Architecture before altering the contract.

## 26. Allowed paths

Primary:
- `crates/gridworks-sim/**`
- `crates/gridworks-sim/tests/**`
- `tests/deterministic/**`
- `packages/schemas/**`
- `packages/content/**` only for generic facility/component type fixtures/config if necessary
- Rust workspace/lock files if required
- CI/build files only if required
- WP-003 architecture/evidence docs

Minimal supportive changes elsewhere must be explained.

## 27. Out of scope

Do not implement:
- causal failure generation;
- symptoms/diagnostic confidence;
- inspection actions;
- repair/rebuild/bypass mechanics;
- maintenance schedules;
- recipes;
- resource inventories;
- production settlement;
- economic ledger;
- marketplace/contracts;
- company/ownership gameplay;
- real estate gameplay;
- skills/managers effects;
- logistics economics;
- UI/Godot GDExtension;
- deployment.

## 28. Evidence

Create concise evidence documenting:
- graph model;
- edge semantics;
- capacity unit/fixed-point convention;
- cycle policy;
- evaluation algorithm;
- aggregate fixture;
- deterministic results;
- complexity notes;
- schema changes;
- dependencies.

## 29. Done-when

Architecture can verify:

1. exact parent;
2. explicit Facility/System/Component types;
3. all six dependency edge types represented;
4. invalid cycles/references rejected;
5. condition and operability remain conceptually separate;
6. deterministic HARD/CAPACITY evaluation works;
7. partial operation works;
8. bottleneck identification works;
9. alternate/parallel path behavior is tested;
10. aggregate fixture proves realistic graph behavior without WP-004/WP-011 scope creep;
11. snapshot/replay determinism remains intact;
12. no external I/O or duplicate simulation logic;
13. all CI/security gates pass;
14. complete GitHub handback exists.

## 30. Handback

Post to the WP-003 issue and draft PR:
- branch;
- exact parent/head;
- commit chain;
- files changed;
- dependencies;
- graph/evaluation design summary;
- capacity representation;
- commands/tests;
- aggregate fixture evidence;
- schema changes;
- CI links/results;
- deviations/risks;
- architecture conflicts;
- explicit no-host/deployment statement.

## 31. Stop condition

After handback, Engineering stops.

Do not start WP-004 without separate authorization.

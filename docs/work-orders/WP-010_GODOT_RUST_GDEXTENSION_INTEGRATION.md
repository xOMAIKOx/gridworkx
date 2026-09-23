# GRIDWORKS — WP-010 Godot ↔ Rust GDExtension Integration

**Status:** AUTHORIZED FOR ENGINEERING  
**Architecture owner:** Architecture  
**Implementation model:** Luna Max if integration investigation is required; Luna XHigh once the dependency/toolchain path is proven  
**Review:** GPT-5.6 Sol-class Architecture review  
**Exact required parent:** `03c0ad7971804897a74489f8f109756f98f6e219`

## 1. Purpose

Implement the accepted Godot ↔ Rust native simulation boundary using Godot GDExtension.

WP-010 proves that the existing Godot 4.x client can invoke the canonical `gridworks-sim` Rust engine through a narrow, versioned, deterministic native-extension boundary without duplicating simulation logic in GDScript, C#, Go or the bridge crate.

The result must be suitable as the integration foundation for WP-011, not a throwaway proof-of-concept.

## 2. Controlling references

Engineering must read and obey:

- `AGENTS.md`
- `AI_AGENT_COLLABORATION_PROTOCOL.md`
- Product Philosophy / Game Design
- DP2 Core Systems Specification
- DP3 Engineering Architecture & Work Packages
- `docs/decisions/ADR-001_GODOT_RUST_INTEGRATION.md`
- `docs/decisions/ADR-003_SCHEMA_AND_SERIALIZATION.md`
- accepted WP-002 Rust simulation kernel
- accepted WP-003/WP-004/WP-005 simulation-domain extensions
- accepted WP-006 persistence boundary
- accepted WP-009 API boundary
- Architecture handoff issue #18

ADR-001 is controlling:
- Godot integrates the canonical Rust simulation through GDExtension;
- the FFI is coarse-grained and versioned;
- no simulation mathematics is duplicated outside `gridworks-sim`;
- engine upgrades require compatibility evidence.

## 3. Exact ancestry

Branch from exactly:

`03c0ad7971804897a74489f8f109756f98f6e219`

This is the accepted WP-009 merge commit.

Verify:
- exact parent;
- clean working tree;
- no unrelated changes.

Do not silently rebase onto a later `main`.

The Architecture work-order commit on `main` is intentionally not the Engineering implementation parent.

## 4. Runtime / host boundary

Repository-only.

No:
- ERIS or other host mutation;
- live service deployment;
- systemd installation/start/restart;
- nginx/reverse proxy;
- DNS/firewall/WireGuard;
- database provisioning;
- NATS;
- provider/OIDC configuration;
- public port activation;
- Docker/Podman/Compose/Kubernetes/OCI/containerd/nerdctl.

Building and executing the native extension inside GitHub Actions / repository test runners is authorized.

## 5. Existing pinned client/toolchain baseline

Current repository baseline includes:
- Godot CI: **4.7.2 stable**, Linux x86_64;
- canonical Rust simulation crate: `crates/gridworks-sim`;
- repository Rust toolchain currently selected in CI;
- Godot project: `apps/game`.

Engineering must not silently change the repository-wide Rust or Godot versions.

If the maintained Rust GDExtension binding needed for Godot 4.7.2 cannot compile against the current Rust toolchain:
1. stop that path;
2. post exact compatibility evidence to the WP-010 issue;
3. request an Architecture decision before changing the repository-wide toolchain.

Do not work around a compatibility problem by hand-rolling Godot's C ABI.

## 6. Integration architecture

Create a dedicated wrapper crate, preferred:

`crates/gridworks-godot`

or a clearly equivalent name.

It should:
- build as `cdylib` for Godot;
- optionally also build as `rlib` for Rust-side testing;
- depend on `gridworks-sim`;
- depend on a pinned, maintained Godot Rust GDExtension binding;
- contain **integration/translation only**;
- contain no canonical simulation mathematics.

Do not convert `gridworks-sim` itself into the Godot extension crate.

## 7. One simulation authority

The following remain implemented only by `gridworks-sim`:

- deterministic advancement;
- command execution;
- component/facility calculations;
- dependency/bottleneck rules;
- fault effects;
- intervention outcomes;
- material transformations;
- inventory rules;
- RNG;
- snapshot/digest behavior;
- operational KPIs.

The bridge may:
- deserialize;
- validate boundary sizes/types;
- call canonical Rust APIs;
- map stable errors;
- serialize results.

The bridge may not recompute simulation outcomes.

## 8. Stateless FFI posture

Initial WP-010 boundary should be **stateless with respect to canonical simulation state**.

Preferred pattern:

`snapshot in → operation → result/snapshot out`

Do not hold authoritative `SimulationState` behind a Godot object pointer across calls.

Rationale:
- avoids lifetime/ownership hazards;
- keeps replay inputs explicit;
- simplifies deterministic parity tests;
- avoids hidden mutable shared state;
- preserves server-validator reuse.

The Godot wrapper object may hold non-authoritative integration metadata only.

## 9. Godot class

Expose one intentionally small Godot-visible class, preferred:

`GridworksSimBridge`

Use a `RefCounted`-style extension class or equivalent lightweight object.

Do not expose one Godot object per facility/system/component.

Do not mirror the Rust object graph into Godot classes.

## 10. Boundary representation

For WP-010, use **UTF-8 JSON** at the Godot↔Rust boundary.

This is consistent with ADR-003 and keeps the first integration inspectable.

Do not introduce a binary codec in this WP unless actual measured payload/performance evidence demonstrates a need and Architecture approves it.

Canonical snapshot and command structures remain governed by repository schemas and Rust serde contracts.

Godot must treat canonical simulation snapshots as boundary data, not reimplement their semantics.

## 11. Bridge contract version

Define an explicit bridge contract version, e.g.:

`godot-rust-bridge-0.1.0`

The bridge metadata response must expose at least:
- bridge contract version;
- simulation schema version;
- rules version;
- kernel revision;
- RNG algorithm identifier where useful.

Do not infer compatibility only from library filename.

## 12. Required coarse-grained operations

The GDExtension boundary must support the following equivalent operations.

Names may differ slightly if idiomatic, but the surface must remain small.

### Metadata

`bridge_info_json()`

Returns safe bridge/simulation version metadata.

### Create deterministic snapshot

`create_snapshot_json(seed)`

- positive deterministic seed;
- returns canonical simulation snapshot + digest.

### Advance by elapsed interval

`advance_json(snapshot_json, elapsed_ms)`

- authoritative elapsed interval supplied by caller;
- no wall-clock reads;
- returns transition/events + resulting snapshot + digest.

### Advance to authoritative time

`advance_to_json(snapshot_json, authoritative_time_ms)`

- canonical monotonic-time validation remains in Rust simulation.

### Execute command batch

`execute_commands_json(snapshot_json, commands_json)`

- commands are versioned canonical simulation commands;
- coarse-grained batch crossing;
- commands execute in deterministic order;
- if any command fails, return an error rather than a partially successful authoritative result;
- returned result includes final snapshot/digest and emitted events/evidence.

### Facility summary/evaluation

`evaluate_facility_json(snapshot_json, facility_id)`

- calls the canonical Rust facility/failure evaluation path;
- returns a presentation-safe serialized evaluation;
- no GDScript calculation of effective capacity/bottlenecks.

### Digest validation

Provide a small operation that can:
- calculate the canonical digest from a snapshot; and/or
- validate an expected digest.

Do not duplicate the digest algorithm in Godot.

## 13. Dynamic rules loading

Do not fake a runtime rules/content loader if the current canonical simulation crate does not yet support one.

WP-010 must expose current rules/schema metadata.

Actual signed rules/content bundle loading remains governed by ADR-004 and later authorized work.

## 14. Response envelope

All bridge calls should return a stable versioned envelope rather than leaking Rust/Godot exceptions or debug strings.

Success equivalent:

```json
{
  "bridge_version": "godot-rust-bridge-0.1.0",
  "schema_version": "schema-0.1.0",
  "rules_version": "rules-0.1.0",
  "ok": true,
  "result": {}
}
```

Error equivalent:

```json
{
  "bridge_version": "godot-rust-bridge-0.1.0",
  "schema_version": "schema-0.1.0",
  "rules_version": "rules-0.1.0",
  "ok": false,
  "error": {
    "code": "simulation.version_mismatch",
    "message": "snapshot or command version is not supported"
  }
}
```

Exact schema may differ, but it must be repository-owned and strict.

## 15. Stable error mapping

Map canonical Rust errors to stable boundary codes.

At minimum distinguish:
- malformed JSON;
- payload too large;
- invalid seed;
- schema/rules/kernel mismatch;
- invalid snapshot/state;
- malformed/unsupported command;
- duplicate command/idempotency;
- invalid/backward time movement;
- arithmetic/validation failure;
- unknown facility;
- digest mismatch;
- internal bridge failure.

Do not return raw Rust `Debug` output as the public Godot-facing error contract.

## 16. Panic safety

No panic may unwind uncontrolled across the GDExtension boundary.

Requirements:
- exported bridge methods must convert expected failures to the stable envelope;
- unexpected panic paths must fail safely and not return partial authoritative state;
- Godot headless tests must prove malformed input cannot crash the process.

Avoid `unwrap()` / `expect()` in normal exported bridge paths.

## 17. No hidden wall clock

The bridge and canonical simulation must not use:
- device wall clock;
- Godot frame time as authoritative simulation time;
- OS time as authoritative input.

Elapsed/time values must be explicit method inputs.

GDScript presentation timing may interpolate visuals but cannot mutate canonical simulation time directly.

## 18. Deterministic randomness

Do not use Godot RNG for canonical simulation decisions.

All canonical randomness remains the existing seeded Rust RNG serialized in simulation state.

Same:
- snapshot;
- seed/state;
- command sequence;
- authoritative elapsed input

must produce identical:
- events;
- resulting snapshot;
- digest.

## 19. Payload size bounds

Protect the native boundary from malformed/accidental huge local payloads.

Define documented bounds for:
- snapshot payload;
- command-batch payload;
- facility ID / metadata strings.

Reasonable initial defaults are acceptable, e.g. multi-MiB snapshot and sub-MiB command batch, but Engineering must record exact limits.

Oversize input returns a stable bridge error; it must not allocate unboundedly.

## 20. Schema catalogue

Add repository-owned schemas as required, preferably including:

- bridge response/envelope;
- command-batch input;
- transition/result;
- bridge metadata.

Reuse existing:
- `simulation-command.schema.json`;
- `simulation-state.schema.json`;
- facility/failure/material schemas.

Do not copy entire simulation schemas into the Godot project.

## 21. Cross-language contract drift

CI must verify:
- bridge schema files parse;
- bridge version is explicit;
- exported method inventory matches the documented bridge contract;
- Godot smoke tests expect the same API surface;
- Rust bridge tests use the same fixtures/contracts.

Do not maintain unrelated duplicate contract inventories.

## 22. Godot project integration

Add the `.gdextension` descriptor under the Godot project, e.g.:

`apps/game/native/gridworks_sim.gdextension`

Generated native libraries must live in a predictable ignored path, e.g.:

`apps/game/bin/<platform>/<profile>/`

Do not commit compiled `.so`, `.dll`, `.dylib` or mobile native binaries.

## 23. Linux CI proof target

WP-010 must prove the extension on the repository CI host:

- Linux;
- x86_64;
- pinned Godot 4.7.2 stable.

The CI build should:
1. build the Rust extension;
2. stage/copy the generated `.so` into the Godot project ignored binary path;
3. load the extension through the actual `.gdextension` descriptor;
4. execute deterministic headless integration tests.

This is repository validation, not host deployment.

## 24. Mobile packaging scope

Android/iOS production packaging is **not** required in WP-010.

However, the directory/descriptor/build structure must not make future mobile targets impossible.

Document:
- expected future Android/iOS native artifact slots;
- that each target will require separate Rust target builds and Godot export integration.

Do not download SDK/NDK/Xcode tooling in this WP.

## 25. GDExtension dependency

Pin the Rust Godot/GDExtension dependency through Cargo lockfiles.

Record:
- package/version;
- why it is compatible with Godot 4.7.2;
- required minimum Rust version.

If compatibility is uncertain, prove it in the branch before broad implementation.

Do not use an unmaintained abandoned binding.

## 26. Rust workspace

Add the bridge crate to the root Cargo workspace.

All existing:
- Rust tests;
- rustfmt;
- clippy/validation gates

must remain green.

No canonical sim public API should be weakened merely to make GDExtension bindings easier.

If a small pure serialization/helper API is required in `gridworks-sim`, keep it generic and Godot-independent.

## 27. GDScript wrapper

Add a thin presentation-side wrapper, preferred:

`apps/game/scripts/sim_bridge.gd`

Responsibilities may include:
- instantiate/check GDExtension class;
- invoke coarse operations;
- parse the stable bridge envelope;
- present typed Godot-side integration errors;
- expose bridge availability/version to UI code.

It must not:
- calculate simulation results;
- mutate snapshot JSON directly to fake results;
- implement RNG;
- implement bottleneck/failure/material rules;
- silently fall back to a GDScript simulation.

## 28. No silent fallback

If the native extension:
- is missing;
- fails ABI/API compatibility;
- rejects contract versions;

the client integration layer must fail explicitly.

Do not fall back to alternate local simulation mathematics.

Presentation-only fallback screens/messages are fine; authoritative simulation fallback is not.

## 29. Godot headless test

Add a dedicated headless Godot test script/scene that loads the real native extension.

It must verify at least:
- GDExtension class is registered;
- required method inventory exists;
- bridge metadata versions match expected repository constants;
- deterministic snapshot creation works;
- advance works;
- command batch works;
- facility evaluation works on a canonical fixture;
- digest returned by bridge matches expected Rust/golden digest;
- malformed input returns stable error rather than crash;
- version mismatch returns stable error;
- invalid/backward time returns stable error.

The test process must return non-zero on failure.

## 30. Deterministic golden fixtures

Create one small shared deterministic fixture set under a repository test path, e.g.:

`tests/deterministic/fixtures/wp010/`

Include only stable text fixtures, not generated native binaries.

At minimum:
- initial snapshot/seed description;
- command batch;
- expected event/result digest;
- expected final snapshot digest.

The same fixtures must be consumed by:
- a Rust integration test; and
- the headless Godot/GDExtension test.

This proves parity across the actual FFI path.

## 31. Cross-boundary determinism test

Required proof:

For the exact same fixture:
1. pure Rust path executes canonical engine;
2. Godot calls the GDExtension;
3. compare resulting digest;
4. compare relevant emitted event/result semantics;
5. compare canonical resulting snapshot.

Differences are a failure.

Do not generate the expected digest independently in GDScript.

## 32. Repetition / state leakage test

Run the same bridge operation repeatedly from fresh identical inputs.

Results must remain identical.

No prior Godot call may influence a later stateless simulation result.

This catches hidden static/global mutable state.

## 33. Command-batch atomicity

For a command batch:
- deserialize/validate all boundary structure;
- execute against a working copy;
- if command N fails, return an error;
- do not emit a success result containing a partially advanced canonical snapshot.

Add a regression where a later command in the batch fails.

## 34. Snapshot compatibility

Bridge snapshot loading must reject:
- unknown fields where canonical Rust serde already rejects them;
- unsupported schema version;
- unsupported rules version;
- unsupported kernel revision;
- invalid deterministic RNG state;
- invalid facility/failure/material state.

Do not auto-correct incompatible snapshots silently.

## 35. Digest behavior

Canonical digest remains implemented in `gridworks-sim`.

Bridge only calls that API.

Tests must prove:
- pure Rust digest == GDExtension digest;
- altered snapshot/command produces different digest where expected;
- a caller-supplied wrong expected digest returns a stable mismatch error.

## 36. Facility evaluation projection

If the bridge exposes facility summary/evaluation:
- serialize a bounded presentation-safe DTO;
- do not return raw Rust/Godot pointers;
- use semantic IDs;
- include enough information for WP-011 presentation to show effective capacity / operational state / bottleneck evidence without reproducing calculations.

Do not implement WP-011 UI here.

## 37. Godot presentation shell

Preserve the existing navigation/presentation shell.

WP-010 may:
- add bridge wrapper;
- add integration-test scene/script;
- optionally expose safe bridge version/health status.

Do not build the aggregate-plant gameplay UI in WP-010.

That is WP-011.

## 38. No API/network coupling

The Rust GDExtension must not:
- call the Go API;
- open sockets;
- use HTTP;
- read PostgreSQL/SQLite;
- publish NATS.

Godot client networking and offline sync remain separate boundaries.

WP-010 proves local simulation integration only.

## 39. No persistence coupling

The bridge must not read/write:
- player save files;
- SQLite;
- PostgreSQL;
- arbitrary filesystem state

as part of canonical simulation execution.

Snapshots are explicit input/output values.

Storage is owned by future client persistence/sync logic.

## 40. Threading

Initial bridge may run synchronously on the calling Godot thread.

Do not introduce worker threads merely for architecture.

If Engineering introduces threading:
- simulation calls must remain deterministic;
- no shared mutable state;
- Godot API thread-safety rules must be respected;
- tests must prove shutdown/lifetime safety.

A simple synchronous WP-010 is preferred.

## 41. Memory ownership

No raw pointer ownership protocol across FFI.

Use Godot/Rust binding-supported value types and copied/owned payloads.

No manual `malloc/free` API exposed to GDScript.

## 42. Native artifact hygiene

Add generated native-library paths to `.gitignore`.

CI must prove no compiled native binary is committed.

Do not commit:
- `target/`;
- `.so`;
- `.dll`;
- `.dylib`;
- Godot import caches.

## 43. Build helper

A repository script such as:

`ops/scripts/build-godot-extension.sh`

may:
- verify tool versions;
- build the bridge crate;
- copy the native library to the Godot ignored path.

Requirements:
- POSIX shell;
- `set -eu`;
- no sudo;
- no package installation;
- no host-service mutation;
- no container invocation;
- predictable paths;
- non-zero on missing artifact/tool.

## 44. CI

Extend the existing Godot/Rust CI without creating an unrelated pipeline.

Required sequence:
- repository validation;
- Rust workspace test/format/clippy as existing policy requires;
- build `gridworks-godot`;
- stage native library;
- launch pinned Godot 4.7.2 headlessly;
- run bridge integration test;
- ensure normal project parse still succeeds.

Do not disable existing gates to make GDExtension pass.

## 45. Security / robustness

Required:
- bounded input;
- strict JSON decode;
- stable error envelope;
- panic-safe FFI boundary;
- no secret/auth data in bridge;
- no network/database access;
- no arbitrary file loading;
- no unsafe raw-pointer API exposed to Godot.

If `unsafe` is required by the Rust GDExtension trait/macros, isolate/document it and do not expand unsafe surface unnecessarily.

## 46. Dependency audit

Evidence must list all new Rust dependencies.

Do not add:
- async runtime;
- HTTP client;
- database driver;
- logging stack;
- binary serialization framework

unless directly required by the bridge and explicitly justified.

Keep the integration crate small.

## 47. Performance posture

WP-010 is correctness-first.

Record a simple non-gating smoke measurement for:
- snapshot decode/encode;
- one command batch;
- one facility evaluation

on CI/dev hardware if convenient.

Do not introduce caching/global mutable state merely to improve a microbenchmark.

Binary serialization/performance optimization is later if measured JSON boundary cost is unacceptable.

## 48. Tests — Rust bridge unit/integration

Required:
- bridge metadata versions;
- valid snapshot creation;
- malformed JSON;
- unsupported versions;
- advance/advance_to;
- deterministic batch execution;
- later-command batch failure is atomic;
- duplicate command;
- wrong digest;
- facility evaluation;
- payload size rejection;
- panic/error mapping behavior where testable;
- repeated identical input gives identical output.

## 49. Tests — Godot headless

Required:
- extension loads under pinned Godot 4.7.2;
- class registration;
- method inventory;
- bridge metadata;
- shared deterministic fixture;
- digest equality;
- malformed/error envelope;
- no GDScript simulation fallback;
- clean non-zero failure behavior.

## 50. Tests — compatibility

Add a compatibility assertion that fails if:
- bridge API method inventory changes unexpectedly;
- bridge version changes without fixture/schema update;
- Godot 4.7.2 cannot load the compiled extension;
- schema/rules/kernel metadata mismatch.

Breaking bridge changes later require an explicit bridge version bump.

## 51. Structural tests

Add repository/static checks for:
- bridge crate is in workspace;
- crate depends on `gridworks-sim`;
- `.gdextension` descriptor exists;
- no native binary committed;
- Godot wrapper contains no obvious duplicated simulation constants/formulas;
- bridge crate contains no network/database dependencies;
- generated artifact path is gitignored;
- test fixture paths exist.

Avoid brittle source-text checks where executable tests are stronger.

## 52. Allowed paths

Primary:
- `crates/gridworks-godot/**`
- bounded generic helper edits in `crates/gridworks-sim/**` if justified
- `Cargo.toml`
- `Cargo.lock`
- `apps/game/**`
- `packages/schemas/**`
- `tests/deterministic/**`
- `tests/structural/**`
- `ops/scripts/**`
- `.gitignore`
- `.github/workflows/**`
- `Makefile` / repository validation scripts as required
- `docs/evidence/**`
- `docs/architecture/**` only for bridge-contract documentation if needed.

Do not modify WP-006/007/008 migrations or WP-009 API behavior.

## 53. Explicit exclusions

Do not implement:
- WP-011 aggregate-plant gameplay vertical slice;
- new simulation domain mathematics;
- client API authentication/network sync;
- SQLite offline persistence;
- realtime/WebSocket;
- inventory/economy settlement API additions;
- manager/recruitment/market/contracts;
- mobile Android/iOS export toolchains;
- app-store packaging;
- production deployment.

## 54. Evidence

Create:

`docs/evidence/WP-010_GODOT_RUST_GDEXTENSION.md`

Document:
- exact GDExtension binding + version;
- Godot compatibility;
- Rust minimum/toolchain compatibility;
- bridge crate/package layout;
- exported Godot class/method inventory;
- bridge contract/schema version;
- payload representation and bounds;
- error mapping;
- panic-safety approach;
- deterministic golden fixture;
- pure Rust ↔ Godot parity result;
- CI build/staging path;
- native artifact ignore policy;
- test commands/results;
- mobile-packaging deferred scope;
- deviations/risks/blockers;
- explicit no-host/deployment statement.

## 55. Done-when

Architecture can verify:

1. exact parent preserved;
2. dedicated Rust GDExtension wrapper crate exists;
3. `gridworks-sim` remains canonical and Godot-independent;
4. no simulation math is duplicated in bridge/GDScript;
5. one small Godot-visible bridge class exists;
6. FFI is coarse-grained and versioned;
7. snapshot/command boundary is schema-backed;
8. canonical state is explicit snapshot input/output, not hidden mutable FFI state;
9. bridge uses no network/database/wall-clock authority;
10. deterministic Rust RNG remains authoritative;
11. stable success/error envelope exists;
12. malformed/incompatible input fails safely;
13. command-batch failure is atomic;
14. digest is produced only by canonical Rust simulation;
15. Godot 4.7.2 loads the real compiled extension headlessly;
16. shared fixture proves pure Rust == GDExtension state/events/digest;
17. repeated calls prove no hidden cross-call state leakage;
18. compiled native artifacts are not committed;
19. CI builds/stages/tests the native extension;
20. existing Rust/Godot/repository/security gates remain green;
21. no WP-011/mobile packaging/deployment scope creep;
22. complete GitHub handback exists.

## 56. Handback

Post to the WP-010 issue and draft PR:

- branch;
- exact parent/head;
- commit chain;
- changed files/dependencies;
- pinned Godot Rust binding/version;
- Rust toolchain compatibility result;
- bridge crate and class names;
- exact method inventory;
- bridge/schema/rules/kernel versions;
- payload limits;
- error mapping table;
- panic-safety approach;
- deterministic fixture path;
- pure Rust expected digest;
- Godot bridge digest;
- parity result;
- headless Godot command/result;
- CI links;
- native artifact hygiene;
- mobile-packaging deferment;
- deviations/risks/conflicts;
- explicit no-host/deployment statement.

Keep the PR draft.

## 57. Stop condition

After handback, Engineering stops.

Do not start WP-011 without separate owner authorization.

Do not merge the WP-010 PR.

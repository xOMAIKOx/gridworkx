# WP-010 Godot ↔ Rust GDExtension Evidence

## Boundary and authority

`crates/gridworks-godot` is a dedicated `cdylib`/`rlib` wrapper over `gridworks-sim`. It exposes one `RefCounted` Godot class, `GridworksSimBridge`, with coarse stateless UTF-8 JSON operations: metadata, snapshot creation, explicit elapsed advance, command batch, digest validation and facility summary.

The bridge contains translation, bounds and stable error-envelope logic only. Simulation transitions, RNG, snapshot serialization, digest calculation and facility evaluation remain in `gridworks-sim`. No network, database, filesystem persistence, device wall clock, Godot RNG or raw pointer protocol is used.

## Toolchain and binding

- Rust toolchain pin: `1.80.0`.
- Godot editor: pinned CI Godot `4.7.2`.
- Binding: exact `godot = "=0.2.4"`.
- Bridge contract: `godot-rust-bridge-0.1.0`.
- Generated libraries are staged under ignored `apps/game/bin/`; no compiled native artifact is committed.

## Boundary safety

Snapshot input is bounded to 4 MiB and command batches to 1 MiB. Malformed JSON, unsupported versions, invalid seeds/time, duplicate commands and simulation errors return stable JSON error envelopes. Batch execution uses a working state and returns no successful partial result after a later command failure.

## Deterministic parity

`tests/deterministic/fixtures/wp010/fixture.json` is consumed by the Rust `wp010_fixture` integration test, which writes an ignored canonical result artifact, and by the real headless Godot test. The Godot test loads the actual `.gdextension`, verifies class/method/metadata registration, runs snapshot creation/advance/batch/digest/facility operations, compares snapshot/digest/events against the Rust artifact, repeats the operation for state-leakage detection and verifies malformed input rejection.

## Build and runtime proof

`ops/scripts/build-godot-extension.sh` builds the bridge, generates the pure-Rust fixture result, stages the `.so` and fixture inputs into the ignored Godot binary path, and rejects tracked native artifacts. CI then parses the Godot project and runs `wp010_headless_test.gd` with Godot 4.7.2.

Future Android/iOS artifact slots remain documentation/build-target concerns; no SDK/NDK/Xcode tooling or packaging was added.

## Scope

No WP-011 gameplay/UI, API/network sync, persistence, realtime, marketplace, contracts, managers, mobile packaging, deployment, host mutation, service activation, ports or container runtime work was performed.

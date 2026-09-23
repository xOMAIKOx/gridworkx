# WP-010 Godot ↔ Rust GDExtension Evidence

## Boundary and authority

`crates/gridworks-godot` is a dedicated `cdylib` wrapper over `gridworks-sim`. It exposes one `RefCounted` Godot class, `GridworksSimBridge`, with a coarse stateless UTF-8 JSON boundary. All canonical transitions, RNG, snapshot serialization, digest calculation and facility evaluation remain in `gridworks-sim`.

No network, database, filesystem persistence, device wall clock, Godot RNG, raw pointer protocol or simulation reimplementation is used.

## Toolchain and binding

- Rust toolchain: exact repository pin `1.80.0`.
- Godot editor/runtime: pinned CI `4.7.2`.
- Binding: exact `godot = "=0.2.4"`.
- Bridge contract: `godot-rust-bridge-0.1.0`.
- Native artifacts are staged below ignored `apps/game/bin/`; no compiled library is committed.

## Exported bridge methods

`GridworksSimBridge` exports exactly:

- `bridge_metadata()`
- `create_snapshot(seed)`
- `advance_snapshot(snapshot_json, elapsed_ms)`
- `advance_to_snapshot(snapshot_json, authoritative_time_ms)`
- `execute_command_batch(snapshot_json, command_batch_json)`
- `digest_snapshot(snapshot_json)`
- `validate_digest(snapshot_json, expected_digest)`
- `evaluate_facility(snapshot_json, facility_id)`

The facility method evaluates only the supplied canonical snapshot. Facility identifiers are bounded to 1–128 UTF-8 bytes. Snapshot payloads are bounded to 4 MiB and command batches to 1 MiB.

## Versioned envelope and errors

Every native success and error envelope includes:

- `bridge_version`
- `schema_version`
- `rules_version`
- `ok`
- `operation`
- `result` or `error`

The schema is `packages/schemas/godot-bridge.schema.json`.

Stable boundary mappings include:

| Condition | Code |
| --- | --- |
| malformed/unknown-field JSON | `bridge.invalid_json` |
| oversized snapshot or batch | `bridge.payload_too_large` |
| non-positive seed | `bridge.invalid_seed` |
| invalid RNG state | `bridge.invalid_rng` |
| unsupported schema/rules | `bridge.version_mismatch` |
| unsupported kernel/state validation | `bridge.invalid_state` |
| malformed command | `bridge.invalid_command` |
| unsupported command type | `bridge.unsupported_command` |
| duplicate command/idempotency key | `bridge.duplicate_command` |
| backward/invalid time | `bridge.invalid_time` |
| arithmetic overflow | `bridge.arithmetic` |
| unknown facility | `bridge.unknown_facility` |
| digest mismatch | `bridge.digest_mismatch` |
| unexpected panic/encoding failure | `bridge.internal` |

Rust diagnostic strings and `Debug` representations are not exposed.

## Panic and atomicity safety

Every exported bridge method executes through a narrow `catch_unwind(AssertUnwindSafe(...))` guard. Unexpected panics become `bridge.internal` envelopes and cannot unwind across the GDExtension boundary.

Command batches execute against a working `SimulationState`; a later command failure returns only an error envelope and no partial snapshot. Duplicate command identity and idempotency are rejected by the canonical kernel.

## Presentation wrapper

`apps/game/scripts/sim_bridge.gd` is the thin Godot-side `GridworksSimBridgeClient` wrapper. It instantiates/checks the native class, exposes all coarse methods, validates versioned envelopes and returns typed availability/method/envelope errors. It contains no simulation math, RNG or fallback state logic. The headless parity scene exercises this wrapper rather than calling the native object directly.

## Committed deterministic golden

The shared fixture is:

`tests/deterministic/fixtures/wp010/fixture.json`

It pins seed `1234`, elapsed time `5000`, authoritative time `7500`, a non-empty adjust-register plus seeded-pulse command batch, event semantics and the accepted canonical `after_elapsed_digest`, `after_command_digest` and `final_digest` values directly in the committed JSON fixture. The Rust test asserts those exact values; they are intentionally not duplicated in prose or evidence logs.

The pure Rust fixture test asserts these committed values and writes ignored parity artifacts. Godot consumes the same committed fixture and compares snapshots, events and digests against those Rust artifacts and pinned values.

## Discovery and runtime proof

`ops/scripts/build-godot-extension.sh` stages the real Linux `.so` and fixture artifacts but does not manufacture Godot project data. CI performs the editor discovery scan with a virtual display and Godot 4.7.2 compatibility rendering mode so the process exits cleanly; any non-zero discovery status fails the gate. CI performs the authorized discovery sequence:

```sh
timeout 30s xvfb-run -a godot --editor --path apps/game --quit-after 2 --rendering-method gl_compatibility --rendering-driver opengl3 --audio-driver Dummy
test -f apps/game/.godot/extension_list.cfg
grep -cFx 'res://native/gridworks_sim.gdextension' apps/game/.godot/extension_list.cfg
```

The generated list is asserted to contain exactly the WP-010 descriptor. A fresh normal runtime process then executes `res://scenes/wp010_test.tscn`, proving ClassDB registration, wrapper availability, every exported method, golden parity, negative/error matrix and non-crashing malformed input behavior.

## Verification scope

The final CI handback records the full repository foundation gates plus the editor discovery scan, normal headless parse and real deterministic parity scene. Mobile packaging remains deferred to the relevant later work package; no mobile SDK/NDK/Xcode work is part of WP-010.

No WP-011 gameplay/UI, API/network sync, persistence, realtime, marketplace, contracts, managers, mobile packaging, deployment, host mutation, service activation, ports or container-runtime work was performed.

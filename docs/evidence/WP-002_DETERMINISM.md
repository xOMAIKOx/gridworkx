# WP-002 Determinism Evidence

## Contract

The canonical entry points are `execute(&mut SimulationState, &Command)` and `advance(&mut SimulationState, authoritative_elapsed_ms)`. `advance_to` provides explicit server-time ordering checks. Neither entry point reads system time, timezone, filesystem state, process environment or external services.

The state carries schema version, rules version, kernel revision, operational time, RNG algorithm/seed/state/draw count, generic proof state and ordered command idempotency identities. Commands carry command ID, command type, schema/rules context, externally supplied effective time, idempotency key and deterministic payload.

## RNG

- Algorithm: `xorshift64star-v1`.
- Seed: explicit non-zero `u64` supplied to `SimulationState::new`.
- State and draw count: serialized in every snapshot.
- Entropy: none; no OS, thread or implicit seed source.
- Selection: deterministic rejection-sampling index selection over fixed proof payload options.

## Digest and serialization

- Snapshot representation: serde JSON with explicit `simulation.snapshot` envelope and aligned schemas under `packages/schemas/`.
- Digest: SHA-256 over the canonical serde JSON bytes of the versioned snapshot envelope, rendered as lowercase hexadecimal.
- State fields use fixed-width integer types and ordered vectors. No unordered maps, floating-point values, locale formatting, filesystem ordering or concurrency are used by the authoritative kernel.

## Proof coverage

The unit and integration tests prove:

- valid register transition and explicit time advancement;
- malformed, unsupported, version-mismatched, duplicate and invalid-time rejection;
- deterministic seeded pulse results;
- snapshot JSON round-trip;
- same snapshot/command stream/seed/time/version produces identical events, state and digest;
- changed command order or seed changes replay evidence;
- monotonic generated elapsed-time sequence.

Run with:

```text
cargo fmt --all -- --check
cargo test --workspace
cargo clippy --workspace --all-targets --all-features -- -D warnings
```

The proof domain is generic kernel test/demo state only. It does not represent facilities, components, failures, production, inventory, economy or any later gameplay work package.

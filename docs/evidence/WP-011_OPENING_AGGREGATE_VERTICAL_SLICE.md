# WP-011 Opening Aggregate-Plant Vertical Slice Evidence

## Scope and authority

WP-011 is a repository-only opening recovery slice built on the accepted WP-003/004/005 Rust foundations and WP-010 `GridworksSimBridgeClient`. Rust remains the sole canonical state/transition authority; Godot presents and dispatches commands only.

No economy, credits, revenue, skills, managers, contracts, network, persistence, database, deployment, host or container scope was added.

## Canonical scenario

- Scenario ID: `scenario.opening.aggregate`
- Seed: `1234` in the committed golden
- Rust builder: `opening_aggregate_scenario(seed)`
- Facility: `facility.aggregate_plant_fixture`
- Critical fault: `fault.motor_bearing_seizure` on `component.feed_conveyor`, instance `fault.instance.opening.feed_motor`
- Visible non-critical fault: `fault.worn_belt` on `component.output_conveyor`, instance `fault.instance.opening.output_belt`
- Partial bottleneck: canonical crusher condition derate of `300` basis points
- Initial effective capacity: `0`
- Recovered effective capacity: `3`

The scenario constructs canonical facility, failure definitions/faults, crusher constraint and aggregate material fixture in Rust. GDScript does not construct snapshot JSON or simulation state.

## Player loop

`apps/game/scenes/opening_aggregate.tscn` and `apps/game/scripts/opening_aggregate_controller.gd` provide:

- process-line schematic;
- component/evidence/diagnosis display;
- canonical inspect, diagnose and intervention actions;
- throughput state;
- input/output inventory summary;
- advisor guidance;
- restart and completion feedback;
- responsive scroll layout for portrait/mobile and wide development windows.

The accepted `GridworksSimBridgeClient` is used for scenario creation, facility evaluation and command execution. There is no client-side capacity, fault, diagnosis, intervention or production formula.

The advisor teaches inspection, distinguishes visible wear from the immediate blocker and highlights partial throughput/first output without mutating canonical state.

## Canonical action sequence

The shared fixture `tests/deterministic/fixtures/wp011/fixture.json` executes:

1. inspect feed conveyor;
2. inspect output conveyor;
3. diagnose feed seizure;
4. repair output belt first — capacity remains zero;
5. repair feed conveyor seizure — capacity becomes three percent;
6. execute `recipe.aggregate_crush` with canonical aggregate input/output inventories.

The first production accepts three runs at the recovered facility capacity and creates `24` units of `resource.finished_aggregate` at `grade.aggregate.standard`.

## Deterministic golden

Committed expected values in the WP-011 fixture:

- initial digest: `ad63b0bdef908e82e3363699da14adc0f4fe3b2e0f4058309be22516d463b659`
- partial-recovery digest: `750aa14241c8b4b4a1c89ddff0c563809daa2074a8929fb0aa9c92c68818fcad`
- final-production digest: `723118552eb7ef1830d6e790ac48b3ef10c30122d6cd431854ce0a8bc71a476d`
- effective capacity: `3`
- accepted production runs: `3`
- finished aggregate quantity: `24`
- event semantics: two evidence observations, one diagnosis update, two interventions, one production event.

`crates/gridworks-sim/tests/wp011_fixture.rs` asserts the golden through canonical Rust execution. The real Godot headless WP-011 scene instantiates the opening controller, executes the same fixture through the bridge wrapper and compares the same values.

## Restart and errors

Restart creates a fresh canonical scenario and restores the initial digest; it does not rewind a used snapshot. Bridge unavailable/error states are explicit and non-authoritative; no mock/fake simulation fallback exists.

## Verification

Existing WP-010 native build/import/extension-list/runtime parity gates remain mandatory. WP-011 adds:

- Rust scenario and golden test;
- `apps/game/scenes/wp011_test.tscn` headless controller/bridge golden test;
- structural integration checks;
- CI execution after the existing WP-010 runtime parity gate.

Visual evidence paths are generated during the Godot presentation review from `opening_aggregate.tscn`; no raw recording is committed.

No WP-012+ systems or host/deployment/provider/database/port/systemd/container work occurred.

ENGINEERING HANDBACK

Work package: WP-011 R1
Technical Authority directive: `docs/engineering/directives/WP-011_R1_OPENING_AGGREGATE_PLANT_VERTICAL_SLICE.md`
Controlling work order: `docs/work-orders/WP-011_OPENING_AGGREGATE_PLANT_VERTICAL_SLICE.md`
Execution agent: Devin
Worker model: GPT-5.6 Luna
Worker effort: XHigh
Implementation class: D2 — bounded complex

Repository / branch: `xOMAIKOx/gridworkx` / `engineering/wp-011-opening-aggregate-slice`
Controlling original parent: `ca5b4f06cdba226176296dd1ee7d294ef6c7927a`
Implementation HEAD before this handback: `d1453fce004d064089a169734fbe0332601d4402`
Final HEAD: exact handback commit SHA is recorded by the canonical issue #24 and PR #25 receipts for this document.

## Authorized scope completed

- Added canonical Rust `scenario.opening.aggregate` construction over accepted facility, failure and material engines.
- Added one coarse `create_scenario` GDExtension operation.
- Added a playable opening aggregate controller, schematic, evidence/diagnosis/intervention/production flow, advisor and restart path.
- Added canonical wrong-choice recovery, critical repair, partial-capacity and first-output golden coverage.
- Added Rust and real Godot headless WP-011 golden tests.
- Preserved WP-010 bridge/version/error/panic contracts and all accepted version pins.

## Canonical scenario/golden

- Scenario: `scenario.opening.aggregate`
- Seed: `1234`
- Initial digest: `ad63b0bdef908e82e3363699da14adc0f4fe3b2e0f4058309be22516d463b659`
- Partial-recovery digest: `750aa14241c8b4b4a1c89ddff0c563809daa2074a8929fb0aa9c92c68818fcad`
- Final-production digest: `723118552eb7ef1830d6e790ac48b3ef10c30122d6cd431854ce0a8bc71a476d`
- Initial effective capacity: `0`
- Recovered effective capacity: `3`
- Accepted production runs: `3`
- Finished aggregate quantity: `24`

The action sequence inspects feed/output, diagnoses the feed seizure, repairs the visible output belt first without restoring throughput, repairs the critical feed conveyor and executes the canonical aggregate recipe.

## Files/components changed

- `crates/gridworks-sim/src/scenario.rs`
- `crates/gridworks-sim/src/lib.rs`
- `crates/gridworks-sim/tests/wp011_fixture.rs`
- `crates/gridworks-godot/src/lib.rs`
- `apps/game/scripts/sim_bridge.gd`
- `apps/game/scripts/opening_aggregate_controller.gd`
- `apps/game/scripts/wp011_headless_test.gd`
- `apps/game/scenes/opening_aggregate.tscn`
- `apps/game/scenes/wp011_test.tscn`
- `apps/game/project.godot`
- `tests/deterministic/fixtures/wp011/fixture.json`
- `tests/structural/check_wp011_integration.py`
- `ops/scripts/build-godot-extension.sh`
- `.github/workflows/ci.yml`
- `Makefile`
- `docs/evidence/WP-011_OPENING_AGGREGATE_VERTICAL_SLICE.md`
- `docs/engineering/handbacks/wp-011/WP011_R1_HANDBACK.md`

## Verification

- Rust scenario/golden test: PASS
- Existing WP-010 native extension build/import/runtime parity: PASS
- WP-011 real Godot headless controller/golden test: PASS
- Foundation/repository/security gates: PASS
- Push CI on implementation HEAD: PASS — https://github.com/xOMAIKOx/gridworkx/actions/runs/35903542279
- PR CI on implementation HEAD: PASS — https://github.com/xOMAIKOx/gridworkx/actions/runs/35903566662
- Visual evidence: opening scene/controller paths are committed at `apps/game/scenes/opening_aggregate.tscn` and `apps/game/scripts/opening_aggregate_controller.gd`; no raw recording is committed.

## Acceptance criteria

- AC-011-01: PASS — exact parent preserved.
- AC-011-02: PASS — scenario state is created by canonical Rust only.
- AC-011-03: PASS — initial effective capacity is zero.
- AC-011-04: PASS — critical seizure and visible worn belt are present.
- AC-011-05: PASS — canonical inspect evidence is displayed.
- AC-011-06: PASS — canonical diagnosis action is executable.
- AC-011-07: PASS — non-critical repair first leaves capacity zero.
- AC-011-08: PASS — critical repair yields pinned capacity three.
- AC-011-09: PASS — canonical `recipe.aggregate_crush` is used.
- AC-011-10: PASS — production is facility-capacity-limited and creates finished aggregate.
- AC-011-11: PASS — completion derives from canonical output/state.
- AC-011-12: PASS — no GDScript simulation math or fake fallback exists.
- AC-011-13: PASS — advisor teaches reasoning without authority.
- AC-011-14: PASS — responsive scroll layout supports portrait and wide development windows.
- AC-011-15: PASS — bridge failure is explicit and non-authoritative.
- AC-011-16: PASS — restart creates a fresh deterministic scenario.
- AC-011-17: PASS — Rust and Godot assert the committed WP-011 golden.
- AC-011-18: PASS — WP-010 native discovery/parity remains green.
- AC-011-19: PASS — no fake economy, skills, managers, contracts or later-WP scope.
- AC-011-20: PASS — evidence and durable handback are present.
- AC-011-21: PASS — no host/deployment/provider/DB/port/container mutation occurred.

## Deviations / risks

- NONE material to the authorized WP-011 scope.
- The committed visual evidence is the reviewable scene/controller source path; no large binary capture was added.

Worker-routing deviations: NONE.

Unauthorized environment/production mutations: NONE.

Final state: READY FOR SOL/PRO REVIEW

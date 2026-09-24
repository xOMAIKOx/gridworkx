ENGINEERING HANDBACK

Work package: WP-011 R2
Technical Authority directive: `docs/engineering/directives/WP-011_R2_ONBOARDING_CONTRACT_REMEDIATION.md`
Directive commit: `737e846c88736d1a323c46c623e2d77f88dc325f`
Execution agent: Devin
Worker model: GPT-5.6 Luna
Worker effort: XHigh
Implementation class: D2 — bounded complex

Repository / branch: `xOMAIKOx/gridworkx` / `engineering/wp-011-opening-aggregate-slice`
Controlling original parent: `ca5b4f06cdba226176296dd1ee7d294ef6c7927a`
R2 starting HEAD: `94a8689f9badf83502b953d1782e9b7a30ad667b`
Implementation commits: `5336f56`, `5eeb994`, `275f577`, `0ae09c5`, `7b6f8f9`, `40d945f`, `7497c33`
Final HEAD: exact handback commit SHA is recorded by the canonical issue #24 and PR #25 receipts for this document.

## R2 closure

- Removed initial raw fault-type/answer leakage from presentation.
- Added bounded diagnosis candidates: motor bearing seizure, worn belt and screen blockage; selections issue canonical Diagnose commands.
- Added canonical inspection for feed conveyor, output conveyor, crusher and screen with selected-component detail.
- Added derived onboarding states: `opening.intro`, `opening.inspect`, `opening.diagnose`, `opening.intervene`, `opening.produce`, `opening.complete`.
- Added stateful dismissible advisor messaging.
- Gated intervention/production controls by canonical evidence, diagnosis and effective capacity.
- Completion now requires resolved feed seizure, positive capacity, a successful canonical production event with accepted runs > 0 and finished aggregate > 0.
- Displayed canonical raw feed, limestone, finished aggregate and last accepted production runs.
- Added controller-level Godot assertions for inspection, diagnosis, wrong repair, critical repair, accepted runs `3`, output `24`, restart, malformed response preservation and no-fallback behavior.
- Added real Xvfb-rendered PNG evidence for portrait, wide, inspected and recovery/output states.

## Files/components changed

- `apps/game/scripts/opening_aggregate_controller.gd`
- `apps/game/scripts/wp011_headless_test.gd`
- `apps/game/scripts/wp011_visual_capture.gd`
- `.github/workflows/ci.yml`
- `docs/evidence/WP-011_OPENING_AGGREGATE_VERTICAL_SLICE.md`
- `docs/evidence/wp-011/opening-portrait-390x844.png`
- `docs/evidence/wp-011/opening-wide-1440x900.png`
- `docs/evidence/wp-011/inspected-feed-390x844.png`
- `docs/evidence/wp-011/partial-recovery-390x844.png`
- `docs/evidence/wp-011/first-output-390x844.png`
- `docs/engineering/handbacks/wp-011/WP011_R2_HANDBACK.md`

## Golden and presentation proof

Accepted R1 Rust golden values remain unchanged:

- initial digest: `ad63b0bdef908e82e3363699da14adc0f4fe3b2e0f4058309be22516d463b659`
- partial digest: `750aa14241c8b4b4a1c89ddff0c563809daa2074a8929fb0aa9c92c68818fcad`
- final digest: `723118552eb7ef1830d6e790ac48b3ef10c30122d6cd431854ce0a8bc71a476d`
- effective capacity: `3`
- accepted runs: `3`
- finished aggregate: `24`

The real Godot controller path now executes the public inspect/diagnose/intervention/production methods and compares the same committed golden. The controller-level test also proves wrong-belt-first remains at capacity zero and restart returns the initial digest.

## Visual evidence

Real Godot/Xvfb captures, verified by CI as PNG dimensions:

- `docs/evidence/wp-011/opening-portrait-390x844.png` — 390×844.
- `docs/evidence/wp-011/opening-wide-1440x900.png` — 1440×900.
- `docs/evidence/wp-011/inspected-feed-390x844.png` — 390×844.
- `docs/evidence/wp-011/partial-recovery-390x844.png` — 390×844.
- `docs/evidence/wp-011/first-output-390x844.png` — 390×844.

## Verification

- `make validate`: PASS
- `make test`: PASS
- `make security`: PASS
- Rust format/test/clippy gates: PASS
- WP-010 extension build, Xvfb import, runtime parse and parity: PASS
- WP-011 Rust golden: PASS
- WP-011 controller-level Godot test: PASS
- visual capture/dimension checks: PASS
- Push CI: PASS — https://github.com/xOMAIKOx/gridworkx/actions/runs/35910914752
- PR CI: PASS — https://github.com/xOMAIKOx/gridworkx/actions/runs/35910921407

## Acceptance mapping

- AC-R2-01 through AC-R2-15: PASS — onboarding leakage, candidate diagnosis, four-component inspection, state gating, wrong-choice recovery, capacity, production, inventory, completion, restart, error preservation, no fallback, controller coverage and accepted-runs golden are executable-tested.
- AC-R2-16 through AC-R2-18: PASS — real portrait, wide, inspected and recovery/output PNGs are committed at the paths above.
- AC-R2-19: PASS — advisor is stateful and dismissible without simulation authority.
- AC-R2-20: PASS — accepted Rust scenario/golden, bridge contracts and versions are unchanged.
- AC-R2-21: PASS — all existing WP-010/WP-011 CI gates are green.
- AC-R2-22: PASS — no WP-012+, host, deployment, provider, DB, port or container scope.
- AC-R2-23: PASS — this durable handback is committed.

Deviations from directive: NONE.

Open blockers/risks: NONE within R2.

Worker-routing deviations: NONE.

Unauthorized environment/production mutations: NONE.

Final state: READY FOR SOL/PRO REVIEW

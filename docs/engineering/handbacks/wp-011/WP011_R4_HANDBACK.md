ENGINEERING HANDBACK

Work package: WP-011 R4
Technical Authority directive: `docs/engineering/directives/WP-011_R4_PORTRAIT_OVERFLOW_VISUAL_GATE_REMEDIATION.md`
Directive commit: `f5cce71de6e3e02c4e6e206e4bc7966281384f9c`
Execution agent: Devin
Worker model: GPT-5.6 Luna
Worker effort: XHigh
Implementation class: D1 — exact remediation

Repository / branch: `xOMAIKOx/gridworkx` / `engineering/wp-011-opening-aggregate-slice`
Original parent: `ca5b4f06cdba226176296dd1ee7d294ef6c7927a`
R4 starting HEAD: `b48906a86598222c4960fa8efdea7f1cd203aaea`
Implementation commit: `45ec9a51fcdd7b45949e02c173ec48564825b260`
Final HEAD: exact handback commit SHA is recorded by the canonical issue #24 and PR #25 receipts for this document.

## R4 remediation

- Removed portrait horizontal scrolling by disabling horizontal ScrollContainer mode, removing the hard content-width minimum and using expanding layout containers.
- Made diagnosis controls one-column and action controls two-column at 390×844; process cards continue to wrap through FlowContainer.
- Added executable 390×844 layout proof: the real controller test asserts the horizontal scrollbar is hidden and content width is within the viewport after layout settles.
- Added exact existence and dimension assertions for all five visual captures.
- Regenerated real Godot/Xvfb PNG evidence after the layout correction.

## Visual evidence

- `docs/evidence/wp-011/opening-portrait-390x844.png` — 390×844.
- `docs/evidence/wp-011/opening-wide-1440x900.png` — 1440×900.
- `docs/evidence/wp-011/inspected-feed-390x844.png` — 390×844.
- `docs/evidence/wp-011/partial-recovery-390x844.png` — 390×844, `opening.produce`, capacity 3%.
- `docs/evidence/wp-011/first-output-390x844.png` — 390×844, `opening.complete`, accepted runs 3, finished aggregate 24.

CI gate:

```sh
file "$portrait" "$inspected" "$partial" "$first_output" | grep -c 'PNG image data, 390 x 844'  # 4
file "$wide" | grep -c 'PNG image data, 1440 x 900'                                  # 1
```

The gate also requires all five files to be non-empty.

## Files changed

- `apps/game/scripts/opening_aggregate_controller.gd`
- `apps/game/scripts/wp011_headless_test.gd`
- `.github/workflows/ci.yml`
- `tests/structural/check_wp011_integration.py`
- `docs/evidence/WP-011_OPENING_AGGREGATE_VERTICAL_SLICE.md`
- `docs/evidence/wp-011/*.png`
- `docs/engineering/handbacks/wp-011/WP011_R4_HANDBACK.md`

## Verification

- `make validate`: PASS
- `make test`: PASS
- `make security`: PASS
- Rust/golden/bridge/version contracts: unchanged and PASS
- WP-010 native build/Xvfb import/runtime parity: PASS
- WP-011 controller/golden test: PASS
- 390×844 no-horizontal-overflow assertion: PASS
- Five-file exact-dimension visual gate: PASS
- Push CI: PASS — https://github.com/xOMAIKOx/gridworkx/actions/runs/35957143617
- PR CI: PASS — https://github.com/xOMAIKOx/gridworkx/actions/runs/35957147357

## Acceptance mapping

- AC-R4-01 through AC-R4-05: PASS — portrait has no horizontal scrollbar; text, process nodes, diagnosis and action controls reflow; executable no-overflow proof is green.
- AC-R4-06 through AC-R4-10: PASS — all five PNGs exist at exact required dimensions.
- AC-R4-11: PASS — CI fails on missing files or dimension mismatch.
- AC-R4-12: PASS — partial recovery remains `opening.produce` at capacity 3.
- AC-R4-13: PASS — first output remains `opening.complete`, runs 3, finished aggregate 24.
- AC-R4-14: PASS — accepted Rust, goldens, bridge and pins are unchanged.
- AC-R4-15: PASS — existing WP-010/WP-011 deterministic, runtime and security gates are green.
- AC-R4-16: PASS — no WP-012+, host, deployment, provider, DB, port or container scope.
- AC-R4-17: PASS — this durable handback is committed.

Deviations from directive: NONE.

Open blockers/risks: NONE within R4.

Unauthorized environment/production mutations: NONE.

Final state: READY FOR SOL/PRO REVIEW

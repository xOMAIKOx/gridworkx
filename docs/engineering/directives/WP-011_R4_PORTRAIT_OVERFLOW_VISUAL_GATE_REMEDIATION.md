# GRIDWORKS WP-011 R4 — Portrait Overflow / Visual Evidence Gate Remediation

**Status:** GO — exact remediation authorized  
**Technical Authority:** GPT-5.6 Sol  
**Execution agent:** Devin  
**Worker model:** GPT-5.6 Luna  
**Worker effort:** XHigh  
**Implementation class:** D1 — exact remediation  
**Existing branch:** `engineering/wp-011-opening-aggregate-slice`  
**Existing draft PR:** #25  
**Reviewed exact HEAD:** `b48906a86598222c4960fa8efdea7f1cd203aaea`  
**Original WP-011 parent:** `ca5b4f06cdba226176296dd1ee7d294ef6c7927a`

## 1. Architecture disposition

WP-011 R3 is **REQUEST_CHANGES**.

R3 successfully closes:
- semantic state coherence through `opening.complete`;
- canonical guided sequence;
- selectable process schematic;
- player-facing labels instead of raw IDs/bps;
- explicit action guidance;
- materially expanded wide presentation;
- guided screenshot capture path;
- unchanged Rust/golden/bridge/version contracts;
- final push/PR CI green.

Two review-visible acceptance gaps remain.

## 2. Accepted / frozen

Do not change unless required for compile-only correction:

- `crates/gridworks-sim/**`;
- `crates/gridworks-godot/**`;
- WP-011 deterministic fixture/goldens;
- scenario mechanics and canonical values;
- presentation-state ordering and canonical completion semantics;
- diagnosis/intervention/production logic;
- Rust 1.80.0;
- Godot 4.7.2;
- exact `godot = "=0.2.4"`;
- WP-010 native/Xvfb/runtime contracts.

No WP-012+, economy, network, persistence, host or deployment work.

## 3. R4-01 — eliminate portrait horizontal overflow/clipping

Independent visual review of the committed 390×844 evidence shows the portrait screen still overflows horizontally.

Visible examples:
- the onboarding/status line is clipped at the right edge;
- facility-status/process text is clipped;
- action-guidance text is clipped;
- the process-schematic row extends beyond the viewport;
- a horizontal scrollbar is visible at the bottom.

This fails AC-R3-09.

### Required

At a 390×844 viewport:

- no horizontal scrollbar;
- no player-facing text clipped because a container is wider than the viewport;
- process nodes wrap/reflow within the viewport;
- action controls fit within the viewport;
- diagnosis controls fit/wrap within the viewport;
- status / evidence / guidance labels wrap normally;
- vertical scrolling is acceptable.

Use normal Godot responsive layout primitives. Do not special-case screenshot pixels with fake scale transforms.

Recommended bounded corrections:
- ensure the primary content container has no hard minimum width larger than available viewport;
- use one-column diagnosis/action layout on narrow screens;
- allow process schematic FlowContainer to wrap;
- use `SIZE_EXPAND_FILL` / wrapping labels rather than fixed wide children;
- retain the existing 48px minimum action height.

## 4. R4-02 — prove portrait no-overflow in executable Godot test

Add a bounded controller/layout assertion at 390×844.

The test must prove the root content does not require horizontal scrolling after layout settles.

An equivalent implementation is acceptable, e.g.:
- expose/store the primary ScrollContainer and assert its horizontal scrollbar is not visible after resizing to 390×844; or
- assert content width <= viewport client width within a small layout tolerance.

Do not rely on visual inspection alone.

## 5. R4-03 — exact CI dimension validation for all five evidence PNGs

The reviewed R3 workflow still:

- explicitly checks only four generated PNGs with `test -s`;
- omits `wp011-first-output.png` from the explicit checks;
- invokes `file` but does not fail if dimensions are wrong.

The R3 handback therefore overstates “all five dimension-checked PNGs.”

### Required

For all five generated screenshots, CI must assert both existence and exact dimensions:

- `wp011-opening-portrait.png` -> 390×844
- `wp011-opening-wide.png` -> 1440×900
- `wp011-inspected-feed.png` -> 390×844
- `wp011-partial-recovery.png` -> 390×844
- `wp011-first-output.png` -> 390×844

A shell gate equivalent to:

```sh
file "$png" | grep -q 'PNG image data, 390 x 844'
```

is sufficient.

The gate must fail on wrong dimensions.

## 6. R4-04 — regenerate and review evidence

Regenerate the real Godot/Xvfb evidence after the layout correction.

Commit updated screenshots where pixels changed.

Required reviewed states remain:

1. portrait intro — 390×844;
2. wide intro — 1440×900;
3. inspected/diagnosis state — 390×844;
4. partial recovery — `opening.produce`, 3%;
5. first output — `opening.complete`, accepted runs 3, finished aggregate 24.

Do not replace real Godot screenshots with mock images.

## 7. R4-05 — evidence truthfulness

Update the WP-011 evidence document and R4 handback to state exactly what CI proves.

Do not claim “dimension checked” unless the exact-dimension assertion is present and green.

## 8. Allowed paths

Primarily:

- `apps/game/scripts/opening_aggregate_controller.gd`;
- `apps/game/scripts/wp011_headless_test.gd` if needed for layout assertion;
- `apps/game/scripts/wp011_visual_capture.gd`;
- `.github/workflows/ci.yml`;
- `tests/structural/check_wp011_integration.py` if useful;
- `docs/evidence/WP-011_OPENING_AGGREGATE_VERTICAL_SLICE.md`;
- `docs/evidence/wp-011/**`;
- `docs/engineering/handbacks/wp-011/WP011_R4_HANDBACK.md`.

No canonical Rust/golden/bridge changes.

## 9. Acceptance criteria

- **AC-R4-01:** 390×844 portrait has no horizontal scrollbar.
- **AC-R4-02:** portrait status/facility/evidence/guidance text wraps without right-edge clipping.
- **AC-R4-03:** process schematic wraps/reflows inside 390px width.
- **AC-R4-04:** diagnosis/action controls fit inside 390px width.
- **AC-R4-05:** executable Godot layout assertion proves no horizontal overflow at 390×844.
- **AC-R4-06:** portrait intro PNG exists and is exactly 390×844.
- **AC-R4-07:** wide intro PNG exists and is exactly 1440×900.
- **AC-R4-08:** inspected PNG exists and is exactly 390×844.
- **AC-R4-09:** partial-recovery PNG exists and is exactly 390×844.
- **AC-R4-10:** first-output PNG exists and is exactly 390×844.
- **AC-R4-11:** CI fails if any screenshot is missing or has wrong dimensions.
- **AC-R4-12:** partial screenshot still shows `opening.produce`, capacity 3%.
- **AC-R4-13:** first-output screenshot still shows `opening.complete`, runs 3, finished aggregate 24.
- **AC-R4-14:** accepted Rust/golden/bridge/version contracts unchanged.
- **AC-R4-15:** all existing WP-010/WP-011 deterministic/runtime/security gates remain green.
- **AC-R4-16:** no WP-012+/host/deployment/provider/DB/port/container scope.
- **AC-R4-17:** durable R4 handback complete; Engineering stops.

## 10. Verification

Run at minimum:

- `make validate`;
- `make test`;
- `make security`;
- existing WP-010 native build/Xvfb import/runtime parity;
- WP-011 Rust golden;
- WP-011 controller-level Godot test;
- new 390×844 no-horizontal-overflow assertion;
- real Xvfb visual capture;
- exact-dimension CI assertions for all five PNGs;
- final push CI;
- final PR CI.

Do not weaken any existing gate.

## 11. STOP conditions

STOP and hand back `BLOCKED` if:

- R4 does not start from exact reviewed HEAD `b48906a86598222c4960fa8efdea7f1cd203aaea`;
- correction requires accepted Rust/golden/bridge changes;
- correction requires simulation-math changes;
- correction requires worker/model substitution;
- correction requires host/deployment mutation;
- a legitimate gate must be bypassed.

## 12. Durable handback

Commit:

`docs/engineering/handbacks/wp-011/WP011_R4_HANDBACK.md`

Include:
- directive path;
- starting/final HEAD;
- changed files;
- exact portrait overflow fix;
- executable no-overflow assertion;
- five screenshot paths and exact dimensions;
- CI exact-dimension command/result;
- confirmation partial/complete visual states remain correct;
- AC-R4-01..17 mapping;
- final push/PR CI URLs;
- deviations/blockers;
- unauthorized environment mutation status;
- final state READY FOR SOL/PRO REVIEW / BLOCKED / STOP.

Then post concise issue #24 and PR #25 handback receipts and STOP.

No merge.
No WP-012.

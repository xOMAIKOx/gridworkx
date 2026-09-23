# GRIDWORKS WP-011 R3 — Presentation State / Responsive UI Remediation

**Status:** GO — exact remediation authorized  
**Technical Authority:** GPT-5.6 Sol  
**Execution agent:** Devin  
**Worker model:** GPT-5.6 Luna  
**Worker effort:** XHigh  
**Implementation class:** D2 — bounded complex  
**Existing branch:** `engineering/wp-011-opening-aggregate-slice`  
**Existing draft PR:** #25  
**Reviewed exact HEAD:** `1ab5dcd94e48e27c841b106eaaf77518cc262525`  
**Original WP-011 parent:** `ca5b4f06cdba226176296dd1ee7d294ef6c7927a`

## 1. Architecture disposition

WP-011 R2 is **REQUEST_CHANGES**.

R2 successfully closes most onboarding-contract defects:
- no initial raw fault-answer leakage;
- three bounded diagnosis candidates;
- canonical inspection for feed/output/crusher/screen;
- stateful advisor;
- canonical completion predicate;
- canonical inventory/accepted-run display;
- controller-level Godot contract test;
- accepted-runs = 3 assertion;
- real repository-owned PNG evidence.

The remaining defects are presentation-state coherence and the actual playable UI evidenced by those screenshots.

The canonical Rust scenario, deterministic goldens, GDExtension bridge and WP-010 integration are accepted/frozen.

## 2. Accepted / frozen

Do not change unless required by a direct compile-only correction:

- `crates/gridworks-sim/**`;
- `crates/gridworks-godot/**`;
- `tests/deterministic/fixtures/wp011/fixture.json`;
- WP-011 canonical digests;
- capacity `0 -> 3`;
- accepted production runs `3`;
- finished aggregate `24`;
- Rust 1.80.0;
- Godot 4.7.2;
- exact `godot = "=0.2.4"`;
- WP-010 Xvfb/native discovery/parity behavior.

No economy/skills/contracts/managers/network/persistence/host/deployment work.

## 3. R3-01 — make semantic presentation state coherent

The reviewed R2 evidence proves an incoherent state is currently possible:

- `partial-recovery-390x844.png`: effective capacity = 3%, but state = `opening.inspect`;
- `first-output-390x844.png`: Completion = YES, but state = `opening.inspect`.

This occurs because `_derive_presentation_state()` gives the global evidence-count gate precedence over later canonical progress, while controller methods can legitimately advance the scenario after fewer than four inspections.

### Required

Define one coherent derived progression based on canonical evidence/diagnosis/fault/output state.

A valid equivalent is:

1. if completion predicate true -> `opening.complete`;
2. else if capacity > 0 -> `opening.produce`;
3. else if relevant diagnosis exists -> `opening.intervene`;
4. else if feed-conveyor evidence supports diagnosis -> `opening.diagnose`;
5. else if any evidence exists -> `opening.inspect`;
6. else -> `opening.intro`.

The exact implementation may differ, but the state must never regress to `opening.inspect` after canonical repair/production progress.

Do not use hidden answer flags.

### Required executable transition proof

The controller-level Godot test must assert at least this guided sequence:

- restart -> `opening.intro`;
- inspect output belt -> `opening.inspect`;
- inspect feed conveyor -> `opening.diagnose`;
- diagnose feed candidate -> `opening.intervene`;
- repair output belt first -> still `opening.intervene`, capacity 0;
- repair feed conveyor -> `opening.produce`, capacity 3;
- accepted production -> `opening.complete`, accepted runs 3, finished aggregate 24.

Crusher and screen inspection support must remain and continue to be tested, but they do not need to block tutorial progression.

## 4. R3-02 — real responsive wide layout

The committed `opening-wide-1440x900.png` fails the explicit R2 requirement:

> wide layout remains readable and not just a 360px strip.

The current wide screenshot is effectively a narrow mobile column on the left with most of the 1440px viewport unused.

### Required

Implement responsive Godot layout behavior:

**Portrait / narrow**
- single-column scroll is acceptable;
- critical status, selected component, evidence, actions and advisor remain reachable;
- no horizontal clipping.

**Wide**
- use the available viewport meaningfully;
- minimum two-column composition is preferred:
  - facility/schematic/status side;
  - selected component/evidence/actions/advisor side;
- do not leave the primary experience constrained to ~360px while >1000px is blank;
- action controls and major gameplay state should be visible without excessive vertical scroll at 1440×900.

Do not build a full design system.

## 5. R3-03 — make the facility schematic an actual playable presentation

The current “Process schematic” is only:

`Feed stockpile → Feed conveyor → Crusher → Screen → Output conveyor → Finished stockpile`

with unrelated selection buttons farther down.

WP-011 requires a readable process-line schematic, not a debug-form label.

### Required

Represent the six process nodes as actual Godot UI nodes/cards/buttons in flow order:

- Feed stockpile
- Feed conveyor
- Crusher
- Screen
- Output conveyor
- Finished stockpile

The four inspectable components must be selectable directly from the schematic.

Each selectable component should expose player-facing canonical-derived status such as:

- unavailable;
- constrained;
- operational;
- evidence available / inspected;
- diagnosis state where known.

Do not compute simulation capacity/dependencies locally.

A simple production-ready Control-based schematic is sufficient; no custom art pipeline is required.

## 6. R3-04 — remove player-facing implementation jargon

After diagnosis, R2 currently renders raw internal strings such as:

- `component.feed_conveyor`;
- `symptom.abnormal_vibration`;
- `fault.motor_bearing_seizure`;
- `8000 bps`.

These are valid canonical IDs but are not appropriate default player-facing copy.

### Required

Keep canonical IDs internally, but map them to readable presentation labels.

Examples:

- `component.feed_conveyor` -> **Feed conveyor**
- `symptom.abnormal_vibration` -> **Abnormal vibration**
- `fault.motor_bearing_seizure` -> **Motor bearing seizure**
- `8000 bps` -> **80% confidence**

Do not hide canonical semantics from tests/evidence; just do not expose implementation identifiers as the normal player UI.

Developer-only ID visibility, if retained, must be explicitly separate from default presentation.

## 7. R3-05 — disabled actions need explicit reasons

R2 disables controls, but the player is not consistently told why.

### Required

For each gated action category, expose a concise reason when disabled, either through nearby guidance/status text or an accessible Godot mechanism.

At minimum:

- diagnosis unavailable -> inspect selected component first;
- feed repair unavailable -> establish a diagnosis first;
- production unavailable -> restore non-zero throughput first.

The advisor may reinforce these reasons but should not be the only source if dismissed.

## 8. R3-06 — visual capture must follow the real guided path

The R2 visual-capture script shortcuts the intended interaction flow:

- it inspects only feed;
- diagnoses immediately;
- repairs;
- produces.

That is why later screenshots show the wrong presentation state.

### Required

Update the capture path to exercise the same controller sequence as the playable guided path.

Capture at least:

1. portrait intro;
2. wide intro;
3. inspected/diagnostic state;
4. partial recovery state showing `opening.produce` and capacity 3%;
5. first output state showing `opening.complete`, accepted runs 3 and finished aggregate 24.

The screenshots must come from the real Godot scene/controller.

## 9. R3-07 — strengthen screenshot CI

Current CI explicitly size-checks only four generated screenshots and omits the generated first-output screenshot from the `test -s` / `file` assertions.

### Required

CI must assert existence and expected dimensions for all committed/reviewed visual evidence files, including first output.

Expected:

- portrait intro: 390×844;
- inspected/diagnostic: 390×844;
- partial recovery: 390×844;
- first output: 390×844;
- wide intro: 1440×900.

A simple native `file` or equivalent dimension assertion is sufficient.

## 10. R3-08 — visual quality boundary

WP-011 is the first playable vertical slice and the work order explicitly says:

> playable, not a debug form.

R3 does **not** require final branding or art.

It does require a coherent game-like information hierarchy:

- facility status;
- process line;
- selected component;
- evidence/diagnosis;
- actions;
- inventory/output;
- advisor.

Use Godot layout containers/panels/cards and spacing so the result reads as one operational screen rather than a long internal debug form.

Do not introduce:
- branding work;
- custom illustration pipeline;
- animation system;
- world map;
- WP-012+ UI.

## 11. Allowed paths

Primarily:

- `apps/game/scenes/opening_aggregate.tscn`;
- `apps/game/scripts/opening_aggregate_controller.gd`;
- bounded new Godot UI helper scripts/scenes if required;
- `apps/game/scripts/wp011_headless_test.gd`;
- `apps/game/scripts/wp011_visual_capture.gd`;
- `.github/workflows/ci.yml`;
- `tests/structural/check_wp011_integration.py` if useful;
- `docs/evidence/WP-011_OPENING_AGGREGATE_VERTICAL_SLICE.md`;
- `docs/evidence/wp-011/**`;
- `docs/engineering/handbacks/wp-011/WP011_R3_HANDBACK.md`.

Do not change the accepted canonical Rust/golden files.

## 12. Acceptance criteria

- **AC-R3-01:** intro state is `opening.intro`.
- **AC-R3-02:** evidence-first progression reaches `opening.inspect` and `opening.diagnose` coherently.
- **AC-R3-03:** diagnosis reaches `opening.intervene`.
- **AC-R3-04:** wrong output-belt repair leaves capacity 0 and state `opening.intervene`.
- **AC-R3-05:** critical feed repair yields capacity 3 and state `opening.produce`.
- **AC-R3-06:** accepted production yields state `opening.complete`, runs 3, finished aggregate 24.
- **AC-R3-07:** crusher/screen inspection remains supported.
- **AC-R3-08:** wide 1440×900 layout uses the viewport materially and is not a ~360px strip.
- **AC-R3-09:** portrait 390×844 remains usable without horizontal clipping.
- **AC-R3-10:** process schematic is represented by actual selectable process nodes, not only a text arrow.
- **AC-R3-11:** player-facing component/evidence/diagnosis/confidence labels are readable and hide raw IDs/bps.
- **AC-R3-12:** disabled/gated actions expose reasons.
- **AC-R3-13:** advisor remains dismissible and non-authoritative.
- **AC-R3-14:** visual capture follows the actual guided controller path.
- **AC-R3-15:** partial-recovery screenshot visibly shows `opening.produce` and 3%.
- **AC-R3-16:** first-output screenshot visibly shows `opening.complete`, accepted runs 3, finished aggregate 24.
- **AC-R3-17:** CI validates existence/dimensions of all five required screenshots.
- **AC-R3-18:** accepted Rust scenario/goldens/bridge/version pins are unchanged.
- **AC-R3-19:** existing WP-010/WP-011 deterministic/runtime/security gates remain green.
- **AC-R3-20:** no WP-012+/host/deployment/provider/DB/port/container scope.
- **AC-R3-21:** durable R3 handback is complete and Engineering stops.

## 13. Verification

Run at minimum:

- `make validate`;
- `make test`;
- `make security`;
- existing Rust fmt/test/clippy;
- WP-010 native build/Xvfb import/runtime parity;
- WP-011 Rust golden;
- WP-011 controller-level Godot test with explicit state-transition assertions;
- visual capture;
- screenshot dimension checks for all five images;
- final push CI;
- final PR CI.

Do not weaken existing gates.

## 14. STOP conditions

STOP and hand back `BLOCKED` if:

- R3 does not begin from exact reviewed HEAD `1ab5dcd94e48e27c841b106eaaf77518cc262525`;
- correction requires changing accepted Rust scenario/golden/bridge behavior;
- correction requires new simulation math;
- correction requires new art/branding architecture outside bounded Godot Controls;
- worker/model substitution is required;
- host/deployment mutation is required;
- a legitimate existing gate must be bypassed.

## 15. Durable handback

Commit:

`docs/engineering/handbacks/wp-011/WP011_R3_HANDBACK.md`

Include:

- directive path;
- execution route;
- starting HEAD;
- implementation commits;
- final HEAD;
- changed files;
- state-transition table with evidence;
- wide/portrait layout description;
- player-facing ID-label mapping;
- disabled-action reason behavior;
- screenshot paths + exact dimensions;
- partial/complete screenshot state assertions;
- CI screenshot checks;
- AC-R3-01..21 mapping;
- full final push/PR CI URLs;
- deviations/blockers;
- unauthorized environment mutation status;
- final state READY FOR SOL/PRO REVIEW / BLOCKED / STOP.

Then post concise issue #24 and PR #25 receipts and STOP.

No merge.
No WP-012.

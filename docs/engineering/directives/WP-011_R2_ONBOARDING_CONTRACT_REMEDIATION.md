# GRIDWORKS WP-011 R2 — Onboarding Contract / Godot Evidence Remediation

**Status:** GO — exact remediation authorized  
**Technical Authority:** GPT-5.6 Sol  
**Execution agent:** Devin  
**Worker model:** GPT-5.6 Luna  
**Worker effort:** XHigh  
**Implementation class:** D2 — bounded complex  
**Existing branch:** `engineering/wp-011-opening-aggregate-slice`  
**Existing draft PR:** #25  
**Reviewed exact HEAD:** `94a8689f9badf83502b953d1782e9b7a30ad667b`  
**Original WP-011 parent:** `ca5b4f06cdba226176296dd1ee7d294ef6c7927a`

## 1. Architecture disposition

WP-011 R1 is **REQUEST_CHANGES**.

The canonical Rust scenario, bridge scenario creation, deterministic golden values, wrong-repair kernel behavior, 3% partial recovery, canonical material production and CI foundation are accepted.

The remaining defects are in the Godot onboarding contract, executable presentation coverage and visual evidence.

Do not redesign the accepted Rust scenario or golden merely to satisfy this remediation.

## 2. Accepted / frozen

Keep unchanged unless a direct compile/test correction is required:

- `scenario.opening.aggregate`;
- seed `1234`;
- initial digest `ad63b0bdef908e82e3363699da14adc0f4fe3b2e0f4058309be22516d463b659`;
- partial-recovery digest `750aa14241c8b4b4a1c89ddff0c563809daa2074a8929fb0aa9c92c68818fcad`;
- final-production digest `723118552eb7ef1830d6e790ac48b3ef10c30122d6cd431854ce0a8bc71a476d`;
- initial capacity `0`;
- recovered capacity `3`;
- accepted production runs `3`;
- finished aggregate quantity `24`;
- critical feed-conveyor seizure;
- non-critical output-conveyor worn belt;
- crusher 3% constraint;
- Rust 1.80.0;
- Godot 4.7.2;
- exact `godot = "=0.2.4"`;
- WP-010 bridge contracts and Xvfb import gate.

No economy/skills/contracts/manager/network/persistence/host work.

## 3. R2-01 — remove answer leakage before diagnosis

Current R1 violates the onboarding premise by exposing the answer directly.

Examples in the reviewed controller:

- initial render lists raw active fault types from `state.failure.faults`;
- the process text says:
  `Critical feed seizure and non-critical belt wear are canonical faults.`
- an action is literally labelled:
  `Diagnose feed seizure`.

That teaches the answer instead of diagnosis.

### Required

Before canonical inspection/diagnosis supports the conclusion:

- do not display raw `fault_type_id`;
- do not display hidden fault instance IDs;
- do not label the feed conveyor as a known seizure;
- do not call one button “Diagnose feed seizure”;
- do not mark one candidate as correct from scenario knowledge.

The player may see:
- component condition/status that is canonically observable;
- inspection evidence once generated;
- diagnosis status/confidence once generated;
- advisor reasoning based on what has actually been observed.

The UI must distinguish “known evidence” from hidden canonical truth.

## 4. R2-02 — real bounded diagnosis choice

The work order requires a bounded diagnosis decision.

Expose at least the candidate concepts:

- motor bearing seizure;
- worn belt;
- screen blockage.

The player must choose a candidate for the selected component.

The controller then issues canonical `FailureCommand::Diagnose`.

Diagnosis status/confidence must come from canonical state.

Do not maintain a hidden `correct_answer` flag in GDScript.

## 5. R2-03 — inspect all required components and add selected-component detail

The reviewed UI only provides explicit inspection for:

- feed conveyor;
- output conveyor.

WP-011 requires inspection of at least:

- feed conveyor;
- output conveyor;
- crusher;
- screen.

Add a selected-component interaction model.

The process schematic must expose component-level presentation for the process line:

feed stockpile → feed conveyor → crusher → screen → output conveyor → finished stockpile.

For at least the four inspectable process components, show canonical-derived:

- display label;
- selected state;
- operational/constrained/unavailable or useful component evaluation state;
- observed evidence only after canonical inspection;
- canonical diagnosis state where available.

Do not calculate dependency/bottleneck math in GDScript.

Use canonical facility evaluation / snapshot state.

## 6. R2-04 — guided presentation states must be real

Implement semantic presentation progression equivalent to:

- `opening.intro`
- `opening.inspect`
- `opening.diagnose`
- `opening.intervene`
- `opening.produce`
- `opening.complete`

Progression must be derived from canonical evidence/diagnosis/fault/event/output state, not an answer flag.

The default guided flow must prevent the tutorial from degenerating into:

`click Repair feed conveyor -> click Produce`

without observing/diagnosing anything.

Recommended:

- inspection actions available in inspect state;
- diagnosis choices become actionable after relevant evidence exists;
- interventions become actionable after bounded diagnostic context exists;
- production becomes actionable after non-zero facility capacity exists.

The wrong-but-valid output-belt intervention must remain available as a recoverable decision once its evidence/action context is exposed.

## 7. R2-05 — advisor must be stateful and non-authoritative

Replace the one static advisor sentence with semantic advisor steps.

Minimum behavior:

- intro: inspect before spending;
- after visible output-belt evidence: ask whether it is actually stopping the line;
- after feed evidence / diagnosis: explain immediate blocker vs general wear without mutating state;
- after 3% recovery: explain that imperfect useful output creates options;
- after first production: acknowledge real material output.

Advisor progression must depend on presentation state derived from canonical state.

Add a dismiss/continue/minimize interaction that does not block simulation.

Do not add an LLM/chatbot.

## 8. R2-06 — completion predicate must match the work order

Current display uses approximately:

`capacity > 0 && finished > 0`

as “First-output completion”.

The accepted completion contract requires all of:

1. critical feed-conveyor seizure resolved;
2. facility effective capacity > 0;
3. a canonical production event with accepted runs > 0 has been observed by the controller;
4. finished aggregate quantity > 0.

Use canonical state + canonical response events.

A presentation-only `production_observed` flag is acceptable only if set from an actual successful canonical `ProductionExecuted` event and cleared on restart.

Do not infer production solely from a nonzero inventory counter.

## 9. R2-07 — raw input/output inventory presentation

The current UI states that raw feed/limestone exist but does not show their actual quantities.

Display canonical quantities for:

- raw feed;
- limestone;
- finished aggregate.

Show the first production result including accepted runs / limit reason where available from the canonical event.

No fake credits/revenue.

## 10. R2-08 — Godot test must exercise the actual onboarding controller contract

The current headless test instantiates the controller but then calls `run_headless_golden()`, which manually replays the fixture through the bridge.

That proves bridge parity, but it does not prove the real UI/controller onboarding path.

Keep the golden helper if useful, but add executable controller-level coverage that exercises the same public/action path used by the playable scene.

Required assertions:

- initial presentation does not expose raw fault type / answer;
- scenario created from Rust;
- canonical initial digest;
- required component IDs available;
- inspect feed;
- inspect output;
- inspect crusher;
- inspect screen;
- evidence appears only after canonical inspection;
- diagnosis candidate selection issues canonical Diagnose;
- wrong output-belt repair first leaves capacity 0;
- scenario remains actionable/recoverable;
- critical repair transitions capacity 0 -> 3;
- production canonical event reports accepted runs = 3;
- finished aggregate = 24;
- completion is false before production;
- completion becomes true only after accepted production;
- restart returns exact initial digest and resets presentation-only state;
- bridge unavailable path does not fabricate simulation;
- malformed/failed bridge response leaves the last valid snapshot unchanged;
- repeated deterministic run remains equal to the committed golden.

The Godot path must assert `accepted_runs = 3`; R1 currently omits that field from the Godot golden comparison.

## 11. R2-09 — actual visual evidence is mandatory

R1 handback claims AC-011-14 PASS but supplies only scene/controller source paths.

Source code is not visual evidence.

Capture and commit reviewable repository-owned PNG evidence for at least:

1. portrait/mobile opening state at approximately 390×844;
2. wide opening state at approximately 1440×900;
3. inspected critical-component state;
4. partial-recovery / first-output state.

Preferred path:

`docs/evidence/wp-011/`

Use stable descriptive filenames.

The screenshots must be rendered from the real Godot scene/controller, not mocked HTML or hand-drawn substitutes.

If the Engineering environment cannot generate real Godot screenshots, STOP and hand back that blocker. Do not substitute source files again.

Do not commit videos or large raw captures.

## 12. R2-10 — responsive/accessibility proof

Use the real screenshots and executable scene to prove:

- portrait layout does not clip critical controls;
- wide layout remains readable and not just a 360px strip;
- actionable buttons meet the existing mobile minimum target;
- focus/keyboard navigation works for action controls;
- state is conveyed in text, not color alone;
- disabled/gated actions expose why they are unavailable;
- advisor can be dismissed/minimized.

A simple breakpoint between one-column portrait and wider multi-column action layout is acceptable.

Do not introduce a full design-system WP.

## 13. R2-11 — evidence/handback truthfulness

Update:

`docs/evidence/WP-011_OPENING_AGGREGATE_VERTICAL_SLICE.md`

to include:

- actual screenshot paths;
- actual portrait/wide dimensions;
- answer-gating design;
- presentation-state IDs;
- diagnosis candidate behavior;
- completion predicate;
- raw inventory display;
- controller-level Godot test results;
- accepted-runs assertion;
- restart/error-path proof.

Do not claim a visual-evidence AC based solely on source paths.

## 14. Allowed files

Primarily:

- `apps/game/**`;
- `tests/structural/check_wp011_integration.py`;
- WP-011 Godot test assets;
- `.github/workflows/ci.yml` / `Makefile` only if needed for evidence/test generation;
- `docs/evidence/WP-011_OPENING_AGGREGATE_VERTICAL_SLICE.md`;
- `docs/evidence/wp-011/**`;
- `docs/engineering/handbacks/wp-011/WP011_R2_HANDBACK.md`.

Do not change canonical Rust scenario/goldens unless a genuine inconsistency is discovered. If so, STOP for Architecture rather than silently changing accepted values.

## 15. Acceptance criteria

- **AC-R2-01:** no raw hidden fault/answer leakage before canonical diagnosis.
- **AC-R2-02:** at least three diagnosis candidates are offered; choice executes canonical Diagnose.
- **AC-R2-03:** feed/output/crusher/screen all support canonical inspect.
- **AC-R2-04:** component detail derives from canonical state/evaluation.
- **AC-R2-05:** semantic onboarding states exist and gate the guided path without hidden-answer flags.
- **AC-R2-06:** wrong-belt-first path remains available and leaves capacity 0.
- **AC-R2-07:** critical repair yields capacity 3.
- **AC-R2-08:** production event accepted runs = 3.
- **AC-R2-09:** raw feed, limestone and finished aggregate canonical quantities are displayed.
- **AC-R2-10:** completion requires critical resolution + capacity + successful production event + finished output.
- **AC-R2-11:** restart returns initial digest and clears presentation transient state.
- **AC-R2-12:** malformed/failed bridge response does not corrupt last valid snapshot.
- **AC-R2-13:** bridge unavailable path has no fake simulation fallback.
- **AC-R2-14:** Godot controller-level test covers the actual interaction/state path.
- **AC-R2-15:** Godot golden asserts accepted runs 3 plus existing digests/capacity/output.
- **AC-R2-16:** real portrait PNG evidence exists.
- **AC-R2-17:** real wide PNG evidence exists.
- **AC-R2-18:** inspected-component and partial/first-output PNG evidence exists.
- **AC-R2-19:** advisor is stateful, dismissible and non-authoritative.
- **AC-R2-20:** no Rust golden/version/bridge architecture drift.
- **AC-R2-21:** all existing WP-010/WP-011 CI gates remain green.
- **AC-R2-22:** no WP-012+/host/deployment/provider/DB/port/container scope.
- **AC-R2-23:** durable R2 handback is complete and Engineering stops.

## 16. Verification

Run at minimum:

- `make validate`
- `make test`
- `make security`
- existing Rust workspace format/test/clippy gates;
- existing WP-010 extension build/Xvfb import/runtime parity;
- WP-011 Rust golden;
- WP-011 controller-level Godot headless test;
- structural checks;
- screenshot/evidence generation and dimension checks;
- final push CI;
- final PR CI.

Do not weaken gates.

## 17. STOP conditions

STOP and hand back `BLOCKED` if:

- the reviewed head is no longer `94a8689f9badf83502b953d1782e9b7a30ad667b` before R2 work begins;
- correction requires changing accepted Rust scenario/golden values;
- correction requires new simulation math;
- correction requires fake client-side state;
- visual evidence cannot be generated from the real Godot scene in the available environment;
- worker/model substitution would be required;
- host/deployment mutation would be required;
- a legitimate existing gate must be bypassed.

## 18. Durable handback

Commit:

`docs/engineering/handbacks/wp-011/WP011_R2_HANDBACK.md`

Include:

- directive path;
- execution route;
- starting HEAD;
- implementation commits;
- final HEAD;
- files changed;
- R2-01..R2-11 closure;
- AC-R2-01..23 mapping;
- presentation state model;
- diagnosis candidates;
- completion predicate;
- actual screenshot paths/dimensions;
- controller-level Godot test evidence;
- accepted-runs proof;
- restart/error-path proof;
- full final push/PR CI URLs;
- deviations/blockers;
- unauthorized environment mutation status;
- final state READY FOR SOL/PRO REVIEW / BLOCKED / STOP.

Then post only concise issue #24 and PR #25 receipts and STOP.

No merge.
No WP-012.

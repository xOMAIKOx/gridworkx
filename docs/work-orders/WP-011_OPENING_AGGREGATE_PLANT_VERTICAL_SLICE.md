# GRIDWORKS — WP-011 Opening Aggregate-Plant Vertical Slice

**Status:** AUTHORIZED DESIGN / CONTROLLING WORK ORDER  
**Technical Authority:** GPT-5.6 Sol  
**Engineering route:** Devin + GPT-5.6 Luna XHigh  
**Implementation class:** D2 — bounded complex  
**Exact Engineering parent:** `ca5b4f06cdba226176296dd1ee7d294ef6c7927a`

## 1. Purpose

WP-011 delivers the first genuinely playable GRIDWORKS loop.

It turns the accepted WP-002–WP-010 foundations into an advisor-guided aggregate-plant recovery experience that teaches:

**Observe → Diagnose → Decide → Act → Observe consequences → Adapt**

The player inherits a barely functioning aggregate plant, inspects evidence, distinguishes an immediate critical fault from visible non-critical degradation, performs a minimal intervention, recovers a small amount of useful throughput, and produces the first real material output.

This is a production-quality vertical slice foundation, not a throwaway tutorial mock.

## 2. Controlling product intent

The slice must preserve the Product Philosophy / DP2 laws:

- inspect before spending;
- failures are causal and diagnosable;
- information exists before the decision;
- wrong choices are allowed and recoverable;
- the UI must not simply reveal the answer;
- dependencies matter;
- partial operation is valuable;
- the player should celebrate the first useful output;
- advisor assistance teaches reasoning but does not become the simulation authority.

The desired first-session teaching sequence is:

1. inspect before acting;
2. distinguish critical failure from non-critical wear;
3. reason about dependencies;
4. choose a minimal intervention;
5. recover approximately 2–5% useful throughput;
6. observe real material flow;
7. leave the player with a meaningful next decision.

## 3. Exact ancestry

Engineering branches from exactly:

`ca5b4f06cdba226176296dd1ee7d294ef6c7927a`

This is the owner-authorized WP-010 merge commit.

The later WP-011 work-order/directive commit on `main` is intentionally not the Engineering implementation parent.

No silent rebase.

## 4. Existing canonical foundations

WP-011 must reuse rather than replace:

- WP-003 facility/system/component graph;
- `aggregate_plant_fixture()`;
- WP-004 causal failure/inspection/diagnosis/intervention model;
- `aggregate_fault_definitions()`;
- WP-005 material/recipe/inventory model;
- `aggregate_material_fixture()`;
- WP-010 `GridworksSimBridge`;
- WP-010 `GridworksSimBridgeClient`;
- committed deterministic snapshot/digest behavior;
- Godot 4.7.2 + Rust 1.80.0 + exact `godot = "=0.2.4"`.

## 5. Opening scenario — canonical state

Add one canonical scenario builder in Rust, preferably in a dedicated scenario module:

`scenario.opening.aggregate`

The scenario must be created by canonical Rust code, not by GDScript editing snapshot JSON.

Preferred public API equivalent:

`opening_aggregate_scenario(seed) -> Result<SimulationState, KernelError>`

The scenario state must contain:

- `facility.aggregate_plant_fixture`;
- `aggregate_fault_definitions()`;
- `aggregate_material_fixture()`;
- initial raw-feed and limestone inventories from the accepted fixture;
- deterministic seed supplied by the caller.

### 5.1 Critical fault

Activate:

- fault instance: `fault.instance.opening.feed_motor`
- fault type: `fault.motor_bearing_seizure`
- component: `component.feed_conveyor`
- severity: canonical full seizure

Initial facility effective capacity must be zero.

### 5.2 Non-critical visible fault

Activate a visible but non-immediate fault on the output side:

- fault instance: `fault.instance.opening.output_belt`
- fault type: `fault.worn_belt`
- component: `component.output_conveyor`

This is intended to tempt premature repair but must not be the immediate cause of the zero-throughput state.

### 5.3 Partial-recovery bottleneck

Configure an accepted graph-level constrained component state so that after resolving the seized feed conveyor, the facility recovers approximately **3% effective capacity**, not 100%.

Preferred implementation:
- constrain `component.crusher` to the accepted graph-equivalent of 3% local capacity.

Do not invent new failure mathematics merely to hit the onboarding number.

The exact resulting facility evaluation must be pinned by tests.

## 6. Scenario creation bridge

Extend the existing GDExtension with one coarse scenario creation operation, preferred:

`create_scenario(scenario_id, seed)`

Requirements:

- accepted ID: `scenario.opening.aggregate`;
- unknown scenario ID -> stable boundary error;
- positive seed only;
- returns the same versioned bridge envelope as WP-010;
- result contains canonical snapshot + digest;
- scenario construction remains in `gridworks-sim`, not the Godot wrapper.

Do not add one bridge method per tutorial step.

## 7. Presentation architecture

Add an opening-slice scene and controller under `apps/game`.

Preferred structure:

- `scenes/opening_aggregate.tscn`
- `scripts/opening_aggregate_controller.gd`
- bounded UI/view scripts/components as needed
- optional data-driven onboarding copy/content under `apps/game/content/onboarding/`

The controller owns presentation flow only.

Canonical simulation state remains the bridge snapshot.

Do not duplicate facility capacity, fault, intervention, production or inventory calculations in GDScript.

## 8. Presentation state model

Use stable semantic presentation states equivalent to:

- `opening.intro`
- `opening.inspect`
- `opening.diagnose`
- `opening.intervene`
- `opening.produce`
- `opening.complete`

Progression should be driven by canonical snapshot/events/results where practical, not a hidden scripted answer flag.

The scenario must survive a wrong-but-valid intervention and remain recoverable.

## 9. Advisor

Include one advisor/foreman panel.

The advisor teaches reasoning, not answers by default.

Required intent:

- opening: “inspect before spending”;
- visible worn belt: prompt player to ask whether it is actually stopping the line;
- after evidence points to the feed conveyor: explain that immediate failure and general wear are different;
- after partial restart: emphasize that imperfect production creates options.

Exact copy may be refined by Engineering, but it must remain concise and operational.

Do not implement a chatbot/LLM advisor in WP-011.

## 10. Facility schematic

Create a readable 2D/2.5D process-line schematic representing:

feed stockpile → feed conveyor → crusher → screen → output conveyor → finished stockpile.

Each component must expose at least:

- display name;
- operational/constrained/unavailable state;
- effective capacity or useful derived status supplied by canonical evaluation;
- relevant observed evidence after inspection;
- whether an unresolved fault is known/diagnosed where canonical state supports it.

The schematic is presentation only.

No GDScript dependency calculation.

## 11. Inspection

The player must be able to inspect at least:

- `component.feed_conveyor`;
- `component.output_conveyor`;
- `component.crusher`;
- `component.screen`.

Inspection of faulted components must execute canonical `FailureCommand::Inspect` commands through `execute_command_batch`.

Use stable semantic command IDs/idempotency keys.

Do not fabricate telemetry/evidence in GDScript.

The UI may format canonical evidence into readable labels.

## 12. Diagnosis

The player must be able to make a bounded diagnosis decision based on evidence.

At minimum expose candidate concepts sufficient to distinguish:

- motor bearing seizure;
- worn belt;
- screen blockage.

The diagnosis action must execute canonical `FailureCommand::Diagnose`.

A diagnosis may be unresolved/low-confidence if canonical evidence does not support it.

Do not mark a diagnosis “correct” solely because the UI knows the scenario answer.

## 13. Intervention

Expose a small set of canonical interventions appropriate to the opening:

- repair critical feed conveyor fault;
- repair/replace the worn output belt where permitted.

Interventions must execute canonical `FailureCommand::Intervene`.

Required behaviors:

- repairing the worn belt first does not restart the plant;
- repairing the seized feed conveyor resolves the critical fault;
- after critical repair, effective facility capacity becomes the pinned partial-recovery value (~3%);
- no action mutates canonical state if the bridge/kernel rejects it.

No repair cost/material economy is invented in WP-011 because repair-cost settlement is not yet an accepted system.

## 14. Production

Once useful capacity exists, let the player run the accepted aggregate recipe through canonical `MaterialCommand::ExecuteProduction`.

Use:

- recipe: `recipe.aggregate_crush`
- facility: `facility.aggregate_plant_fixture`
- input inventory: `inventory.aggregate_feed`
- output inventory: `inventory.aggregate_finished`

Request enough runs that facility capacity is the limiting factor for the first restart.

The canonical material engine determines:
- accepted runs;
- consumed inputs;
- output quantity;
- limit reason.

The UI must show real resulting inventory/material flow.

No fake output counter.

## 15. Completion

The opening recovery task completes only when canonical state proves:

- critical feed-conveyor seizure is resolved;
- facility has non-zero useful effective capacity;
- at least one canonical production event has executed;
- finished aggregate inventory is greater than zero.

Completion is presentation state, but its predicate derives from canonical state/events.

## 16. What WP-011 does NOT fake

Do not display fake:

- cash/revenue;
- contract completion;
- market sale;
- XP/skill gain;
- manager progression;
- player/company balance;
- leaderboard score.

Those belong to later WPs.

The advisor may say the plant now has “useful output” or “options”; it must not claim money was earned unless an accepted economic transaction actually exists.

## 17. UI requirements

The vertical slice must be playable, not a debug form.

Required layout:

- facility title/status;
- process schematic;
- selected-component detail;
- evidence/diagnosis panel;
- intervention/action panel;
- throughput/status summary;
- raw input + finished output inventory summary;
- advisor panel;
- clear progress/completion feedback.

Touch targets must be mobile-appropriate.

Support at minimum:
- portrait/mobile-sized layout around 390×844;
- desktop/wide development layout around 1440×900.

Do not build world map, portfolio map or later operational dashboards in this WP.

## 18. Accessibility / interaction

At minimum:

- keyboard/focus navigation for actionable controls;
- readable labels independent of color;
- disabled controls have explanatory state;
- status is not conveyed by color alone;
- no tiny icon-only critical actions;
- advisor can be advanced/dismissed without blocking canonical state.

## 19. Scenario reset

Provide a deterministic “Restart opening scenario” developer/test path.

Restart must:
- create a fresh scenario with the configured seed;
- clear presentation-only transient state;
- restore the exact canonical initial digest.

Do not mutate a used snapshot backward by hand.

## 20. Determinism

For a fixed seed and fixed command sequence:

- initial digest is stable;
- inspection/diagnosis/intervention events are stable;
- partial-recovery digest is stable;
- first-production digest and inventories are stable.

Commit a WP-011 deterministic golden fixture that includes:
- seed;
- scenario ID;
- action/command sequence;
- initial digest;
- post-critical-repair digest;
- final production digest;
- expected effective capacity;
- expected accepted production runs;
- expected finished aggregate quantity.

Pure Rust and Godot/GDExtension must assert the same golden.

## 21. Rust scenario tests

Required:

- scenario builds and validates;
- initial effective capacity = 0;
- critical + non-critical faults exist on intended components;
- wrong/non-critical repair first leaves capacity = 0;
- critical repair yields pinned ~3% capacity;
- first production is facility-capacity-limited;
- canonical finished aggregate quantity matches golden;
- repeated seed/sequence produces identical digest/events/inventory;
- scenario unknown ID bridge path rejects safely.

## 22. Godot integration tests

Add a dedicated headless WP-011 test scene/script that executes the real opening controller/bridge path.

Required:

- scene instantiates;
- bridge available;
- opening scenario created from Rust;
- process component IDs match canonical fixture;
- inspect action changes canonical evidence state;
- wrong repair path remains recoverable;
- critical repair moves throughput 0 -> pinned partial value;
- production creates finished aggregate inventory;
- completion predicate becomes true only after production;
- restart returns to initial digest;
- no simulation fallback if bridge unavailable;
- malformed/failed bridge result does not corrupt local snapshot.

## 23. UI structural tests

Add bounded structural checks for:

- opening scene/controller paths exist;
- controller uses `GridworksSimBridgeClient`;
- no local facility-capacity/bottleneck formula exists in GDScript;
- no fake credit/revenue/XP constant or mutation exists;
- semantic component IDs are referenced through scenario/view model rather than duplicated arbitrary state;
- WP-011 test scene is included in CI.

Avoid brittle source-string checks where executable tests are stronger.

## 24. Visual evidence

Capture repository-owned evidence for both:

- mobile/portrait opening state;
- desktop/wide opening state;
- inspected critical component;
- partial-recovery/first-output state.

Use Godot-rendered screenshots or deterministic UI capture where practical.

Do not commit huge raw recordings.

Evidence is for review; it is not a substitute for executable tests.

## 25. Content/data discipline

Scenario gameplay state belongs in Rust.

Presentation copy/config may be data-driven.

Use semantic IDs for advisor steps and component presentation labels.

Do not build the full WP-021 localization system, but avoid scattering opaque text/state literals across many scripts.

## 26. Existing shell

Preserve the accepted presentation shell architecture.

The opening slice may become the default development/start scene for WP-011 if done cleanly, but do not remove the world/portfolio/site/operational presentation-layer concepts from the shell baseline.

## 27. Error handling

If native bridge/scenario creation fails:

- show an explicit non-authoritative error state;
- do not silently substitute mock simulation;
- do not fabricate a successful scenario;
- provide a developer/test retry/restart path.

## 28. Persistence/network boundary

WP-011 is local deterministic gameplay integration.

Do not add:

- API calls;
- WebSocket;
- PostgreSQL;
- SQLite save system;
- NATS;
- authentication;
- account/company fetch;
- server clock integration.

Those require later explicit WPs.

## 29. Economy boundary

WP-006 ledger and WP-009 API remain intact, but WP-011 does not invent public economy endpoints or local wallet logic.

First material output is sufficient for this vertical slice.

## 30. No later-WP scope

Do not implement:

- WP-012 skills/managers;
- WP-013 recruitment;
- WP-014 contracts/Work Exchange;
- WP-015 marketplace/trade settlement;
- WP-016 lifecycle/opportunity;
- WP-017 sale/acquisition;
- WP-018 messaging;
- WP-019 consortium;
- WP-020 notifications;
- WP-021 full localization;
- WP-022 admin;
- WP-023 challenges;
- WP-024 leaderboards;
- WP-025 economy telemetry/fairness.

## 31. Allowed paths

Primary:

- `crates/gridworks-sim/**`
- `crates/gridworks-godot/**`
- `apps/game/**`
- `tests/deterministic/**`
- `tests/structural/**`
- `packages/schemas/**` if a bounded scenario/view schema is justified
- `.github/workflows/**`
- `Makefile` / repository validation scripts as needed
- `docs/evidence/**`
- `docs/engineering/handbacks/wp-011/**`

Do not modify WP-006–WP-009 persistence/API semantics.

## 32. Toolchain/runtime boundary

Keep:

- Rust 1.80.0;
- Godot 4.7.2;
- exact `godot = "=0.2.4"`;
- native/systemd estate policy.

Repository-only implementation.

No ERIS/host mutation, deployment, ports, proxy, DB provisioning or containers.

## 33. CI

Extend existing CI rather than creating pipeline theatre.

Required:

- all existing foundation/security/Rust/Godot WP-010 gates remain green;
- WP-011 Rust golden tests;
- WP-011 headless Godot controller/vertical-slice tests;
- existing Xvfb import discovery remains green;
- normal Godot project parse remains green;
- native extension parity remains green.

## 34. Evidence

Create:

`docs/evidence/WP-011_OPENING_AGGREGATE_VERTICAL_SLICE.md`

Document:

- exact scenario seed/ID;
- canonical initial conditions/faults;
- partial-recovery bottleneck;
- component/recipe/inventory IDs;
- UI scene/controller structure;
- advisor step IDs;
- deterministic golden path/values;
- wrong-choice recovery proof;
- critical repair proof;
- first material-output proof;
- Godot headless result;
- visual evidence paths;
- CI links;
- explicit deferred economy/skills/contracts scope;
- no-host/deployment statement.

## 35. Done-when

Architecture can independently verify:

1. exact parent preserved;
2. canonical Rust scenario builder exists;
3. initial plant is genuinely blocked by the feed-conveyor seizure;
4. visible non-critical wear is present;
5. player can inspect canonical evidence;
6. player can make a canonical diagnosis;
7. wrong-but-valid repair does not magically solve the plant;
8. critical repair resolves the zero-throughput fault;
9. throughput becomes a pinned small partial value near 3%;
10. production uses the accepted material engine;
11. first finished aggregate appears in canonical inventory;
12. completion derives from canonical state;
13. Godot contains no duplicate simulation math;
14. advisor teaches reasoning without becoming authority;
15. mobile and wide layouts are usable;
16. bridge failure has no fake fallback;
17. restart is deterministic;
18. Rust and Godot assert the same committed golden;
19. all existing CI gates remain green;
20. no fake economy/skills/contracts/later-WP scope exists;
21. evidence and durable Engineering handback are complete.

## 36. Stop

After handback, Engineering stops.

Do not merge WP-011.
Do not start WP-012 without owner authorization.

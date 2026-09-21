# GRIDWORKS — WP-004 Causal Failures, Diagnosis and Interventions

**Status:** AUTHORIZED FOR ENGINEERING  
**Architecture owner:** Architecture  
**Implementation model:** GPT-5.6 Luna XHigh  
**Review:** GPT-5.6 Sol-class Architecture review  
**Exact required parent:** `2924c544680b1814ab258c29ec92b64015cfe367`

## 1. Purpose

Implement the canonical deterministic failure, diagnosis and intervention layer on top of the accepted WP-003 Facility/System/Component graph.

WP-004 must establish the causal machinery that lets GRIDWORKS represent:
- component condition separate from operability;
- latent faults and observed symptoms;
- evidence and diagnostic confidence;
- player inspection/testing actions;
- repair / rebuild / replace / bypass / reroute / derate / shut-down style interventions;
- deterministic state transitions and consequences.

It must not implement recipes, resource inventories, economic settlement, contracts, skills/managers, UI, or production accounting.

## 2. Controlling references

Engineering must read and obey:
- `AGENTS.md`
- `AI_AGENT_COLLABORATION_PROTOCOL.md`
- Product Philosophy
- DP2 Core Systems Specification
- DP3 Engineering Architecture
- ADR-003 Schema/Serialization
- ADR-004 Content/Rules Distribution
- ADR-008 World Time
- accepted WP-002 deterministic kernel
- accepted WP-003 facility/component graph

If implementation conflicts with accepted architecture, stop that path and raise evidence in the WP-004 issue.

## 3. Exact ancestry

Branch from exactly:

`2924c544680b1814ab258c29ec92b64015cfe367`

Verify:
- exact parent;
- clean working tree;
- no unrelated changes.

Do not silently rebase to a newer `main` head.

## 4. Runtime / host constraints

Repository-only.

No additional ERIS host package installation is authorized by default.

Prohibited:
- Docker/Podman/Compose/OCI/Kubernetes;
- deployment;
- systemd installation;
- database/service provisioning;
- nginx/DNS/firewall/WireGuard changes;
- external port activation.

## 5. Core product invariants

Implement these directly:

> Every failure has a cause.

> Low skill or incomplete inspection means incomplete evidence/lower confidence, not fabricated telemetry.

> Condition is not the same as operability.

> Partial operation is valid and strategically meaningful.

> Mistakes matter, but systems remain recoverable.

> Information should exist before a decision.

WP-004 must preserve these principles in code and tests.

## 6. Condition model

Introduce a deterministic condition representation for components that is distinct from:
- current available/unavailable state;
- enabled/disabled state;
- effective capacity.

A component may be:
- worn but operational;
- degraded but still available;
- healthy but disabled;
- failed due to a specific fault;
- temporarily bypassed/derated.

Use deterministic fixed-point/integer representation.

Do not use authoritative floating-point condition values.

Document the chosen scale.

## 7. Fault model

A Fault must be an explicit causal object/state.

At minimum:
- stable fault instance ID;
- fault type ID;
- affected component ID;
- severity;
- latent/active/resolved status;
- deterministic onset/effective time context;
- causal effects;
- observable symptom references;
- rules/content version compatibility as required.

The kernel must not create “generic broken” state without an attributable fault where a failure is meant to exist.

## 8. Fault type/content contract

Support data-driven fault definitions with stable IDs such as:

`fault.motor_bearing_seizure`

Do not hardcode translated prose.

A fault definition should be able to express at least:
- applicable component types/categories;
- severity;
- capacity effect;
- availability effect;
- potential symptoms;
- inspection/test evidence;
- permitted interventions;
- deterministic resolution effects.

Keep the first implementation bounded and generic.

## 9. Symptoms and evidence

Model observable symptoms separately from underlying faults.

Examples:
- abnormal vibration;
- elevated temperature;
- no rotation;
- low flow;
- high current;
- blockage indication.

A symptom is evidence, not the diagnosis itself.

Required:
- stable symptom IDs;
- observed/not observed/unknown states;
- source/test context;
- deterministic timestamp/effective time if stateful;
- no invented observations.

Unknown must remain distinct from false/absent.

## 10. Diagnostic evidence model

Introduce deterministic evidence items produced by:
- visual inspection;
- measurement/test;
- system telemetry where valid;
- component state already observable by design.

Evidence must be factual from authoritative state.

Do not let a low-information actor receive fake values.

## 11. Diagnosis model

Support diagnostic hypotheses over known candidate faults.

At minimum:
- candidate fault type/instance reference;
- evidence used;
- confidence represented deterministically;
- unresolved/confirmed/rejected state as appropriate.

Confidence must be based on available evidence, not random fabricated certainty.

Use fixed-point basis points or another explicit integer scale.

Do not implement LLM-based diagnosis in authoritative gameplay.

## 12. Skill boundary

WP-004 does not implement the player skill system itself; that is WP-012.

However, the diagnostic model must expose a clean future input for:
- information reveal;
- confidence calibration;
- test interpretation quality.

For WP-004, use an explicit deterministic `diagnostic_capability` or similar proof parameter only if required by tests.

Do not create full player skill progression.

## 13. Inspections/tests

Implement generic deterministic actions for at least:

- inspect;
- test/measure.

Exact taxonomy may be content-driven.

An inspection/test:
- targets a valid component/facility context;
- consumes authoritative time input if appropriate;
- produces evidence;
- may reveal a symptom/fault indicator;
- must not mutate unrelated state;
- must be replayable/idempotent under the accepted kernel contract.

WP-004 does not need economic cost/resource consumption yet.

## 14. Interventions

Support generic deterministic interventions sufficient to prove the architecture:

- inspect;
- test;
- clean;
- calibrate;
- lubricate;
- patch;
- repair;
- rebuild;
- replace;
- bypass;
- reroute;
- derate;
- shut down;
- outsource;
- salvage/dismantle.

Not every intervention needs full deep mechanics in WP-004, but the intervention type contract must support this accepted list.

Minimum executable mechanics should cover:
- repair;
- replace;
- bypass;
- derate;
- shut down;
- restore/re-enable as appropriate.

Do not implement economic/resource costs yet.

## 15. Intervention preconditions

An intervention must validate:
- target exists;
- intervention is allowed for current component/fault state;
- required fault/evidence condition where appropriate;
- time does not move backwards;
- command/version/idempotency rules remain valid.

Invalid interventions return typed errors and do not mutate state.

## 16. Repair semantics

Repair must:
- target a known repairable fault;
- change fault state deterministically;
- update condition/capacity/availability as defined by rules;
- preserve audit/replay evidence;
- not magically improve unrelated components.

Repair does not necessarily mean “back to factory new.”

## 17. Replace semantics

Replace must:
- deterministically resolve applicable active faults on the replaced component;
- reset condition according to defined replacement baseline;
- restore availability/capacity subject to graph dependencies;
- preserve stable component identity unless Architecture explicitly approves component-instance identity split.

For WP-004, retain component ID stability.

## 18. Bypass/reroute semantics

Bypass/reroute must integrate with the WP-003 graph without destroying its structural determinism.

A bypass may:
- disable a dependency edge or mark a component/path as bypassed through authoritative state;
- preserve alternate-path behavior;
- alter achievable facility capacity.

Do not mutate graph topology in an ad hoc way that breaks snapshot/replay.

Prefer explicit state/config around active path/dependency behavior.

## 19. Derate semantics

Derate must allow a component/system/facility to continue at intentionally reduced capacity.

This is a key partial-operation strategy.

Use deterministic fixed-point capacity limits.

A derated asset remains distinguishable from one constrained by a fault.

## 20. Shutdown semantics

Shutdown must allow deliberate safe unavailability:
- component or bounded system/facility target;
- explicit operator intervention;
- not misclassified as a fault.

This reinforces condition vs operability separation.

## 21. Recovery semantics

The model must support recovery to service after:
- repair;
- replacement;
- bypass/reroute;
- removal of deliberate shutdown;
- derate adjustment.

Recovery must remain causal and deterministic.

## 22. Aggregate-plant proof scenarios

Extend the aggregate fixture with failure/diagnostic/intervention proofs.

Required scenarios:

### A. Seized feed conveyor motor
Accepted onboarding fault:
- motor bearing seizure is active;
- feed conveyor cannot operate;
- downstream plant output becomes zero;
- fault has causal symptom/evidence;
- repair or replacement restores bounded production.

### B. Worn belt but still operational
- condition is degraded;
- component remains available;
- no false “failed” classification;
- capacity/reliability hook may be constrained without WP-005 production logic.

### C. Partially blocked screen
- plant remains operational;
- facility throughput is materially reduced;
- bottleneck evidence points to the affected path/component.

### D. Bypass/alternate path
Where structurally valid, bypass one component/path and prove bounded degraded operation rather than total failure.

### E. Deliberate shutdown
- healthy component can be intentionally unavailable;
- no fault is fabricated.

## 23. Fault generation

WP-004 should support deterministic activation of known faults through commands/test fixtures.

Do not build a broad stochastic wear/failure scheduler yet unless necessary to prove seeded deterministic onset.

If random fault activation is included for proof:
- use WP-002 seeded RNG only;
- same seed/state/time -> same fault result;
- no OS entropy.

Do not implement large-scale degradation scheduling beyond WP-004 scope.

## 24. Causal propagation

Fault effects must flow through WP-003 graph semantics.

Examples:
- seized motor sets affected component effective availability/capacity;
- HARD dependencies propagate outage downstream;
- partial blockage constrains CAPACITY path;
- bypass can restore alternate path.

Do not duplicate graph evaluation logic in the failure module.

## 25. Canonical state

Add authoritative state needed for:
- component condition;
- active/resolved fault instances;
- symptom/evidence state;
- intervention state such as deliberate shutdown, derate and bypass where required.

All must:
- serialize canonically;
- participate in snapshots/digests;
- remain insertion-order deterministic;
- validate strictly.

## 26. Schema contracts

Add/update canonical schemas under `packages/schemas` for:
- condition/fault state;
- diagnosis/evidence/intervention payloads if externally meaningful;
- updated simulation state.

Rust and schema validation must remain aligned.

No duplicate schema authority.

## 27. Typed errors

At minimum:
- unknown component/fault/evidence target;
- duplicate fault instance;
- incompatible fault/component type;
- invalid condition;
- invalid severity/confidence;
- unsupported intervention;
- intervention precondition failure;
- already resolved fault;
- invalid derate;
- invalid bypass target/path;
- invalid shutdown/recovery transition;
- arithmetic overflow;
- malformed command/version/time/idempotency through existing kernel paths.

Normal bad input must not panic.

## 28. Determinism

No:
- unordered hash iteration;
- floating-point authoritative math;
- wall clock;
- filesystem;
- network;
- database;
- thread-scheduling-dependent outcomes.

Canonical ordering must extend to all new vectors/state.

Logically equivalent reordered fault/evidence state must serialize/digest identically.

## 29. Command integration

Extend the accepted command model cleanly for WP-004 proof actions.

Likely commands:
- activate fault;
- inspect component;
- perform test;
- apply intervention;
- set/clear deliberate shutdown;
- set derate.

Exact names/types are Engineering-owned.

Requirements:
- versioned;
- externally supplied authoritative time;
- idempotent/replay-safe;
- typed result events;
- no external side effects.

## 30. Event model

Emit deterministic events for important transitions, for example:
- fault activated;
- evidence observed;
- diagnosis updated;
- intervention applied;
- fault resolved;
- component shutdown/restored;
- derate changed;
- bypass changed.

Events must use stable semantic identifiers, not translated prose.

## 31. No economic/resource settlement

WP-004 interventions do not yet consume:
- spare parts;
- money/credits;
- labour contracts;
- manager time;
- inventory.

Those belong to later work packages.

Keep hooks/identifiers extensible, but do not implement settlement.

## 32. Testing

Required tests include:

### Condition/operability
- worn-but-operational;
- healthy-but-shutdown;
- failed-due-to-fault;
- derated-but-operational.

### Faults
- activation;
- duplicate rejection;
- incompatible target rejection;
- resolution.

### Evidence/diagnosis
- unknown vs observed vs absent;
- deterministic evidence reveal;
- confidence bounded and deterministic;
- no fabricated evidence.

### Interventions
- repair success;
- replace success;
- invalid repair rejection;
- bypass/alternate path;
- derate;
- shutdown/recovery.

### Aggregate scenarios
- seized feed conveyor motor;
- worn belt;
- blocked screen partial operation;
- alternate/bypass recovery.

### Determinism
- reordered new state -> identical canonical JSON/digest;
- replay identical;
- negative replay/version/idempotency cases remain valid.

## 33. CI

All existing gates remain green:
- `cargo fmt --check`;
- `cargo test --workspace`;
- `cargo clippy --workspace --all-targets --all-features -- -D warnings`;
- schema/content validation;
- secret scan;
- forbidden-runtime scan;
- Godot parse gate.

## 34. Complexity / performance

Document algorithmic complexity for:
- fault lookup;
- evidence lookup;
- intervention resolution;
- graph re-evaluation after state change.

No external rule engine.

No premature event-sourcing framework.

## 35. Architecture constraints carried forward

From WP-002:
- do not allow `executed_commands` to become unbounded production history.

From WP-003:
- explicit facility output topology remains authoritative;
- multiple facility outputs currently represent alternate authoritative routes using maximum effective output, not additive multi-product throughput;
- canonical graph ordering must remain stable.

Do not silently reinterpret these contracts.

## 36. Allowed paths

Primary:
- `crates/gridworks-sim/**`
- `crates/gridworks-sim/tests/**`
- `tests/deterministic/**`
- `packages/schemas/**`
- `packages/content/**` for generic fault/symptom/intervention definitions/fixtures
- Rust workspace/lock files if necessary
- WP-004 evidence/architecture docs
- CI/build files only if required

Minimal supportive changes elsewhere must be explained.

## 37. Out of scope

Do not implement:
- production recipes;
- inventories;
- material flows;
- economic ledger;
- marketplace/trade;
- contracts/Work Exchange;
- player skills/managers progression;
- company/ownership gameplay;
- real estate gameplay;
- localization UI text;
- Godot GDExtension;
- deployment;
- live services.

## 38. Evidence

Create concise evidence documenting:
- condition model;
- fault model;
- symptom/evidence model;
- diagnosis confidence representation;
- intervention semantics;
- aggregate scenarios;
- canonical ordering;
- deterministic results;
- complexity;
- schema changes;
- dependencies.

## 39. Done-when

Architecture can verify:

1. exact parent;
2. condition separate from operability;
3. causal fault instances exist;
4. symptoms/evidence separate from faults;
5. unknown evidence remains distinct from absent;
6. deterministic diagnosis/confidence;
7. repair/replace/bypass/derate/shutdown semantics work;
8. aggregate onboarding scenarios work;
9. fault effects propagate through WP-003 graph rather than duplicated logic;
10. canonical serialization/digest remains insertion-order stable;
11. no fabricated telemetry/evidence;
12. no WP-005+ scope creep;
13. CI/security gates pass;
14. complete GitHub handback exists.

## 40. Handback

Post to the WP-004 issue and draft PR:
- branch;
- exact parent/head;
- commit chain;
- files/dependencies;
- condition/fault/evidence/diagnosis design summary;
- intervention semantics;
- aggregate fixture evidence;
- schema changes;
- deterministic/replay evidence;
- CI links/results;
- deviations/risks/conflicts;
- explicit no-host/deployment statement.

## 41. Stop condition

After handback, Engineering stops.

Do not start WP-005 without separate authorization.

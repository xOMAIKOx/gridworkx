# GRIDWORKS — WP-005 Resources, Recipes and Inventory

**Status:** AUTHORIZED FOR ENGINEERING  
**Architecture owner:** Architecture  
**Implementation model:** GPT-5.6 Luna XHigh  
**Review:** GPT-5.6 Sol-class Architecture review  
**Exact required parent:** `70c67fe987abcd045584fee75d17836c28684a85`

## 1. Purpose

Implement the deterministic resource, recipe, inventory and material-transformation layer on top of the accepted simulation kernel, facility graph and causal failure model.

WP-005 establishes the authoritative machinery for:
- resource identities and grades;
- inventory locations/stores;
- quantities and capacity limits;
- deterministic recipe definitions;
- input consumption and output production;
- bounded production execution;
- shortages, partial runs and capacity constraints;
- material-state serialization and replay.

WP-005 must **not** implement monetary settlement, double-entry accounting, marketplace trading, contracts, company finance, logistics economics or player-facing UI.

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
- accepted WP-004 causal failure/diagnosis/intervention layer

If implementation conflicts with accepted architecture, stop that path and raise evidence in the WP-005 issue.

## 3. Exact ancestry

Branch from exactly:

`70c67fe987abcd045584fee75d17836c28684a85`

Verify:
- exact parent;
- clean working tree;
- no unrelated changes.

Do not silently rebase onto a newer `main` head.

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

## 5. Product invariants

Implement these directly:

> The world behaves causally.

> Material cannot be created or destroyed except through an explicit accepted recipe/process rule.

> Partial operation remains useful.

> Production is constrained by physical inputs, facility capacity and recipe rules.

> Authoritative quantities use deterministic integer/fixed-point arithmetic.

> Gameplay content uses stable semantic IDs rather than translated prose.

## 6. Resource model

Introduce a canonical resource definition/state model.

At minimum support:
- stable resource ID, e.g. `resource.limestone`;
- resource category/type;
- canonical quantity unit;
- grade/quality dimension where applicable;
- optional density/conversion metadata only if required by current recipes;
- deterministic content/rules version compatibility.

Do not embed localized names/descriptions as authoritative identifiers.

## 7. Quantity representation

Use deterministic fixed-width integer arithmetic.

Select and document a canonical quantity representation suitable for later large industrial volumes.

Requirements:
- no authoritative floating point;
- explicit unit/scaling convention;
- checked arithmetic;
- deterministic rounding rules where conversion is unavoidable;
- typed overflow/underflow errors.

Do not silently truncate quantities.

## 8. Resource grades / quality

The architecture must permit materially distinct grades/qualities where they matter.

Examples:
- ore grade;
- product specification;
- purity;
- moisture;
- quality class.

WP-005 does not need a universal chemistry engine.

Use an explicit bounded representation such as:
- semantic grade ID;
- integer quality basis points;
- or a small typed property set.

Document the chosen rule.

Resources of incompatible grade/specification must not be silently merged.

## 9. Inventory model

Implement authoritative inventory stores.

At minimum:
- stable inventory/store ID;
- facility/site ownership/reference;
- permitted resource constraints where applicable;
- capacity;
- current resource lots/stacks/balances;
- deterministic ordering;
- validation.

Inventory may represent:
- stockpile;
- silo;
- tank;
- warehouse/store;
- work-in-process buffer;
- finished-goods storage.

Keep this generic and data-driven.

## 10. Inventory lots vs aggregate balances

Choose a bounded model that can later support:
- grade/spec differences;
- provenance where needed;
- spoilage/age later if required.

For WP-005, either:
- explicit lots; or
- aggregate balances keyed by resource+grade/spec.

The model must not force all future resources into one undifferentiated scalar.

Document the choice.

## 11. Capacity

Inventory stores must have deterministic capacity constraints.

Support:
- total capacity;
- optional resource-specific capacity if required;
- rejection or bounded partial transfer when destination capacity is insufficient.

Do not silently overfill stores.

## 12. Recipe model

Introduce data-driven deterministic recipe definitions.

At minimum:
- stable recipe ID;
- applicable facility/system/component/process type;
- input resources and required quantities;
- output resources and produced quantities;
- optional by-products;
- nominal batch/run basis;
- minimum operating/capacity constraints where relevant;
- rules/content version.

Examples later:
- limestone feed -> crushed aggregate;
- ore -> concentrate;
- grain -> processed food;
- crude -> refined products.

Keep WP-005 generic; the aggregate fixture is the proof vertical.

## 13. Conservation / transformation

Every material transition must be attributable to a recipe or explicit transfer.

A recipe execution must:
1. validate required inputs;
2. calculate executable amount;
3. consume the exact accepted input amount;
4. produce deterministic outputs;
5. respect destination capacity;
6. emit deterministic events;
7. leave no partial mutation on rejected execution.

Where a recipe intentionally changes mass/volume due to waste, yield or conversion, that relationship must be explicit in the recipe definition.

No hidden material creation.

## 14. Yield

Support deterministic yield/loss.

Use fixed-point basis points or another explicit integer scale.

Examples:
- 10000 = 100% yield;
- 9500 = 95%;
- 300 = 3%.

If output yield is less than input-equivalent material, the missing quantity must be represented as explicit waste/loss/by-product where the recipe definition requires material-accounting visibility.

Do not invent generalized thermodynamic accounting beyond scope.

## 15. Production execution

Implement a deterministic production/run operation that ties together:
- recipe;
- source inventory/input stores;
- destination output stores;
- facility effective capacity from WP-003/WP-004 projection;
- authoritative elapsed/effective time where required.

Production must not duplicate graph/failure evaluation.

Use the existing facility evaluation as an input to determine bounded executable production.

## 16. Time basis

WP-005 may support production over an explicit authoritative duration or a discrete bounded run/batch.

Requirements:
- no wall clock;
- no client/device time;
- no hidden time advancement;
- explicit authoritative elapsed time if time-based;
- deterministic replay.

Do not implement the full offline catch-up policy yet.

## 17. Partial production

Partial production is a first-class outcome.

If:
- facility capacity is 3%;
- only part of required inputs are present;
- destination capacity is limited;

then production may execute a bounded partial amount if the recipe allows it.

Do not collapse partial production into failure.

Return deterministic evidence explaining the limiting factor.

## 18. Production limiting factors

Expose typed limiting reasons such as:
- facility capacity;
- missing input;
- insufficient input quantity;
- destination capacity;
- recipe minimum threshold;
- deliberate shutdown/fault-derived zero capacity;
- invalid grade/spec.

Do not emit user-facing prose from the kernel.

## 19. Transfers

Implement deterministic inventory transfer semantics:
- source store;
- destination store;
- resource/grade;
- requested quantity;
- accepted quantity;
- rejected/partial result;
- typed reason.

Transfers must:
- never create material;
- never drive source below zero;
- never overfill destination;
- preserve grade/spec.

WP-005 does not implement transport time/cost/logistics routing.

## 20. Aggregate-plant proof chain

Extend the accepted aggregate fixture to prove:

`raw feed / limestone → aggregate processing → finished aggregate`

At minimum include:
- raw-feed stockpile inventory;
- finished-goods stockpile inventory;
- one deterministic aggregate recipe;
- facility effective capacity as a production constraint;
- input consumption;
- output creation;
- partial output when the screen/facility is constrained;
- zero output when the feed conveyor seizure drives facility capacity to zero.

Do not implement sales or credits yet.

## 21. Required scenarios

### A. Healthy production
- sufficient feed;
- sufficient destination capacity;
- healthy facility;
- exact deterministic input consumption/output creation.

### B. Input shortage
- recipe requests more than source inventory;
- bounded partial production or typed rejection according to recipe contract;
- no negative inventory.

### C. Destination full
- output store capacity constrains execution;
- no overfill.

### D. Fault-constrained production
- use WP-004 blocked-screen/seized-motor state;
- production amount tracks evaluated facility capacity;
- no duplicate failure math.

### E. Grade mismatch
- incompatible input grade/spec rejected or excluded;
- no silent merging/substitution.

### F. Transfer conservation
- source decrease equals destination increase for direct transfer;
- total conserved.

## 22. Canonical state

Extend `SimulationState` with authoritative material state as required.

All new state must:
- serialize canonically;
- participate in snapshots/digests;
- remain stable under insertion-order differences;
- validate strictly.

Canonicalize:
- resource definitions;
- recipes;
- inventories;
- inventory balances/lots;
- transfers/production state only if persisted.

Do not canonicalize away semantically meaningful sequence where chronology matters.

## 23. Commands

Extend the accepted command envelope cleanly.

Likely proof commands:
- register/load material content state where appropriate;
- transfer inventory;
- execute production run.

Exact names/types are Engineering-owned.

Requirements:
- versioned;
- explicit effective time;
- idempotent/replay-safe;
- typed events/results;
- no external side effects.

Do not expose internal helper calculations as public commands unnecessarily.

## 24. Events

Emit deterministic events for:
- inventory transfer;
- production executed;
- inputs consumed;
- outputs produced;
- partial/rejected production where useful.

Events use stable IDs and quantities, not translated prose.

## 25. Schemas/content

Add/update canonical schemas for:
- resources;
- recipes;
- inventory/material state;
- production/transfer command payloads;
- simulation state.

Add minimal data-driven aggregate resource/recipe definitions in `packages/content`.

Actual content files must be exercised by validation, not merely have schemas registered.

Rust and schemas must remain aligned.

## 26. Typed errors

At minimum:
- unknown resource;
- unknown recipe;
- unknown inventory;
- incompatible resource/grade;
- insufficient quantity;
- destination capacity exceeded;
- invalid quantity/unit;
- duplicate IDs;
- invalid recipe;
- invalid transfer;
- invalid production target;
- facility unavailable;
- arithmetic overflow/underflow;
- schema/version/idempotency errors through existing kernel paths.

Rejected commands must not mutate state.

## 27. Determinism

No:
- float authoritative math;
- unordered hash iteration;
- system time;
- filesystem;
- network;
- database;
- thread-dependent outcomes.

Use stable ordering and checked fixed-width arithmetic.

Equivalent material states supplied in different insertion order must serialize/digest identically.

## 28. Failure integration

WP-004 failure state constrains production only through accepted facility evaluation/projection.

Do not duplicate:
- fault logic;
- condition logic;
- capacity propagation.

Example:
- screen blockage -> WP-004 projects facility state -> WP-003 computes effective capacity -> WP-005 uses that capacity to limit recipe execution.

## 29. No economics yet

WP-005 must not introduce:
- Credits;
- prices;
- cost basis;
- revenue;
- double-entry ledger;
- invoices;
- contracts;
- market orders;
- valuation.

Those belong to WP-006 and later packages.

Material events should be structured so WP-006 can later journal economically relevant consequences without rewriting production semantics.

## 30. No logistics economics/routing

Transfers are local deterministic inventory operations only.

Do not implement:
- trucks;
- rail;
- ships;
- route times;
- freight cost;
- cargo risk;
- transport contracts.

Those remain later scope.

## 31. Carried architecture constraints

From WP-002:
- do not turn `executed_commands` into unbounded production history.

From WP-003:
- explicit facility outputs remain authoritative;
- multiple facility outputs currently mean alternate routes, not additive multi-product throughput.

From WP-004:
- low information cannot leak hidden state;
- fault severity is stored but effect magnitude remains definition-driven unless explicitly changed;
- intervention semantics remain causal.

Do not silently reinterpret these contracts.

## 32. Testing

Required:

### Resource/inventory
- valid resource definitions;
- duplicate/invalid IDs rejected;
- quantity arithmetic;
- capacity enforcement;
- grade/spec separation;
- canonical ordering.

### Transfers
- full transfer;
- partial/rejected transfer;
- source underflow prevention;
- destination overflow prevention;
- conservation.

### Recipes
- valid recipe;
- missing input;
- insufficient input;
- grade mismatch;
- output capacity limit;
- deterministic yield/by-product handling.

### Production
- healthy aggregate run;
- partial run;
- blocked-screen constrained run;
- seized-motor zero run;
- no input mutation on rejected command.

### Determinism
- reordered material state -> identical canonical JSON/digest;
- replay identical;
- schema fixtures valid/invalid.

## 33. CI

All existing gates remain green:
- `cargo fmt --check`;
- `cargo test --workspace`;
- `cargo clippy --workspace --all-targets --all-features -- -D warnings`;
- schema/content validation;
- secret scan;
- forbidden-runtime scan;
- Godot parse gate.

## 34. Performance posture

Avoid obviously pathological behavior.

Document complexity for:
- resource lookup;
- inventory lookup;
- transfer;
- recipe execution;
- production limiting calculation.

No external rules engine.

No premature database/index design.

## 35. Allowed paths

Primary:
- `crates/gridworks-sim/**`
- `crates/gridworks-sim/tests/**`
- `tests/deterministic/**`
- `packages/schemas/**`
- `packages/content/**`
- Rust workspace/lock files if required
- WP-005 evidence/architecture docs
- CI/build files only if required

Minimal supportive changes elsewhere must be explained.

## 36. Out of scope

Do not implement:
- PostgreSQL persistence;
- economic ledger;
- pricing;
- trade/marketplace;
- contracts;
- company finances;
- player skills/managers;
- logistics routing/economics;
- UI/Godot GDExtension;
- deployment;
- live services.

## 37. Evidence

Create concise evidence documenting:
- quantity/unit representation;
- resource/grade model;
- inventory model;
- recipe model;
- yield/conservation rules;
- production limiting logic;
- aggregate proof chain;
- canonical ordering;
- complexity;
- schemas/content;
- deterministic results.

## 38. Done-when

Architecture can verify:

1. exact parent;
2. deterministic resource/quantity model;
3. grade/spec separation;
4. bounded inventory capacity;
5. deterministic transfers with conservation;
6. data-driven recipes;
7. explicit material transformation/yield;
8. facility/failure capacity constrains production through existing evaluators;
9. partial production works;
10. aggregate fixture consumes feed and produces finished aggregate;
11. insertion-order-stable serialization/digest;
12. no economic/ledger/logistics scope creep;
13. all CI/security gates pass;
14. complete GitHub handback exists.

## 39. Handback

Post to the WP-005 issue and draft PR:
- branch;
- exact parent/head;
- commit chain;
- files/dependencies;
- quantity/unit model;
- resource/grade design;
- inventory/recipe design;
- transfer and production semantics;
- aggregate proof evidence;
- schema/content changes;
- determinism/replay evidence;
- CI links/results;
- deviations/risks/conflicts;
- explicit no-host/deployment statement.

## 40. Stop condition

After handback, Engineering stops.

Do not start WP-006 without separate authorization.

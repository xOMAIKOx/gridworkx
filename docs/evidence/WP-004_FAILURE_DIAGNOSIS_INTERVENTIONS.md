# WP-004 Failure, Diagnosis and Intervention Evidence

## Condition and operability

`ComponentConditionState` stores condition and derate in integer basis points (`10000 = 100%`) separately from WP-003 component availability/enabled state and graph-derived effective capacity. A deliberate shutdown is explicit state, not a fault. Missing condition state means healthy/fully derated baseline; active fault effects are projected onto a cloned facility before delegating to the existing WP-003 evaluator.

## Causal faults

`FaultDefinition` is data-shaped and versioned by stable IDs such as `fault.motor_bearing_seizure`. It declares applicable component types, severity, capacity/availability effects, symptom IDs and permitted interventions. `FaultInstance` carries a stable instance ID, causal fault type, target component, severity, latent/active/resolved status, onset time, symptoms and rules version. No generic broken state is created.

The aggregate proof definitions include:

- `fault.motor_bearing_seizure` — conveyor unavailable, symptoms vibration/no rotation;
- `fault.worn_belt` — conveyor remains available at 70% capacity, vibration symptom;
- `fault.screen_blockage` — screen remains available at 30% capacity, blockage/low-flow symptoms.

## Symptoms, evidence and diagnosis

Symptoms are separate from faults and evidence. `Observation` distinguishes `Observed`, `Absent` and `Unknown`; low diagnostic capability produces `Unknown`, never fabricated telemetry. Evidence records source, component, symptom, observation ID and authoritative effective time. Diagnosis records candidate fault, evidence IDs, deterministic confidence basis points and unresolved/confirmed/rejected status.

Inspection and measurement commands are replayable and idempotent through the WP-002 command envelope. Confidence is a deterministic function of matching observed/unknown/absent evidence and diagnostic capability; no LLM or random certainty is involved.

## Interventions

The typed intervention contract supports Inspect, Test, Clean, Calibrate, Lubricate, Patch, Repair, Rebuild, Replace, Bypass, Reroute, Derate, Shutdown, Restore, Outsource and SalvageDismantle. Executable semantics cover:

- repair: resolves the specified active fault and improves condition without unrelated changes;
- replace: restores condition/derate baseline, resolves active faults and retains component ID;
- bypass/reroute: records explicit bypass state and removes only the matching dependency edge during graph projection;
- derate: sets an intentional fixed-point capacity limit distinct from fault constraint;
- shutdown/restore: controls deliberate availability without fabricating a fault;
- clean/patch/lubricate/rebuild: bounded deterministic condition/fault recovery hooks.

Invalid target, fault, evidence, intervention, bypass, derate and recovery transitions return typed errors without mutating state. No resource, money, labor, inventory or production settlement is consumed.

## Aggregate scenarios

Executable tests prove:

1. seized feed conveyor motor: active causal fault drives output to zero, inspection observes causal evidence, repair resolves and restores output;
2. worn belt: degraded condition remains operational at partial capacity;
3. blocked screen: plant remains operational but constrained;
4. low-information test: unknown evidence remains unknown and diagnosis remains below confirmation;
5. deliberate shutdown and restore: output changes without a fabricated fault;
6. failure command execution participates in kernel replay/snapshot/digest state.

Failure effects are projected into the existing WP-003 `Facility::evaluate()` path; no duplicate graph evaluator exists.

## Determinism and complexity

Failure vectors are canonicalized by stable IDs before snapshots/digests. Fault/evidence/diagnosis lookup is linear over bounded fixture-scale vectors; graph re-evaluation remains the existing ordered topological evaluation, `O(V + E)` after projection. No floating point, wall clock, filesystem, network, database, hash iteration or thread scheduling is used.

Canonical schemas were added for failure state and failure command payloads. Host prerequisites and runtime dependencies were unchanged.

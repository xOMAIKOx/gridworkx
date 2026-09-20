# WP-003 Facility Graph Evidence

## Model and hierarchy

`Facility` is the deterministic graph container with stable facility ID, facility type, World/Region/Site context references, nominal capacity, systems, components and directed dependency edges. `System` groups component IDs and declares deterministic output components. `Component` carries stable ID/type, nominal capacity, availability, enabled state and a fixed-point capacity factor. Evaluation is derived as `FacilityEvaluation`; it is not authoritative translated prose or a second simulation engine.

The accepted hierarchy is preserved as context references:

```text
World → Region → Site/Parcel → Facility → System → Component
```

WP-003 implements the Facility/System/Component graph only. It does not implement world simulation.

## Edge semantics

`DependencyType` serializes explicitly as:

- `HARD`: every ungrouped upstream dependency must be available; grouped edges are alternate paths and the group is available when any upstream is available.
- `CAPACITY`: ungrouped edges constrain by the minimum upstream effective capacity; grouped edges use the maximum available alternate route.
- `QUALITY`, `RELIABILITY`, `COST`: structural hooks only in WP-003 and do not alter capacity.
- `OPTIONAL`: structural/non-mandatory edge and never becomes a HARD dependency.

Edges are directed upstream → downstream and may carry a stable `path_group` for alternate routes.

## Capacity and bottleneck algorithm

All capacities use unsigned integer basis points: `10000 = 100%`, `300 = 3%`. A component local capacity is `nominal_capacity × capacity_factor_bps / 10000`, with unavailable/disabled components at zero. Components are evaluated in deterministic topological order. HARD and CAPACITY edge groups propagate availability and constraints; each system evaluates its declared outputs, and facility capacity is the minimum effective capacity across systems and the facility nominal basis.

Operational state is `Operational`, `Constrained` or `Unavailable`; non-100% output is not collapsed into failure. Bottleneck evidence returns stable component IDs, effective capacity and a typed reason (`ComponentUnavailable`, `ComponentCapacity`, `HardDependencyUnavailable`, `CapacityDependency` or `SystemCapacity`).

Validation rejects invalid IDs, duplicate facility/system/component IDs, orphan or multiply assigned components, missing references, duplicate edges, self-edges, invalid capacity/factors, empty graphs and dependency cycles. BTree-backed ordering and fixed-width integer arithmetic avoid insertion-order, hash-order, floating-point and platform-width nondeterminism.

## Aggregate fixture

`aggregate_plant_fixture()` represents the test-only line:

```text
feed stockpile → feed conveyor → crusher → screen → output conveyor → finished stockpile
```

The unit tests prove a healthy 100% line, hard unavailability propagation, 3% partial operation with identified crusher bottleneck, optional-element non-blocking behavior, parallel alternate-path behavior, typed cycle/reference/duplicate errors, deterministic evaluation under reordered storage, and snapshot/digest preservation through `SimulationState`.

## Scope boundary

The graph has no failure causes, symptoms, repairs, maintenance, recipes, inventory, economics, UI, Godot dependency, network, database, filesystem or wall-clock access. Those remain later work-package boundaries.

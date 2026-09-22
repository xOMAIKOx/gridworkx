# WP-005 Resources, Recipes and Inventory Evidence

## Quantity and grade model

WP-005 uses checked `u64` base quantities (`QUANTITY_SCALE = 1`) with no authoritative floating point. Resource definitions carry stable resource IDs, category, unit, rules version and explicit grade IDs. Inventory balances are keyed by `(resource_id, grade_id)`, so incompatible grades cannot merge silently.

Recipe yields use integer basis points (`10000 = 100%`). Input/output quantities use checked multiplication and division; no silent truncation or negative quantity is permitted.

## Inventory model

`InventoryStore` is a bounded generic store with stable ID, optional facility reference, capacity, permitted resources and ordered material lots. Lots are aggregate balances keyed by resource and grade for WP-005. Canonicalization sorts resources, recipes, stores, recipe inputs/outputs and lots by stable IDs.

Transfers compute accepted quantity as the minimum of requested quantity, source availability and destination free capacity. They emit a deterministic partial/accepted result and mutate neither side when accepted quantity is zero.

## Recipe and production model

`RecipeDefinition` is data-shaped with stable recipe ID, facility type targets, typed inputs, outputs, yields and rules version. `MaterialState::produce` first obtains facility capacity from the existing WP-004 projection/WP-003 evaluator, then bounds executable runs by facility capacity, each input quantity and destination capacity. It consumes inputs and produces outputs atomically after all limits are calculated.

The aggregate fixture defines:

```text
raw feed + limestone → finished aggregate + aggregate waste
```

The fixture proves healthy production, input shortage, destination capacity, grade mismatch, transfer conservation and failure-constrained production. No credits, prices, ledger, transport economics or contracts are involved.

## Determinism and scope

Material state participates in the WP-002 snapshot/digest through `SimulationState`, and material commands/events use the existing version/time/idempotency envelope. No graph or failure math is duplicated: production consumes the existing facility evaluation result.

Schemas and content are exercised by the repository validation tool. No host packages, services, databases, deployment, UI or container runtime work was performed.

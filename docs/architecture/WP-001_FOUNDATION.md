# WP-001 Foundation Boundaries

## Scope

WP-001 establishes the repository, build/test, schema/content, domain-ownership and native-runtime foundations described by the accepted DP3 architecture and ADR-001 through ADR-009. It does not implement gameplay, the WP-002 simulation kernel, deployment, or host service configuration.

## Authority boundaries

| Concern | Foundation owner |
|---|---|
| Deterministic simulation rules | `crates/gridworks-sim` only |
| Presentation and input | `apps/game` |
| Server-owned identity, ownership and settlement | `services/` boundaries |
| Cross-language contracts | `packages/schemas` |
| Versioned rules/content/localization | `packages/content`, `packages/localization`, `packages/visual-assets` |
| PostgreSQL migration conventions | `db/migrations` |
| Native process conventions | `ops/systemd` and `ops/config` |

The Rust crate contains no network, database or device-clock dependencies. Go services own orchestration boundaries and do not reproduce simulation mathematics.

## Domain scaffolding

The content catalogue and schema registry reserve explicit identifiers for identity/account, player profile, company/ownership, world/region/site, facility/component, resources/inventory, construction/maintenance, skills/managers, trade/marketplace, contracts/work exchange, consortium/JV, accounting/ledger, business lifecycle, social/messaging, notifications, localization, real estate/property, GRIDWORKS system ownership, authority/treasury ledgers, starter-path selection, operating state, bootstrap policy, presentation layers, transport modes, time and climate.

These are ownership and contract locations, not fake gameplay implementations.

## Versioning and time

All external contract envelopes carry schema, rules/content and release versions. Authoritative timestamps are supplied by server boundaries. The initial operational time ratio is represented as data (`1 in-game day = 1 real hour`); clients cannot make their wall clock authoritative.

## Runtime

The four process templates are native systemd examples for `gridworks-api`, `gridworks-worker`, `gridworks-realtime` and `gridworks-sim-validator`. They use non-root identities, external environment files, loopback/non-privileged defaults, restart policy and hardening settings. They are not installed, enabled or started by WP-001.

# WP-008 Company, Ownership and Business Identity Evidence

## Domain hierarchy

WP-008 adds a transport-independent Go company domain under `services/internal/company` while consuming WP-007 `player_id`-style owner principals. Account, Player and Company remain separate. Company and group IDs are opaque server-generated identities; provider subjects and account IDs are never company owners.

Supported company types are Holding, Operating and System. Company/group public identity is separated from ownership/history internals through strict schemas. System ownership uses an explicit `system` principal and does not fabricate a player.

## Names and normalization

Company names reuse WP-007 NFKC/case-fold/confusable subset semantics through the shared identity normalization package. Company display names may contain spaces while canonical/skeleton keys use the shared normalized comparison behavior. Reserved baseline names are represented in `packages/content/config/reserved-company-names.json` and protected by the domain store. No second Unicode algorithm is introduced.

## Ownership and history

Ownership uses exact basis points (`10000 = 100%`). Company creation creates one full owner and an append-oriented ownership event. Transfer/reallocation closes active projections, appends new active rows and appends an immutable history event while preserving the exact total. Insufficient shares, invalid totals and conflicting owners reject.

Groups are first-class identities. Active company→group assignment is exclusive and idempotent; group/company history is represented by durable membership rows and mutation receipts. The initial model does not implement recursive company ownership, so parent/subsidiary cycles are structurally excluded rather than deferred through a later recursive graph.

System-owned companies are explicit and can be created with `principal.gridworks.system`. No money movement or ledger duplication exists in WP-008.

## Idempotency and concurrency boundary

Company creation, group creation, ownership transfer and group assignment share one logical mutation receipt namespace with mutation type/request digest/result reference. Same key plus same payload replays; contradictory reuse conflicts. Repository domain operations stage all random IDs before committing in-memory state. PostgreSQL migration uniqueness, partial active-membership indexes, exact-share deferred trigger and idempotency primary key provide the later transactional/concurrency boundary; no process-local mutex is treated as durable correctness.

## PostgreSQL migration and schemas

`db/migrations/0004_wp008_company_ownership.sql` defines companies, groups, current ownership, group memberships, append-only ownership history, company mutation receipts and reserved names. It includes canonical/skeleton uniqueness, active group exclusivity, exact basis-point checks, player/system owner type fields, and a deferred active ownership total guard.

Strict company state, mutation-command and public-company-profile schemas are registered and exercised by the content validator. Public profile fields exclude ownership/audit/auth internals.

## Verification and exclusions

Go tests cover company creation/idempotency, exact initial ownership, transfer/history, insufficient share rejection, group assignment, system ownership, reserved names, generation failure and shared normalization behavior. Structural migration validation checks all required WP-008 tables/constraints and rejects destructive/provisioning SQL.

No WP-009 API, acquisition/sale settlement, marketplace, real estate, JV/Consortium, money movement, ledger duplication, deployment, host/database provisioning or container runtime work was performed.

## REQUEST_CHANGES remediation evidence

- Ownership history no longer has global owner-tuple uniqueness; only active ownership projection is unique, allowing historical re-entry.
- SQL ownership references are guarded for company/group existence, player principal existence and the explicit `principal.gridworks.system`; Go validates the same principal contract.
- Ownership identity fields are frozen after creation, active totals lock the authoritative company/group row during deferred validation, and transfer mutations preserve exact `10000` basis points.
- Full-share source exits omit zero-share active rows. Invalid principals and exact-total violations reject.
- Added detach/reassign group mutations with idempotent replay, closed prior membership rows and exclusive active-parent behavior; SQL freezes membership identity fields.
- Go and SQL share a durable company/group name-claim namespace and reserved policy. Company-create/group-create IDs are always server-generated.
- Company commands are strict per-mutation schemas and reject caller-supplied create IDs or unrelated fields.
- Ownership history is SQL append-only and structural validation asserts the guard.

## R10–R14 remediation evidence

- Ownership projection identity is frozen and only one-way lifecycle closure is permitted; authoritative entity rows are locked before ownership mutation and deferred exact-total validation.
- Group membership reassign stages replacement IDs before closing the old row; failed randomness leaves the old active membership unchanged. Closed membership identity is immutable.
- Shared company/group name claims are append-only and name identity fields are frozen because rename is out of scope. Company-specific names use the shared normalization algorithm with a consistent 3–80 policy and reserved baseline.
- Player owner resolution uses an explicit in-memory/repository player registry seam rather than accepting syntactically plausible IDs alone; system ownership remains the exact GRIDWORKS principal.
- Ownership transfer accepts typed company or group entity references and applies the same exact-share/idempotency/history rules to both.

## R15–R18 remediation evidence

- Ownership and membership SQL projections now reject DELETE and allow only one-way close; identity/effective-start fields are frozen and closed rows cannot be rewritten.
- Company/group display, canonical and skeleton name identity is fully frozen; Go validates the display value itself at 3–80 runes before canonical whitespace normalization, matching SQL and schemas.
- Company ownership consumes an injected `PlayerResolver` interface; the company domain no longer registers or mints player validity. The WP-007 identity repository exposes the resolver contract.
- The legacy company-only ownership transfer path is removed. Typed `entity_type` + `entity_id` is required for all ownership transfers, with group partial/full/replay/conflict/invalid-entity regressions.

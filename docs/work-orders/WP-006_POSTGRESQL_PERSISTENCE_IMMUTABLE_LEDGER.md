# GRIDWORKS — WP-006 PostgreSQL Persistence and Immutable Economic Ledger

**Status:** AUTHORIZED FOR ENGINEERING  
**Architecture owner:** Architecture  
**Implementation model:** GPT-5.6 Luna XHigh  
**Review:** GPT-5.6 Sol-class Architecture review  
**Exact required parent:** `106bf35e3d22ab1ef32425513157e4ae7fb438ef`

## 1. Purpose

Implement the first authoritative persistence and economic-ledger foundation for GRIDWORKS.

WP-006 establishes:
- PostgreSQL 18 schema and migrations;
- durable simulation snapshot persistence contracts;
- immutable/double-entry-inspired economic ledger primitives;
- idempotent posting boundaries;
- auditable balances derived from journal entries;
- deterministic persistence mappings for the accepted WP-002 through WP-005 state;
- repository-level integration tests that do not require unauthorized estate host changes.

This work package must create the accounting/persistence foundation required by later company, marketplace, contract and finance work without prematurely implementing those products.

## 2. Controlling references

Engineering must read and obey:
- `AGENTS.md`
- `AI_AGENT_COLLABORATION_PROTOCOL.md`
- Product Philosophy
- DP2 Core Systems Specification
- DP3 Engineering Architecture
- ADR-003 Schema/Serialization
- ADR-004 Content/Rules Distribution
- ADR-005 Initial Server Topology
- accepted WP-002 deterministic kernel
- accepted WP-003 facility/component graph
- accepted WP-004 causal failure layer
- accepted WP-005 resources/recipes/inventory

If implementation conflicts with accepted architecture, stop that path and raise evidence in the WP-006 issue.

## 3. Exact ancestry

Branch from exactly:

`106bf35e3d22ab1ef32425513157e4ae7fb438ef`

Verify:
- exact parent;
- clean working tree;
- no unrelated changes.

Do not silently rebase onto newer `main`.

## 4. Runtime / host constraints

Repository-only.

No PostgreSQL server installation, database creation, service start, role creation, credential deployment or ERIS host modification is authorized by default.

Prohibited:
- Docker/Podman/Compose/OCI/Kubernetes;
- deployment;
- systemd installation;
- nginx/DNS/firewall/WireGuard changes;
- external port activation;
- production/dev database provisioning on estate hosts.

Repository-local SQL, migrations, fixtures, pure/unit tests and CI-safe validation are authorized.

If a live PostgreSQL prerequisite is genuinely required beyond the existing repository/CI environment, stop and post a blocker with exact evidence before any host change.

## 5. PostgreSQL target

Target PostgreSQL 18 semantics.

Repository artifacts must remain compatible with the estate’s accepted PG18 target.

Do not depend on PG16-specific behavior merely because a local client may exist.

Use standard PostgreSQL features unless a PG18-specific feature is deliberate, justified and tested.

## 6. Persistence architecture

Persistence is **not** the simulation authority.

Rust remains authoritative for:
- simulation math;
- graph/failure/material transitions;
- deterministic replay.

PostgreSQL persists accepted state and economic records.

Do not reimplement simulation rules in SQL, Go or database triggers.

## 7. Persistence boundaries

Implement schema contracts for at least:

### Simulation snapshots
Persist:
- snapshot ID;
- logical aggregate/world/tenant owner reference where available;
- schema version;
- rules version;
- operational time;
- deterministic state digest;
- canonical serialized snapshot payload;
- created/accepted metadata;
- previous snapshot reference where useful.

### Command/idempotency receipt
Persist enough durable identity to prevent duplicate economic/persistence application later:
- command ID;
- idempotency key;
- command type;
- accepted/rejected status;
- state digest before/after where useful;
- effective time;
- unique constraints.

Do **not** blindly persist the entire current in-memory `executed_commands` vector as an ever-growing snapshot history.

WP-006 is the point where durable idempotency ownership should begin moving to the persistence boundary.

## 8. Resolve carried `executed_commands` concern

Architecture has carried this from WP-002 onward:

> `executed_commands` must not become unbounded production history in long-lived simulation snapshots.

WP-006 must define and implement a bounded production-safe approach.

Required:
- durable command/idempotency uniqueness in PostgreSQL;
- a bounded or explicitly non-production in-memory replay history in Rust;
- snapshot size must not grow forever merely because commands were executed;
- deterministic replay tests must remain valid.

Engineering may:
- introduce a bounded replay receipt window in simulation state; or
- split replay/test receipts from durable persistence identity.

Do not silently remove replay protection.

Document the chosen design.

## 9. Economic ledger invariant

Hard product invariant:

> Economic value movements are recorded in an immutable, auditable, double-entry-inspired journal from the first economic engineering phase.

WP-006 must implement the ledger foundation now.

Do not implement a mutable “balance column as truth” model.

## 10. Ledger model

Implement canonical repository/domain types and SQL schema for:

### Ledger account
At minimum:
- stable account ID;
- owner/entity reference;
- account type/class;
- currency/asset code;
- status;
- created metadata.

### Journal transaction
At minimum:
- stable transaction ID;
- transaction type;
- effective time;
- external/idempotency reference;
- source command/event reference;
- description/code/reference fields using stable IDs, not localized prose;
- immutable posting status.

### Journal entry / line
At minimum:
- transaction ID;
- account ID;
- debit/credit or signed amount representation;
- exact integer amount;
- currency/asset;
- line sequence;
- metadata/reference.

## 11. Double-entry balance rule

Every posted monetary transaction must balance exactly.

Preferred invariant:

`sum(debits) == sum(credits)`

or an equivalent signed-entry invariant where total signed amount = 0.

Requirements:
- checked integer arithmetic;
- same currency within a balancing unit unless explicit FX transaction structure is later added;
- no authoritative floating point;
- imbalance rejected before commit;
- transaction posting atomic.

## 12. Monetary representation

Choose and document deterministic monetary units.

Preferred:
- signed/unsigned 64-bit integer minor units;
- currency ID/ISO-like code stored separately.

Examples:
- USD cents;
- EUR cents;
- game Credits smallest unit.

No floats/decimals as authoritative application math.

If PostgreSQL `NUMERIC` is used for storage, Rust application boundaries must still use explicit checked integer semantics unless Architecture approves otherwise.

## 13. Currency scope

WP-006 establishes the ledger contract; it does not implement FX markets.

Support multiple currency/asset IDs structurally.

A single journal transaction must not accidentally balance USD debit against EUR credit.

Cross-currency exchange requires explicit later transaction semantics.

## 14. Account classes

Support a minimal extensible class set sufficient for later company/player/system economics, for example:
- Asset
- Liability
- Equity
- Revenue
- Expense
- Clearing/System

Do not hardcode full accounting products yet.

Account ownership should support future:
- player;
- company;
- GRIDWORKS system principal;
- authority/treasury;
- contract/JV where needed.

Use generic owner/entity references if those domains are not yet implemented.

## 15. Immutable posting

Posted journal transactions and lines must be append-only.

Required:
- application code does not UPDATE posted lines;
- correction model is reversal/compensating transaction, not mutation;
- SQL constraints/triggers may enforce immutability if simple and deterministic;
- do not put simulation/business decision logic in triggers.

If using triggers for append-only protection, keep them strictly integrity-focused.

## 16. Derived balances

Balances are derived from journal entries.

You may add:
- query/view helpers;
- cached/materialized balance projections;

but they must be reconstructible from immutable journal lines.

A cached balance must never become the sole source of truth.

## 17. Ledger idempotency

Posting the same economic transaction twice must be prevented.

Use explicit unique keys such as:
- external transaction/idempotency ID;
- source command ID;
- source event ID;
- owner/scope where required.

Duplicate posting must return a deterministic duplicate/idempotency result, not create extra money.

## 18. Atomic posting

A journal transaction and all of its lines must commit atomically.

No state where:
- transaction row exists without all lines;
- only one side of the double-entry exists;
- balance projection changes without journal evidence.

Repository-level API/transaction abstractions should make partial posting difficult/impossible.

## 19. Initial economic proof

Add a bounded synthetic proof transaction, not marketplace/company gameplay.

Examples:
- system grants 1,000 Credits from a system issuance/clearing account to a synthetic player wallet;
- material production cost proof is **not required** yet.

Prove:
- transaction balances;
- debit/credit entries are immutable;
- derived account balances are correct;
- duplicate posting rejected;
- reversal produces opposite entries without mutating original transaction.

Do not infer a full monetary policy from the synthetic proof.

## 20. GRIDWORKS/system principal

The architecture already requires GRIDWORKS/system participation later.

WP-006 may define generic system-owned ledger accounts needed for proofs.

Do not implement:
- buyer-of-last-resort behavior;
- taxes;
- procurement;
- stabilization;
- marketplace liquidity.

Those are later economic behaviors.

## 21. Material/economic separation

WP-005 material quantities are **not currency**.

Do not mix:
- inventory quantity;
- physical resource units;
- monetary ledger amount.

Future inventory valuation can reference both domains, but WP-006 must preserve the distinction.

## 22. Persistence of material/failure/facility state

Define repository persistence mapping for accepted simulation snapshots, including WP-003/004/005 state.

Preferred model:
- canonical JSON/JSONB snapshot as replay artifact;
- selected normalized metadata/index columns for retrieval/integrity;
- avoid prematurely normalizing every component/fault/material lot into relational tables unless justified.

The deterministic canonical snapshot remains the authoritative replay payload.

## 23. Snapshot integrity

Persist and validate:
- digest;
- schema version;
- rules version;
- operational time.

On load:
- payload must deserialize through current supported schema/migration path;
- digest mismatch must reject;
- unsupported version must reject or route through an explicit migration path.

No silent snapshot repair.

## 24. Snapshot write semantics

A snapshot write must be atomic and idempotent where keyed.

Requirements:
- duplicate snapshot ID/digest behavior explicit;
- no partial snapshot metadata/payload rows;
- previous/current linkage cannot create invalid chains.

## 25. Repository API boundaries

Implement repository/domain abstractions that later Go services can consume without owning simulation/accounting rules.

Suggested boundaries:
- SnapshotRepository
- CommandReceiptRepository
- LedgerRepository / JournalRepository

Exact module names are Engineering-owned.

Avoid generic “repository framework” overengineering.

## 26. Language ownership

Rust owns:
- simulation state serialization/digest;
- ledger domain invariants if implemented in Rust.

Go may later orchestrate API/database transactions, but WP-006 should not duplicate accounting validation in multiple languages.

If ledger domain primitives live in Rust, SQL constraints still independently protect integrity.

If a minimal Go persistence adapter is added, it must not reimplement balance math.

## 27. SQL migrations

Add versioned SQL migrations under the repository’s accepted migration structure.

Migrations must:
- be deterministic;
- be rerunnable only according to migration tooling conventions;
- use explicit constraints/indexes;
- avoid destructive reset patterns;
- contain no credentials;
- avoid host-specific paths.

Required tables should include equivalent concepts for:
- simulation snapshots;
- command/idempotency receipts;
- ledger accounts;
- journal transactions;
- journal entries.

Names are Engineering-owned but must be clear.

## 28. SQL constraints

Use database constraints where appropriate:
- primary keys;
- foreign keys;
- unique idempotency keys;
- non-zero/positive amount rules as chosen;
- valid line sequencing;
- transaction/account currency compatibility where practical;
- immutable posted-row protection if implemented safely.

Do not rely solely on application code for basic relational integrity.

## 29. Database transaction isolation

Document expected transaction boundaries and isolation assumptions.

At minimum:
- journal post uses one DB transaction;
- idempotency receipt + journal posting ordering is defined;
- concurrent duplicate requests cannot double-post.

Do not implement distributed transactions.

## 30. Concurrency

Repository/domain design must account for concurrent posting attempts.

Use:
- unique constraints;
- transaction semantics;
- deterministic conflict handling.

Do not rely on process-local mutexes as the correctness mechanism.

## 31. Reversal

Implement a minimal reversal/compensation model.

A reversal:
- creates a new journal transaction;
- references the original;
- posts exact opposite lines;
- cannot mutate original entries;
- itself is idempotent.

No arbitrary delete/update of posted transactions.

## 32. Audit metadata

Persist bounded audit metadata:
- created/posted time from authoritative service/database context;
- source command/event reference;
- actor/principal reference when available;
- correlation ID where available.

Do not put sensitive secrets into journal metadata.

WP-006 does not implement the full user/audit product.

## 33. PostgreSQL tests without unauthorized host changes

Use repository/CI-safe strategies.

Acceptable:
- SQL parser/shape tests;
- migration static validation;
- domain/unit tests;
- existing CI PostgreSQL service only if already provided and compliant;
- ephemeral test database only if CI natively provides it without Docker/container invocation by project code.

Not acceptable:
- installing/starting PostgreSQL on ERIS;
- invoking Docker/Podman locally;
- modifying estate services.

If no live DB exists in CI, prove SQL invariants as far as possible statically and clearly document the remaining live-PG gate for a later authorized environment.

## 34. Schema validation

Extend repository validation to verify:
- migration files exist and are ordered;
- no destructive/forbidden patterns;
- required ledger tables/constraints are represented;
- simulation/ledger schemas remain aligned with domain types where applicable.

Do not create a fake SQL parser that gives false confidence.

## 35. Ledger/domain tests

Required unit/domain tests:

### Balanced posting
- valid 2-line transaction passes;
- valid multi-line transaction passes.

### Imbalance
- unbalanced transaction rejects with no mutation.

### Currency
- mixed-currency accidental balancing rejects.

### Amounts
- zero/invalid amounts rejected per chosen representation;
- overflow rejected.

### Duplicate/idempotency
- duplicate external/source key rejects.

### Reversal
- original remains unchanged;
- reversal posts opposite entries;
- net derived balance returns appropriately.

### Derived balances
- account balance derives correctly from immutable lines.

## 36. Persistence tests

Required:
- canonical snapshot round-trip through persistence representation;
- digest preserved;
- version/digest mismatch rejected;
- duplicate snapshot behavior deterministic;
- command receipt uniqueness;
- bounded in-memory executed-command strategy tested;
- canonical state remains replayable.

## 37. Integration proof

Provide an executable repository-level proof tying together:

1. deterministic simulation state exists;
2. canonical snapshot is serialized + digest computed;
3. persistence record is created/validated;
4. synthetic balanced journal transaction is created;
5. duplicate journal post is rejected;
6. reversal is represented;
7. no simulation/material mutation is caused by ledger storage itself.

A live DB is preferred only if an already-authorized CI/test facility exists. Otherwise use domain + migration evidence and clearly state live DB execution remains a later environment gate.

## 38. No marketplace/company implementation yet

Do not implement:
- company balance sheets;
- revenue from real player transactions;
- marketplace settlement;
- subscriptions;
- invoices;
- contracts;
- taxes;
- stock exchange;
- dividends;
- loans;
- interest;
- valuation.

WP-006 creates primitives only.

## 39. Security

Required:
- parameterized SQL in any adapter code;
- no string-built SQL from untrusted values;
- no credentials committed;
- least-privilege assumptions documented;
- ledger IDs/amounts validated;
- immutable journal records protected.

Secret and forbidden-runtime scans remain green.

## 40. Canonical IDs

Use stable semantic IDs where domain-facing.

Database primary keys may use UUID/text/native identifiers, but external/domain identity semantics must be explicit and stable.

Do not let database surrogate IDs become the only externally meaningful identity.

## 41. Time

Simulation effective time and database/audit creation time are distinct.

Preserve:
- simulation operational/effective time from kernel;
- journal effective time;
- persistence created/posted timestamps.

Do not substitute DB wall-clock time for simulation time.

## 42. Carried architecture constraints

From WP-002:
- deterministic replay;
- bounded `executed_commands` production strategy now required.

From WP-003:
- explicit facility output topology.

From WP-004:
- causal failures and information boundaries.

From WP-005:
- checked material arithmetic;
- material inventory and money remain distinct;
- production is atomic.

Do not silently reinterpret accepted contracts.

## 43. Allowed paths

Primary:
- `crates/gridworks-sim/**` for bounded command-history adjustments or ledger domain primitives if appropriate;
- new persistence/ledger domain crate/module if justified;
- `services/**` only for minimal persistence adapter/library scaffolding;
- `db/**`, `migrations/**` or existing repository migration path;
- `packages/schemas/**`;
- `tests/**`;
- `docs/evidence/**`;
- CI/build validation files if required.

No host/deployment files unless purely repository templates already within accepted scope and directly necessary.

## 44. Out of scope

Do not implement:
- WP-007 identity;
- WP-008 company ownership;
- WP-009 API foundation;
- marketplace/trade;
- contracts;
- real estate;
- managers/skills;
- finance products;
- UI;
- deployment;
- live PostgreSQL provisioning.

## 45. Evidence

Create concise evidence documenting:
- persistence boundaries;
- snapshot schema/mapping;
- command/idempotency strategy;
- resolution of unbounded `executed_commands`;
- monetary representation;
- ledger account/transaction/entry model;
- balancing invariant;
- immutability/reversal;
- transaction/concurrency assumptions;
- SQL migration design;
- tests;
- live-DB limitations if any;
- security constraints.

## 46. Done-when

Architecture can verify:

1. exact parent;
2. PG18-targeted migrations exist;
3. snapshot persistence contract exists;
4. digest/version integrity is enforced;
5. durable command/idempotency uniqueness exists in schema/design;
6. unbounded snapshot command history is resolved;
7. immutable ledger account/journal/entry model exists;
8. double-entry balance invariant is enforced;
9. amounts use checked deterministic integer semantics;
10. duplicate posting cannot create money;
11. reversal is append-only;
12. derived balances come from journal lines;
13. material quantities remain distinct from money;
14. no marketplace/company/finance scope creep;
15. CI/security gates pass;
16. complete GitHub handback exists.

## 47. Handback

Post to the WP-006 issue and draft PR:
- branch;
- exact parent/head;
- commit chain;
- changed files/dependencies;
- persistence architecture;
- migration list;
- snapshot/idempotency design;
- executed-command retention solution;
- ledger model;
- money/currency representation;
- balancing/idempotency/reversal semantics;
- test evidence;
- CI links/results;
- live-PG limitation/blocker if applicable;
- deviations/risks/conflicts;
- explicit no-host/deployment statement.

## 48. Stop condition

After handback, Engineering stops.

Do not start WP-007 without separate authorization.

# GRIDWORKS — WP-008 Company, Ownership and Business Identity Foundation

**Status:** AUTHORIZED FOR ENGINEERING  
**Architecture owner:** Architecture  
**Implementation model:** GPT-5.6 Luna XHigh  
**Review:** GPT-5.6 Sol-class Architecture review  
**Exact required parent:** `f347de36d3d1ec4095a63b4c4e88d1c7d8e97e1a`

## 1. Purpose

Implement the server-authoritative company and ownership foundation that sits between Player identity and later operational/business systems.

WP-008 establishes:
- persistent company identities independent from human accounts;
- player-controlled company ownership;
- company groups / holding-company hierarchy;
- operating-company identity;
- ownership-share records and auditable ownership history;
- company profile/public identity foundations;
- company-name uniqueness and normalization;
- company lifecycle/status foundations;
- deterministic/idempotent ownership mutations;
- relational integrity suitable for later sale/acquisition/JV/market domains.

This work package does **not** implement business acquisition flows, marketplace settlement, finance products, public HTTP APIs, JVs, consortiums, facility ownership transfer, real estate, or deployment.

## 2. Controlling references

Engineering must read and obey:
- `AGENTS.md`
- `AI_AGENT_COLLABORATION_PROTOCOL.md`
- Product Philosophy and Game Design
- DP2 Core Systems Specification
- DP3 Engineering Architecture & Work Packages
- accepted WP-006 ledger/persistence foundation
- accepted WP-007 identity/player-profile foundation

Controlling architecture invariants include:
- Player, Account and Company are separate domains.
- Company identity persists independently from the player.
- Company ownership can change without account transfer.
- Company ownership is server-authoritative.
- Business/ownership history is permanent and auditable.
- PostgreSQL is authoritative for transactional ownership state.

## 3. Exact ancestry

Branch from exactly:

`f347de36d3d1ec4095a63b4c4e88d1c7d8e97e1a`

Verify:
- exact parent;
- clean working tree;
- no unrelated changes.

Do not silently rebase onto a newer `main`.

The Architecture work-order commit on `main` is intentionally not the Engineering parent unless explicitly stated later.

## 4. Runtime / host constraints

Repository-only.

No:
- live PostgreSQL provisioning;
- host changes;
- systemd/nginx/DNS/firewall/WireGuard;
- service deployment;
- container runtime;
- NATS or external provider changes;
- public listener/port activation.

Repository migrations, schemas, transport-independent domain code, tests and CI validation are authorized.

## 5. Domain hierarchy

Preserve the accepted hierarchy:

`Account → PlayerProfile → Ownership/Control → CompanyGroup / Company`

Company is an economic/legal game entity, not an authentication identity.

Definitions:

- `company_id`: immutable company identity.
- `group_id`: immutable holding/group identity where used.
- `player_id`: human/game principal identity from WP-007.
- ownership/control records reference principals; they do not alter account identity.
- company profile data is separate from player profile data.

## 6. Company types

Minimum initial types:
- Holding / Group;
- Operating Company;
- System / GRIDWORKS-owned company or principal representation where structurally useful.

Do not add speculative legal forms by country.

A company type must not imply ownership percentages or operating status.

## 7. Company identity

Each company must have:
- immutable company ID;
- company type;
- stable canonical name identity;
- display name;
- status;
- optional short description;
- optional logo/emblem asset reference;
- optional industry/category semantic ID;
- optional headquarters region reference;
- public/private visibility state;
- created/updated metadata.

Do not require later-domain data such as balance sheet, facilities, market valuation, contracts, managers or reputation implementation.

References/placeholders are acceptable where later domains will own the values.

## 8. Company names

Company names are unique within the logical world.

Requirements:
- Unicode normalization;
- Unicode case-folding;
- confusable-aware canonical/skeleton uniqueness, reusing or sharing the accepted WP-007 normalization infrastructure where appropriate;
- reject controls/invisible abuse;
- bounded length;
- reserved/system names;
- deterministic comparison key;
- display/original form may differ from canonical uniqueness key.

Do not create a second incompatible normalization algorithm.

Prefer a shared normalization package or clearly shared versioned logic.

## 9. Reserved company names

Protect at least:
- GRIDWORKS / official system representations;
- administrator/support/moderation impersonation names;
- operator-reserved names where product policy requires.

Use configurable/data-driven policy where practical.

Do not over-reserve ordinary commercial words.

## 10. Holding groups

Support a first-class company-group/holding structure.

Minimum semantics:
- one group may contain zero or more companies;
- a company may belong to at most one parent group at a time;
- group identity persists independently;
- group membership history is auditable;
- parent-child cycles are impossible;
- no arbitrary recursive nesting unless Architecture evidence demonstrates a requirement.

Preferred initial model:
- Player-controlled Group/Holding;
- Operating Companies directly under a Group;
- no deeper recursive tree for Release-1 foundation unless needed.

If Engineering chooses nested groups, prove cycle prevention and explain why.

## 11. Ownership model

Ownership is a first-class durable relation.

Minimum ownership record:
- ownership ID;
- company ID or group ID;
- owner principal type;
- owner principal ID;
- ownership share / basis-points or equivalent exact integer representation;
- voting/control share field only if needed structurally;
- effective-from;
- effective-to nullable;
- source mutation/reference;
- created metadata.

Do not use floating-point percentages.

Preferred:
- basis points or another exact integer scale;
- total active ownership must equal exactly 100% / fixed scale for entities requiring fully allocated ownership.

If unallocated/system treasury ownership is needed, model it explicitly rather than allowing silent totals below 100%.

## 12. Owner principals

Initial owner principal types may include:
- Player;
- Company;
- GRIDWORKS/System principal.

Do not implement JV or Consortium ownership yet.

A player ownership reference must point to a valid WP-007 `player_id`.

A company-owned subsidiary may be allowed if it supports the accepted holding model, but cycles must be impossible.

## 13. Control model

Do not equate authentication account with economic control.

Model control as ownership-derived or explicit authority state.

Minimum requirement:
- domain can determine whether a `player_id` currently controls a company/group;
- control survives external identity linking because WP-007 preserves `player_id`;
- account transfer is not required when ownership changes.

Avoid complex corporate governance voting mechanics in WP-008.

## 14. Ownership creation

Initial company creation operation should create:
- company/group identity;
- initial ownership allocation;
- company profile;
- ownership-history entry;
- idempotency receipt;

atomically.

A new player-owned company should default to a fully allocated initial owner unless explicit system-owned creation is requested by an internal trusted path.

No orphan company without an ownership state.

## 15. Ownership mutation types

WP-008 may support only bounded foundation mutations such as:
- create company;
- create holding/group;
- attach/detach operating company to/from group;
- change company profile metadata;
- transfer or reallocate ownership shares between already-valid principals;
- deactivate/archive company if policy allows.

Do not implement sale price, payment settlement, acquisition negotiation or marketplace flows.

Ownership transfer here is domain state mechanics only.

WP-017 will own sale/acquisition transfer orchestration.

## 16. Ownership transfer invariant

Any ownership-share mutation must:
- be server-authoritative;
- be idempotent;
- be transactional;
- preserve exact share totals;
- preserve history;
- never mutate historical records in place if append-only event/history design is chosen;
- never orphan an entity with invalid ownership;
- reject insufficient ownership / impossible transfer;
- reject self/cycle structures that violate hierarchy constraints.

No money moves in WP-008.

If later ownership sale requires money, settlement must use WP-006 ledger under later work packages.

## 17. Ownership history

Ownership history is permanent/auditable.

Minimum history/event record:
- event ID;
- entity ID;
- mutation type;
- from-owner principal if applicable;
- to-owner principal if applicable;
- share amount;
- effective time;
- source/idempotency reference;
- actor/player reference where appropriate;
- created/audit metadata.

History must not be silently rewritten.

Current ownership may be represented by active rows/projection, but reconstructability/auditability must be retained.

## 18. Group membership history

Company ↔ Group assignment changes must also be auditable.

Record:
- membership ID;
- company ID;
- group ID;
- effective-from;
- effective-to;
- source mutation;
- created metadata.

A company cannot be active in two groups simultaneously.

## 19. Company status

Minimum statuses should cover:
- Active;
- Suspended/Mothballed placeholder only if structurally useful;
- Archived/Closed only if needed;
- System-owned state should not be encoded as a status.

Do not implement business operational lifecycle from WP-016.

Company status is entity-state, not facility production state.

## 20. Public/private company profile boundary

Add strict schema/domain separation for:
- internal company record;
- public company profile.

Public profile may include:
- company ID;
- name;
- logo/emblem asset reference;
- industry/category;
- HQ/region reference;
- public/private state;
- optional short description.

Do not expose:
- private ownership-control internals beyond intentionally public ownership summary;
- audit/internal mutation references;
- security metadata;
- private account/player auth data;
- hidden system flags.

## 21. Ownership visibility

Architecture must allow later public ownership display without making private internal relations automatically public.

Use an explicit visibility/public-summary boundary.

Do not assume all ownership structures are public.

## 22. Company ticker

Tickers are later.

WP-008 may reserve an optional field/interface but must not implement exchange/ticker allocation.

Do not create stock-market semantics.

## 23. Company profile metrics

Do not implement financial/operational metrics calculations.

Public/internal profile structures may reserve semantic references for later:
- reputation;
- operating metrics;
- financial metrics;
- consortium memberships.

No fake or placeholder computed values are required.

## 24. Persistence / PostgreSQL 18 migration

Add a PG18-targeted migration for company/ownership foundations.

Expected tables or equivalent:
- companies;
- company_groups;
- company_profiles or profile fields on companies;
- company_ownership;
- company_group_membership;
- ownership_history / ownership_events;
- company_mutation_receipts;
- reserved_company_names if persisted/config-backed.

Use explicit:
- PK/FK;
- owner-principal integrity strategy;
- exact share constraints;
- active-ownership uniqueness;
- total-share enforcement;
- group-membership uniqueness;
- cycle-prevention design where applicable;
- idempotency uniqueness;
- timestamps/status checks;
- canonical/skeleton company-name uniqueness.

No credentials/roles/provisioning SQL.

## 25. Total ownership enforcement

The durable database contract must prevent invalid active ownership totals.

Required:
- exact scale;
- no floating point;
- active shares cannot exceed 100%;
- mutations that are intended to leave fully allocated ownership must commit only at exactly 100%;
- concurrency cannot allow two simultaneous writes to over-allocate.

Use PostgreSQL transaction/locking/constraint mechanisms appropriate to the chosen schema.

Do not rely on a process-local mutex as the durable correctness boundary.

## 26. Ownership concurrency

Document and prove the concurrency model for:
- company creation idempotency;
- simultaneous ownership transfer/reallocation;
- simultaneous group assignment;
- name claim races.

Correctness must rely on DB uniqueness/locking/transaction semantics.

Repository/in-memory domain may use mutexes for test isolation, but must model the same durable invariants.

## 27. Parent/subsidiary cycles

If Company may own Company:
- reject direct self-ownership;
- reject indirect ownership cycles where they would make control/hierarchy invalid.

If Group membership is modeled separately:
- prevent group/company membership cycles;
- prevent a company from being both parent and child in contradictory ways.

If the initial design avoids recursive company ownership, document that bounded decision explicitly.

## 28. System / GRIDWORKS ownership

The architecture must support non-player ownership by the GRIDWORKS system principal.

Requirements:
- system principal is explicit;
- player identity is not fabricated for GRIDWORKS;
- system-owned company assets can later enter liquidity/bootstrap flows;
- ownership is still auditable.

Do not implement buyback/resale behavior in WP-008.

## 29. Identity integration

WP-008 must consume accepted WP-007 identities.

Requirements:
- owner player references valid `player_id`;
- ownership persists across guest→linked conversion because player ID is immutable;
- provider subjects/emails/session IDs never become owner principals;
- no duplicated player/account identity tables.

## 30. Ledger integration boundary

Do not move money.

Where an ownership mutation may later be associated with economic settlement:
- store source/reference hooks only;
- do not create a fake cash settlement;
- do not duplicate ledger logic.

WP-006 remains economic value authority.

## 31. Idempotency

All company/ownership mutations must use a single logical company-domain idempotency namespace.

Receipt must include:
- idempotency key;
- mutation type;
- deterministic request digest;
- safe result reference(s);
- created metadata.

Same key + same mutation/payload => deterministic replay.

Same key + different mutation type/payload => typed idempotency conflict.

Concurrent duplicates cannot create duplicate companies/ownership events.

## 32. Mutation atomicity

Company creation:
- identity;
- profile;
- initial ownership;
- initial history;
- receipt

must commit atomically.

Ownership change:
- current ownership projection;
- history/event;
- group/control updates if relevant;
- receipt

must commit atomically.

No partial state.

## 33. Typed errors

At minimum:
- invalid idempotency;
- idempotency conflict;
- company not found;
- group not found;
- owner principal not found;
- company name invalid;
- company name unavailable;
- company name reserved;
- confusable company-name collision;
- invalid ownership share;
- ownership total mismatch;
- insufficient share;
- duplicate ownership relation;
- invalid group membership;
- group membership conflict;
- ownership/control cycle;
- forbidden system-principal mutation;
- invalid company status/profile.

Errors must not leak auth-private identity data.

## 34. Company name normalization

Prefer reuse of the accepted identity normalization primitives.

If shared code is refactored from WP-007:
- behavior must remain backward compatible for player handles;
- tests for WP-007 must remain green;
- company-specific length/character policy may differ from handles;
- normalization/case-fold/confusable algorithm version should be shared where possible.

Do not silently fork Unicode semantics.

## 35. Tests — company creation

Required:
- create company with stable opaque ID;
- creates initial owner at exact full ownership;
- profile/name state created;
- history created;
- receipt created;
- replay is deterministic;
- contradictory idempotency reuse rejects;
- generation failure causes no partial state;
- duplicate name/canonical/skeleton rejects.

## 36. Tests — ownership

Required:
- exact share total accepted;
- over-allocation rejected;
- under-allocation rejected when operation contract requires fully allocated state;
- valid transfer/reallocation preserves 100%;
- insufficient share rejected;
- same transfer replay idempotent;
- contradictory reuse rejected;
- history remains append-only;
- player ownership survives account protection/linking assumptions via stable player ID.

## 37. Tests — groups

Required:
- create holding/group;
- assign operating company;
- duplicate active assignment rejected/idempotent as designed;
- second simultaneous parent/group conflict rejected by durable contract;
- detach/reassign records history;
- cycle prevention if recursive structure exists.

## 38. Tests — system principal

Required:
- explicit GRIDWORKS/system principal may own a company;
- no fake player/account needed;
- public/private serialization does not leak internal system control flags.

## 39. Tests — public/private schemas

Required:
- strict internal company schema;
- strict public company-profile schema;
- public profile rejects private ownership/audit/auth fields;
- unknown fields reject;
- name normalization/canonical/skeleton fields not exposed publicly unless explicitly intended.

## 40. Tests — migration/static validation

Validate:
- expected company/ownership tables;
- exact-share representation;
- player-owner FK/integrity strategy;
- company-name canonical/skeleton uniqueness;
- idempotency PK/uniqueness;
- active group membership uniqueness;
- total-share enforcement mechanism;
- immutable/auditable history;
- no credentials/destructive/provisioning SQL.

Live PG execution is not required unless an already-authorized PG18 test environment exists.

## 41. Repository architecture

Preferred:
- transport-independent company/ownership domain under `services/internal/company` or equivalent;
- share reusable normalization utilities rather than copy-paste where practical;
- migrations under `db/migrations`;
- strict schemas under `packages/schemas`;
- config under `packages/content/config`;
- structural tests under `tests/structural`;
- evidence under `docs/evidence`.

Do not create generic enterprise DDD/framework abstractions.

## 42. No API implementation yet

WP-009 owns the Go API foundation.

WP-008 must not add:
- HTTP handlers;
- routers;
- public endpoints;
- WebSocket;
- network listeners.

Domain/service/repository interfaces only.

## 43. No acquisition/business lifecycle implementation

Do not implement:
- opportunity board;
- due diligence;
- business listings;
- buy/sell negotiation;
- acquisition price;
- sale settlement;
- GRIDWORKS buyer-of-last-resort;
- real estate transactions;
- company failure/turnaround generation.

Those belong to later WPs.

## 44. No JV / Consortium implementation

Do not implement:
- JVs;
- consortium membership;
- JV share structures;
- consortium governance.

Reserve extensible owner-principal/entity references if needed, but keep implementation out of scope.

## 45. Security

Requirements:
- server-authoritative mutations;
- no client-provided authoritative IDs where server generation is required;
- parameterized SQL in later adapters;
- no auth/session/provider secrets in company records;
- no implicit ownership based on request-supplied player IDs without later authenticated-principal binding;
- internal system-principal paths distinguishable from player paths;
- audit/event payloads contain no secrets.

## 46. Evidence

Create concise evidence documenting:
- Player/Account/Company separation;
- company/group identity model;
- company name normalization version;
- ownership share scale;
- ownership/control principal model;
- company creation transaction;
- ownership mutation/idempotency;
- group membership model;
- ownership/history immutability;
- cycle prevention;
- system principal ownership;
- public/private company profile boundary;
- PG18 constraints/indexes/concurrency approach;
- tests;
- live-PG limitation;
- exclusions.

## 47. Done-when

Architecture can verify:

1. exact parent;
2. company identity independent from account/player identity;
3. player ownership references immutable `player_id`;
4. company names are Unicode/case/confusable aware and unique;
5. holding/group + operating company foundation exists;
6. ownership uses exact integer shares;
7. active ownership cannot over/under allocate contrary to contract;
8. ownership mutations are atomic/idempotent;
9. ownership history is permanent/auditable;
10. group membership is exclusive/auditable;
11. cycles are impossible or structurally excluded;
12. system principal ownership is supported explicitly;
13. public/private company profile schemas are separated;
14. PostgreSQL migration enforces durable integrity/concurrency boundaries;
15. no ledger duplication or money movement;
16. no WP-009/API, acquisition, marketplace, JV/Consortium or deployment scope creep;
17. CI/security gates pass;
18. complete GitHub handback exists.

## 48. Handback

Post to the WP-008 issue and draft PR:
- branch;
- exact parent/head;
- commit chain;
- changed files/dependencies;
- company/group model;
- ownership share representation;
- principal model;
- company-name normalization/confusable strategy;
- idempotency model;
- group/cycle model;
- system-principal model;
- migration list and integrity constraints;
- public/private schema design;
- concurrency semantics;
- tests and CI;
- live-PG limitations;
- deviations/risks/conflicts;
- explicit no-host/deployment statement.

## 49. Stop condition

After handback, Engineering stops.

Do not start WP-009 without separate authorization.

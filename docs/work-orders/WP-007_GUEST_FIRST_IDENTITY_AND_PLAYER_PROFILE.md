# GRIDWORKS — WP-007 Guest-First Identity and Player Profile Foundation

**Status:** AUTHORIZED FOR ENGINEERING  
**Architecture owner:** Architecture  
**Implementation model:** GPT-5.6 Luna XHigh  
**Review:** GPT-5.6 Sol-class Architecture review  
**Exact required parent:** `8306ac89cc25c8d7e193e1b93f51d3fe9c12ecf9`

## 1. Purpose

Implement the server-authoritative identity and player-profile foundation required by ADR-002.

WP-007 establishes:
- immutable GRIDWORKS account and player identities;
- server-issued guest-first accounts;
- guest session / credential domain contracts;
- external-identity linking scaffolding;
- recovery/protection state;
- player profile persistence;
- unique handle normalization/confusable defenses;
- privacy / DM / notification preference foundations;
- durable audit/idempotency boundaries for identity mutations.

This work package does **not** implement the public/mobile HTTP API, live Apple/Google/OIDC integration, company ownership, messaging, or deployment.

## 2. Controlling references

Engineering must read and obey:
- `AGENTS.md`
- `AI_AGENT_COLLABORATION_PROTOCOL.md`
- `docs/decisions/ADR-002_IDENTITY_AND_ACCOUNT_LINKING.md`
- Product Philosophy
- DP2 Core Systems Specification
- DP3 Engineering Architecture
- accepted WP-006 persistence/ledger boundary

If implementation conflicts with ADR-002 or accepted architecture, stop that path and raise evidence in the WP-007 issue.

## 3. Exact ancestry

Branch from exactly:

`8306ac89cc25c8d7e193e1b93f51d3fe9c12ecf9`

Verify:
- exact parent;
- clean working tree;
- no unrelated changes.

Do not silently rebase onto a newer `main`.

## 4. Runtime / host constraints

Repository-only.

No host identity provider, PostgreSQL server, Redis/Valkey, NATS, service, DNS, proxy or systemd changes are authorized.

Prohibited:
- Docker/Podman/Compose/Kubernetes/OCI;
- deployment;
- live DB/role provisioning;
- external OIDC client creation;
- Apple/Google production configuration;
- secrets or credentials committed to the repository.

Repository migrations, schemas, domain libraries, pure tests and CI validation are authorized.

## 5. Identity hierarchy

Preserve the accepted distinction:

`Account → PlayerProfile → later Ownership/Control → Company`

Definitions:

- `account_id`: private immutable authentication/account identity.
- `player_id`: immutable game identity.
- `unique_handle`: public unique identifier.
- `display_name`: mutable, non-unique public name.

Account is not a Company.
PlayerProfile is not a Company.
Provider identity is not the GRIDWORKS account identity.

WP-008 will own company/control relationships.

## 6. Guest-first invariant

First launch must be able to create a server-issued guest identity without a conventional registration flow.

A guest issuance operation must create exactly one:
- immutable account;
- immutable player;
- initial player profile;
- guest session/credential record.

Guest issuance must be:
- server generated;
- transactionally atomic;
- replay/idempotency safe;
- independent of client-generated authoritative IDs;
- free of provider-specific identity assumptions.

## 7. Account status model

Define an explicit account lifecycle.

Minimum states should cover:
- Guest / Unprotected;
- Protected / Linked;
- Suspended;
- RecoveryRestricted / equivalent if justified;
- Closed/Disabled only if needed structurally.

Do not encode moderation or commercial entitlements into authentication status.

A player may continue normal gameplay as a guest until a product policy later requires protection for a specific sensitive action.

## 8. Stable identifiers

Use durable opaque account/player identifiers.

Requirements:
- globally collision-resistant;
- server generated;
- not derived from email, provider subject, handle or device ID;
- immutable after creation;
- SQL unique/PK constraints;
- no personally identifying data encoded in IDs.

Exact ID representation is Engineering-owned.

If UUID is used, choose a PostgreSQL-compatible representation and document ordering/privacy implications.

## 9. Guest session / credential model

Define a secure guest session credential boundary.

Requirements:
- raw long-lived session/refresh secret is never stored in plaintext;
- database stores a one-way digest/verifier plus metadata;
- session record has stable ID;
- account association;
- created/last-used/expiry/revocation metadata;
- credential rotation/revocation can be represented;
- concurrent valid sessions can be supported later;
- short-lived app access credentials can be layered above this without changing account identity.

Do not implement production JWT signing infrastructure unless strictly necessary for the domain proof.

A repository-safe opaque-token proof is acceptable.

## 10. Secret/token handling

For any token proof:
- generate with cryptographically secure randomness in production-facing code;
- store only verifier/hash;
- compare securely;
- never log raw tokens;
- tests may use deterministic fixtures only where clearly isolated.

Do not invent custom encryption.

Use standard library / well-established primitives already available in the repository where possible.

## 11. External identity link model

Define durable link records for providers such as:
- Apple;
- Google;
- Email/OIDC;
- Estate OIDC.

A link record must contain:
- link ID;
- account ID;
- provider type;
- provider issuer / authority where relevant;
- provider subject;
- normalized provider identity key;
- linked timestamp;
- verification/proof metadata reference as appropriate;
- status/revocation fields if required.

Do not store raw OAuth access/refresh tokens as the identity link itself.

Provider `subject` is not the account ID.

## 12. Link uniqueness

Enforce:
- one external identity key cannot be linked to two GRIDWORKS accounts;
- duplicate delivery of the same valid link operation is idempotent;
- linking an already-linked external identity to a different account rejects;
- linking a second provider to the same account is allowed;
- multiple identities from the same provider are allowed only if explicitly distinguishable and policy-compatible.

Document the exact uniqueness key, e.g. issuer + subject.

## 13. Guest → protected linking

The linked identity attaches to the **existing** guest account.

Must preserve:
- account ID;
- player ID;
- profile;
- later company ownership references;
- economic/ledger references;
- facilities/inventory;
- reputation;
- purchases/entitlements;
- messages/consortium/challenge history when those domains exist.

WP-007 proves preservation of account/player/profile IDs and provides schema references suitable for later domains.

Do not create a second player and migrate state.

## 14. Proof-of-control boundary

ADR-002 requires proof of control of both:
- the guest session;
- the external identity.

WP-007 must model this as an explicit linking precondition.

Repository proof may use a verified `ExternalIdentityAssertion`/equivalent supplied by a trusted future auth adapter.

Do **not** implement or fake Apple/Google signature validation inside this work package.

The domain function must not accept arbitrary unverified provider claims.

## 15. Account merge policy

Normal self-service account merge is out of scope.

If the external identity is already linked to another progressed account:
- reject with a typed conflict;
- do not automatically merge;
- do not transfer economic state;
- preserve an explicit future support/recovery path.

## 16. Player profile model

Minimum persisted profile fields:
- player_id;
- account_id;
- unique_handle;
- display_name;
- avatar_asset_id nullable;
- bio nullable/bounded;
- locale;
- timezone preference;
- reputation summary placeholder/reference;
- public-achievement visibility;
- privacy settings;
- DM permissions;
- notification preference foundation;
- moderation state/reference where structurally required;
- created/updated metadata.

Keep reputation, achievements, moderation and notification domains bounded; use references/settings rather than implementing those later systems now.

## 17. Handle rules

`unique_handle` is public and globally unique within the logical world.

Requirements:
- case-insensitive uniqueness;
- Unicode normalization;
- Unicode/confusable awareness;
- reserved/system names;
- banned impersonation-like forms;
- bounded length;
- deterministic canonical comparison key;
- display/original handle may be preserved separately from canonical uniqueness key if useful.

At minimum protect obvious collisions such as:
- `GridWorks` vs `gridworks`;
- composed vs decomposed Unicode;
- visually confusable Latin/Cyrillic forms where feasible;
- whitespace/invisible-control abuse.

Do not rely solely on PostgreSQL `lower()` for Unicode/confusable safety.

## 18. Confusable strategy

Choose and document a deterministic strategy.

Acceptable approaches:
- Unicode NFKC + casefold + confusable skeleton using a pinned Unicode data/version source;
- a well-defined library with pinned version and tests.

Requirements:
- no online service;
- deterministic;
- versioned behavior;
- test vectors;
- future migration/version implications documented.

If the current dependency set cannot provide a robust confusable skeleton without adding a justified dependency, Engineering may add a pinned library and must document why.

## 19. Reserved handles

Reserve at least system/security-sensitive names and close variants:
- GRIDWORKS;
- admin / administrator;
- moderator / support;
- system / official;
- CortexSSG or operator-branded names if appropriate.

Exact reserved list should be data-driven/configurable rather than deeply hard-coded.

Do not reserve arbitrary common player names unnecessarily.

## 20. Display names

`display_name`:
- is mutable;
- need not be unique;
- must be Unicode-safe;
- bounded in length;
- must reject control/invisible abuse;
- is not an authentication credential;
- does not grant namespace ownership.

Do not reuse unique-handle canonicalization as a reason to make display names globally unique.

## 21. Locale/timezone

Profile must persist:
- locale preference;
- timezone preference.

Validation:
- locale uses a bounded BCP-47-compatible representation or canonical locale key;
- timezone uses IANA TZ identifiers or a stable equivalent;
- invalid free-form values reject.

WP-021 will own full localization packs. WP-007 stores the preference only.

## 22. Privacy and DM policy foundation

Persist bounded settings sufficient for later social WP:
- profile visibility;
- DM permission policy;
- optional discoverability/searchability flag;
- block-list relation may be schema-reserved but full blocking behavior belongs later.

Example DM policies:
- Everyone;
- Contacts/Consortium later;
- Nobody.

Do not implement messaging in WP-007.

## 23. Notification preference foundation

Persist per-category/channel preference structure suitable for later WP-020.

Do not build push delivery.

At minimum support stable categories/channel keys or a structured preference record without translated prose.

## 24. Authentication vs public identity

Never expose private authentication data as profile state.

Public profile APIs later must not include:
- provider subjects;
- session digests;
- recovery tokens;
- raw email unless a future privacy surface explicitly permits it;
- account security metadata.

Design schemas so private/public boundaries are clear.

## 25. Database schema/migration

Add a PG18-targeted migration for identity/profile foundations.

Expected tables or equivalent:
- accounts;
- players;
- player_profiles;
- account_sessions / guest_sessions;
- account_identity_links;
- identity_mutation_receipts / idempotency records if not reusing a safe generic receipt boundary;
- reserved_handles / handle policy data if persisted;
- player_preferences or bounded JSON/settings table if justified.

Use explicit:
- PK/FK;
- uniqueness;
- status checks;
- timestamps;
- session verifier uniqueness where appropriate;
- external identity uniqueness;
- handle canonical/skeleton uniqueness.

No credentials/roles/provisioning SQL.

## 26. Transactional guest issuance

Document and implement repository/domain semantics where guest account, player, profile and initial session are one transaction.

If any insert fails:
- no orphan account;
- no orphan player;
- no profile without account/player;
- no issued secret without durable verifier.

Use a repository abstraction / domain transaction contract, not a live DB requirement.

## 27. Idempotency

Guest issuance and linking mutations must be replay-safe.

Requirements:
- explicit request/idempotency key;
- duplicate same request returns/reconstructs the same durable result where safe;
- same idempotency key with contradictory payload rejects;
- concurrent duplicate submissions cannot create two accounts or two links.

Use DB uniqueness as the correctness boundary, not process-local locks.

## 28. Session revocation

Model:
- active;
- revoked;
- expired.

A revoked/expired session verifier must fail authentication.

Guest account linking must not silently invalidate all sessions unless policy explicitly says so; document chosen behavior.

Provide a session-rotation/revocation foundation for later API work.

## 29. Recovery scaffolding

Implement domain/schema scaffolding for account recovery/protection state, not full email/provider recovery flows.

Support future concepts such as:
- protected account;
- verified linked identity;
- recovery challenge/reference;
- support-assisted recovery marker.

Do not store plaintext recovery secrets.

## 30. Security events / audit

Persist or model auditable identity mutations:
- guest account created;
- identity linked;
- session revoked/rotated;
- handle changed;
- security-sensitive preference changes where appropriate.

Use stable semantic event/action codes.

Do not store secrets in audit payloads.

## 31. Public/private DTO/schema separation

Add strict schemas/domain types for at least:
- private account record;
- public player profile;
- internal player profile;
- guest issuance result;
- verified identity-link request;
- session metadata;
- profile update command.

Public profile schema must not accidentally permit private auth fields.

## 32. No API implementation yet

WP-009 owns the Go API foundation.

WP-007 may implement:
- identity/profile domain packages;
- repository interfaces;
- migrations;
- pure service/domain functions;
- schema contracts.

Do not implement:
- HTTP routers;
- REST endpoints;
- WebSocket;
- public server listeners;
- nginx/systemd deployment.

If Go domain code is introduced, keep it transport-independent.

## 33. No live provider integration

Do not:
- create Apple Service IDs;
- create Google OAuth clients;
- call provider discovery endpoints;
- fetch JWKS;
- perform live OIDC redirects;
- add production client secrets.

Define adapter interfaces / verified assertion inputs only.

Provider-specific implementation belongs later when API/auth integration is explicitly authorized.

## 34. No company ownership yet

WP-008 owns company/ownership.

WP-007 may reserve/reference future principal IDs but must not create:
- companies;
- company shares;
- company ownership transfer;
- company groups;
- JVs.

## 35. SQLite/client cache

Do not make local SQLite authoritative for identity.

Client-side secure-storage integration is not required in WP-007.

Document the contract that long-lived session material must eventually be stored in platform secure storage.

## 36. Typed errors

At minimum:
- duplicate account/player ID;
- duplicate guest issuance idempotency;
- invalid/expired/revoked session;
- invalid token/verifier;
- external identity already linked;
- cross-account link conflict;
- unverified external assertion;
- invalid handle;
- handle unavailable;
- reserved handle;
- confusable collision;
- invalid display name;
- invalid locale/timezone;
- invalid account/profile state;
- duplicate mutation/idempotency conflict.

Errors must not leak sensitive verifier/provider details.

## 37. Tests — guest identity

Required:
- new guest issuance creates stable account + player + profile + session;
- no duplicate identities on replay;
- contradictory idempotency payload rejects;
- no orphaned state on staged failure;
- guest session verifier authenticates;
- wrong token fails;
- revoked/expired session fails;
- raw token is absent from persisted domain record.

## 38. Tests — linking

Required:
- verified external identity links to existing guest account;
- account ID/player ID remain unchanged;
- same link replay is idempotent;
- same provider identity cannot link to a second account;
- unverified assertion rejects;
- linking a second distinct provider to same account succeeds;
- automatic two-account merge is rejected.

## 39. Tests — handles/profile

Required:
- casefold collision;
- Unicode normalization collision;
- at least one Latin/Cyrillic confusable collision;
- reserved handle rejection;
- valid Unicode handle;
- display-name non-uniqueness;
- control/invisible abuse rejection;
- valid/invalid locale;
- valid/invalid timezone;
- public profile serialization excludes private account/session/link fields.

## 40. Tests — SQL/schema

Required repository/static validation:
- migration ordered;
- expected tables/constraints/indexes present;
- provider uniqueness present;
- handle canonical/skeleton uniqueness present;
- session verifier uniqueness/revocation fields present;
- no plaintext-token column;
- no destructive/provisioning SQL;
- schemas reject unknown/private fields on public surfaces.

Live PG execution is not required without an already-authorized test environment.

## 41. Concurrency assumptions

Document:
- guest issuance concurrent duplicate handling;
- external link race handling;
- handle-claim race handling.

Correctness must rely on transactional DB uniqueness/locking semantics, not only pre-checks.

## 42. Performance posture

Use indexes for:
- account/player lookup;
- session verifier lookup;
- external issuer+subject;
- handle canonical/skeleton;
- profile player lookup.

Avoid premature search infrastructure.

PostgreSQL remains sufficient.

## 43. Allowed paths

Primary:
- new identity/profile domain package under `services/**` or a dedicated shared package;
- `db/migrations/**`;
- `packages/schemas/**`;
- `packages/content/config/**` for reserved handle policy if justified;
- `tests/**`;
- `docs/evidence/**`;
- build/CI files only when required.

Minimal supportive changes elsewhere must be explained.

## 44. Out of scope

Do not implement:
- WP-008 company ownership;
- WP-009 public API;
- live OAuth/OIDC provider integration;
- email delivery;
- messaging/social behavior;
- notifications delivery;
- purchases/entitlements;
- moderation workflow;
- deployment;
- host/database provisioning.

## 45. Evidence

Create concise evidence documenting:
- identity hierarchy;
- guest issuance transaction;
- ID strategy;
- session secret/verifier design;
- external identity link uniqueness;
- proof-of-control boundary;
- guest→linked preservation;
- merge rejection policy;
- handle normalization/confusable algorithm + Unicode data/library version;
- profile/public-private boundaries;
- SQL constraints/indexes;
- concurrency/idempotency;
- tests;
- live-provider/live-PG limitations.

## 46. Done-when

Architecture can verify:

1. exact parent;
2. guest-first account/player/profile/session domain exists;
3. immutable account/player IDs;
4. raw session secrets are never persisted;
5. session verification/revocation/expiry works;
6. verified external link attaches to existing account;
7. cross-account identity collision rejects;
8. normal self-service account merge is absent;
9. player IDs survive guest→linked conversion;
10. Unicode/case/confusable-aware unique handles exist;
11. public profile is separated from auth-private data;
12. locale/timezone/privacy/DM/notification foundations persist;
13. PG18 migration has required integrity/uniqueness constraints;
14. concurrency/idempotency boundaries are documented/tested;
15. no WP-008/WP-009/provider/deployment scope creep;
16. all CI/security gates pass;
17. complete GitHub handback exists.

## 47. Handback

Post to WP-007 issue and draft PR:
- branch;
- exact parent/head;
- commit chain;
- changed files/dependencies;
- ID model;
- guest issuance/session design;
- token/verifier handling;
- external link model;
- handle/confusable strategy;
- profile/public-private design;
- migration list and integrity constraints;
- concurrency/idempotency behavior;
- test evidence;
- CI links/results;
- live provider/PG limitations;
- deviations/risks/conflicts;
- explicit no-host/deployment statement.

## 48. Stop condition

After handback, Engineering stops.

Do not start WP-008 without separate authorization.

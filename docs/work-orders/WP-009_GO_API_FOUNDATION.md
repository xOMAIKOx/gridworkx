# GRIDWORKS — WP-009 Go API Foundation

**Status:** AUTHORIZED FOR ENGINEERING  
**Architecture owner:** Architecture  
**Implementation model:** GPT-5.6 Luna XHigh  
**Review:** GPT-5.6 Sol-class Architecture review  
**Exact required parent:** `f93c6526f9ce8743406606265456b283c8d8787c`

## 1. Purpose

Implement the first production-shaped, versioned Go HTTP API boundary for GRIDWORKS.

WP-009 must turn the accepted WP-006/WP-007/WP-008 server-authoritative domains into a coherent mobile/public API foundation without:
- duplicating Rust simulation mathematics;
- exposing auth-private data;
- bypassing idempotency/ownership invariants;
- inventing later-WP business flows;
- using in-memory state as a production persistence path.

The result must establish the transport, contract, authentication, error, OpenAPI, PostgreSQL-adapter and service-composition patterns that later API groups extend.

## 2. Controlling references

Engineering must read and obey:
- `AGENTS.md`
- `AI_AGENT_COLLABORATION_PROTOCOL.md`
- Product Philosophy / Game Design
- DP2 Core Systems Specification
- DP3 Engineering Architecture & Work Packages
- `docs/decisions/ADR-002_IDENTITY_AND_ACCOUNT_LINKING.md`
- `docs/decisions/ADR-003_SCHEMA_AND_SERIALIZATION.md`
- `docs/decisions/ADR-005_INITIAL_SERVER_TOPOLOGY.md`
- accepted WP-006 persistence/ledger contract
- accepted WP-007 identity/profile contract
- accepted WP-008 company/ownership contract
- Architecture handoff issue #18

If implementation discovers an architecture conflict, stop that path and raise evidence on the WP-009 issue.

## 3. Exact ancestry

Branch from exactly:

`f93c6526f9ce8743406606265456b283c8d8787c`

Verify:
- exact parent;
- clean working tree;
- no unrelated changes.

Do not silently rebase to a later `main`.

## 4. Runtime and host boundary

Repository-only.

No:
- ERIS or other host mutation;
- PostgreSQL server/role/database provisioning;
- systemd installation/start/restart;
- nginx/reverse-proxy changes;
- DNS/firewall/WireGuard changes;
- port activation;
- NATS provisioning;
- Apple/Google/OIDC client creation;
- deployment;
- Docker/Podman/Compose/Kubernetes/OCI/containerd/nerdctl.

Repository code, config templates, systemd-unit source adjustments required by the binary contract, schemas, OpenAPI, migrations, tests and CI are allowed.

## 5. Process topology

WP-009 implements the **gridworks-api** process only.

Initial process boundaries remain:
- `gridworks-api` — HTTP API;
- `gridworks-worker` — later/background;
- `gridworks-realtime` — later WebSocket/realtime;
- `gridworks-sim-validator` — Rust deterministic validation boundary.

Do not collapse realtime/worker/sim-validator responsibilities into the API process.

Do not create new microservices.

## 6. Bind and ingress model

The API process continues to bind loopback by default:

`127.0.0.1:18080`

External HTTPS is terminated by the shared reverse-proxy layer when deployment is separately authorized.

The API binary must not:
- self-provision certificates;
- bind public interfaces by default;
- assume direct Internet exposure.

Repository config may allow an explicit bind override.

## 7. API versioning

Public/mobile routes must live under:

`/api/v1`

Operational endpoints remain unversioned:
- `GET /healthz`
- `GET /readyz`
- `GET /version`

Breaking public contract changes require a new API version rather than silently changing v1 semantics.

## 8. Canonical representation

Per ADR-003:
- external representation is JSON;
- contracts are repository-owned;
- HTTP contracts are described by OpenAPI;
- externally meaningful payloads are schema-versioned or unambiguously associated with a versioned contract;
- semantic IDs are used instead of translated prose/status copy.

Do not introduce protobuf/gRPC as the public/mobile API in WP-009.

## 9. OpenAPI

Add a canonical OpenAPI 3.1 contract for the actual WP-009 routes.

Preferred location:

`packages/schemas/openapi/gridworks-api-v1.yaml`

or an equivalently clear repository-owned schema path.

Requirements:
- route/method definitions;
- request/response schemas;
- auth scheme;
- idempotency header where required;
- error envelope;
- status codes;
- examples;
- explicit API version;
- no routes for unimplemented later domains.

CI must parse/validate the OpenAPI document and detect handler/contract drift through an explicit route inventory or equivalent test.

Do not publish speculative endpoints as if implemented.

## 10. HTTP implementation posture

Use a small modular Go HTTP layer.

The Go standard library is preferred unless a routing dependency provides clear value.

Avoid:
- framework-heavy abstractions;
- reflection magic;
- service-locator/global mutable state;
- generated boilerplate that obscures authorization;
- generic CRUD scaffolding.

Handler modules should be explicit and domain-oriented.

## 11. API composition

Create a clear application composition boundary, e.g.:

- `services/api/internal/httpapi`;
- `services/api/internal/app`;
- or equivalent.

Handlers must depend on narrow application/domain interfaces, not global concrete stores.

Test composition may use in-memory/fake dependencies.

Normal production composition must not silently instantiate WP-007/WP-008 in-memory stores as authoritative state.

## 12. Production persistence path

WP-009 must establish the PostgreSQL-backed application persistence path for the stateful routes it exposes.

Use the accepted WP-007/WP-008 tables and invariants.

Preferred:
- Go `database/sql`;
- PostgreSQL driver such as pinned `pgx/v5/stdlib`;
- explicit SQL;
- no ORM.

Requirements:
- transactional operations;
- context deadlines/cancellation;
- parameterized queries;
- no SQL string concatenation for user values;
- no secret/DSN logging;
- no process-local lock as durable correctness;
- accepted database constraints/triggers remain authoritative.

Do not rewrite accepted WP-006/007/008 migrations merely to suit API code.

If an additive API-support schema change is genuinely required, add a new ordered migration and explain it. Do not mutate accepted historical migrations unless Architecture explicitly approves.

## 13. Avoid policy duplication

PostgreSQL adapters must persist the accepted domain rules; they must not create a second divergent implementation of:
- identity normalization;
- session credential policy;
- company name normalization;
- ownership arithmetic;
- idempotency semantics.

Extract/reuse pure domain validation and canonicalization where necessary.

The API must not duplicate Rust facility/simulation mathematics at all.

## 14. Database configuration

Normal API startup must require explicit durable database configuration for stateful routes.

At minimum support:
- database DSN/config from environment or equivalent runtime config;
- connection open;
- bounded pool settings;
- readiness ping.

Requirements:
- credentials never logged;
- configuration errors fail closed;
- no automatic schema creation/migration on normal API process startup;
- no fallback to in-memory production state.

An explicit test/dev composition may use in-memory stores, but it must be unmistakably non-production and must not be the default binary behavior.

## 15. Server lifecycle

Refactor the existing minimal runtime server as necessary to support:
- injected handler;
- `http.Server` timeouts;
- graceful shutdown;
- request context cancellation;
- bounded header size;
- clean SIGTERM/SIGINT handling in `gridworks-api`.

At minimum configure sensible:
- ReadHeaderTimeout;
- ReadTimeout or equivalent body/read control;
- WriteTimeout;
- IdleTimeout;
- MaxHeaderBytes.

Do not add aggressive timeouts that break valid mobile requests; document chosen values.

## 16. Operational endpoints

### `GET /healthz`

Liveness only.

Must not expose:
- secrets;
- DSN;
- provider config;
- internal topology details beyond safe role/version metadata.

### `GET /readyz`

Readiness.

For normal durable composition, readiness should fail when required stateful dependencies are unavailable.

Do not perform destructive checks.

### `GET /version`

Return safe build/API version metadata.

No Git credentials, environment dump or internal secrets.

## 17. Request identity / correlation

Every request must have a server-controlled request ID.

Requirements:
- generate a request ID when absent;
- if an inbound correlation/request value is accepted, validate length/character set before reflecting it;
- include request ID in responses/errors;
- make it available to logs/context;
- never treat it as an authorization or idempotency credential.

Correlation IDs for future domain events should have a clean seam.

## 18. Common JSON behavior

For API JSON requests:
- require appropriate `Content-Type: application/json`;
- reject unsupported media types;
- apply bounded body size;
- reject malformed JSON;
- reject unknown fields for mutation contracts;
- reject trailing multiple JSON values;
- return stable machine-readable errors.

Suggested default body limit for ordinary v1 JSON commands: 1 MiB or smaller if Engineering justifies it.

File/media upload is out of scope.

## 19. Error contract

Define a stable error envelope.

At minimum:
- machine code;
- safe human message;
- request ID;
- optional field/path details;
- server timestamp if common response metadata uses one.

Never expose:
- SQL text;
- stack traces;
- DSNs;
- token digests;
- provider subjects;
- raw session tokens;
- internal filesystem paths.

Map domain errors deliberately.

Suggested status classes:
- 400 malformed request;
- 401 missing/invalid/expired/revoked session;
- 403 authenticated but forbidden;
- 404 hidden/not found;
- 409 idempotency/name/state conflict;
- 413 body too large;
- 415 unsupported content type;
- 422 semantically invalid command where appropriate;
- 500 generic internal error.

## 20. Panic recovery

A handler panic must:
- be recovered at the HTTP boundary;
- produce a generic 500 error with request ID;
- not leak stack/internal state to the client;
- allow server process continuity where safe.

Server logs may capture diagnostic stack data but must redact request secrets.

## 21. Security headers / cache behavior

Set appropriate API headers, at minimum:
- `Content-Type`;
- `X-Content-Type-Options: nosniff`.

Auth/session responses must use:
- `Cache-Control: no-store`.

Do not add browser-specific CORS permissiveness by default.

Native mobile clients do not require wildcard CORS.

## 22. Authentication transport

Use bearer session authentication for authenticated v1 routes.

Expected:

`Authorization: Bearer <session-secret>`

Requirements:
- compute/lookup verifier/digest; do not persist raw bearer secret;
- invalid/expired/revoked token => 401;
- derive authenticated Account/Player principal server-side;
- do not accept account/player actor IDs from the request as authority;
- raw bearer token must never appear in logs/errors/audit payloads.

Extend the accepted identity repository with a bounded session-authentication query if required.

## 23. Authenticated principal

Create a narrow request principal/context type containing only what handlers need, e.g.:
- account ID;
- player ID;
- account protection/status if required for policy;
- current session ID.

Do not include:
- raw bearer token;
- token digest in public DTOs;
- provider subject;
- provider access/refresh token.

## 24. Public WP-009 route set

Implement only routes backed by already accepted domains.

Required public routes:

### Authentication/session
- `POST /api/v1/auth/guest`
- `DELETE /api/v1/auth/session`

### Player/profile
- `GET /api/v1/me`
- `PATCH /api/v1/me/profile`
- `GET /api/v1/players/{player_id}`

### Company
- `POST /api/v1/companies`
- `GET /api/v1/companies/{company_id}`
- `POST /api/v1/company-groups`

Additional tiny read endpoints may be added only when directly required to support these contracts and must be documented.

Do not expose generic table CRUD.

## 25. Guest issuance endpoint

`POST /api/v1/auth/guest`

Requirements:
- no prior auth;
- requires `Idempotency-Key`;
- server generates account/player/session IDs and raw session secret;
- creates account/player/profile/session atomically;
- persists only session verifier/digest;
- initial response returns raw session secret once;
- response includes safe account/player/profile/session metadata;
- `Cache-Control: no-store`.

### Guest replay semantics

Durable receipt state does not contain the raw session secret.

Therefore a replay of an already committed guest issuance must **not** pretend it can recover the original secret.

Define and document deterministic replay semantics that:
- return the same safe identity/result references;
- do not return the raw token again;
- clearly indicate that the credential is unavailable on durable replay / recovery is required.

Do not store plaintext session secret merely to make HTTP idempotency cosmetically identical.

## 26. Session authentication query

Add/reuse a repository operation that:
- hashes/verifies the presented bearer token;
- resolves active session;
- checks expiry and revocation;
- resolves account + player;
- returns the safe authenticated principal.

Successful auth may update `last_used_at` using a documented strategy.

Do not make every read depend on provider identity.

## 27. Session revocation endpoint

`DELETE /api/v1/auth/session`

Requirements:
- authenticated current session only;
- revokes the presented session;
- idempotent/replay-safe behavior;
- does not revoke unrelated sessions unless a later policy explicitly adds that operation;
- no token in response.

## 28. External identity linking

Do **not** expose a public HTTP endpoint that accepts client-supplied:

`verified: true`

or otherwise trusts an arbitrary client assertion.

WP-007 requires a trusted verified-assertion boundary.

Live Apple/Google/OIDC verification is not authorized in WP-009.

You may define an internal adapter/application seam for a future verified external identity provider, but no public linking route becomes active until proof verification is implemented and authorized.

## 29. `GET /api/v1/me`

Authenticated.

Return a self-safe view containing:
- player identity;
- profile fields;
- safe account protection/status information if useful;
- no session verifier;
- no provider subject;
- no recovery/security internals.

This may be richer than the public player profile because it is the authenticated owner’s own resource, but keep auth secrets excluded.

## 30. Profile update

`PATCH /api/v1/me/profile`

Implement a bounded WP-007 profile mutation.

Allowed fields:
- display name;
- bio;
- locale;
- timezone;
- visibility;
- DM policy;
- discoverability;
- notification-preference foundation.

Out of scope:
- handle rename;
- avatar binary upload;
- provider/auth changes;
- moderation state mutation.

Requirements:
- authenticated player derived from session;
- strict patch schema;
- validation reuses WP-007 validators;
- `updated_at` is server-owned;
- mutation uses the identity-domain idempotency namespace;
- requires `Idempotency-Key`;
- same key/same payload replays;
- contradictory reuse conflicts;
- no partial update on failure.

If an additive migration is needed for profile mutation receipt result metadata, use a new migration rather than rewriting WP-007 migration.

## 31. Public player profile

`GET /api/v1/players/{player_id}`

Return only the accepted `PublicPlayerProfile`.

Requirements:
- no account ID;
- no session/security state;
- no provider identity;
- honor profile visibility;
- private/non-visible profiles should fail in a privacy-preserving manner;
- unknown and non-public should not create unnecessary account-existence leakage.

Document exact 404/visibility behavior.

## 32. Company creation

`POST /api/v1/companies`

Authenticated.

Public request may include only safe creation fields such as:
- company type;
- name.

Optional already-accepted public profile fields may be included only if the WP-008 domain supports them without inventing new ownership/lifecycle semantics.

Critical authorization:
- initial Player owner is **derived from authenticated `player_id`**;
- client must not choose an arbitrary owner principal;
- public API must not permit creation as GRIDWORKS system principal;
- server generates company ID;
- requires `Idempotency-Key`;
- exact accepted WP-008 name normalization/reserved-name/idempotency contract applies.

Do not add money, fees, share sales or acquisition semantics.

## 33. Company type restrictions

Public/mobile creation must not expose privileged/system entity creation.

At minimum:
- operating company: allowed;
- holding/company type: allowed if product baseline supports player-created holding entities;
- system type: forbidden on public route.

System-company creation remains an internal/trusted capability.

## 34. Public company profile

`GET /api/v1/companies/{company_id}`

Return only the accepted public company profile.

Do not expose:
- private ownership audit rows;
- idempotency receipts;
- internal name skeleton/canonical key unless intentionally public;
- account/auth/provider data;
- hidden system flags.

Honor public/private company visibility.

## 35. Holding/group creation

`POST /api/v1/company-groups`

Authenticated.

Requirements:
- name only plus already accepted safe profile inputs if any;
- owner is authenticated player, server-derived;
- no system-owner selection;
- server-generated group ID;
- shared company/group name-claim namespace;
- requires `Idempotency-Key`;
- returns a safe group creation result.

A public group-read endpoint is optional only if Engineering adds a strict public group DTO and documents it. Do not expose raw ownership internals.

## 36. Ownership-transfer API is deliberately not public in WP-009

Do **not** expose WP-008’s raw ownership-transfer primitive as a public/mobile route yet.

Reason:
- WP-017 owns sale/acquisition orchestration;
- future settlements may require ledger/payment invariants;
- exposing a raw transfer now could bypass commercial and authorization policy.

The typed ownership-transfer domain remains available internally for later authorized orchestration.

Similarly, do not expose GRIDWORKS buyer-of-last-resort behavior in WP-009.

## 37. Group membership mutation API

Do not expose raw assign/detach/reassign publicly unless Engineering can prove an accepted authorization policy from existing product/design documents.

WP-009 does not invent corporate-governance/control thresholds.

The domain mechanics remain internal and can be surfaced by a later explicitly authorized work package.

## 38. Idempotency header

State-changing API operations that require idempotency use:

`Idempotency-Key`

Requirements:
- bounded length;
- reject empty/control-character values;
- do not log as a secret-bearing header dump;
- pass the exact semantic key into the appropriate identity/company mutation namespace;
- same key + same mutation/payload => replay;
- same key + different mutation/payload => 409 conflict;
- no process-local idempotency map as durable production boundary.

GET routes do not require idempotency keys.

## 39. Server-owned time

Client time is never authoritative.

For mutations:
- use server UTC time;
- return RFC3339/ISO-8601 UTC timestamps where exposed;
- ignore/reject client attempts to set authoritative created/updated/effective timestamps.

Do not implement simulation-time logic in Go API.

## 40. Stable IDs

Path/request IDs must be syntax-validated before use.

Do not infer authorization from ID shape.

Do not expose sequential database IDs.

Do not let clients choose authoritative account/player/company/group/session IDs.

## 41. PostgreSQL identity adapter

For exposed identity routes, implement the durable repository/application operations required for:
- guest issuance;
- session authentication;
- current session revocation;
- self profile read;
- public profile read;
- bounded profile update/idempotency.

Use accepted tables:
- `accounts`;
- `players`;
- `player_profiles`;
- `guest_sessions`;
- `identity_mutation_receipts`.

External identity-link persistence can remain unexposed unless needed by an internal seam.

Raw session secret must never be written to PostgreSQL.

## 42. PostgreSQL company adapter

For exposed company/group routes, implement durable operations required for:
- player-owned company creation;
- public company read;
- player-owned group creation.

Use accepted WP-008 constraints/triggers:
- shared name claims;
- owner/player integrity;
- exact ownership;
- append-only history;
- idempotency receipts.

Do not bypass triggers with alternate tables or manual “temporary” state.

Public handlers must not accept the system principal as a caller-selected owner.

## 43. Player resolver in durable composition

WP-008’s `PlayerResolver` must resolve through the durable identity authority in normal API composition.

Do not replace it with a local registration map.

PostgreSQL-backed identity lookup is the production resolver.

In-memory/fake resolvers remain test-only.

## 44. Transaction boundaries

Guest issuance transaction must atomically cover:
- account;
- player;
- profile;
- session verifier;
- receipt.

Profile update transaction must atomically cover:
- profile change;
- receipt.

Company/group creation transaction must atomically cover:
- entity;
- name claim via accepted trigger;
- initial ownership;
- history;
- receipt.

If any part fails:
- rollback all state;
- do not return a successful HTTP mutation.

## 45. DB error mapping

Map known PostgreSQL constraint/conflict conditions to typed domain/application errors without leaking SQL.

At minimum distinguish:
- duplicate/idempotency conflict;
- handle/name collision;
- reserved name;
- unknown player principal;
- invalid state;
- generic persistence failure.

Do not key correctness solely off fragile human-readable PostgreSQL error strings when a SQLSTATE/constraint identifier is available.

## 46. Request/response DTO separation

HTTP DTOs must not directly serialize arbitrary internal structs.

Use explicit API DTOs.

Examples:
- guest issuance response;
- self profile response;
- public profile response;
- company create request/result;
- public company response;
- group create request/result;
- error envelope.

This prevents future internal fields from silently becoming public.

## 47. Logging

Use structured or consistently parseable logs.

At minimum include:
- request ID;
- method;
- route/template where safe;
- status;
- duration;
- safe principal IDs after authentication where appropriate.

Never log:
- Authorization header;
- raw session secret;
- token digest;
- database DSN/password;
- provider access/refresh tokens;
- full request body by default.

## 48. Metrics / telemetry seam

Do not add a full metrics stack.

It is acceptable to add a small instrumentation seam/counters interface for later observability.

No Prometheus/OTel infrastructure deployment is authorized.

Avoid hard-wiring product analytics into handlers.

## 49. Realtime boundary

WP-009 does not implement WebSocket.

Do not add DM/market/facility realtime handlers to `gridworks-api`.

`gridworks-realtime` remains the accepted WebSocket process.

HTTP may later return URLs/capability metadata, but no speculative realtime contract is required now.

## 50. NATS boundary

No NATS dependency is required merely to satisfy WP-009.

If domain-event publication is architecturally necessary for an exposed mutation:
- define a transactional/outbox-safe seam;
- do not publish before DB commit;
- do not create a dual-write correctness bug.

Prefer leaving actual event publication for the work package that owns the consuming domain unless directly required.

No NATS provisioning or connection to estate infrastructure is authorized.

## 51. Rust/simulation boundary

Go API must not calculate:
- facility output;
- bottlenecks;
- production transforms;
- wear;
- failure effects;
- deterministic advancement;
- simulation KPIs.

Future simulation commands are forwarded/validated through accepted Rust boundaries.

WP-009 should not add facility/world simulation routes simply to fill out the API catalogue.

## 52. Later-domain route exclusions

Do not implement routes for:
- facilities/world;
- inventory/material mutations;
- construction/repair;
- managers/recruitment;
- contracts;
- market/trading;
- Consortium/JV;
- business opportunities/acquisition;
- messaging;
- notifications delivery;
- challenges/leaderboards;
- real-estate transaction flows;
- finance/exchange.

Those route groups are added by the work packages that own those domains.

## 53. No raw database administration endpoints

Do not add:
- generic query endpoint;
- debug SQL endpoint;
- migration endpoint;
- environment/config dump endpoint;
- admin backdoor.

Admin/back-office is WP-022.

## 54. Dependency policy

New Go dependencies must be:
- justified;
- pinned through `go.mod/go.sum`;
- actively maintained;
- minimal.

No web framework is required to complete this WP.

If `pgx` or an OpenAPI validator/parser is added, record the exact version and purpose in evidence.

## 55. Test strategy — routing and HTTP contract

Required:
- operational endpoint behavior;
- `/api/v1` route versioning;
- wrong method => 405 with Allow where appropriate;
- unknown route => safe 404;
- wrong content type => 415;
- malformed JSON => 400;
- unknown JSON field => reject;
- multiple JSON documents/trailing payload => reject;
- oversized body => 413;
- request ID always present;
- panic recovery => safe 500;
- auth responses => no-store;
- no wildcard CORS by default.

## 56. Tests — authentication/session

Required:
- guest issue first response returns raw token once;
- persisted record contains verifier/digest only;
- guest replay returns same safe identity refs and does not recover raw token;
- valid bearer resolves correct account/player;
- wrong token => 401;
- expired => 401;
- revoked => 401;
- revoke current session invalidates subsequent use;
- auth middleware never trusts caller-supplied player/account ID.

## 57. Tests — profile

Required:
- self read authenticated;
- public profile contains only public schema fields;
- private profile hidden according to documented policy;
- profile update validates locale/timezone/display/bio/privacy fields;
- profile update uses server `updated_at`;
- idempotent replay;
- contradictory idempotency reuse;
- unknown/private fields rejected;
- handle/auth/provider/security mutation cannot be smuggled through patch.

## 58. Tests — company/group

Required:
- authenticated player creates operating company and becomes initial owner;
- holding creation if public policy allows it;
- system company creation rejected publicly;
- caller cannot choose another player as initial owner;
- duplicate/case/confusable/reserved names map to stable conflict errors;
- public company response excludes private ownership/audit/auth state;
- private company visibility honored;
- authenticated player creates group as themselves;
- group/company shared-name namespace conflict preserved;
- client cannot choose system owner;
- raw ownership-transfer route is absent.

## 59. Tests — PostgreSQL adapters

Without a live PG18 environment, provide repository-level evidence using appropriate SQL adapter tests/mocks/static checks for:
- transaction begin/commit/rollback boundaries;
- parameterized queries;
- no raw token persistence;
- digest lookup;
- receipt/idempotency use;
- server-derived player ownership;
- constraint/SQLSTATE mapping;
- context cancellation.

Do not pretend mocked SQL proves PostgreSQL trigger semantics.

The later live-PG environment gate remains explicit.

## 60. Tests — OpenAPI

Required:
- OpenAPI parses/validates;
- every implemented public v1 route appears in OpenAPI;
- no documented unimplemented route;
- security scheme applied to authenticated operations;
- idempotency header documented on required mutations;
- error response schema reused consistently.

## 61. Tests — graceful server behavior

Required:
- server handler test does not require public port;
- graceful shutdown path is testable;
- configured server timeouts are non-zero;
- loopback remains default bind;
- runtime errors propagate to process exit rather than being silently swallowed.

## 62. Security scan expectations

Existing:
- gitleaks;
- forbidden-runtime scan;
- ShellCheck;
- Go vet/tests;
- schema validation;
- Godot/Rust gates

must remain green.

No secret fixtures that resemble deployable credentials.

## 63. Performance posture

Do not optimize prematurely.

Required:
- bounded request bodies;
- DB contexts/deadlines;
- connection pooling;
- no unbounded per-request goroutine fan-out;
- no full-table scans for session token lookup (digest is indexed/unique);
- indexed ID/profile/company reads.

No Valkey/cache required.

## 64. Allowed paths

Primary:
- `services/api/**`
- `services/internal/runtime/**`
- `services/internal/identity/**` for bounded API/repository seams
- `services/internal/company/**` for bounded API/repository seams
- new `services/internal/postgres/**` or equivalent durable adapter package
- `packages/schemas/**`
- `db/migrations/**` only for new additive WP-009 support migration if justified
- `tests/**`
- `docs/evidence/**`
- `docs/api/**` if used
- `go.mod/go.sum`
- CI/build files only when required.

Minimal supportive edits elsewhere must be explained.

## 65. Existing migrations are accepted history

Treat merged WP-006/WP-007/WP-008 migrations as accepted historical migrations.

Do not edit:
- `0002_wp006_persistence_ledger.sql`
- `0003_wp007_identity_profile.sql`
- `0004_wp008_company_ownership.sql`

If WP-009 needs additive schema support, create the next ordered migration and document exactly why.

## 66. Evidence

Create:

`docs/evidence/WP-009_GO_API_FOUNDATION.md`

Document:
- route inventory;
- versioning;
- OpenAPI path;
- authentication/session flow;
- guest replay/raw-token semantics;
- principal derivation;
- public/private DTO boundaries;
- profile update/idempotency;
- company/group ownership derivation;
- PostgreSQL adapter design;
- transaction boundaries;
- error mapping;
- server timeout/shutdown config;
- request-ID/log redaction;
- test evidence;
- dependency additions;
- live-PG/provider/reverse-proxy limitations;
- exclusions.

## 67. Done-when

Architecture can verify:

1. exact parent preserved;
2. modular `gridworks-api` HTTP foundation exists;
3. loopback `:18080` default preserved;
4. `/api/v1` explicit versioning exists;
5. valid OpenAPI matches implemented routes;
6. strict JSON/body/error behavior exists;
7. request IDs + panic recovery + server timeouts/graceful shutdown exist;
8. guest issuance is durable/idempotent and raw token stored only as digest;
9. durable replay does not recover the raw session secret;
10. bearer session authentication derives account/player server-side;
11. current-session revocation works;
12. self/public profile boundaries are safe;
13. profile update is bounded/idempotent;
14. public company creation derives owner from auth player;
15. public group creation derives owner from auth player;
16. system ownership cannot be selected through public endpoints;
17. raw ownership transfer is not publicly exposed;
18. PostgreSQL is the normal durable production path;
19. no in-memory production fallback;
20. accepted domain normalization/idempotency/ownership invariants are reused, not forked;
21. no Rust simulation math is reimplemented in Go;
22. no WebSocket/NATS/deployment/later-domain scope creep;
23. CI/security gates pass;
24. complete GitHub handback exists.

## 68. Handback

Post to the WP-009 issue and draft PR:
- branch;
- exact parent/head;
- commit chain;
- changed files/dependencies;
- public route inventory;
- OpenAPI path/version;
- handler/application/repository package layout;
- server lifecycle/timeouts;
- request/error contract;
- authentication/session flow;
- raw-token/replay semantics;
- profile update semantics;
- company/group authorization derivation;
- PostgreSQL adapter + transaction design;
- any additive migration;
- SQLSTATE/constraint mapping strategy;
- tests and CI;
- live-PG/provider/reverse-proxy limitations;
- deviations/risks/conflicts;
- explicit no-host/deployment statement.

## 69. Stop condition

After handback, Engineering stops.

Do not start WP-010 without separate owner authorization.

Do not merge the WP-009 PR.

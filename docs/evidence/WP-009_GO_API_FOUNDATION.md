# WP-009 Go API Foundation Evidence

## Composition and runtime

`gridworks-api` is composed in `services/api/cmd/gridworks-api/main.go` with `database/sql` and the pinned pgx/v5 stdlib driver. Normal production startup requires `GRIDWORKS_DATABASE_URL`; absence or failed PostgreSQL readiness is fatal. No WP-007/WP-008 in-memory store is instantiated as an authoritative production fallback. Test fakes are confined to API unit tests.

The default bind is `127.0.0.1:18080`. The server has bounded read-header/read/write/idle timeouts, request panic recovery, graceful SIGINT/SIGTERM shutdown and no public-interface default.

## Contract and route inventory

`packages/schemas/openapi/gridworks-api-v1.json` is the repository-owned OpenAPI 3.1 contract. `tests/structural/check_wp009_openapi.py` validates the exact implemented route inventory:

- `/healthz`, `/readyz`, `/version`;
- `POST /api/v1/auth/guest`;
- `DELETE /api/v1/auth/session`;
- `GET /api/v1/me`;
- `PATCH /api/v1/me/profile`;
- `GET /api/v1/players/{player_id}`;
- `POST /api/v1/companies`;
- `GET /api/v1/companies/{company_id}`;
- `POST /api/v1/company-groups`.

No ownership-transfer, provider-link, simulation, marketplace or later-domain route is published.

## HTTP/security contract

The HTTP boundary generates/validates request IDs, applies JSON content-type and bounded-body checks, rejects unknown/trailing JSON fields, emits a stable error envelope, sets `nosniff`, uses `no-store` for auth/stateful responses and recovers panics without exposing internals.

Bearer authentication hashes the presented secret and resolves the account/player/session through PostgreSQL. Handlers receive only a server-derived principal. Raw bearer tokens, digests, provider subjects, SQL/DSNs and private auth data do not enter DTOs or errors.

## PostgreSQL adapter

`services/api/internal/postgres` implements the API repository interfaces using explicit parameterized SQL and transactions over accepted WP-006/WP-007/WP-008 tables. Guest issuance covers account/player/profile/session verifier/receipt atomically; durable replay returns safe references without a raw token. Profile updates use the identity receipt namespace. Company/group creation derives ownership from the authenticated player, inserts through accepted name/ownership/history triggers and records company-domain receipts.

No migration was added or modified. No live PostgreSQL provisioning or execution occurred.

## Scope and limitations

Provider linking is intentionally not public. Profile handle/avatar/auth changes, raw ownership transfer, group governance, realtime, workers, simulation, acquisition, marketplace, finance, admin and later-WP routes are not implemented. No host, database, role, service, port, reverse-proxy, provider or container runtime mutation occurred.

## REQUEST_CHANGES remediation evidence

R1–R12 remediation extends the API foundation with JSON-safe guest receipt references and explicit durable replay credential-unavailable semantics; transactional company/group receipt claim/replay using accepted WP-008 digest inputs; propagated CSPRNG errors; SQLSTATE/constraint-aware safe error mapping; centralized profile validation and JSONB preference encoding; restricted account-state authentication; replay-safe current-session revoke; nullable public DTO handling; shared WP-008 name normalization; bounded DB-pool/server hardening; OpenAPI schema/operation validation; and expanded HTTP/adapter regression evidence.

The API remains PostgreSQL-only in normal composition. No live PostgreSQL, provider, host, reverse-proxy, deployment, port or container runtime action occurred.

## R13–R21 remediation evidence

- Identity/profile and company/group first-use idempotency serializes through transaction-scoped PostgreSQL advisory locks derived from the domain namespace plus idempotency key before receipt lookup or mutation. Same-payload replay and contradictory-key conflict paths are exercised through the real repository methods with sqlmock.
- Profile notification preferences use an explicit `$8::jsonb` write contract while preserving nil/no-change semantics.
- Stable method-aware dispatch is generated from the same Go route registry used to register handlers and produce `RouteInventory()`; unknown paths remain JSON 404 and known paths with the wrong method return JSON 405 plus `Allow`, including parameterized routes.
- OpenAPI 3.1 is semantically validated by pinned `kin-openapi v0.123.0` and compared against the Go route inventory. Successful response schemas/examples, stable errors, bearer security and idempotency metadata remain repository-owned.
- HTTP regression coverage includes media type, malformed/trailing/oversized JSON, request-ID generation/reflection, panic recovery, CORS absence, no-store, authentication states, replay-safe revoke, profile policy/idempotency, company/group authorization/conflicts/visibility, and forbidden later/raw routes.
- PostgreSQL adapter regression coverage executes real guest/profile/company/group repository methods using sqlmock, verifies advisory-lock ordering, transaction replay/conflict, raw-token exclusion, server-derived player ownership, explicit JSONB SQL, nullable public reads, representative SQLSTATE mapping and parent-context cancellation. Tests explicitly do not claim to prove PostgreSQL trigger behavior.
- Repository entropy is injectable for tests while production uses `crypto/rand.Reader`. Failure coverage proves rollback/no partial commit at guest generation plus company/group entity, ownership and history-ID stages.
- API runtime construction is factored through `services/api/internal/apiruntime`: loopback remains the default, HTTP timeouts/header limits are non-zero, graceful shutdown is testable without a public port, and non-`ErrServerClosed` serve failures propagate.
- External request IDs are restricted to ASCII alphanumeric plus `._:-`, non-empty and <=96 characters. Unsafe printable/control values are replaced before reflection/logging; Authorization, bodies and DSNs are not access-log fields.

The API remains PostgreSQL-only in normal composition. No live PostgreSQL, provider, host, reverse-proxy, deployment, port or container runtime action occurred.

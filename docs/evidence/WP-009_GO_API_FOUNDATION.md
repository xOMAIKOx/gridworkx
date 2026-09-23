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

## R13–R20 remediation evidence

- First-use idempotency now serializes through transaction-scoped PostgreSQL advisory locks per domain namespace/key before receipt lookup or mutation.
- Profile preferences use an explicit `$8::jsonb` write contract; create/replay display results use persisted trimmed display values; repository SQLSTATE mapping and profile domain errors remain stable API errors.
- Guest receipts are strictly decoded with unknown-field rejection and replay joins account/player/session/profile consistency.
- Explicit method-aware dispatch provides stable 404/405 behavior and `Allow`; route, timeout, header, request, panic and SQL-mock adapter tests cover the hardening boundary.
- OpenAPI is semantically validated by pinned kin-openapi and compared against the Go route inventory; concrete response schemas/examples, error responses and mutation metadata are present.
- PostgreSQL adapter tests execute actual repository methods through sqlmock, including receipt locking/query order and entropy-failure rollback; live PostgreSQL trigger behavior remains an environment gate.
- Repository entropy is injectable; production uses `crypto/rand.Reader` and failure tests prove no mutation SQL follows entropy failure.

## R21–R22 remediation evidence

- Domain errors now map explicitly to safe 4xx API envelopes: company name validation/conflict, identity validation/conflict and accepted state/principal failures no longer rely on generic 500 handling.
- Guest receipt decoding is strict about unknown fields and trailing JSON; replay verifies relational account/player/session/profile consistency and uses the same nullable-normalized profile shape as first issuance.

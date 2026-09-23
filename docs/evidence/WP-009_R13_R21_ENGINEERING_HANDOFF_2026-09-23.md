# WP-009 R13–R21 Engineering Handback

**Repository:** `xOMAIKOx/gridworkx`  
**Branch:** `engineering/wp-009-go-api-foundation`  
**Existing PR:** #20  
**Exact required parent:** `f93c6526f9ce8743406606265456b283c8d8787c`  
**Authorized source:** issue #19 Architecture comment `5792358224`  
**Scope:** R13–R21 only

The exact final implementation SHA and CI URLs are recorded in the GitHub issue/PR handback posted after the final branch gates complete.

## Implemented corrections

### R13 — company/group durable idempotency

Company/group mutations use transaction-scoped PostgreSQL advisory locks derived from the company-domain namespace plus idempotency key before receipt lookup or durable entity mutation. The single WP-008 company receipt namespace remains authoritative. Repository tests execute real create/replay/conflict paths and assert lock/query/mutation ordering.

### R14 — profile-update idempotency

Profile updates use the identity-domain advisory-lock namespace before receipt lookup or profile mutation. Same-key/same-payload replay and contradictory reuse are exercised against the real repository method with transactional sqlmock expectations.

### R15 — JSONB preferences

Notification preferences persist through an explicit `$8::jsonb` SQL cast while nil retains no-change semantics. Repository tests assert the JSONB SQL shape.

### R16 — stable 404/405

HTTP route registration and method handling now derive from one Go route registry. Known paths with the wrong method return the stable JSON error envelope, HTTP 405 and `Allow`; unknown paths return JSON 404. Parameterized player/company paths use exact segment matching so extra path segments remain 404.

### R17 — OpenAPI semantics and parity

Pinned `github.com/getkin/kin-openapi v0.123.0` performs OpenAPI 3.1 semantic validation. The repository contract includes concrete success schemas/examples and stable error/security/idempotency metadata. The parity test compares OpenAPI paths/methods to the same Go route inventory used for handler registration.

### R18 — acceptance test matrix

Added executable HTTP evidence for:
- wrong-method 405 + `Allow` and unknown-route 404;
- invalid media type, malformed JSON, trailing JSON and oversized body;
- request-ID generation/preservation and panic-safe 500;
- no wildcard CORS and no-store auth responses;
- valid/invalid/expired/revoked/restricted account behavior;
- replay-safe revoke followed by failed normal auth;
- self/public/private profile behavior and public DTO leakage checks;
- valid profile update, server-owned time, same-key replay, contradictory conflict, unknown fields and validation failures;
- operating/holding company creation, system rejection, authenticated owner derivation, alternate-owner rejection, conflict mapping, company visibility, group creation and absence of raw ownership/provider-link routes.

Added executable PostgreSQL-adapter evidence for:
- begin/commit/rollback and advisory-lock ordering;
- guest safe receipt persistence/replay and raw-token exclusion;
- profile JSONB mutation/replay/conflict;
- company/group server-derived owner creation/replay/conflict;
- nullable public reads and representative SQLSTATE mapping;
- parent-context cancellation.

Mocks do not claim to prove PostgreSQL trigger semantics; live PG18 remains a separate environment gate.

### R19 — CSPRNG rollback proof

Repository entropy is injectable for tests and defaults to `crypto/rand.Reader` in production. Failure tests cover guest generation plus company/group entity, ownership and history-event generation stages and require rollback/no committed partial mutation.

### R20 — testable runtime seam

`services/api/internal/apiruntime` owns production HTTP server configuration and run/shutdown behavior. Tests prove:
- default bind `127.0.0.1:18080`;
- non-zero read/read-header/write/idle timeouts and header limit;
- graceful shutdown without opening a public port;
- non-`http.ErrServerClosed` serve failures propagate;
- shutdown failures propagate.

The main binary uses this seam with signal cancellation.

### R21 — deterministic safe request-ID logging

External request IDs are accepted only when non-empty, <=96 characters and composed of ASCII letters/digits plus `.`, `_`, `:`, `-`. Unsafe printable/control values are replaced by a server-generated ID before reflection/logging. Access logs do not contain Authorization headers, request bodies or DSNs.

## Scope compliance

No WP-010 or later-WP route/service was added. No live PostgreSQL contact/provisioning, provider configuration, host mutation, deployment, port activation, reverse-proxy/systemd change or Docker/Podman/Compose/Kubernetes/OCI runtime work was performed.

Engineering stops after the final GitHub handback for Architecture review.

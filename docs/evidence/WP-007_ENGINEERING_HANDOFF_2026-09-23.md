# WP-007 Engineering Handback — Guest-First Identity and Player Profile Foundation

**Repository:** `xOMAIKOx/gridworkx`
**Branch:** `engineering/wp-007-guest-identity-profile-rerun`
**Exact required parent:** `8306ac89cc25c8d7e193e1b93f51d3fe9c12ecf9`
**Exact HEAD:** `bd12e6fe7c381de3ad42340bef1ab676929672ae`
**Draft PR:** pending creation from this branch
**Work order:** `docs/work-orders/WP-007_GUEST_FIRST_IDENTITY_AND_PLAYER_PROFILE.md`

## Scope delivered

Implemented only the authorized WP-007 guest-first identity and player-profile foundation:

- immutable opaque `account_id` and `player_id` contracts;
- atomic server-issued guest account/player/profile/session issuance;
- verifier-only session-token persistence;
- verified external identity-linking seam with proof-of-control metadata;
- existing-account guest-to-linked preservation;
- no normal self-service merge;
- deterministic NFKC/case-fold/confusable-subset handles;
- reserved-handle policy configuration;
- public/private profile schema separation;
- locale, timezone, visibility, DM and notification-preference foundations;
- PostgreSQL migration 0003 with integrity, uniqueness and idempotency boundaries.

## Changed-file inventory

The parent-to-HEAD compare contains 15 files:

```text
Makefile
apps/tools/src/validate-content.mjs
db/migrations/0003_wp007_identity_profile.sql
docs/evidence/WP-007_GUEST_FIRST_IDENTITY_PROFILE.md
packages/content/config/reserved-handles.json
packages/schemas/catalog.json
packages/schemas/identity-command.schema.json
packages/schemas/identity-public-profile.schema.json
packages/schemas/identity-state.schema.json
services/go.mod
services/go.sum
services/internal/identity/identity.go
services/internal/identity/identity_test.go
tests/structural/check_repository.py
tests/structural/check_wp007_migration.py
```

## Identity / security contract

- Account identity, player identity and future company/control identity remain separate.
- Guest issuance creates one account, player, profile and session under one idempotency key.
- CSPRNG failures propagate before state mutation; failure injection proves no partial state/panic.
- Raw guest/session tokens are returned only at issuance and never enter receipt/durable state.
- Session records store verifier/digest metadata and support expiry/revocation.
- External linking requires a verified assertion and attaches to the existing account.
- Issuer/subject provider identity cannot become the GRIDWORKS account ID.
- Same external identity cannot link to two accounts; independent-account self-service merge is rejected.
- Guest issuance and external linking share one global idempotency receipt namespace with deterministic request digests and contradictory-reuse rejection.
- Proof references are retained for linked identities.
- Profile hierarchy uses composite player/account relational integrity.
- Public profile schema is strict and excludes account/session/provider/security fields.
- Handle algorithm is `gridworks-unicode-15.1-confusable-subset-v1`: NFKC + Unicode case folding plus the documented Latin/Cyrillic subset, not full UTS #39.

## Verification

Commands executed in a fresh worktree at the exact parent-derived branch:

```bash
make validate
make test
make format-check
make security
make godot-check
go test ./services/...
go vet ./services/...
```

Results:

- repository structure/content/schema/WP-006/WP-007/systemd checks: PASS;
- Rust tests: 47 passed;
- Rust replay integration: 2 passed;
- Go identity/runtime tests: PASS;
- admin typecheck/build/test: PASS (2 tests);
- format checks: PASS;
- gitleaks: no leaks found;
- forbidden-runtime scan: PASS;
- ShellCheck: PASS;
- Godot headless validation: PASS.

A fresh `npm ci --ignore-scripts` was used for the existing committed dependency lockfile. npm reported four audit advisories in the existing dependency graph; no dependency manifest or lockfile was changed by WP-007.

## CI / PR

The draft PR is opened against the repository’s normal review flow from the exact-parent-derived branch. CI evidence will be recorded in the PR after GitHub completes the run.

## Boundaries / deviations

- No live PostgreSQL contact or mutation.
- No PostgreSQL provisioning, roles, databases or ports.
- No provider calls, Apple/Google/OIDC client creation or secrets.
- No host, systemd, nginx, DNS, deployment or runtime mutation.
- No Docker/Podman/Compose/Kubernetes/OCI runtime.
- No WP-008, WP-009, messaging, notifications delivery or public HTTP API.
- No Parley collaboration.

Engineering stops for Architecture review of exact HEAD `bd12e6fe7c381de3ad42340bef1ab676929672ae`.

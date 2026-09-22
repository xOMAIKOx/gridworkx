# WP-007 Guest-First Identity and Player Profile Evidence

## Identity hierarchy and guest issuance

The transport-independent Go domain in `services/internal/identity` preserves `Account → PlayerProfile → Player`. Account and player IDs are opaque server-generated random values and never derive from handles, provider subjects, email or device IDs.

`InMemoryStore.IssueGuest` is the repository-safe transaction seam: one idempotency key creates one account, player, profile and opaque guest session result. The session record stores only a SHA-256 token digest; the raw token exists only in the issuance result for secure handoff. Replaying the same request returns the original result and does not create duplicate identities.

The repository interface is the boundary for a later PostgreSQL transaction implementation. No live database or provider was used in WP-007.

## Sessions and linking

Session lookup uses constant-time digest comparison and rejects revoked/expired sessions. External identity linking accepts only a verified `ExternalIdentityAssertion`; provider issuer/subject remains a link key and never becomes the GRIDWORKS account ID. A verified link attaches to the existing guest account and transitions it to protected status while preserving account/player/profile identity. Existing links replay idempotently; linking the same external identity to another account rejects rather than merging accounts.

No Apple/Google/OIDC validation, redirects, provider credentials or live integration was implemented.

## Handles and public/private boundaries

Handles use NFKC normalization, Unicode case folding and the deterministic pinned GRIDWORKS strategy label `gridworks-unicode-15.1-confusable-subset-v1`, with the supported explicit Latin/Cyrillic confusable subset, invisible/control rejection, reserved skeletons and canonical/skeleton uniqueness. This is a versioned GRIDWORKS subset and does not claim full UTS #39 coverage. Display names are mutable, non-unique and separately validated.

`PublicPlayerProfile` intentionally excludes account status, session verifier, provider identity, recovery and security metadata. Locale, IANA timezone, visibility, DM policy, discoverability and notification preference foundations are represented without implementing messaging or delivery.

## PostgreSQL migration

`db/migrations/0003_wp007_identity_profile.sql` defines accounts, players, profiles, guest sessions, external identity links, idempotency receipts and reserved handles. It includes primary/foreign keys, status checks, token-digest uniqueness, issuer/subject uniqueness, canonical/skeleton handle uniqueness and no raw-token column. The migration is transaction-wrapped and contains no credentials, roles, provider configuration or provisioning operations.

## Verification and scope

Go domain tests cover atomic/idempotent guest issuance, raw-secret absence, session revocation, verified-link proof, guest identity preservation, cross-account conflict, normalization/case/confusable/reserved/invisible handle cases, locale/timezone/display-name validation and public/private separation. Structural validation checks the migration contract and strict identity schemas.

No WP-008 company ownership, WP-009 HTTP/API, live provider integration, messaging, notification delivery, deployment, host/database provisioning or container runtime work was performed.

## REQUEST_CHANGES remediation evidence

- Link mutations now require an explicit idempotency key and deterministic request digest. Same key plus same payload replays the link; contradictory reuse rejects with `ErrIdempotencyConflict`.
- Guest idempotency receipts store only account/player/session references and request digest. A duplicate issuance never recovers or persists the raw session token; raw token material exists only in the one-time result.
- Secure ID generation uses an injectable `io.Reader`; every randomness error propagates before state/handle mutation. Failure-injection tests prove no partial identity state.
- External link proof references are retained in the domain record and require non-empty verified proof metadata.
- Player/profile SQL uses a composite `(player_id, account_id)` foreign key, preventing cross-account profile association.
- Added a strict public-profile schema and fixture; public fields cannot carry account, session, provider or security fields.
- Canonical handle comparison uses pinned `golang.org/x/text v0.21.0` NFKC plus Unicode case-folding and an explicitly versioned GRIDWORKS Latin/Cyrillic confusable subset. The reserved baseline is represented in `packages/content/config/reserved-handles.json`; broader UTS #39 coverage is a future versioned algorithm change, not implied by this subset.

## Unified identity mutation receipt namespace

Guest issuance and external-link mutations now share one logical idempotency receipt map keyed by the same mutation key namespace. Each receipt retains mutation type, request digest and safe result references only. Same key plus same mutation/payload replays deterministically; same key plus a different mutation type or payload returns `ErrIdempotencyConflict`. The cross-mutation regression test proves a guest-issuance key cannot be reused for external linking and that no link state is mutated. Raw session tokens remain excluded from receipt state.

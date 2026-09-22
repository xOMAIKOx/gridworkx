# WP-007 Guest-First Identity and Player Profile Evidence

## Identity hierarchy and guest issuance

The transport-independent Go domain in `services/internal/identity` preserves `Account → PlayerProfile → Player`. Account and player IDs are opaque server-generated random values and never derive from handles, provider subjects, email or device IDs.

`InMemoryStore.IssueGuest` is the repository-safe transaction seam: one idempotency key creates one account, player, profile and opaque guest session result. The session record stores only a SHA-256 token digest; the raw token exists only in the issuance result for secure handoff. Replaying the same request returns the original result and does not create duplicate identities.

The repository interface is the boundary for a later PostgreSQL transaction implementation. No live database or provider was used in WP-007.

## Sessions and linking

Session lookup uses constant-time digest comparison and rejects revoked/expired sessions. External identity linking accepts only a verified `ExternalIdentityAssertion`; provider issuer/subject remains a link key and never becomes the GRIDWORKS account ID. A verified link attaches to the existing guest account and transitions it to protected status while preserving account/player/profile identity. Existing links replay idempotently; linking the same external identity to another account rejects rather than merging accounts.

No Apple/Google/OIDC validation, redirects, provider credentials or live integration was implemented.

## Handles and public/private boundaries

Handles use NFKC normalization, lower-casing and a deterministic pinned strategy label `unicode-15.1-skeleton-v1`, with explicit Latin/Cyrillic confusable mapping, invisible/control rejection, reserved skeletons and canonical/skeleton uniqueness. Display names are mutable, non-unique and separately validated.

`PublicPlayerProfile` intentionally excludes account status, session verifier, provider identity, recovery and security metadata. Locale, IANA timezone, visibility, DM policy, discoverability and notification preference foundations are represented without implementing messaging or delivery.

## PostgreSQL migration

`db/migrations/0003_wp007_identity_profile.sql` defines accounts, players, profiles, guest sessions, external identity links, idempotency receipts and reserved handles. It includes primary/foreign keys, status checks, token-digest uniqueness, issuer/subject uniqueness, canonical/skeleton handle uniqueness and no raw-token column. The migration is transaction-wrapped and contains no credentials, roles, provider configuration or provisioning operations.

## Verification and scope

Go domain tests cover atomic/idempotent guest issuance, raw-secret absence, session revocation, verified-link proof, guest identity preservation, cross-account conflict, normalization/case/confusable/reserved/invisible handle cases, locale/timezone/display-name validation and public/private separation. Structural validation checks the migration contract and strict identity schemas.

No WP-008 company ownership, WP-009 HTTP/API, live provider integration, messaging, notification delivery, deployment, host/database provisioning or container runtime work was performed.

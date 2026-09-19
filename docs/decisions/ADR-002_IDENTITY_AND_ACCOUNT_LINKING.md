# ADR-002 — Guest-First Identity and Account Linking

**Status:** ACCEPTED  
**Date:** 2026-09-19  
**Applies to:** WP-007 and all account/security work

## Context

GRIDWORKS should start like a game, not like enterprise SaaS. A player must be able to enter the opening scenario immediately, while still allowing durable account protection and cross-device recovery.

## Decision

GRIDWORKS uses a **server-issued guest account first**, followed by explicit account protection/linking.

On first launch:

1. client requests a guest identity;
2. server creates immutable `account_id` and `player_id`;
3. player begins gameplay;
4. later UX prompts the player to protect/link the account.

Supported linking targets should include:

- Sign in with Apple;
- Google;
- email/OIDC;
- estate identity integration where appropriate.

The linked identity attaches to the existing account. It does **not** create a new player or migrate economic state between separate records.

## Invariants

Guest-to-linked conversion preserves:

- player ID;
- company ownership;
- facilities;
- inventory;
- manager history;
- reputation;
- purchases/entitlements;
- messages;
- consortium membership;
- challenge history.

Account linking requires proof of control of both the guest session and external identity.

Account merge between two independently progressed GRIDWORKS accounts is **not** a normal self-service path and requires explicit support/recovery policy.

## Authentication architecture

The mobile client receives short-lived application access credentials. Long-lived refresh/session material must be stored using platform secure storage.

The game API trusts GRIDWORKS application identity/session claims, not raw third-party provider tokens for every domain request.

## Names

Authentication identity is separate from public identity.

- `account_id` — private authentication identity;
- `player_id` — immutable game identity;
- `unique_handle` — public unique identifier;
- `display_name` — mutable/non-unique public name.

## Rationale

This preserves low-friction onboarding while avoiding irreversible guest-device lock-in.

## Consequences

- WP-007 includes guest issuance, linking, recovery scaffolding and uniqueness rules.
- Company ownership references `player_id/account-controlled principals`, never a provider-specific subject directly.
- Account protection reminders must not block core gameplay unfairly.

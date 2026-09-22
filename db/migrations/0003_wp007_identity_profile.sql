BEGIN;

SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '30s';

CREATE TABLE IF NOT EXISTS gridworks.accounts (
    account_id text PRIMARY KEY,
    status text NOT NULL CHECK (status IN ('guest', 'protected', 'suspended', 'recovery_restricted')),
    authority_version bigint NOT NULL CHECK (authority_version > 0),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gridworks.players (
    player_id text PRIMARY KEY,
    account_id text NOT NULL UNIQUE REFERENCES gridworks.accounts(account_id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gridworks.player_profiles (
    player_id text PRIMARY KEY REFERENCES gridworks.players(player_id) ON DELETE RESTRICT,
    account_id text NOT NULL REFERENCES gridworks.accounts(account_id) ON DELETE RESTRICT,
    handle_display text NOT NULL,
    handle_canonical text NOT NULL,
    handle_skeleton text NOT NULL,
    display_name text NOT NULL,
    avatar_asset_id text,
    bio text,
    locale text NOT NULL,
    timezone text NOT NULL,
    visibility text NOT NULL CHECK (visibility IN ('public', 'private')),
    dm_policy text NOT NULL CHECK (dm_policy IN ('everyone', 'nobody')),
    discoverable boolean NOT NULL DEFAULT true,
    notification_preferences jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (length(handle_display) BETWEEN 3 AND 32),
    CHECK (length(display_name) BETWEEN 1 AND 80),
    CHECK (length(locale) BETWEEN 2 AND 16),
    CHECK (length(timezone) BETWEEN 1 AND 128)
);
CREATE UNIQUE INDEX IF NOT EXISTS player_profiles_handle_canonical_uq ON gridworks.player_profiles(handle_canonical);
CREATE UNIQUE INDEX IF NOT EXISTS player_profiles_handle_skeleton_uq ON gridworks.player_profiles(handle_skeleton);
CREATE INDEX IF NOT EXISTS player_profiles_account_idx ON gridworks.player_profiles(account_id);

CREATE TABLE IF NOT EXISTS gridworks.guest_sessions (
    session_id text PRIMARY KEY,
    account_id text NOT NULL REFERENCES gridworks.accounts(account_id) ON DELETE RESTRICT,
    token_digest char(64) NOT NULL UNIQUE CHECK (token_digest ~ '^[0-9a-f]{64}$'),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_used_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    CHECK (expires_at > created_at)
);
CREATE INDEX IF NOT EXISTS guest_sessions_account_idx ON gridworks.guest_sessions(account_id);
CREATE INDEX IF NOT EXISTS guest_sessions_active_idx ON gridworks.guest_sessions(expires_at, revoked_at);

CREATE TABLE IF NOT EXISTS gridworks.account_identity_links (
    link_id text PRIMARY KEY,
    account_id text NOT NULL REFERENCES gridworks.accounts(account_id) ON DELETE RESTRICT,
    provider text NOT NULL,
    issuer text NOT NULL,
    subject text NOT NULL,
    identity_key text NOT NULL UNIQUE,
    linked_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    proof_reference text,
    revoked_at timestamptz,
    UNIQUE (issuer, subject)
);
CREATE INDEX IF NOT EXISTS account_identity_links_account_idx ON gridworks.account_identity_links(account_id);

CREATE TABLE IF NOT EXISTS gridworks.identity_mutation_receipts (
    idempotency_key text PRIMARY KEY,
    mutation_type text NOT NULL,
    request_digest char(64) NOT NULL CHECK (request_digest ~ '^[0-9a-f]{64}$'),
    result_ref text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gridworks.reserved_handles (
    handle_skeleton text PRIMARY KEY,
    reason_code text NOT NULL,
    active boolean NOT NULL DEFAULT true
);

COMMIT;

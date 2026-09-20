BEGIN;

SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '30s';

CREATE SCHEMA IF NOT EXISTS gridworks;

CREATE TABLE IF NOT EXISTS gridworks.schema_metadata (
    schema_id text PRIMARY KEY,
    schema_version text NOT NULL,
    applied_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    source text NOT NULL
);

CREATE TABLE IF NOT EXISTS gridworks.domain_contracts (
    domain_id text PRIMARY KEY,
    owner_boundary text NOT NULL,
    contract_state text NOT NULL CHECK (contract_state = 'foundation'),
    schema_version text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gridworks.content_manifests (
    manifest_id text PRIMARY KEY,
    schema_version text NOT NULL,
    rules_version text NOT NULL,
    content_version text NOT NULL,
    signature_required boolean NOT NULL DEFAULT true,
    checksum text,
    published_at timestamptz
);

CREATE TABLE IF NOT EXISTS gridworks.system_principals (
    principal_id text PRIMARY KEY,
    display_label text NOT NULL,
    principal_kind text NOT NULL,
    political_governance_enabled boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gridworks.domain_reservations (
    reservation_id text PRIMARY KEY,
    domain_id text NOT NULL REFERENCES gridworks.domain_contracts(domain_id),
    reservation_kind text NOT NULL,
    schema_version text NOT NULL,
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb
);

INSERT INTO gridworks.schema_metadata (schema_id, schema_version, source)
VALUES ('schema.wp001.foundation', 'schema-0.1.0', 'docs/work-orders/WP-001_REPOSITORY_AND_NATIVE_RUNTIME_BASELINE.md')
ON CONFLICT (schema_id) DO NOTHING;

COMMIT;

BEGIN;

SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '30s';

CREATE TABLE IF NOT EXISTS gridworks.companies (
    company_id text PRIMARY KEY,
    company_type text NOT NULL CHECK (company_type IN ('holding', 'operating', 'system')),
    name_display text NOT NULL,
    name_canonical text NOT NULL,
    name_skeleton text NOT NULL,
    status text NOT NULL CHECK (status IN ('active', 'suspended', 'archived')),
    visibility text NOT NULL CHECK (visibility IN ('public', 'private')),
    description text,
    logo_asset_id text,
    industry_id text,
    headquarters_region_id text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS companies_name_canonical_uq ON gridworks.companies(name_canonical);
CREATE UNIQUE INDEX IF NOT EXISTS companies_name_skeleton_uq ON gridworks.companies(name_skeleton);

CREATE TABLE IF NOT EXISTS gridworks.company_groups (
    group_id text PRIMARY KEY,
    name_display text NOT NULL,
    name_canonical text NOT NULL,
    name_skeleton text NOT NULL,
    status text NOT NULL CHECK (status IN ('active', 'suspended', 'archived')),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS company_groups_name_canonical_uq ON gridworks.company_groups(name_canonical);
CREATE UNIQUE INDEX IF NOT EXISTS company_groups_name_skeleton_uq ON gridworks.company_groups(name_skeleton);

CREATE TABLE IF NOT EXISTS gridworks.company_ownership (
    ownership_id text PRIMARY KEY,
    entity_id text NOT NULL,
    entity_type text NOT NULL CHECK (entity_type IN ('company', 'group')),
    owner_type text NOT NULL CHECK (owner_type IN ('player', 'system')),
    owner_id text NOT NULL,
    share_bps integer NOT NULL CHECK (share_bps > 0 AND share_bps <= 10000),
    active boolean NOT NULL DEFAULT true,
    effective_from timestamptz NOT NULL,
    effective_to timestamptz,
    UNIQUE (entity_id, entity_type, owner_type, owner_id) DEFERRABLE INITIALLY IMMEDIATE
);
CREATE UNIQUE INDEX IF NOT EXISTS company_ownership_active_owner_uq ON gridworks.company_ownership(entity_id, entity_type, owner_type, owner_id) WHERE active;
CREATE INDEX IF NOT EXISTS company_ownership_entity_idx ON gridworks.company_ownership(entity_id, entity_type, active);

CREATE TABLE IF NOT EXISTS gridworks.company_group_membership (
    membership_id text PRIMARY KEY,
    company_id text NOT NULL REFERENCES gridworks.companies(company_id) ON DELETE RESTRICT,
    group_id text NOT NULL REFERENCES gridworks.company_groups(group_id) ON DELETE RESTRICT,
    active boolean NOT NULL DEFAULT true,
    effective_from timestamptz NOT NULL,
    effective_to timestamptz,
    source_ref text NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS company_group_membership_active_company_uq ON gridworks.company_group_membership(company_id) WHERE active;
CREATE UNIQUE INDEX IF NOT EXISTS company_group_membership_active_pair_uq ON gridworks.company_group_membership(company_id, group_id) WHERE active;

CREATE TABLE IF NOT EXISTS gridworks.ownership_history (
    event_id text PRIMARY KEY,
    entity_id text NOT NULL,
    mutation_type text NOT NULL,
    from_owner_type text,
    from_owner_id text,
    to_owner_type text,
    to_owner_id text,
    share_bps integer NOT NULL CHECK (share_bps > 0 AND share_bps <= 10000),
    effective_time timestamptz NOT NULL,
    source_ref text NOT NULL,
    actor_player_id text,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gridworks.company_mutation_receipts (
    idempotency_key text PRIMARY KEY,
    mutation_type text NOT NULL,
    request_digest char(64) NOT NULL CHECK (request_digest ~ '^[0-9a-f]{64}$'),
    result_ref text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gridworks.reserved_company_names (
    name_skeleton text PRIMARY KEY,
    reason_code text NOT NULL,
    active boolean NOT NULL DEFAULT true
);

CREATE OR REPLACE FUNCTION gridworks.validate_active_ownership_total()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE target_entity text; target_type text; total_bps bigint;
BEGIN
    target_entity := COALESCE(NEW.entity_id, OLD.entity_id);
    target_type := COALESCE(NEW.entity_type, OLD.entity_type);
    SELECT COALESCE(sum(share_bps), 0) INTO total_bps FROM gridworks.company_ownership WHERE entity_id = target_entity AND entity_type = target_type AND active;
    IF total_bps <> 10000 THEN RAISE EXCEPTION 'active ownership must total exactly 10000 basis points'; END IF;
    IF TG_OP = 'DELETE' THEN RETURN OLD; ELSE RETURN NEW; END IF;
END
$$;

DROP TRIGGER IF EXISTS company_ownership_total_guard ON gridworks.company_ownership;
CREATE CONSTRAINT TRIGGER company_ownership_total_guard
    AFTER INSERT OR UPDATE OR DELETE ON gridworks.company_ownership
    DEFERRABLE INITIALLY DEFERRED
    FOR EACH ROW EXECUTE FUNCTION gridworks.validate_active_ownership_total();

COMMIT;

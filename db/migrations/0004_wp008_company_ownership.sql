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
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (length(name_display) BETWEEN 3 AND 80)
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
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK (length(name_display) BETWEEN 3 AND 80)
);
CREATE UNIQUE INDEX IF NOT EXISTS company_groups_name_canonical_uq ON gridworks.company_groups(name_canonical);
CREATE UNIQUE INDEX IF NOT EXISTS company_groups_name_skeleton_uq ON gridworks.company_groups(name_skeleton);

CREATE TABLE IF NOT EXISTS gridworks.company_name_claims (
    name_canonical text PRIMARY KEY,
    name_skeleton text NOT NULL UNIQUE,
    entity_id text NOT NULL,
    entity_type text NOT NULL CHECK (entity_type IN ('company', 'group')),
    claimed_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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
    UNIQUE (ownership_id)
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

INSERT INTO gridworks.reserved_company_names (name_skeleton, reason_code)
VALUES ('gridworks', 'system'), ('official', 'system'), ('system', 'system'), ('admin', 'security'), ('administrator', 'security'), ('support', 'security'), ('moderator', 'security')
ON CONFLICT (name_skeleton) DO NOTHING;

CREATE OR REPLACE FUNCTION gridworks.claim_company_name()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF EXISTS (SELECT 1 FROM gridworks.reserved_company_names WHERE name_skeleton = NEW.name_skeleton AND active) THEN
        RAISE EXCEPTION 'company/group name is reserved';
    END IF;
    INSERT INTO gridworks.company_name_claims (name_canonical, name_skeleton, entity_id, entity_type)
    VALUES (NEW.name_canonical, NEW.name_skeleton, NEW.company_id, 'company')
    ON CONFLICT DO NOTHING;
    IF NOT EXISTS (SELECT 1 FROM gridworks.company_name_claims WHERE name_canonical = NEW.name_canonical AND entity_id = NEW.company_id) THEN
        RAISE EXCEPTION 'company/group name is unavailable';
    END IF;
    RETURN NEW;
END
$$;

DROP TRIGGER IF EXISTS company_name_claim_guard ON gridworks.companies;
CREATE TRIGGER company_name_claim_guard BEFORE INSERT ON gridworks.companies
FOR EACH ROW EXECUTE FUNCTION gridworks.claim_company_name();

CREATE OR REPLACE FUNCTION gridworks.claim_group_name()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF EXISTS (SELECT 1 FROM gridworks.reserved_company_names WHERE name_skeleton = NEW.name_skeleton AND active) THEN
        RAISE EXCEPTION 'company/group name is reserved';
    END IF;
    INSERT INTO gridworks.company_name_claims (name_canonical, name_skeleton, entity_id, entity_type)
    VALUES (NEW.name_canonical, NEW.name_skeleton, NEW.group_id, 'group')
    ON CONFLICT DO NOTHING;
    IF NOT EXISTS (SELECT 1 FROM gridworks.company_name_claims WHERE name_canonical = NEW.name_canonical AND entity_id = NEW.group_id) THEN
        RAISE EXCEPTION 'company/group name is unavailable';
    END IF;
    RETURN NEW;
END
$$;

DROP TRIGGER IF EXISTS company_group_name_claim_guard ON gridworks.company_groups;
CREATE TRIGGER company_group_name_claim_guard BEFORE INSERT ON gridworks.company_groups
FOR EACH ROW EXECUTE FUNCTION gridworks.claim_group_name();

CREATE OR REPLACE FUNCTION gridworks.prevent_ownership_history_mutation()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'ownership history is append-only';
END
$$;

DROP TRIGGER IF EXISTS ownership_history_append_only ON gridworks.ownership_history;
CREATE TRIGGER ownership_history_append_only
    BEFORE UPDATE OR DELETE ON gridworks.ownership_history
    FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_ownership_history_mutation();

CREATE OR REPLACE FUNCTION gridworks.freeze_group_membership_identity()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.company_id IS DISTINCT FROM OLD.company_id
       OR NEW.group_id IS DISTINCT FROM OLD.group_id
       OR NEW.source_ref IS DISTINCT FROM OLD.source_ref THEN
        RAISE EXCEPTION 'group membership identity is immutable';
    END IF;
    RETURN NEW;
END
$$;

DROP TRIGGER IF EXISTS company_group_membership_identity_guard ON gridworks.company_group_membership;
CREATE TRIGGER company_group_membership_identity_guard
    BEFORE UPDATE ON gridworks.company_group_membership
    FOR EACH ROW EXECUTE FUNCTION gridworks.freeze_group_membership_identity();

CREATE OR REPLACE FUNCTION gridworks.prevent_company_name_rewrite()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.name_canonical IS DISTINCT FROM OLD.name_canonical OR NEW.name_skeleton IS DISTINCT FROM OLD.name_skeleton THEN
        RAISE EXCEPTION 'company/group name identity is immutable';
    END IF;
    RETURN NEW;
END
$$;
DROP TRIGGER IF EXISTS company_name_identity_guard ON gridworks.companies;
CREATE TRIGGER company_name_identity_guard BEFORE UPDATE ON gridworks.companies FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_company_name_rewrite();
DROP TRIGGER IF EXISTS company_group_name_identity_guard ON gridworks.company_groups;
CREATE TRIGGER company_group_name_identity_guard BEFORE UPDATE ON gridworks.company_groups FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_company_name_rewrite();

CREATE OR REPLACE FUNCTION gridworks.prevent_company_name_claim_mutation()
RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'company name claims are append-only'; END $$;
DROP TRIGGER IF EXISTS company_name_claims_append_only ON gridworks.company_name_claims;
CREATE TRIGGER company_name_claims_append_only BEFORE UPDATE OR DELETE ON gridworks.company_name_claims FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_company_name_claim_mutation();

CREATE OR REPLACE FUNCTION gridworks.lock_ownership_entity_before_mutation()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE target_entity text; target_type text;
BEGIN
    target_entity := COALESCE(NEW.entity_id, OLD.entity_id); target_type := COALESCE(NEW.entity_type, OLD.entity_type);
    IF target_type = 'company' THEN PERFORM 1 FROM gridworks.companies WHERE company_id = target_entity FOR UPDATE;
    ELSE PERFORM 1 FROM gridworks.company_groups WHERE group_id = target_entity FOR UPDATE; END IF;
    IF TG_OP = 'DELETE' THEN RETURN OLD; ELSE RETURN NEW; END IF;
END
$$;
DROP TRIGGER IF EXISTS company_ownership_entity_lock ON gridworks.company_ownership;
CREATE TRIGGER company_ownership_entity_lock BEFORE INSERT OR UPDATE OR DELETE ON gridworks.company_ownership FOR EACH ROW EXECUTE FUNCTION gridworks.lock_ownership_entity_before_mutation();

CREATE OR REPLACE FUNCTION gridworks.validate_company_ownership_refs()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.entity_type = 'company' AND NOT EXISTS (SELECT 1 FROM gridworks.companies WHERE company_id = NEW.entity_id) THEN
        RAISE EXCEPTION 'ownership entity company does not exist';
    END IF;
    IF NEW.entity_type = 'group' AND NOT EXISTS (SELECT 1 FROM gridworks.company_groups WHERE group_id = NEW.entity_id) THEN
        RAISE EXCEPTION 'ownership entity group does not exist';
    END IF;
    IF NEW.owner_type = 'player' AND NOT EXISTS (SELECT 1 FROM gridworks.players WHERE player_id = NEW.owner_id) THEN
        RAISE EXCEPTION 'ownership player principal does not exist';
    END IF;
    IF NEW.owner_type = 'system' AND NEW.owner_id <> 'principal.gridworks.system' THEN
        RAISE EXCEPTION 'unknown system ownership principal';
    END IF;
    RETURN NEW;
END
$$;

DROP TRIGGER IF EXISTS company_ownership_reference_guard ON gridworks.company_ownership;
CREATE TRIGGER company_ownership_reference_guard
    BEFORE INSERT OR UPDATE ON gridworks.company_ownership
    FOR EACH ROW EXECUTE FUNCTION gridworks.validate_company_ownership_refs();

CREATE OR REPLACE FUNCTION gridworks.freeze_company_ownership_identity()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.ownership_id IS DISTINCT FROM OLD.ownership_id OR NEW.entity_id IS DISTINCT FROM OLD.entity_id OR NEW.entity_type IS DISTINCT FROM OLD.entity_type
       OR NEW.owner_type IS DISTINCT FROM OLD.owner_type OR NEW.owner_id IS DISTINCT FROM OLD.owner_id
       OR NEW.share_bps IS DISTINCT FROM OLD.share_bps OR NEW.effective_from IS DISTINCT FROM OLD.effective_from THEN
        RAISE EXCEPTION 'ownership identity is immutable; append a new projection row';
    END IF;
    IF OLD.active = false THEN
        IF NEW.active IS DISTINCT FROM OLD.active OR NEW.effective_to IS DISTINCT FROM OLD.effective_to THEN
            RAISE EXCEPTION 'closed ownership projection is immutable';
        END IF;
    ELSIF NEW.active = true OR OLD.effective_to IS NOT NULL OR NEW.effective_to IS NULL THEN
        IF NEW.active IS DISTINCT FROM OLD.active OR NEW.effective_to IS DISTINCT FROM OLD.effective_to THEN
            RAISE EXCEPTION 'ownership projection may only close once';
        END IF;
    END IF;
    RETURN NEW;
END
$$;
DROP TRIGGER IF EXISTS company_ownership_identity_guard ON gridworks.company_ownership;
CREATE TRIGGER company_ownership_identity_guard BEFORE UPDATE ON gridworks.company_ownership FOR EACH ROW EXECUTE FUNCTION gridworks.freeze_company_ownership_identity();

CREATE OR REPLACE FUNCTION gridworks.validate_active_ownership_total()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE target_entity text; target_type text; total_bps bigint;
BEGIN
    target_entity := COALESCE(NEW.entity_id, OLD.entity_id);
    target_type := COALESCE(NEW.entity_type, OLD.entity_type);
    IF target_type = 'company' THEN
        PERFORM 1 FROM gridworks.companies WHERE company_id = target_entity FOR UPDATE;
    ELSE
        PERFORM 1 FROM gridworks.company_groups WHERE group_id = target_entity FOR UPDATE;
    END IF;
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

CREATE OR REPLACE FUNCTION gridworks.prevent_referenced_identity_delete()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF TG_TABLE_NAME = 'companies' AND EXISTS (SELECT 1 FROM gridworks.company_ownership WHERE entity_type = 'company' AND entity_id = OLD.company_id) THEN RAISE EXCEPTION 'company is referenced by ownership'; END IF;
    IF TG_TABLE_NAME = 'company_groups' AND EXISTS (SELECT 1 FROM gridworks.company_ownership WHERE entity_type = 'group' AND entity_id = OLD.group_id) THEN RAISE EXCEPTION 'group is referenced by ownership'; END IF;
    IF TG_TABLE_NAME = 'players' AND EXISTS (SELECT 1 FROM gridworks.company_ownership WHERE owner_type = 'player' AND owner_id = OLD.player_id) THEN RAISE EXCEPTION 'player is referenced by ownership'; END IF;
    RETURN OLD;
END
$$;
DROP TRIGGER IF EXISTS companies_ownership_delete_guard ON gridworks.companies;
CREATE TRIGGER companies_ownership_delete_guard BEFORE DELETE ON gridworks.companies FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_referenced_identity_delete();
DROP TRIGGER IF EXISTS company_groups_ownership_delete_guard ON gridworks.company_groups;
CREATE TRIGGER company_groups_ownership_delete_guard BEFORE DELETE ON gridworks.company_groups FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_referenced_identity_delete();
DROP TRIGGER IF EXISTS players_ownership_delete_guard ON gridworks.players;
CREATE TRIGGER players_ownership_delete_guard BEFORE DELETE ON gridworks.players FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_referenced_identity_delete();

COMMIT;

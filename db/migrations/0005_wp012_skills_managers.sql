BEGIN;

SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '30s';

CREATE TABLE IF NOT EXISTS gridworks.player_skills (
    player_id text NOT NULL REFERENCES gridworks.players(player_id) ON DELETE RESTRICT,
    skill_id text NOT NULL,
    cumulative_xp bigint NOT NULL DEFAULT 0 CHECK (cumulative_xp >= 0),
    proficiency_bps integer NOT NULL DEFAULT 0 CHECK (proficiency_bps BETWEEN 0 AND 10000),
    progression_version text NOT NULL,
    last_source_event_id text,
    last_event_time timestamptz,
    repetition_key text,
    repetition_count integer NOT NULL DEFAULT 0 CHECK (repetition_count >= 0),
    PRIMARY KEY (player_id, skill_id)
);

CREATE TABLE IF NOT EXISTS gridworks.player_skill_events (
    source_event_id text PRIMARY KEY,
    player_id text NOT NULL REFERENCES gridworks.players(player_id) ON DELETE RESTRICT,
    skill_id text NOT NULL,
    activity_kind text NOT NULL,
    base_xp bigint NOT NULL CHECK (base_xp >= 0),
    awarded_xp bigint NOT NULL CHECK (awarded_xp >= 0),
    activity_class text NOT NULL CHECK (activity_class IN ('trivial', 'meaningful')),
    repetition_key text NOT NULL,
    occurrence_time timestamptz NOT NULL,
    rules_version text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS player_skill_events_player_idx ON gridworks.player_skill_events(player_id, occurrence_time, source_event_id);

CREATE TABLE IF NOT EXISTS gridworks.managers (
    manager_id text PRIMARY KEY,
    display_name text NOT NULL CHECK (length(display_name) BETWEEN 1 AND 120),
    rarity text NOT NULL CHECK (rarity IN ('Bronze', 'Silver', 'Gold', 'Platinum')),
    specialization_id text NOT NULL,
    total_xp bigint NOT NULL DEFAULT 0 CHECK (total_xp >= 0),
    level integer NOT NULL DEFAULT 1 CHECK (level > 0),
    workload_bps integer NOT NULL DEFAULT 0 CHECK (workload_bps BETWEEN 0 AND 10000),
    fatigue_bps integer NOT NULL DEFAULT 0 CHECK (fatigue_bps BETWEEN 0 AND 10000),
    morale_bps integer NOT NULL DEFAULT 10000 CHECK (morale_bps BETWEEN 0 AND 10000),
    status text NOT NULL CHECK (status IN ('Active', 'Inactive')),
    created_at timestamptz NOT NULL,
    effective_from timestamptz NOT NULL
);

CREATE TABLE IF NOT EXISTS gridworks.manager_skills (
    manager_id text NOT NULL REFERENCES gridworks.managers(manager_id) ON DELETE RESTRICT,
    skill_id text NOT NULL,
    proficiency_bps integer NOT NULL DEFAULT 0 CHECK (proficiency_bps BETWEEN 0 AND 10000),
    potential_bps integer NOT NULL CHECK (potential_bps BETWEEN 0 AND 10000),
    PRIMARY KEY (manager_id, skill_id),
    CHECK (proficiency_bps <= potential_bps)
);

CREATE TABLE IF NOT EXISTS gridworks.manager_traits (
    manager_id text NOT NULL REFERENCES gridworks.managers(manager_id) ON DELETE RESTRICT,
    trait_id text NOT NULL,
    source_ref text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (manager_id, trait_id)
);

CREATE TABLE IF NOT EXISTS gridworks.manager_progression_events (
    source_event_id text PRIMARY KEY,
    manager_id text NOT NULL REFERENCES gridworks.managers(manager_id) ON DELETE RESTRICT,
    activity_kind text NOT NULL,
    awarded_xp bigint NOT NULL CHECK (awarded_xp >= 0),
    occurrence_time timestamptz NOT NULL,
    rules_version text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gridworks.manager_employment_history (
    employment_id text PRIMARY KEY,
    manager_id text NOT NULL REFERENCES gridworks.managers(manager_id) ON DELETE RESTRICT,
    company_id text NOT NULL REFERENCES gridworks.companies(company_id) ON DELETE RESTRICT,
    role text NOT NULL,
    state text NOT NULL CHECK (state IN ('active', 'closed')),
    effective_from timestamptz NOT NULL,
    effective_to timestamptz,
    source_ref text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CHECK ((state = 'active' AND effective_to IS NULL) OR (state = 'closed' AND effective_to IS NOT NULL))
);
CREATE UNIQUE INDEX IF NOT EXISTS manager_employment_active_uq ON gridworks.manager_employment_history(manager_id) WHERE state = 'active';
CREATE INDEX IF NOT EXISTS manager_employment_company_idx ON gridworks.manager_employment_history(company_id, state, manager_id);

CREATE TABLE IF NOT EXISTS gridworks.manager_facility_assignments (
    assignment_id text PRIMARY KEY,
    manager_id text NOT NULL REFERENCES gridworks.managers(manager_id) ON DELETE RESTRICT,
    company_id text NOT NULL REFERENCES gridworks.companies(company_id) ON DELETE RESTRICT,
    facility_id text NOT NULL,
    assignment_role text NOT NULL,
    active boolean NOT NULL DEFAULT true,
    effective_from timestamptz NOT NULL,
    effective_to timestamptz,
    source_ref text NOT NULL,
    CHECK ((active AND effective_to IS NULL) OR ((NOT active) AND effective_to IS NOT NULL))
);
CREATE UNIQUE INDEX IF NOT EXISTS manager_assignment_primary_active_uq ON gridworks.manager_facility_assignments(manager_id) WHERE active AND assignment_role = 'primary';
CREATE INDEX IF NOT EXISTS manager_assignment_company_idx ON gridworks.manager_facility_assignments(company_id, active, manager_id);

CREATE TABLE IF NOT EXISTS gridworks.manager_mutation_receipts (
    idempotency_key text PRIMARY KEY,
    mutation_type text NOT NULL,
    request_digest char(64) NOT NULL CHECK (request_digest ~ '^[0-9a-f]{64}$'),
    result_ref text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION gridworks.prevent_wp012_history_mutation()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'WP-012 history is append-only';
END
$$;
DROP TRIGGER IF EXISTS player_skill_events_append_only ON gridworks.player_skill_events;
CREATE TRIGGER player_skill_events_append_only BEFORE UPDATE OR DELETE ON gridworks.player_skill_events FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_wp012_history_mutation();
DROP TRIGGER IF EXISTS manager_progression_events_append_only ON gridworks.manager_progression_events;
CREATE TRIGGER manager_progression_events_append_only BEFORE UPDATE OR DELETE ON gridworks.manager_progression_events FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_wp012_history_mutation();

CREATE OR REPLACE FUNCTION gridworks.prevent_wp012_employment_rewrite()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.employment_id IS DISTINCT FROM OLD.employment_id OR NEW.manager_id IS DISTINCT FROM OLD.manager_id OR NEW.company_id IS DISTINCT FROM OLD.company_id OR NEW.effective_from IS DISTINCT FROM OLD.effective_from OR NEW.source_ref IS DISTINCT FROM OLD.source_ref THEN
        RAISE EXCEPTION 'employment identity is immutable';
    END IF;
    IF OLD.state = 'closed' OR (OLD.state = 'active' AND NEW.state = 'active') THEN
        IF NEW.state IS DISTINCT FROM OLD.state OR NEW.effective_to IS DISTINCT FROM OLD.effective_to THEN
            RAISE EXCEPTION 'employment may only close once';
        END IF;
    END IF;
    RETURN NEW;
END
$$;
DROP TRIGGER IF EXISTS manager_employment_append_close ON gridworks.manager_employment_history;
CREATE TRIGGER manager_employment_append_close BEFORE UPDATE ON gridworks.manager_employment_history FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_wp012_employment_rewrite();
CREATE OR REPLACE FUNCTION gridworks.prevent_wp012_employment_delete()
RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'employment history is append-close only'; END $$;
DROP TRIGGER IF EXISTS manager_employment_no_delete ON gridworks.manager_employment_history;
CREATE TRIGGER manager_employment_no_delete BEFORE DELETE ON gridworks.manager_employment_history FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_wp012_employment_delete();

CREATE OR REPLACE FUNCTION gridworks.prevent_wp012_assignment_delete()
RETURNS trigger LANGUAGE plpgsql AS $$ BEGIN RAISE EXCEPTION 'assignment history is append-close only'; END $$;
DROP TRIGGER IF EXISTS manager_assignment_no_delete ON gridworks.manager_facility_assignments;
CREATE TRIGGER manager_assignment_no_delete BEFORE DELETE ON gridworks.manager_facility_assignments FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_wp012_assignment_delete();



CREATE OR REPLACE FUNCTION gridworks.prevent_wp012_employment_rewrite()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.employment_id IS DISTINCT FROM OLD.employment_id OR NEW.manager_id IS DISTINCT FROM OLD.manager_id OR NEW.company_id IS DISTINCT FROM OLD.company_id OR NEW.role IS DISTINCT FROM OLD.role OR NEW.effective_from IS DISTINCT FROM OLD.effective_from OR NEW.source_ref IS DISTINCT FROM OLD.source_ref OR NEW.created_at IS DISTINCT FROM OLD.created_at THEN
        RAISE EXCEPTION 'employment history identity/content is immutable';
    END IF;
    IF OLD.state = 'closed' THEN
        IF NEW.state IS DISTINCT FROM OLD.state OR NEW.effective_to IS DISTINCT FROM OLD.effective_to THEN RAISE EXCEPTION 'closed employment is immutable'; END IF;
    ELSIF NEW.state = 'active' OR NEW.effective_to IS NULL THEN
        IF NEW.state IS DISTINCT FROM OLD.state OR NEW.effective_to IS DISTINCT FROM OLD.effective_to THEN RAISE EXCEPTION 'employment may only close once'; END IF;
    END IF;
    RETURN NEW;
END
$$;

CREATE OR REPLACE FUNCTION gridworks.prevent_wp012_assignment_rewrite()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.assignment_id IS DISTINCT FROM OLD.assignment_id OR NEW.manager_id IS DISTINCT FROM OLD.manager_id OR NEW.company_id IS DISTINCT FROM OLD.company_id OR NEW.facility_id IS DISTINCT FROM OLD.facility_id OR NEW.assignment_role IS DISTINCT FROM OLD.assignment_role OR NEW.effective_from IS DISTINCT FROM OLD.effective_from OR NEW.source_ref IS DISTINCT FROM OLD.source_ref THEN
        RAISE EXCEPTION 'assignment history identity/content is immutable';
    END IF;
    IF OLD.active = false THEN
        IF NEW.active IS DISTINCT FROM OLD.active OR NEW.effective_to IS DISTINCT FROM OLD.effective_to THEN RAISE EXCEPTION 'closed assignment is immutable'; END IF;
    ELSIF NEW.active = true OR NEW.effective_to IS NULL THEN
        IF NEW.active IS DISTINCT FROM OLD.active OR NEW.effective_to IS DISTINCT FROM OLD.effective_to THEN RAISE EXCEPTION 'assignment may only close once'; END IF;
    END IF;
    RETURN NEW;
END
$$;
DROP TRIGGER IF EXISTS manager_assignment_append_close ON gridworks.manager_facility_assignments;
CREATE TRIGGER manager_assignment_append_close BEFORE UPDATE ON gridworks.manager_facility_assignments FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_wp012_assignment_rewrite();

CREATE OR REPLACE FUNCTION gridworks.validate_wp012_assignment_employer()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM gridworks.manager_employment_history WHERE manager_id=NEW.manager_id AND company_id=NEW.company_id AND state='active') THEN
        RAISE EXCEPTION 'assignment company must match active manager employer';
    END IF;
    RETURN NEW;
END
$$;
DROP TRIGGER IF EXISTS manager_assignment_employer_guard ON gridworks.manager_facility_assignments;
CREATE TRIGGER manager_assignment_employer_guard BEFORE INSERT OR UPDATE ON gridworks.manager_facility_assignments FOR EACH ROW EXECUTE FUNCTION gridworks.validate_wp012_assignment_employer();

COMMIT;

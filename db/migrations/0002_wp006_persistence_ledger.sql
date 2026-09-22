BEGIN;

SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '30s';

CREATE TABLE IF NOT EXISTS gridworks.simulation_snapshots (
    snapshot_id text PRIMARY KEY,
    owner_ref text,
    schema_version text NOT NULL,
    rules_version text NOT NULL,
    kernel_revision text NOT NULL,
    operational_time_ms bigint NOT NULL CHECK (operational_time_ms >= 0),
    state_digest char(64) NOT NULL CHECK (state_digest ~ '^[0-9a-f]{64}$'),
    payload jsonb NOT NULL,
    previous_snapshot_id text REFERENCES gridworks.simulation_snapshots(snapshot_id),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    accepted_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS simulation_snapshots_digest_idx
    ON gridworks.simulation_snapshots (state_digest, schema_version, rules_version);

CREATE TABLE IF NOT EXISTS gridworks.command_receipts (
    command_id text PRIMARY KEY,
    idempotency_key text NOT NULL UNIQUE,
    command_type text NOT NULL,
    status text NOT NULL CHECK (status IN ('accepted', 'rejected')),
    state_digest_before char(64) CHECK (state_digest_before IS NULL OR state_digest_before ~ '^[0-9a-f]{64}$'),
    state_digest_after char(64) CHECK (state_digest_after IS NULL OR state_digest_after ~ '^[0-9a-f]{64}$'),
    effective_time_ms bigint NOT NULL CHECK (effective_time_ms >= 0),
    result_event jsonb,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gridworks.ledger_accounts (
    account_id text PRIMARY KEY,
    owner_ref text NOT NULL,
    account_class text NOT NULL CHECK (account_class IN ('Asset', 'Liability', 'Equity', 'Revenue', 'Expense', 'ClearingSystem')),
    currency char(3) NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS gridworks.journal_transactions (
    transaction_id text PRIMARY KEY,
    transaction_type text NOT NULL,
    effective_time_ms bigint NOT NULL CHECK (effective_time_ms >= 0),
    idempotency_key text NOT NULL UNIQUE,
    source_ref text NOT NULL,
    currency char(3) NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    reversal_of text REFERENCES gridworks.journal_transactions(transaction_id),
    status text NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'posted')),
    posted_at timestamptz
);
CREATE UNIQUE INDEX IF NOT EXISTS journal_transactions_one_reversal_idx
    ON gridworks.journal_transactions (reversal_of)
    WHERE reversal_of IS NOT NULL;

CREATE TABLE IF NOT EXISTS gridworks.journal_lines (
    line_id text PRIMARY KEY,
    transaction_id text NOT NULL REFERENCES gridworks.journal_transactions(transaction_id) ON DELETE RESTRICT,
    line_sequence integer NOT NULL CHECK (line_sequence > 0),
    account_id text NOT NULL REFERENCES gridworks.ledger_accounts(account_id) ON DELETE RESTRICT,
    currency char(3) NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    side text NOT NULL CHECK (side IN ('Debit', 'Credit')),
    amount_minor bigint NOT NULL CHECK (amount_minor > 0),
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
    UNIQUE (transaction_id, line_sequence)
);

CREATE OR REPLACE FUNCTION gridworks.prevent_posted_journal_mutation()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF TG_TABLE_NAME = 'journal_transactions'
       AND TG_OP = 'UPDATE'
       AND OLD.status = 'draft'
       AND NEW.status = 'posted' THEN
        IF NEW.transaction_id IS DISTINCT FROM OLD.transaction_id
           OR NEW.transaction_type IS DISTINCT FROM OLD.transaction_type
           OR NEW.effective_time_ms IS DISTINCT FROM OLD.effective_time_ms
           OR NEW.idempotency_key IS DISTINCT FROM OLD.idempotency_key
           OR NEW.source_ref IS DISTINCT FROM OLD.source_ref
           OR NEW.currency IS DISTINCT FROM OLD.currency
           OR NEW.reversal_of IS DISTINCT FROM OLD.reversal_of THEN
            RAISE EXCEPTION 'draft-to-posted transition cannot change accounting identity fields';
        END IF;
        IF NEW.posted_at IS NOT NULL THEN
            RAISE EXCEPTION 'posted_at is trigger-owned';
        END IF;
        NEW.posted_at := CURRENT_TIMESTAMP;
        RETURN NEW;
    END IF;
    RAISE EXCEPTION 'journal records are immutable; use a reversal transaction';
END
$$;

CREATE OR REPLACE FUNCTION gridworks.prevent_direct_posted_transaction()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    IF NEW.status <> 'draft' OR NEW.posted_at IS NOT NULL THEN
        RAISE EXCEPTION 'draft transactions require trigger-owned null posted_at';
    END IF;
    RETURN NEW;
END
$$;

DROP TRIGGER IF EXISTS journal_transactions_draft_guard ON gridworks.journal_transactions;
CREATE TRIGGER journal_transactions_draft_guard
    BEFORE INSERT ON gridworks.journal_transactions
    FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_direct_posted_transaction();

CREATE OR REPLACE FUNCTION gridworks.validate_posted_journal_balance()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE line_count bigint; debit_total bigint; credit_total bigint;
BEGIN
    SELECT count(*), COALESCE(sum(CASE WHEN side = 'Debit' THEN amount_minor ELSE 0 END), 0), COALESCE(sum(CASE WHEN side = 'Credit' THEN amount_minor ELSE 0 END), 0)
      INTO line_count, debit_total, credit_total
      FROM gridworks.journal_lines WHERE transaction_id = NEW.transaction_id;
    IF line_count < 2 OR debit_total <= 0 OR debit_total <> credit_total THEN
        RAISE EXCEPTION 'posted journal transaction is not balanced';
    END IF;
    RETURN NEW;
END
$$;

DROP TRIGGER IF EXISTS journal_transactions_balance_guard ON gridworks.journal_transactions;
CREATE CONSTRAINT TRIGGER journal_transactions_balance_guard
    AFTER UPDATE OF status ON gridworks.journal_transactions
    DEFERRABLE INITIALLY DEFERRED
    FOR EACH ROW WHEN (NEW.status = 'posted')
    EXECUTE FUNCTION gridworks.validate_posted_journal_balance();

DROP TRIGGER IF EXISTS journal_transactions_immutable ON gridworks.journal_transactions;
CREATE TRIGGER journal_transactions_immutable
    BEFORE UPDATE OR DELETE ON gridworks.journal_transactions
    FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_posted_journal_mutation();

DROP TRIGGER IF EXISTS journal_lines_immutable ON gridworks.journal_lines;
CREATE TRIGGER journal_lines_immutable
    BEFORE UPDATE OR DELETE ON gridworks.journal_lines
    FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_posted_journal_mutation();

CREATE OR REPLACE FUNCTION gridworks.prevent_posted_journal_line_insert()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE tx_status text;
BEGIN
    SELECT status INTO tx_status FROM gridworks.journal_transactions WHERE transaction_id = NEW.transaction_id FOR SHARE;
    IF tx_status IS NULL THEN
        RAISE EXCEPTION 'journal transaction does not exist';
    END IF;
    IF tx_status = 'posted' THEN
        RAISE EXCEPTION 'journal lines cannot be inserted after transaction is posted';
    END IF;
    RETURN NEW;
END
$$;

DROP TRIGGER IF EXISTS journal_lines_posted_insert_guard ON gridworks.journal_lines;
CREATE TRIGGER journal_lines_posted_insert_guard
    BEFORE INSERT ON gridworks.journal_lines
    FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_posted_journal_line_insert();

CREATE OR REPLACE FUNCTION gridworks.validate_journal_line_currency()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE tx_currency char(3); account_currency char(3);
BEGIN
    SELECT currency INTO tx_currency FROM gridworks.journal_transactions WHERE transaction_id = NEW.transaction_id;
    SELECT currency INTO account_currency FROM gridworks.ledger_accounts WHERE account_id = NEW.account_id;
    IF tx_currency IS NULL OR account_currency IS NULL OR NEW.currency <> tx_currency OR NEW.currency <> account_currency THEN
        RAISE EXCEPTION 'journal line currency does not match transaction currency';
    END IF;
    RETURN NEW;
END
$$;

DROP TRIGGER IF EXISTS journal_lines_currency_guard ON gridworks.journal_lines;
CREATE TRIGGER journal_lines_currency_guard
    BEFORE INSERT ON gridworks.journal_lines
    FOR EACH ROW EXECUTE FUNCTION gridworks.validate_journal_line_currency();

CREATE OR REPLACE VIEW gridworks.ledger_account_balances AS
SELECT
    account_id,
    currency,
    COALESCE(SUM(CASE WHEN side = 'Debit' THEN amount_minor ELSE -amount_minor END), 0)::bigint AS balance_minor
FROM gridworks.journal_lines
JOIN gridworks.journal_transactions USING (transaction_id)
WHERE journal_transactions.status = 'posted'
GROUP BY account_id, currency;

COMMIT;

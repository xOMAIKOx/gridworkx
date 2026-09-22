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
    accepted_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (state_digest, schema_version, rules_version)
);

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
    posted_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

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
    RAISE EXCEPTION 'journal records are immutable; use a reversal transaction';
END
$$;

DROP TRIGGER IF EXISTS journal_transactions_immutable ON gridworks.journal_transactions;
CREATE TRIGGER journal_transactions_immutable
    BEFORE UPDATE OR DELETE ON gridworks.journal_transactions
    FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_posted_journal_mutation();

DROP TRIGGER IF EXISTS journal_lines_immutable ON gridworks.journal_lines;
CREATE TRIGGER journal_lines_immutable
    BEFORE UPDATE OR DELETE ON gridworks.journal_lines
    FOR EACH ROW EXECUTE FUNCTION gridworks.prevent_posted_journal_mutation();

CREATE OR REPLACE FUNCTION gridworks.validate_journal_line_currency()
RETURNS trigger LANGUAGE plpgsql AS $$
DECLARE tx_currency char(3);
BEGIN
    SELECT currency INTO tx_currency FROM gridworks.journal_transactions WHERE transaction_id = NEW.transaction_id;
    IF tx_currency IS NULL OR NEW.currency <> tx_currency THEN
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
GROUP BY account_id, currency;

COMMIT;

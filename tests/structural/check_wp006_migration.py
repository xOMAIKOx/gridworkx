from pathlib import Path

root = Path(__file__).resolve().parents[2]
migration = root / "db/migrations/0002_wp006_persistence_ledger.sql"
text = migration.read_text()
required = [
    "gridworks.simulation_snapshots",
    "gridworks.command_receipts",
    "gridworks.ledger_accounts",
    "gridworks.journal_transactions",
    "gridworks.journal_lines",
    "simulation_snapshots_digest_idx",
    "journal_transactions_balance_guard",
    "journal_transactions_one_reversal_idx",
    "journal_transactions_balance_guard",
    "journal_lines_currency_guard",
    "journal_lines_posted_insert_guard",
    "draft-to-posted transition cannot change accounting identity fields",
    "status text NOT NULL DEFAULT 'draft'",
    "UNIQUE (transaction_id, line_sequence)",
    "prevent_posted_journal_mutation",
    "ledger_account_balances",
]
for marker in required:
    if marker not in text:
        raise SystemExit(f"WP-006 migration missing required contract: {marker}")
for forbidden in ("DROP SCHEMA", "DROP TABLE", "CREATE ROLE", "ALTER SYSTEM"):
    if forbidden in text.upper():
        raise SystemExit(f"WP-006 migration contains forbidden destructive/host operation: {forbidden}")
if not text.lstrip().startswith("BEGIN;") or not text.rstrip().endswith("COMMIT;"):
    raise SystemExit("WP-006 migration must be transaction-wrapped")
print("WP-006 migration contract shape is valid.")

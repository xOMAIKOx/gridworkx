# WP-006 PostgreSQL Persistence and Immutable Ledger Evidence

## Persistence boundaries

Rust remains authoritative for simulation, graph, failure and material transitions. WP-006 adds repository/domain contracts for persistence without moving simulation or accounting rules into SQL or Go.

`SnapshotRecord` persists the canonical JSON snapshot payload, SHA-256 digest, schema/rules/kernel versions, operational time, owner reference and previous-snapshot linkage. Loading validates the payload through the Rust snapshot parser and rejects version, kernel, operational-time or digest mismatches.

`CommandReceipt` is the durable identity boundary for command ID, idempotency key, command type, status, effective time and before/after state digests. The in-memory `executed_commands` vector is now a bounded 1,024-entry replay window; long-lived duplicate/economic protection belongs to PostgreSQL `command_receipts` uniqueness rather than unbounded snapshots.

## PostgreSQL migration

`db/migrations/0002_wp006_persistence_ledger.sql` targets PostgreSQL 18-compatible standard features and defines:

- `simulation_snapshots`;
- `command_receipts`;
- `ledger_accounts`;
- `journal_transactions`;
- `journal_lines`;
- an immutable `ledger_account_balances` derived view.

The migration is transaction-wrapped, rerunnable under the repository convention, contains no credentials/roles/host paths and performs no provisioning.

## Ledger model

The Rust ledger uses checked signed `i64` minor units with a separate three-letter currency code. `LedgerState` contains stable accounts and immutable journal transactions with ordered debit/credit lines. Every posted transaction requires positive amounts, known accounts, one currency and exact debit/credit equality.

Duplicate idempotency keys return a deterministic duplicate event without adding money. Reversals append an opposite-line transaction referencing the original; the original transaction and lines are never mutated. Balances are derived by summing immutable journal lines, with debit positive and credit negative for the account view.

The synthetic proof transfers 1,000 `CRD` minor units from a system clearing account to a synthetic wallet, proves duplicate rejection, and proves reversal returns both derived balances to zero. Physical WP-005 quantities remain entirely separate from monetary amounts.

## Verification

Rust tests cover balanced/multi-line semantics, imbalance, mixed currency, invalid/overflow amounts, duplicate posting, append-only reversal, derived balances, snapshot persistence integrity, kernel ledger command integration and bounded replay history. The migration shape checker verifies required tables, uniqueness, immutability trigger, derived view, transaction wrapping and absence of destructive/provisioning operations.

No live PostgreSQL server was installed or started on ERIS. Any live-PG execution is limited to an already-authorized CI/test context; no host, service, role, database or port mutation occurred.

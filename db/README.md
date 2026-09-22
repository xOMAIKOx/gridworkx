# Database foundation

`db/migrations` and `db/seeds` define repository-side PostgreSQL conventions for the accepted transactional authority. The WP-001 migration reserves schema metadata, domain contracts, content manifests, system principals and domain reservation metadata without implementing the ledger, gameplay tables or shared database service.

No database was created or modified on ERIS by WP-001. PostgreSQL client validation uses the pre-existing host tooling only when a later authorized test requires it.

WP-006 adds repository-only PG18-targeted snapshot, command-receipt, ledger-account, journal-transaction and journal-line contracts in `0002_wp006_persistence_ledger.sql`. It does not provision or connect to an ERIS database. Balances derive from immutable journal lines; corrections are reversal transactions.

# Database foundation

`db/migrations` and `db/seeds` define repository-side PostgreSQL conventions for the accepted transactional authority. The WP-001 migration reserves schema metadata, domain contracts, content manifests, system principals and domain reservation metadata without implementing the ledger, gameplay tables or shared database service.

No database was created or modified on ERIS by WP-001. PostgreSQL client validation uses the pre-existing host tooling only when a later authorized test requires it.

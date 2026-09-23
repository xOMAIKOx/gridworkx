from pathlib import Path

root = Path(__file__).resolve().parents[2]
migration = root / "db/migrations/0004_wp008_company_ownership.sql"
text = migration.read_text()
required = [
    "gridworks.companies",
    "gridworks.company_groups",
    "gridworks.company_ownership",
    "gridworks.company_group_membership",
    "gridworks.ownership_history",
    "gridworks.company_mutation_receipts",
    "companies_name_canonical_uq",
    "companies_name_skeleton_uq",
    "company_group_membership_active_company_uq",
    "company_ownership_total_guard",
    "company_name_claims",
    "company_name_claim_guard",
    "company_group_name_claim_guard",
    "company_ownership_reference_guard",
    "company_ownership_identity_guard",
    "FOR UPDATE",
    "ownership_history_append_only",
    "company_group_membership_identity_guard",
    "ownership_history_append_only",
    "active ownership must total exactly 10000 basis points",
]
for marker in required:
    if marker not in text:
        raise SystemExit(f"WP-008 migration missing required contract: {marker}")
for forbidden in ("DROP SCHEMA", "DROP TABLE", "CREATE ROLE", "ALTER SYSTEM", "password", "secret"):
    if forbidden.lower() in text.lower():
        raise SystemExit(f"WP-008 migration contains forbidden contract: {forbidden}")
if not text.lstrip().startswith("BEGIN;") or not text.rstrip().endswith("COMMIT;"):
    raise SystemExit("WP-008 migration must be transaction-wrapped")
print("WP-008 migration contract shape is valid.")

from pathlib import Path

root = Path(__file__).resolve().parents[2]
sql = (root / "db/migrations/0005_wp012_skills_managers.sql").read_text()
required_tables = [
    "player_skills", "player_skill_events", "managers", "manager_skills", "manager_traits",
    "manager_progression_events", "manager_employment_history", "manager_facility_assignments", "manager_mutation_receipts",
]
for table in required_tables:
    if f"gridworks.{table}" not in sql:
        raise SystemExit(f"WP-012 migration missing {table}")
for token in ["append-only", "append-close", "manager_employment_active_uq", "manager_assignment_primary_active_uq", "manager_assignment_employer_guard", "manager_employment_append_close", "manager_assignment_append_close", "manager_employment_assignment_guard", "manager_mutation_receipts_immutable", "prevent_wp012_receipt_mutation", "request_digest", "skill_id text NOT NULL"]:
    if token not in sql:
        raise SystemExit(f"WP-012 migration missing integrity contract: {token}")
if "REFERENCES gridworks.players" not in sql or "REFERENCES gridworks.companies" not in sql:
    raise SystemExit("WP-012 migration missing accepted identity/company references")
if "DROP TABLE" in sql or "DROP SCHEMA" in sql:
    raise SystemExit("WP-012 migration contains destructive DDL")
print("WP-012 migration shape is valid.")

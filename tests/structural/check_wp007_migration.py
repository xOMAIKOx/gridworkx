from pathlib import Path

root = Path(__file__).resolve().parents[2]
migration = root / "db/migrations/0003_wp007_identity_profile.sql"
text = migration.read_text()
required = [
    "gridworks.accounts",
    "gridworks.players",
    "gridworks.player_profiles",
    "gridworks.guest_sessions",
    "gridworks.account_identity_links",
    "gridworks.identity_mutation_receipts",
    "player_profiles_handle_canonical_uq",
    "player_profiles_handle_skeleton_uq",
    "account_identity_links_account_idx",
    "token_digest char(64)",
    "UNIQUE (issuer, subject)",
    "identity_mutation_receipts",
    "player_profiles_player_account_fk",
    "proof_reference text NOT NULL",
]
for marker in required:
    if marker not in text:
        raise SystemExit(f"WP-007 migration missing required contract: {marker}")
for forbidden in ("DROP SCHEMA", "DROP TABLE", "CREATE ROLE", "ALTER SYSTEM", "raw_token", "plaintext"):
    if forbidden.lower() in text.lower():
        raise SystemExit(f"WP-007 migration contains forbidden contract: {forbidden}")
if not text.lstrip().startswith("BEGIN;") or not text.rstrip().endswith("COMMIT;"):
    raise SystemExit("WP-007 migration must be transaction-wrapped")
print("WP-007 migration contract shape is valid.")

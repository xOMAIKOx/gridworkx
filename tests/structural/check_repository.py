from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
REQUIRED_PATHS = [
    "apps/game/project.godot",
    "apps/admin/package.json",
    "apps/tools/src/validate-content.mjs",
    "crates/gridworks-sim/Cargo.toml",
    "services/go.mod",
    "packages/schemas/catalog.json",
    "packages/content/domain-catalog.json",
    "packages/localization/catalogue.json",
    "packages/visual-assets/manifests/semantic-assets.json",
    "db/migrations/0001_wp001_foundation.sql",
    "db/migrations/0002_wp006_persistence_ledger.sql",
    "db/migrations/0003_wp007_identity_profile.sql",
    "db/migrations/0004_wp008_company_ownership.sql",
    "docs/evidence/WP-007_GUEST_FIRST_IDENTITY_PROFILE.md",
    "docs/evidence/WP-008_COMPANY_OWNERSHIP_BUSINESS_IDENTITY.md",
    "packages/content/config/reserved-company-names.json",
    "packages/content/config/reserved-handles.json",
    "db/seeds/0001_domain_registry.sql",
    "ops/config/gridworks.env.example",
    "ops/systemd/gridworks-api.service",
    "ops/systemd/gridworks-worker.service",
    "ops/systemd/gridworks-realtime.service",
    "ops/systemd/gridworks-sim-validator.service",
    ".github/workflows/ci.yml",
]
REQUIRED_DIRECTORIES = [
    "apps/game",
    "apps/admin",
    "apps/tools",
    "crates/gridworks-sim",
    "services/api",
    "services/worker",
    "services/realtime",
    "services/sim-validator",
    "packages/schemas",
    "packages/content",
    "packages/localization",
    "packages/visual-assets",
    "db/migrations",
    "db/seeds",
    "ops/systemd",
    "ops/config",
    "tests/integration",
    "tests/deterministic",
    "tests/economy",
    "tests/security",
]
FORBIDDEN_NAMES = {"Dockerfile", "docker-compose.yml", "docker-compose.yaml", "compose.yml", "compose.yaml"}

missing = [path for path in REQUIRED_PATHS if not (ROOT / path).is_file()]
missing_directories = [path for path in REQUIRED_DIRECTORIES if not (ROOT / path).is_dir()]
forbidden = [str(path.relative_to(ROOT)) for path in ROOT.rglob("*") if path.name in FORBIDDEN_NAMES and ".git" not in path.parts]

if missing or missing_directories or forbidden:
    if missing:
        print("missing files:", ", ".join(missing))
    if missing_directories:
        print("missing directories:", ", ".join(missing_directories))
    if forbidden:
        print("forbidden container manifests:", ", ".join(forbidden))
    raise SystemExit(1)

print(f"Repository structure is complete: {len(REQUIRED_PATHS)} files and {len(REQUIRED_DIRECTORIES)} directories checked.")

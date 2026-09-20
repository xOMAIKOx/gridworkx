# GRIDWORKS

GRIDWORKS is a persistent cooperative economic and systems-simulation game. This repository is the native Linux monorepo foundation established by WP-001.

## Architecture boundaries

- `crates/gridworks-sim` is the only canonical Rust simulation boundary. It is deliberately Godot-, network-, database- and wall-clock-independent.
- `services/` is one Go workspace containing the API, worker, realtime and explicit simulation-validator process boundaries.
- `apps/game` is the Godot presentation foundation. It does not reimplement simulation rules.
- `apps/admin` is the React + TypeScript administrative foundation.
- `packages/schemas` owns versioned cross-language contract definitions.
- `packages/content`, `packages/localization` and `packages/visual-assets` own versioned data and semantic identifiers.
- `db/` contains repository-side PostgreSQL migration conventions only; WP-001 does not provision or run a database.
- `ops/systemd` contains native systemd templates/examples only; no unit is installed or enabled by WP-001.

The repository remains a monorepo and uses native Linux/systemd as the target runtime. Docker, Podman, Compose, Kubernetes and OCI runtime dependencies are prohibited.

## Verification

Prerequisites are recorded in [the WP-001 ERIS receipt](docs/evidence/WP-001_ERIS_PREREQUISITE_RECEIPT.md) and the authoritative GitHub issue comment.

```text
make validate       # structure, schemas, Rust, Go and admin foundations
make test           # executable tests for all local foundations
make security       # secret scan, forbidden-runtime scan and shell lint
make godot-check    # headless Godot project parse/editor validation
```

The WP-001 foundation intentionally provides boundaries, contracts, configuration ownership and executable smoke checks. It does not implement gameplay or the WP-002 simulation kernel.

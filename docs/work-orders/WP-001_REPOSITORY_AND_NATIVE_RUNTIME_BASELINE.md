# GRIDWORKS — WP-001 Repository and Native Runtime Baseline

**Status:** READY FOR ENGINEERING / NOT YET ACCEPTED  
**Architecture owner:** Architecture  
**Implementation model:** GPT-5.6 Luna XHigh or equivalent bounded implementation worker  
**Review:** GPT-5.6 Sol-class Architecture review  
**Design baseline before this work order:** `bef77ffba313f87525dd540d5f4ef323c1c9b167`

## 1. Purpose

Create the durable monorepo, build/test, schema/content and native-runtime foundation required by the accepted GRIDWORKS architecture.

WP-001 is a **foundation work package**. It must establish boundaries and conventions without prematurely implementing gameplay.

## 2. Controlling references

Engineering must read and obey:

- `docs/design/GRIDWORKS_PRODUCT_PHILOSOPHY_AND_GAME_DESIGN.md`
- `docs/design/GRIDWORKS_DP2_CORE_SYSTEMS_SPECIFICATION.md`
- `docs/design/GRIDWORKS_DP3_ENGINEERING_ARCHITECTURE_AND_WORK_PACKAGES.md`
- `docs/decisions/ADR-001_GODOT_RUST_INTEGRATION.md`
- `docs/decisions/ADR-002_IDENTITY_AND_ACCOUNT_LINKING.md`
- `docs/decisions/ADR-003_SCHEMA_AND_SERIALIZATION.md`
- `docs/decisions/ADR-004_CONTENT_RULES_DISTRIBUTION.md`
- `docs/decisions/ADR-005_INITIAL_SERVER_TOPOLOGY.md`
- `docs/decisions/ADR-006_GRIDWORKS_OWNED_ASSETS_AND_TREASURY.md`
- `docs/decisions/ADR-007_REAL_ESTATE_AND_STARTING_PATHS.md`
- `docs/decisions/ADR-008_WORLD_TIME_SEASONS_AND_TIMERS.md`
- `docs/decisions/ADR-009_WORLD_BOOTSTRAP_PRESENTATION_RELEASE1.md`

If implementation appears to require contradicting an accepted ADR, stop that decision path and raise the conflict in the WP-001 GitHub issue. Do not silently redesign Architecture in code.

## 3. Runtime constraints

Target runtime is native Linux/systemd.

**Prohibited:**
- Docker;
- Podman;
- Docker Compose;
- OCI runtime dependency;
- Kubernetes;
- container-only build/test assumptions.

WP-001 authorizes a **bounded ERIS development prerequisite audit and toolchain preparation** before any coding begins.

Authorized host scope is **ERIS only**, with repository workspace:

`/srv/gridworx/`

Engineering must first inspect the host and record the current prerequisite/toolchain state. It may then install only missing development prerequisites required to build/test the accepted stack.

Permitted prerequisite preparation may include, as required:
- Git and normal repository tooling;
- Rust toolchain/cargo/rustfmt/clippy;
- Go toolchain;
- Node.js and the selected package manager for the admin/tooling workspace;
- Godot 4.x editor/headless tooling required for project validation;
- native compiler/build essentials and pkg-config;
- PostgreSQL client/development tooling required for migrations/build tests;
- SQLite development/client tooling;
- NATS client/development tooling where required for local validation;
- static analysis, formatting, schema-validation and secret-scanning tools required by WP-001 gates.

Rules:
- inspect before installing;
- prefer distro/native packages where appropriate, otherwise use pinned/checksummed upstream toolchains;
- record package/tool name, version, source and installation command in the WP-001 evidence;
- do not replace a working compatible estate-standard package merely to chase the newest version;
- do not install Docker, Podman, Compose, Kubernetes, OCI runtimes or container-dependent tooling;
- do not install or start application infrastructure services merely because client/dev libraries are needed.

Still prohibited without separate owner authorization:
- changing nginx/reverse proxies;
- changing DNS;
- changing firewall/WireGuard;
- allocating externally reachable application ports;
- creating/modifying production-style systemd application units on ERIS;
- creating shared PostgreSQL databases/users;
- installing/running NATS/PostgreSQL/object-storage daemons for GRIDWORKS;
- deploying GRIDWORKS services;
- changes outside ERIS.

Repository templates/examples for later native deployment remain allowed.

**Hard gate:** no application/repository implementation work begins until the ERIS prerequisite audit is complete, missing authorized prerequisites are installed, and a concise prerequisite receipt is posted to the WP-001 GitHub issue.

## 4. Required monorepo foundation

Create or normalize the following structure where appropriate:

```
apps/
  game/                  # Godot client foundation only
  admin/                 # React + TypeScript admin foundation
  tools/                 # development/content/balance tooling

crates/
  gridworks-sim/         # canonical Rust simulation package foundation

services/
  api/                   # Go API foundation
  worker/                # Go background worker foundation
  realtime/              # Go realtime/WebSocket foundation
  sim-validator/         # explicit validator boundary; may remain thin in WP-001

packages/
  schemas/               # canonical schema catalogue / generation boundary
  content/               # versioned game-content definitions
  localization/          # locale packs, glossary/source catalogues
  visual-assets/         # semantic asset/icon manifest conventions, not production art

db/
  migrations/
  seeds/

docs/
  architecture/
  decisions/
  design/
  work-orders/

ops/
  systemd/
  config/
  scripts/

tests/
  integration/
  deterministic/
  economy/
  security/
```

Do not create a microservice repository split.

## 5. Required technology baselines

Establish buildable/testable foundations for:

- Godot 4.x client project;
- Rust canonical simulation crate;
- Go backend workspace/modules;
- React + TypeScript admin application;
- PostgreSQL 18 migration conventions;
- SQLite client/offline boundary;
- NATS/JetStream configuration contracts;
- S3-compatible object-storage configuration contracts.

Valkey is **not** required and must not be introduced by default.

OpenSearch/dedicated search infrastructure is **not** required.

## 6. Canonical simulation boundary

The Rust crate must be structurally capable of becoming the single canonical simulation/rules engine.

WP-001 must establish:
- crate/package ownership;
- no network access dependency in the simulation core;
- no database dependency in the simulation core;
- no direct device-wall-clock dependency;
- deterministic-test location;
- versioned command/state schema boundary placeholder.

Do not implement the WP-002 simulation kernel in WP-001.

## 7. Schema and semantic-ID foundation

Establish repository conventions for stable semantic identifiers and versioned contracts.

The foundation must be able to represent future namespaces including:

- `resource.*`
- `fault.*`
- `transport.*`
- `industry.*`
- `event.*`
- `contract.*`
- `property.*`
- `region.*`
- `climate.*`

Create validation/tests for malformed or duplicate IDs where practical.

Do not hard-code translated prose into authoritative schemas.

## 8. Domain scaffolding required in WP-001

WP-001 must make these accepted domains visible in package/module/schema organization without implementing their gameplay:

- Identity / Account;
- Player Profile;
- Company / Ownership;
- World / Region / Site;
- Facility / Component;
- Resources / Inventory;
- Construction / Maintenance;
- Skills / Managers;
- Trade / Marketplace;
- Contracts / Work Exchange;
- Consortium / JV;
- Accounting / Ledger;
- Business Lifecycle;
- Social / Messaging;
- Notifications;
- Localization;
- Real Estate / Land / Property;
- GRIDWORKS system-principal ownership;
- Non-player Treasury / Authority ledgers;
- Starter Path selection;
- Business operating state;
- GRIDWORKS buyer-of-last-resort policy;
- World/Portfolio/Site/Operational presentation boundaries;
- transport-mode contracts;
- bootstrap-economy policy;
- time/climate configuration.

Scaffolding means explicit module/schema/config ownership and testable structure, not fake implementations.

## 9. Time and climate foundation

Represent configuration contracts for:

- authoritative UTC timestamps;
- initial ratio: **1 in-game day = 1 real hour**;
- operational simulation clock;
- regional climate calendar/profile;
- timed-job configuration;
- offline-safe policy metadata;
- rules/content version association.

No client wall clock may become authoritative.

## 10. Real-estate / system ownership foundation

Reserve domain/schema locations for:

- land parcels;
- buildings;
- property ownership;
- property listings;
- leases;
- occupancy;
- renovation jobs;
- GRIDWORKS-owned inventory;
- authority principals;
- treasury accounts/transactions;
- system purchase offers;
- valuation policy;
- carrying costs;
- reacquisition rules.

Do not implement the full Release-1 property game in WP-001.

## 11. World/presentation foundation

The client/content structure must support later implementation of:

1. World Map;
2. Player Portfolio Map;
3. Site / Business View;
4. Operational Detail View.

Visual architecture assumptions:
- 2.5D isometric/oblique;
- modular stylized low-poly assets;
- 2D overlays;
- limited/fixed site rotations;
- semantic icon/asset IDs.

WP-001 does **not** require production art or finished UI.

## 12. Bootstrap economy foundation

Create policy/config/schema ownership for GRIDWORKS system participation so later work can support an economy with only one active player.

Reserve concepts for:
- system demand;
- system supply;
- contracts/jobs;
- transport opportunities;
- property/business inventory;
- emergency supply;
- asset liquidity;
- market-depth responsive participation.

Do not implement live market intervention algorithms yet.

## 13. Social/moderation foundation

Repository/domain structure must reserve policy/config ownership for:
- DM permissions;
- block;
- mute;
- report;
- rate limits;
- naming/impersonation/confusable rules;
- moderation audit.

Do not implement social UI in WP-001.

## 14. Native runtime repository artefacts

Create repository-side native runtime templates/conventions for the accepted processes:
- `gridworks-api`;
- `gridworks-worker`;
- `gridworks-realtime`;
- `gridworks-sim-validator`.

Systemd artefacts must be templates/examples only and must:
- use non-root service users;
- support environment/config files outside source;
- define restart behavior;
- use explicit working directories;
- avoid embedding secrets;
- avoid privileged ports by default.

No unit may be installed/enabled/started on a host under WP-001.

## 15. Configuration and secrets

Provide:
- example configuration;
- environment variable documentation;
- fail-fast validation pattern;
- secret placeholders only;
- no committed credentials/tokens/private keys;
- no production endpoint assumptions.

## 16. CI and quality gates

CI must exercise relevant repository foundations, including where applicable:

- formatting/lint;
- Rust build/test;
- Go build/test;
- TypeScript build/typecheck/test;
- schema/content validation;
- deterministic test placeholder/fixture proving the test harness works;
- secret scan;
- forbidden-container-runtime scan;
- repository structural checks.

The exact CI provider may use GitHub Actions unless the repository already dictates otherwise.

## 17. Minimum smoke contracts

WP-001 should prove the skeleton is executable without pretending gameplay exists.

Acceptable smoke evidence includes:
- Rust crate test passes;
- Go services compile and expose only minimal health/version interfaces if needed;
- admin app builds;
- Godot project opens/parses/builds sufficiently for CI-supported validation;
- schemas validate;
- native unit templates validate statically;
- no container manifests/dependencies exist.

Avoid ornamental endpoints or placeholder business logic.

## 18. Allowed paths

WP-001 may create/modify:
- root build/workspace/config files required by the monorepo;
- `apps/**`;
- `crates/**`;
- `services/**`;
- `packages/**`;
- `db/**`;
- `ops/**`;
- `tests/**`;
- `.github/**`;
- documentation directly required by WP-001.

Accepted design/ADR documents are controlling references and should not be substantively rewritten by Engineering.

## 19. Explicitly out of scope

Do not implement:
- production gameplay;
- full facility simulation;
- real economy/market matching;
- actual ledger/accounting engine;
- account provider integrations;
- real-estate gameplay;
- finance/exchange;
- messaging;
- manager/recruitment gameplay;
- production localization catalogue;
- final art;
- detailed transport simulation;
- deployment;
- host configuration.

Those belong to later WPs.

## 20. Required engineering evidence

Engineering must return through the WP-001 GitHub issue/PR with:

- branch;
- exact parent SHA;
- exact HEAD SHA;
- files changed;
- repository tree summary;
- versions/toolchains selected;
- commands executed;
- build/test/lint results;
- secret-scan result;
- forbidden-container scan result;
- evidence paths;
- known deviations;
- unresolved risks;
- architecture conflicts discovered;
- ERIS prerequisite audit receipt, including detected versions, installed prerequisites and commands used;
- explicit statement that no host changes outside the authorized ERIS prerequisite scope were performed.

If a PR is opened, keep it draft until Architecture reviews it.

## 21. Done-when

WP-001 is complete only when Architecture can verify:

1. monorepo foundations exist and are coherent;
2. accepted technology boundaries are represented;
3. Rust is structurally the sole canonical simulation location;
4. cross-language schema/content ownership is explicit;
5. accepted domains introduced through ADR-006..009 are not omitted;
6. ERIS prerequisite audit/preparation completed before coding, and native systemd runtime templates exist without application deployment;
7. build/test/validation gates pass;
8. no container dependency exists;
9. no secrets are committed;
10. no substantive gameplay was prematurely implemented;
11. Engineering handback is complete and reproducible.

## 22. Stop condition

After publishing the handback, Engineering stops.

Engineering does not start WP-002 or any other work package without a separate Architecture instruction.

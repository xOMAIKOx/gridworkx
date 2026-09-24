# WP-012 Skills and Manager Domain Evidence

## Scope and boundary

WP-012 establishes player skills, trusted progression, anti-grind, manager state/potential/traits, workload/fatigue/morale, employment/assignment history, deterministic Rust information-quality capability, PostgreSQL persistence shape and authenticated owner-only read APIs.

Recruitment remains WP-013. No daily draws, candidate generation, rarity odds, recruit/reject/shortlist, Platinum acquisition, manager market, economy, network, host, deployment, database provisioning or container scope was added.

## Versioned content

`packages/content/config/skills-managers.json` is the shared versioned source for:

- 12 player skill families;
- 9 manager skill dimensions;
- Bronze/Silver/Gold/Platinum potential ceilings;
- five deterministic traits;
- anti-grind multipliers `[10000, 5000, 2500, 0]` for repeated trivial activity.

`packages/schemas/skills-manager-content.schema.json` is registered and exercised by the content validator.

## Canonical Rust progression

`crates/gridworks-sim/src/progression.rs` provides deterministic integer logic for:

- monotonic XP and proficiency derivation;
- duplicate source-event replay safety;
- trivial repetition diminishing awards;
- manager potential caps and level derivation;
- deterministic trait vectors;
- bounded fatigue contribution;
- effective diagnostic capability.

The capability formula is clamped to `0..10000`, uses player skill plus optional current manager skill, bounded trait effects and fatigue/morale, and has no rarity or payment input.

The integration proof uses accepted WP-004 inspection semantics: capability can move an observation from unknown to observed without changing the underlying symptom identity; it does not fabricate telemetry.

## Deterministic golden

`tests/deterministic/fixtures/wp012/fixture.json` and `crates/gridworks-sim/tests/wp012_fixture.rs` pin:

- trivial awards `100, 50, 25, 0`;
- duplicate award `0`;
- cumulative XP `175` and proficiency `1750`;
- Gold current technical skill `8000`, potential `8500`;
- Platinum current technical skill `1000`, potential `10000`;
- Gold diagnostic capability `7460` versus fresh Platinum `3000`;
- same `symptom.abnormal_vibration` identity across capability levels.

Gold outperforms fresh Platinum because of current relevant skill and bounded trait history; Platinum retains higher potential. Rarity is not a direct output multiplier.

## Persistence and trusted mutation boundary

`db/migrations/0005_wp012_skills_managers.sql` adds player skill/event, manager/skill/trait/progression, employment, assignment and progression-receipt tables. It includes bounded checks, accepted player/company FKs, append-only event triggers, append-close history protections, active-employer/primary-assignment uniqueness and a distinct idempotency namespace.

`services/internal/progression` owns transport-independent view/domain validation and trusted progression representation. PostgreSQL repository methods use transactions, deterministic source/idempotency locking and row locking for progression/employment paths. Facility references intentionally remain semantic IDs because no authoritative PostgreSQL facility table exists yet.

## Authenticated read APIs

Added read-only routes:

- `GET /api/v1/me/skills` — authenticated self skills only;
- `GET /api/v1/companies/{company_id}/managers` — authenticated active company owner/controller only;
- `GET /api/v1/companies/{company_id}/managers/{manager_id}` — same owner-only boundary.

No public XP, manager creation, employment transfer, facility assignment or recruitment route exists. OpenAPI and the route registry are semantically synchronized; response views do not expose internal receipt/source digests.

API and sqlmock tests cover authenticated reads, owner rejection, deterministic ordering/view validation and repository SQL methods.

## Verification

- Rust WP-012 golden and full `gridworks-sim` tests: PASS.
- Go domain/API/PostgreSQL tests: PASS.
- Migration/content/OpenAPI structural checks: PASS.
- Full push CI: PASS.
- Full PR CI: PASS.
- Existing WP-010 native/Xvfb/runtime parity and WP-011 deterministic/controller/visual gates: PASS through CI.
- Local full-history gitleaks reports one pre-existing WP-010 historical evidence false positive; repository CI security gate passes on the current shallow checkout. No security policy was changed.
- Local host has Rust 1.75 and no Rustup; workspace-wide local Cargo verification is therefore environment-blocked by the accepted Rust 1.80 requirement. CI runs the pinned Rust 1.80 toolchain and is authoritative.
- Live PostgreSQL 18 execution: `BLOCKED_DEPENDENCY` — no authorized native PG18 environment was available. Repository-local migration/repository/sqlmock proof is complete; no fake live-DB claim is made.

No host, deployment, provider, database provisioning, port, proxy, systemd or container mutation occurred.


## R2 integrity remediation

R2 makes `packages/content/config/skills-managers.json` authoritative for player-skill metadata, manager skill IDs, rarity potential/trait capacities, trait semantic vectors and anti-grind policy. Rust embeds and validates the shared content; Go has a generated representation with a repository drift test.

Trusted Go progression now rejects unknown activity classes/inapplicable activity kinds/version drift, uses checked arithmetic, and rejects overflow before persistence. Manager progression is an internal trusted operation with deterministic XP/level/skill-cap behavior, row locking and durable replay/conflict receipts.

Employment and assignment mutations use durable request-digest receipts, append-close replay semantics and PostgreSQL row locks. Migration triggers freeze historical identity/content, forbid reopen/delete, and enforce assignment company equals active employer.

The Rust WP-004 integration proof now exercises diagnosis confidence: low capability yields lower unresolved confidence while higher capability yields higher confidence, with the same underlying symptom identity and no fabricated evidence.

Trait vectors preserve semantic trade-offs across diagnostic/workload/fatigue/morale dimensions and rarity trait-capacity validation prevents universal dominance from becoming a hidden multiplier.

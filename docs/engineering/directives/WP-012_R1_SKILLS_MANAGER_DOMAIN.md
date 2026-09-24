# GRIDWORKS WP-012 R1 — Skills and Manager Domain Engineering Directive

**Status:** GO  
**Technical Authority:** GPT-5.6 Sol  
**Execution agent:** Devin  
**Worker model:** GPT-5.6 Luna  
**Worker effort:** XHigh  
**Implementation class:** D2 — bounded complex  
**Exact Engineering parent:** `466faa05ad24d0dc1af35a1fa3d71a04e8920961`

## 1. Objective

Implement the complete WP-012 skills/manager domain defined by:

`docs/work-orders/WP-012_SKILLS_MANAGER_DOMAIN.md`

This pass must establish the repository-local Rust, PostgreSQL, Go domain/API, schema/content and deterministic-test contracts needed for later recruitment and manager usage.

## 2. Repository / branch

Repository:

`xOMAIKOx/gridworkx`

Engineering branch:

`engineering/wp-012-skills-managers`

Branch from exactly:

`466faa05ad24d0dc1af35a1fa3d71a04e8920961`

The Architecture work-order/directive commit on `main` is intentionally not the Engineering parent.

Open one draft PR against `main` after the first implementation commit.

No silent rebase.

## 3. Execution route

Architecture classification: **D2 bounded complex**.

Execute as:

**Devin + GPT-5.6 Luna XHigh**

Architecture and product laws are fixed. Ordinary implementation/debugging remains inside this route.

If work exposes a genuine unresolved architecture contradiction, STOP and return it to Architecture. Do not silently escalate model/effort.

## 4. Mandatory read-before-edit

Read:

1. `AGENTS.md`
2. `AI_AGENT_COLLABORATION_PROTOCOL.md`
3. `docs/GITHUB_COLLABORATION_TRANSPORT.md`
4. `docs/work-orders/WP-012_SKILLS_MANAGER_DOMAIN.md`
5. Product Philosophy §§10–11
6. DP2 §§13–14
7. DP3 authority/persistence/API sections
8. WP-004/006/007/008/009 accepted work orders/evidence
9. current main / exact parent ancestry

## 5. Required implementation sequence

### R1.1 — versioned content

Add strict versioned definitions for:

- 12 player skill families;
- manager skill dimensions;
- rarity/potential policy;
- manager traits;
- anti-grind/progression policy.

One shared source; do not fork tunable tables across languages.

### R1.2 — Rust progression module

Implement deterministic pure logic for:

- XP threshold/proficiency calculation;
- manager potential cap;
- trait effect resolution;
- bounded fatigue contribution;
- effective diagnostic capability.

Add integration proof against WP-004 failure/diagnosis semantics.

### R1.3 — deterministic fixture

Add the WP-012 fixture including veteran Gold vs fresh Platinum.

Pin exact expected outputs.

### R1.4 — PostgreSQL migration

Add `0005_wp012_skills_managers.sql` with required tables/constraints/history/receipt semantics.

Add structural migration tests.

### R1.5 — Go domain

Add transport-independent progression/manager domain package.

No domain logic only in HTTP handlers.

### R1.6 — PostgreSQL repository

Implement parameterized read/internal-mutation methods, row/advisory locking and replay/idempotency protections.

### R1.7 — API/OpenAPI

Add only:

- `GET /api/v1/me/skills`
- `GET /api/v1/companies/{company_id}/managers`
- `GET /api/v1/companies/{company_id}/managers/{manager_id}`

All manager routes require authenticated active company ownership.

No XP/recruitment/create/transfer public mutation routes.

### R1.8 — full proof

Run all new and existing gates and produce evidence/handback.

## 6. Architectural invariants

- Rust owns simulation-relevant manager/skill calculation.
- PostgreSQL owns server transactional/history authority.
- Go does not duplicate Rust simulation math.
- client does not award XP.
- rarity does not directly multiply output.
- higher skill improves information/confidence, not observed reality.
- fatigue never becomes a paywall/stamina lock.
- recruitment remains WP-013.
- manager market remains later work.
- accepted WP-011 gameplay/goldens remain unchanged.

## 7. Acceptance criteria

- **AC-012-01:** exact parent `466faa05ad24d0dc1af35a1fa3d71a04e8920961`.
- **AC-012-02:** 12 player skill definitions versioned.
- **AC-012-03:** player XP/proficiency deterministic/monotonic.
- **AC-012-04:** repeated trivial XP diminishes per versioned policy.
- **AC-012-05:** duplicate source event cannot award twice.
- **AC-012-06:** no public client XP mutation.
- **AC-012-07:** manager entity/skills/potential/traits/workload/fatigue/morale modeled.
- **AC-012-08:** Bronze/Silver/Gold/Platinum rarity has no direct universal multiplier.
- **AC-012-09:** proficiency cannot exceed potential ceiling.
- **AC-012-10:** veteran Gold fixture outperforms fresh Platinum in relevant diagnostic capability.
- **AC-012-11:** at least one trait effect is executable and deterministic.
- **AC-012-12:** fatigue effect bounded; no action lock/payment input.
- **AC-012-13:** higher capability changes diagnosis confidence without changing observed evidence identity.
- **AC-012-14:** migration contains required skills/manager/history/receipt tables.
- **AC-012-15:** progression/event history append-only.
- **AC-012-16:** employment history append-close + one active employer.
- **AC-012-17:** assignment history append-close + one active primary assignment.
- **AC-012-18:** company/employer refs use accepted WP-008 company identity.
- **AC-012-19:** authenticated `/me/skills` read works.
- **AC-012-20:** owned-company manager roster/detail reads work.
- **AC-012-21:** non-owner/cross-company manager reads fail closed.
- **AC-012-22:** OpenAPI and route registry semantically match.
- **AC-012-23:** no public manager creation/recruitment/transfer/assignment mutation.
- **AC-012-24:** no WP-013 recruitment tables/odds/draws/candidates implemented.
- **AC-012-25:** existing WP-001–011 gates remain green.
- **AC-012-26:** evidence and durable handback complete.
- **AC-012-27:** no unauthorized host/deployment/provider/DB/port/container mutation.

## 8. Mandatory verification

Run at minimum:

- `make validate`
- `make test`
- `make security`
- Rust fmt/test/clippy
- new Rust WP-012 deterministic/integration tests
- migration structural validation
- Go domain unit tests
- Go API route/OpenAPI tests
- PostgreSQL repository sqlmock/accepted repository-method tests
- existing WP-010 native/Xvfb/runtime parity
- existing WP-011 controller/golden/visual gates
- final push CI
- final PR CI

If a live authorized PG18 environment is unavailable, record that honestly; do not fabricate live evidence.

## 9. Mandatory STOP conditions

STOP and hand back `BLOCKED` if:

- branch does not start from exact parent;
- implementing WP-012 requires changing accepted WP-011 golden/scenario behavior;
- success requires a public arbitrary XP endpoint;
- success requires implementing recruitment candidate generation/daily draws;
- success requires manager market/transfer UI;
- success requires inventing a fake facility PostgreSQL authority/FK;
- success requires duplicating Rust simulation math in Go/GDScript;
- success requires changing version pins;
- success requires weakening an accepted security/integrity gate;
- worker/model substitution is required;
- host/deployment mutation is required.

## 10. Required durable handback

Commit:

`docs/engineering/handbacks/wp-012/WP012_R1_HANDBACK.md`

Include:

- directive/work-order paths;
- execution agent/model/effort/class;
- repository/branch;
- exact parent;
- implementation commits;
- final HEAD;
- files/components changed;
- skill registry/progression version;
- anti-grind exact policy/result;
- manager rarity/potential representation;
- trait definitions/effects;
- Gold-vs-Platinum fixture/result;
- diagnostic confidence/no-fabricated-evidence proof;
- migration tables/constraints;
- employment/assignment invariants;
- API routes + authorization matrix;
- OpenAPI proof;
- repository SQL/locking/idempotency proof;
- full validation/test/security results;
- final push/PR CI URLs;
- live PG18 status;
- recruitment/market exclusions;
- AC-012-01..27 mapping;
- deviations/blockers;
- unauthorized environment mutation status;
- final state READY FOR SOL/PRO REVIEW / BLOCKED / STOP.

## 11. GitHub return path

After committing the handback:

1. post one concise WP-012 issue handback pointer;
2. post one concise draft-PR receipt pointing to it;
3. give the Owner only the concise status/pointer;
4. STOP for Architecture review.

Do not merge.
Do not start WP-013.

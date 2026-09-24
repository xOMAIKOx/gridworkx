# GRIDWORKS — WP-012 Skills and Manager Domain

**Status:** AUTHORIZED DESIGN / CONTROLLING WORK ORDER  
**Technical Authority:** GPT-5.6 Sol  
**Engineering route:** Devin + GPT-5.6 Luna XHigh  
**Implementation class:** D2 — bounded complex  
**Exact Engineering parent:** `466faa05ad24d0dc1af35a1fa3d71a04e8920961`

## 1. Purpose

WP-012 implements the canonical **player-skill** and **facility-manager** domains.

It must establish the progression, persistence, deterministic simulation inputs and read APIs required by later recruitment, contracts, Work Exchange and manager-market work without pulling those later systems forward.

Controlling product laws:

- player skill and manager skill are separate systems;
- skill improves through relevant activity;
- repetitive zero-risk spam must not be an optimal progression path;
- manager rarity represents potential / unusual combinations / trait capacity, **not** a universal production multiplier;
- a veteran Gold manager must be able to outperform a newly recruited Platinum manager in relevant circumstances;
- managers improve **information quality**, not merely output numbers;
- skill/manager effects may improve confidence/effectiveness but must never fabricate telemetry or repeal physical constraints;
- fatigue/workload must never become an energy/stamina paywall;
- no manager/recruitment mechanic may create paid power.

## 2. Exact ancestry

Engineering branches from exactly:

`466faa05ad24d0dc1af35a1fa3d71a04e8920961`

This is the owner-authorized WP-011 merge commit.

The later WP-012 Architecture documentation commit on `main` is intentionally **not** the Engineering implementation parent.

No silent rebase.

## 3. Controlling references

Read before implementation:

1. `AGENTS.md`
2. `AI_AGENT_COLLABORATION_PROTOCOL.md`
3. `docs/GITHUB_COLLABORATION_TRANSPORT.md`
4. Product Philosophy §§10–11
5. DP2 Part V §§11–14
6. DP3 canonical Rust / authoritative-state / persistence / WP ledger sections
7. WP-004 failure/diagnosis contracts
8. WP-006 persistence/ledger conventions
9. WP-007 identity/player model
10. WP-008 company/ownership model
11. WP-009 Go API/OpenAPI conventions
12. WP-010 bridge invariants
13. active WP-012 issue / draft PR

## 4. Domain boundary

WP-012 owns:

- player skill definitions/state;
- trusted skill-XP award/progression rules;
- anti-grind/diminishing eligibility;
- manager identity and current progression state;
- manager rarity/potential model;
- manager skill state;
- manager traits and deterministic modifier vectors;
- manager workload/fatigue/morale state;
- manager employment history;
- manager facility-assignment history/projection;
- deterministic skill/manager contribution to information-quality calculations;
- PostgreSQL persistence for the above;
- authenticated read APIs for player skills and owned-company managers;
- internal/repository mutation primitives needed by later trusted modules.

WP-012 does **not** own candidate generation/recruitment.

## 5. Initial player skill registry

Create versioned semantic definitions for the DP2 top-level skill families:

- `skill.mechanical`
- `skill.electrical`
- `skill.process_engineering`
- `skill.agriculture`
- `skill.mining`
- `skill.energy`
- `skill.water`
- `skill.logistics`
- `skill.construction`
- `skill.commerce`
- `skill.finance`
- `skill.management`

Definitions belong in versioned content/config, not scattered transport constants.

Each definition must include at least:

- semantic skill ID;
- player-facing label key/string source;
- applicable activity classes;
- progression curve/version reference.

## 6. Player skill state

Canonical server-owned player skill state must include at least:

- player_id;
- skill_id;
- cumulative XP;
- derived proficiency;
- progression/rules version;
- last progression event reference/time.

Use an integer deterministic representation.

Preferred internal proficiency scale:

- `0..10000` basis points;
- presentation may display `0..100`.

Proficiency must be derived from XP against versioned progression content. It must not be directly client-writable.

## 7. Trusted XP award model

Create a transport-independent domain operation equivalent to:

`ApplyPlayerSkillEvent(player_id, activity_event) -> progression result`

The event must carry server/trusted-domain evidence such as:

- immutable source event ID/reference;
- activity kind;
- skill ID;
- base XP;
- difficulty/eligibility classification;
- repetition/novelty key;
- authoritative occurrence time;
- rules version.

Requirements:

- duplicate source event cannot award twice;
- negative XP is impossible;
- XP/proficiency never regresses;
- proficiency cannot exceed the configured maximum;
- unknown skill/activity/version fails closed;
- no public API accepts arbitrary XP from a client.

## 8. Anti-grind rule

Implement a versioned deterministic anti-grind policy for repeated **trivial** actions.

The exact policy must be data/config driven and pinned by tests.

Initial accepted shape may use diminishing multipliers equivalent to:

- first eligible trivial repetition: 100%;
- second: 50%;
- third: 25%;
- fourth and later in the same repetition series: 0%.

Meaningful/non-trivial activity may use different eligibility.

The server/trusted domain decides the activity classification and repetition key.

A client must not be able to declare its own action “novel” or “high difficulty”.

## 9. Manager entity

Manager state must include at least:

- manager_id;
- display name;
- rarity;
- specialization ID;
- manager total XP;
- derived level;
- current manager skills;
- skill potential ceilings;
- traits;
- workload_bps;
- fatigue_bps;
- morale_bps;
- status;
- created/effective metadata.

Use opaque immutable IDs.

No user-facing code may treat a manager as a commodity object owned like inventory.

## 10. Rarity and potential

Supported rarity:

- Bronze
- Silver
- Gold
- Platinum

Rarity controls only versioned potential dimensions such as:

- potential ceiling ranges;
- possible trait count/trait rarity;
- adaptability/specialization opportunity.

Rarity must **not** apply a generic production/output multiplier.

A Platinum manager with low current skill may perform worse today than a veteran Gold manager with high relevant skill.

Commit a deterministic fixture proving this invariant.

## 11. Manager skill dimensions

Initial manager dimensions:

- `manager_skill.operations`
- `manager_skill.technical`
- `manager_skill.maintenance`
- `manager_skill.safety`
- `manager_skill.leadership`
- `manager_skill.logistics`
- `manager_skill.energy_efficiency`
- `manager_skill.crisis_response`
- `manager_skill.mentoring`

Manager proficiency uses the same deterministic integer discipline as player skills but has a per-skill potential ceiling.

Manager proficiency cannot exceed its ceiling.

Potential ceiling is not itself a current bonus.

## 12. Manager progression

Provide trusted-domain manager XP/progression operations.

Manager progression sources may include:

- operating a facility;
- diagnosing faults;
- managing repairs;
- production;
- crisis recovery;
- mentoring where later supported.

For WP-012, trusted deterministic test events are sufficient.

No public client endpoint may directly award manager XP.

Manager level must be derived from versioned XP thresholds, not client supplied.

## 13. Manager traits

Create versioned semantic trait definitions for at least:

- `trait.aggressive_operator`
- `trait.maintenance_first`
- `trait.cost_cutter`
- `trait.mentor`
- `trait.crisis_specialist`

Traits must be represented as deterministic effect vectors / semantic effects, not prose-only flags.

No trait may be universally dominant.

Where a product axis is not yet implemented (e.g. operating cost, mature wear model), preserve the semantic effect in data/domain contracts but do not invent fake economy or new simulation mathematics.

At least one trait must affect an already accepted axis in executable tests.

## 14. Information-quality integration

Implement one canonical Rust calculation for diagnostic capability using:

- relevant player skill;
- optional assigned manager relevant skill;
- bounded manager trait adjustment;
- bounded fatigue effect;
- rules/content version.

Requirements:

- result is deterministic;
- result is clamped `0..10000`;
- no manager means player skill remains usable;
- fatigue reduces manager contribution but does not lock the player out of actions;
- no premium/membership input exists;
- rarity is not an input to the final capability formula except indirectly through already-earned/current skill/potential history;
- a veteran Gold fixture can exceed a fresh Platinum fixture in resulting capability.

Use this capability through existing failure/diagnosis semantics in a Rust integration test and prove:

- higher capability can increase diagnosis confidence;
- the underlying observed symptom/evidence identity does not change merely because the actor is more skilled;
- no telemetry is fabricated.

Do not duplicate failure math in Go/GDScript.

## 15. Manager workload/fatigue/morale

Persist bounded integer state:

- workload `0..10000`;
- fatigue `0..10000`;
- morale `0..10000`.

WP-012 may apply fatigue only as a bounded contribution modifier.

Do not implement:

- energy/stamina action blocking;
- real-money recovery;
- premium-only recovery;
- punitive log-in timers.

Rest/deputy/rotation staffing mechanics may be introduced later; do not fabricate them here.

## 16. Employment history

Manager employment is server authoritative.

Persist append-close employment history with:

- employment_id;
- manager_id;
- company_id;
- role/state;
- effective_from;
- effective_to;
- source_ref;
- created_at.

Rules:

- one active employer at a time;
- historical rows cannot be deleted;
- closed rows cannot be reopened/rewritten;
- active employment may close once;
- company must exist;
- manager identity remains stable across future transfers.

WP-012 does not implement the manager market.

## 17. Facility assignment history

Persist manager facility assignment projection/history with:

- assignment_id;
- manager_id;
- employing company_id;
- facility_id semantic ID;
- assignment role;
- active;
- effective_from;
- effective_to;
- source_ref.

Rules:

- one active primary facility assignment per manager;
- assignment company must match active employer;
- history is append-close;
- no delete/rewrite.

There is no authoritative PostgreSQL facility table yet, so do **not** invent a fake FK.

Document the deferred facility-reference FK/integrity gate explicitly.

Do not expose an unsafe public assignment mutation until facility ownership can be authoritatively verified.

## 18. Manager fixtures / no production recruitment

WP-012 may create deterministic **test fixtures only** for managers.

Required fixture pair:

- veteran Gold manager with high relevant current skill;
- fresh Platinum manager with higher potential but lower current relevant skill.

The veteran Gold must outperform fresh Platinum in at least the accepted diagnostic-capability scenario.

Do not create production manager rows in migrations/seeds.

## 19. PostgreSQL migration

Add:

`db/migrations/0005_wp012_skills_managers.sql`

Required table families:

- `gridworks.player_skills`
- `gridworks.player_skill_events`
- `gridworks.managers`
- `gridworks.manager_skills`
- `gridworks.manager_traits`
- `gridworks.manager_progression_events`
- `gridworks.manager_employment_history`
- `gridworks.manager_facility_assignments`
- `gridworks.manager_mutation_receipts` or one explicitly documented equivalent receipt namespace

Requirements:

- FK player skills/events to accepted `gridworks.players`;
- FK employer/history to accepted `gridworks.companies`;
- append-only progression/event histories;
- append-close employment/assignment semantics;
- one active employer constraint;
- one active primary assignment constraint;
- bounded bps/check constraints;
- no destructive migration;
- no plaintext secrets;
- no DB roles/system config.

## 20. Mutation/idempotency discipline

All server-side manager/skill mutations require deterministic idempotency/source references.

Do not reuse identity/company mutation receipt namespaces.

Use one progression/manager domain namespace with request/source digest where needed.

Concurrent mutation paths must use PostgreSQL row/advisory locking where required to prevent:

- double XP award;
- duplicate active employment;
- conflicting active assignment;
- double manager progression.

## 21. Go domain package

Add a transport-independent package, preferred:

`services/internal/progression`

or equivalent clearly bounded name.

It owns:

- skill/manager view models;
- trusted XP progression operations;
- authorization-independent domain validation;
- manager employment/assignment repository contracts;
- stable domain errors.

Do not put domain rules only inside HTTP handlers.

Do not duplicate Rust simulation mathematics.

## 22. Public API — read surface only

Extend the accepted WP-009 route registry/OpenAPI with authenticated read endpoints:

- `GET /api/v1/me/skills`
- `GET /api/v1/companies/{company_id}/managers`
- `GET /api/v1/companies/{company_id}/managers/{manager_id}`

Requirements:

- authenticated principal required;
- player skill endpoint returns only current player skills;
- company manager endpoints require the authenticated player to be an active owner/controller of the company under accepted WP-008 ownership semantics;
- no cross-company manager disclosure through this owner-only surface;
- no recruitment endpoint;
- no arbitrary XP endpoint;
- no public manager-create endpoint;
- no public manager-employment-transfer endpoint;
- no public facility-assignment mutation until facility ownership authority exists.

Response views may include:

- manager ID/name;
- rarity;
- specialization;
- current level/current skills;
- traits;
- workload/fatigue/morale;
- active employment;
- active facility assignment if any;
- owner-visible potential band.

Do not expose internal receipt/source digests.

## 23. PostgreSQL API adapter

Extend the WP-009 PostgreSQL repository with explicit parameterized SQL for the new read surface and internal domain mutations.

Tests must prove:

- authenticated self skill read;
- owned-company manager roster/detail read;
- non-owner access rejected;
- missing manager not found;
- no cross-company leak;
- deterministic ordering;
- actual repository SQL methods are exercised with sqlmock or accepted repository test pattern;
- transaction/locking behavior on progression/employment paths.

## 24. OpenAPI

Update the single repository OpenAPI 3.1 contract.

Requirements:

- route registry and OpenAPI remain semantically synchronized;
- strict schemas;
- stable error envelopes;
- no undocumented mutation routes;
- no recruitment semantics in WP-012.

## 25. Rust module

Add a dedicated canonical Rust module, preferred:

`crates/gridworks-sim/src/progression.rs`

or equivalent.

It owns simulation-relevant pure deterministic calculations only:

- proficiency threshold evaluation;
- manager potential cap;
- trait modifier vector resolution;
- effective diagnostic capability;
- bounded fatigue contribution.

Server persistence/history remains Go/PostgreSQL authority.

## 26. Cross-language content/config

Add versioned content files under `packages/content/config/` for:

- skill definitions/progression version;
- manager skill definitions;
- manager rarity/potential policy;
- manager trait definitions;
- anti-grind policy.

Add strict schemas under `packages/schemas/` where consistent with existing repository patterns.

Do not embed tunable progression tables independently in Rust and Go.

One versioned source must control shared values.

## 27. Deterministic golden

Commit one WP-012 golden fixture proving at least:

- player skill XP/proficiency progression;
- trivial-repeat diminishing award;
- duplicate source replay awards zero/returns replay-safe result;
- manager skill cap;
- veteran Gold current-skill profile;
- fresh Platinum higher-potential/lower-current profile;
- veteran Gold diagnostic capability > fresh Platinum capability;
- fatigue bounded effect;
- trait effect;
- same underlying observation identity under low/high capability;
- diagnosis confidence differs only through accepted capability semantics.

Preferred path:

`tests/deterministic/fixtures/wp012/fixture.json`

Pure Rust asserts simulation-relevant outputs.

Go/domain tests assert persistence/progression representations against the same semantic fixture where practical.

## 28. Security / abuse requirements

Explicitly test:

- no public XP award;
- no public manager creation;
- no public recruitment;
- no client-supplied rarity/potential mutation;
- no negative/overflow XP;
- duplicate XP source cannot farm XP;
- manager cannot have two active employers;
- manager cannot have conflicting active primary assignments;
- non-owner cannot read private company manager roster;
- fatigue cannot be bypassed by membership/payment field;
- Platinum has no direct bonus multiplier.

## 29. Recruitment boundary — WP-013

WP-012 must **not** implement:

- 5 free / ~8 member daily recruitment opportunities;
- daily refresh;
- candidate generation;
- rarity odds;
- recruit/reject/shortlist;
- Platinum acquisition logic;
- candidate inventory;
- paid/member recruitment differences;
- candidate market/referral.

Those belong to WP-013.

WP-012 may define the manager entity that WP-013 will instantiate.

## 30. Manager market boundary

Do not implement:

- permanent employment transfer marketplace;
- temporary manager lease/contract marketplace;
- manager bidding/pricing;
- referral/trading;
- consortium assignment.

Employment/history primitives may exist internally, but no market behavior.

## 31. Godot boundary

No new production Godot manager/skills UI is required in WP-012.

Do not destabilize the accepted WP-011 opening slice.

If a Godot/Rust bridge proof is useful, it must be test-only or additive and must not change WP-011 goldens.

The authoritative user-facing manager/recruitment UI will follow the server domain/recruitment contracts.

## 32. Existing gates

All existing WP-001–WP-011 gates remain green.

Specifically preserve:

- WP-010 native/Xvfb/extension parity;
- WP-011 deterministic opening golden and visual evidence;
- WP-009 API strictness/security;
- WP-006/007/008 migrations/integrity.

## 33. Live PostgreSQL gate

Repository-local migration/schema/domain/API proof is required now.

Actual PostgreSQL 18 execution remains an environment gate if no authorized native PG18 test environment is available.

Do not weaken repository tests to simulate a live DB claim.

Document `BLOCKED_DEPENDENCY` honestly if live PG18 is unavailable.

This does not block repository-local completion if the work order's repository gates are otherwise met.

## 34. Allowed paths

Primary:

- `crates/gridworks-sim/**`
- `services/internal/**`
- `services/api/**`
- `db/migrations/**`
- `packages/schemas/**`
- `packages/content/config/**`
- `tests/**`
- `.github/workflows/**` only if necessary
- `Makefile` only if necessary
- `docs/evidence/**`
- `docs/engineering/handbacks/wp-012/**`

Do not modify accepted identity/company semantics except additive repository interfaces needed for authorization.

## 35. Evidence

Create:

`docs/evidence/WP-012_SKILLS_MANAGER_DOMAIN.md`

Document:

- skill registry/version;
- progression representation;
- anti-grind policy;
- manager entity;
- rarity/potential invariant;
- manager skill dimensions;
- traits/effect vectors;
- diagnostic capability calculation;
- veteran Gold vs fresh Platinum proof;
- fatigue behavior;
- employment/assignment history rules;
- migration tables/constraints;
- API routes/authorization;
- no-recruitment boundary;
- deterministic golden;
- repository tests/CI;
- live PG18 status;
- no host/deployment statement.

## 36. Done-when

Architecture can independently verify:

1. exact parent preserved;
2. player and manager skills are separate;
3. canonical skill registry is versioned;
4. XP progression is deterministic and monotonic;
5. trivial repetition diminishes;
6. duplicate source cannot award XP twice;
7. client cannot award XP directly;
8. manager rarity is not a direct power multiplier;
9. current skill is capped by potential;
10. veteran Gold can outperform fresh Platinum in a relevant fixture;
11. traits are deterministic trade-off/effect vectors;
12. fatigue is bounded and not an action/paywall lock;
13. higher skill/manager capability changes confidence without fabricating observations;
14. manager employment history is append-close;
15. one active employer invariant holds;
16. one active primary assignment invariant holds;
17. migration constraints/structural tests pass;
18. authenticated skill/manager read APIs exist;
19. company manager reads enforce ownership;
20. OpenAPI and route registry remain synchronized;
21. no recruitment system is implemented;
22. no manager market is implemented;
23. existing WP-001–011 gates remain green;
24. evidence + durable handback are complete;
25. no unauthorized host/deployment work occurred.

## 37. Stop

After handback, Engineering stops.

Do not merge WP-012.
Do not start WP-013 without owner authorization.

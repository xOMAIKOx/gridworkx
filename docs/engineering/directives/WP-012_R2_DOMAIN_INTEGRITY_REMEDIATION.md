# GRIDWORKS WP-012 R2 — Domain Integrity Remediation

**Status:** GO — REMEDIATION AUTHORIZED  
**Technical Authority:** GPT-5.6 Sol  
**Execution agent:** Devin  
**Worker model:** GPT-5.6 Luna  
**Worker effort:** XHigh  
**Implementation class:** D2 — bounded complex remediation  
**Reviewed R1 HEAD:** `4ca762ae70306d6287da364465fdcd805508eead`  
**Engineering branch:** `engineering/wp-012-skills-managers`

## 1. Disposition

WP-012 R1 is **REQUEST_CHANGES**.

The implementation direction is accepted, the repository builds/tests are green, and ancestry is correct, but several controlling WP-012 domain/integrity requirements are not yet satisfied.

Continue on the existing Engineering branch from exact R1 HEAD:

`4ca762ae70306d6287da364465fdcd805508eead`

Do **not** rebase onto `main` or onto this Architecture directive commit.

## 2. R1 accepted / frozen

Preserve unless a change is strictly required by the remediation below:

- exact original Engineering parent `466faa05ad24d0dc1af35a1fa3d71a04e8920961`;
- draft PR #27 and branch;
- WP-012-only scope;
- read-only public API route set:
  - `GET /api/v1/me/skills`
  - `GET /api/v1/companies/{company_id}/managers`
  - `GET /api/v1/companies/{company_id}/managers/{manager_id}`
- PostgreSQL table-family names;
- no public XP/manager-create/recruitment/transfer/assignment mutation routes;
- veteran-Gold / fresh-Platinum product invariant;
- no rarity direct universal output multiplier;
- bounded fatigue and no payment/stamina input;
- semantic facility IDs without a fake PostgreSQL facility FK;
- WP-013 recruitment boundary;
- WP-001–WP-011 accepted behavior and gates;
- no host/deployment/provider/database provisioning/port/proxy/systemd/container work.

## 3. R2-01 — one authoritative progression/content source

R1 creates `packages/content/config/skills-managers.json`, but the same tunable/domain values are independently hard-coded again in Rust and Go.

Examples include:

- Rust `PLAYER_SKILL_IDS`;
- Rust anti-grind multipliers;
- Rust rarity potential caps;
- Rust trait diagnostic effects;
- Go `PlayerSkillIDs`;
- Go anti-grind multipliers;
- Go progression version constants.

This violates the WP-012 rule that one versioned source controls shared values and that tunable progression tables are not independently embedded across languages.

Also, the current content file represents player skills as bare strings, while the controlling work order requires each player skill definition to carry at least:

- semantic skill ID;
- player-facing label key/string source;
- applicable activity classes;
- progression curve/version reference.

### Required remediation

1. Make the versioned content/config the controlling source for shared progression policy.
2. Rust and Go must consume/validate the same pinned semantic content rather than maintain independent tunable copies.
3. It is acceptable to use build-time embedding / generated typed representation if deterministic and repository-local, but generated output must be mechanically derived from the one source and checked for drift.
4. Strengthen the schema so it proves:
   - exactly the required 12 unique player skill definitions;
   - exactly the required manager skill dimensions;
   - the four supported rarity policies;
   - unique trait IDs;
   - required player-skill metadata;
   - bounded progression/trait values;
   - version presence/consistency.
5. Add drift tests proving Rust/Go behavior is controlled by the shared content.

## 4. R2-02 — trusted progression must fail closed and be overflow-safe

The Go progression path currently accepts any non-`trivial` `ActivityClass` as full-value progression and does not validate `activity_kind` against the skill definition.

The Go arithmetic also uses unchecked signed `int64` multiplication/addition, so sufficiently large trusted events can overflow and violate monotonic/non-negative progression.

### Required remediation

- reject unknown activity classes;
- reject unknown/inapplicable activity kinds;
- reject unknown rules/progression versions;
- validate the existing state version before mutation;
- use checked arithmetic for XP award, cumulative XP and derived proficiency;
- return a stable domain overflow/invalid-event error rather than wrap;
- prove negative/overflow XP cannot be persisted;
- add boundary tests for maximum accepted values and overflow rejection.

Rust and Go semantics must remain aligned through the shared content contract.

## 5. R2-03 — manager progression is incomplete

R1 creates manager state and `manager_progression_events`, but no complete trusted manager XP/progression operation exists.

WP-012 requires trusted manager progression, including derived level/current skill progression bounded by potential.

### Required remediation

Implement repository-local trusted manager progression primitives that:

- accept immutable trusted source/event identity;
- apply versioned manager XP progression;
- derive level from the versioned progression policy;
- progress the targeted manager skill deterministically;
- never exceed that skill's potential ceiling;
- persist manager total XP / derived level / skill state atomically;
- append `manager_progression_events`;
- are replay safe and concurrency safe;
- are not exposed as a public client mutation route.

Add deterministic unit/repository tests, including potential-cap behavior.

## 6. R2-04 — durable replay/idempotency contract is not complete

R1 creates `manager_mutation_receipts`, but the implemented manager employment mutations do not use a durable request-digest receipt.

The current advisory lock is only a concurrency primitive, not a durable replay receipt.

Current consequences include:

- replaying `OpenManagerEmployment` can return `ErrActiveEmployment` rather than replaying the original success;
- replaying `CloseManagerEmployment` can return `ErrNoActiveEmployment` rather than replaying the original success;
- same idempotency/source identity with different payload is not proven to conflict;
- `CloseManagerEmployment` accepts `sourceRef` but does not persist/use it;
- persistent player-skill replay currently returns the previous awarded-XP value while the canonical deterministic/Rust contract records duplicate award as zero, creating semantic divergence.

### Required remediation

For every trusted server-side WP-012 mutation:

- establish one explicit domain receipt/source namespace;
- persist mutation type + request/source digest + result reference where replay identity requires a digest;
- same key/source + same semantic request replays deterministically without applying state twice;
- same key/source + different mutation/payload fails with a stable conflict;
- replay result semantics are consistent across Rust/domain/repository representations;
- advisory/row locks remain for concurrency, but do not substitute for durable receipts.

Add direct repository tests for same-request replay and changed-payload conflict.

## 7. R2-05 — employment/assignment append-close integrity is incomplete

The migration does not yet enforce the full append-close contract.

### Employment

The employment update guard leaves fields such as `role` and `created_at` mutable and allows closed rows to be rewritten in fields other than `state/effective_to`.

A close operation must not be able to mutate unrelated historical identity/content.

### Facility assignment

The assignment table currently has a no-delete trigger but no equivalent immutable-identity / close-once update guard.

As written, a row can be rewritten or reopened by update.

The database also does not enforce the WP-012 rule that an assignment's company matches the manager's active employer.

### Required remediation

- freeze all historical identity/content fields that are not explicitly the close transition;
- active employment/assignment may transition to closed exactly once;
- closed rows are immutable;
- rows cannot be reopened;
- delete remains forbidden;
- assignment company must match the manager's active employer at mutation time;
- preserve the one-active-employer and one-active-primary-assignment constraints;
- add trusted repository primitives for assignment open/close if required to prove these invariants;
- test the actual SQL/repository paths, not only token presence in a structural script.

## 8. R2-06 — diagnosis-confidence proof does not prove the required contract

The current Rust integration test exercises `FailureCommand::Inspect` and proves `Unknown -> Observed` while preserving the symptom ID.

That is useful, but the controlling acceptance criterion specifically requires proof that higher capability can increase **diagnosis confidence** without changing/fabricating the underlying observed symptom/evidence identity.

### Required remediation

Add a deterministic WP-004 integration proof that:

1. begins from the same accepted underlying fault/symptom/evidence identity;
2. uses the WP-012 capability result through the accepted diagnosis path;
3. demonstrates lower versus higher diagnosis confidence/effectiveness as supported by the existing WP-004 contract;
4. proves the symptom/evidence identity is unchanged;
5. proves no new telemetry/evidence identity is fabricated.

Do not invent parallel diagnosis mathematics in WP-012.

The existing inspect-visibility proof may remain as additional coverage.

## 9. R2-07 — trait semantics must preserve trade-offs and avoid dominance

The current shared content reduces every trait to a single `diagnostic_capability_bps` number.

This causes two problems:

- several product traits lose their semantic meaning (`cost_cutter`, `mentor`, `aggressive_operator`, etc.);
- on the only represented axis, `crisis_specialist +500` strictly dominates smaller positive diagnostic values, contrary to the requirement that no trait be universally dominant.

### Required remediation

Represent traits as deterministic semantic effect vectors.

- Preserve future/unimplemented axes as semantic data without inventing new economy/simulation mathematics.
- At least one effect must remain executable against an already accepted axis.
- Encode meaningful trade-offs/conditional specialization so no trait is universally better across the modeled semantic vector.
- Enforce rarity trait-capacity constraints in the manager-domain validation where applicable.
- Keep all values controlled by the shared versioned content source.

## 10. R2-08 — verification depth must match the repaired contracts

The current structural migration test largely checks for strings/tokens, and the current API fake repository does not exercise company ownership itself.

R2 must add focused proof for the repaired contracts, including at least:

- unknown activity/class/version rejection;
- XP overflow rejection;
- player replay semantic consistency;
- manager progression + potential cap;
- manager mutation same-request replay;
- changed-payload idempotency conflict;
- employment close-once / closed-row immutability;
- assignment close-once / no-reopen / employer-company match;
- owner roster/detail repository reads;
- missing manager;
- cross-company isolation;
- deterministic ordering;
- actual repository transaction/locking behavior;
- diagnosis-confidence / unchanged-evidence integration proof;
- shared-content drift prevention.

## 11. CI / environment

R1 push, PR and final-handback-head CI are green. Preserve that.

Live PostgreSQL 18 remains an honest `BLOCKED_DEPENDENCY` if no authorized native PG18 environment exists.

Do not weaken repository-local gates and do not claim live PostgreSQL execution without it.

## 12. Required R2 handback

Update/replace the durable handback at:

`docs/engineering/handbacks/wp-012/WP012_R1_HANDBACK.md`

or add an R2 handback alongside it if repository convention prefers, but the final pointer must unambiguously identify the exact R2 HEAD and remediation evidence.

The handback must include:

- R2 implementation commits and exact final HEAD;
- R2-01..R2-08 mapping;
- shared-source mechanism and drift proof;
- overflow/fail-closed proof;
- manager progression proof;
- durable replay/conflict proof;
- employment/assignment integrity proof;
- diagnosis-confidence proof;
- trait-vector/non-dominance proof;
- focused repository/SQL/API test results;
- full validation/test/security results;
- final push and PR CI URLs;
- live PG18 status;
- confirmation that recruitment/WP-013 and environment work remain untouched.

Then post one concise issue #26 pointer and one PR #27 receipt and **STOP for Architecture review**.

Do not merge.
Do not start WP-013.

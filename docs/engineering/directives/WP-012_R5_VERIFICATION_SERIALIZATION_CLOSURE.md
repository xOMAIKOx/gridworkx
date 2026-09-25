# GRIDWORKS WP-012 R5 — Verification and Serialization Closure

**Status:** GO — REMEDIATION AUTHORIZED  
**Technical Authority:** GPT-5.6 Sol  
**Execution agent:** Devin  
**Worker model:** GPT-5.6 Luna  
**Worker effort:** XHigh  
**Implementation class:** D1 — bounded exact remediation  
**Reviewed R4 HEAD:** `1e8a48c02806ee3be8f73d381b28bf2ce8b35299`  
**Engineering branch:** `engineering/wp-012-skills-managers`

## 1. Disposition

WP-012 R4 is **REQUEST_CHANGES**.

Continue on the existing Engineering branch from exactly:

`1e8a48c02806ee3be8f73d381b28bf2ce8b35299`

Do **not** rebase onto `main` or this Architecture directive commit.

R4 materially closes the implementation defects. R5 is limited to exact proof, repository lock ordering, missing mutation-path tests, and handback cleanup.

## 2. Accepted / frozen

Preserve the accepted WP-012 implementation and do not redesign it.

Specifically freeze:

- shared JSON progression/manager policy;
- Rust and Go policy consumption;
- player anti-grind behavior;
- player replay/conflict implementation;
- checked manager skill-gain implementation;
- manager progression history `skill_id`;
- receipt immutability trigger;
- DB-level manager serialization guards;
- assignment↔active-employer invariant;
- Gold-vs-Platinum and diagnosis-confidence proofs;
- API/OpenAPI route boundary;
- WP-013 recruitment exclusion;
- existing WP-001–WP-011 gates;
- no environment/host/container work.

## 3. R5-01 — exact generated rarity policy proof

R4 improved the drift test, but exact key-set proof is incomplete for `SharedRarityTraitCapacity`.

Current proof establishes exact canonical membership/cardinality for `SharedRarityPotential`, then checks canonical trait-capacity values by key. It does not independently reject an extra generated trait-capacity key.

### Required

Prove exact key-set/cardinality equality independently for:

- `SharedRarityPotential`;
- `SharedRarityTraitCapacity`;
- `SharedTraitEffects`.

All canonical values must match and no generated extra key may exist.

No runtime policy change is required.

## 4. R5-02 — complete manager arithmetic boundary proof

The R4 code now checks the final manager skill-gain addition. Keep that implementation.

The new test named for skill-gain overflow currently initializes `TotalXP=10` and awards `MaxInt64`, so it fails earlier at the total-XP overflow guard. It does not exercise the newly repaired non-zero-current-skill addition boundary.

### Required

Add separate tests proving:

1. total-XP overflow fails with no mutation;
2. skill-gain addition overflow is reached independently while total XP itself remains valid;
3. valid large/safe gain still caps at potential;
4. error paths preserve total XP, level and targeted skill.

## 5. R5-03 — repository employment-close must acquire the common manager lock before lifecycle check

The DB triggers now serialize direct writes on the manager row. Preserve them.

The trusted repository close path, however, currently:

1. locks the employment row;
2. counts active assignments;
3. only reaches the manager-row lock later through the UPDATE trigger.

That means the trusted precheck itself is not performed under the same manager serialization point used by assignment-open.

### Required

After resolving the employment's manager ID and before checking active assignments:

- acquire `SELECT ... FROM gridworks.managers WHERE manager_id=$1 FOR UPDATE`;
- then query active assignments;
- then close employment if permitted.

Preserve a deterministic lock order and the DB trigger backstop.

Add a repository test that asserts this manager-row lock occurs before the active-assignment query.

## 6. R5-04 — receipt immutability must be part of the structural gate

R4 added the immutable UPDATE/DELETE trigger for `manager_mutation_receipts`.

The WP-012 structural migration check was not updated to require that guard.

### Required

Make `tests/structural/check_wp012_migration.py` fail if the WP-012 receipt immutability function/trigger is absent.

No live PostgreSQL claim is required.

## 7. R5-05 — complete the mutation-path repository test matrix required by R4

R4 added `progression_mutation_test.go`, but it currently contains only three tests:

- player first application + mismatch conflict;
- manager first application;
- employment-close active-assignment rejection.

This does not satisfy the R4 mutation matrix.

Add direct repository/sqlmock tests for the missing cases.

### Player progression

- same semantic request replay succeeds;
- replay returns canonical cumulative XP/proficiency/repetition fields;
- changed semantic payload with same source conflicts.

### Manager progression

- same semantic request replay succeeds and returns total XP/level;
- changed occurrence time or other semantic payload with same source conflicts;
- first application proves skill identity is persisted in event history.

### Employment

- open same-request replay;
- open changed-payload conflict;
- close same-request replay;
- close changed-payload conflict;
- close rejected while active assignment exists;
- close succeeds after assignments are closed;
- manager-row lock ordering is asserted.

### Assignment

- open same-request replay;
- open changed-payload conflict;
- close same-request replay;
- close changed-payload conflict;
- assignment close succeeds while matching employment is active.

Tests must call the actual repository methods and verify transaction SQL/commit/rollback behavior.

## 8. R5-06 — durable handback cleanup and corrected final-receipt contract

R4's durable handback still contains stale/inconsistent metadata:

- both the R1 directive and R4 directive are listed as `Directive`;
- implementation class says D2 although R4 was D1;
- the main verification section still presents old R1 CI URLs as if current;
- older R3 text still says final values are only in pointers without clearly separating historical from current evidence.

Clean the handback so a fresh Architecture session can distinguish historical rounds from the current closure.

### Important correction to R4-06

A commit cannot practically contain its own commit SHA and CI URLs generated only after that commit without creating a self-referential commit loop.

Therefore the controlling closure contract is:

**Durable handback file**
- records the controlling R5 directive path/commit;
- records immediate R5 parent;
- records original WP-012 parent;
- records all R5 implementation commits that precede the handback commit;
- records exact test/evidence results available before handback commit;
- contains no stale current-round directive/class/CI metadata.

**Canonical GitHub issue/PR receipt**
- records the exact final handback HEAD after it is committed;
- records exact push and PR CI URLs for that final HEAD;
- points to the durable handback;
- is the canonical exact-head closure receipt.

This exception is deliberate and supersedes the impossible self-SHA/self-CI wording in R4-06.

## 9. Verification

Run and preserve all existing gates, including:

- `make validate`;
- `make test`;
- `make security`;
- Go progression/repository tests;
- Rust progression/golden/integration tests;
- WP-010 native/Xvfb/runtime parity;
- WP-011 deterministic/controller/visual gates;
- final push CI;
- final PR CI.

R4 exact final-head CI is already verified green:

- push `36117386581`;
- PR `36117390346`.

Live PostgreSQL 18 remains `BLOCKED_DEPENDENCY` if no authorized native PG18 environment exists.

## 10. STOP

After R5:

1. commit the cleaned durable handback;
2. wait for/check final push + PR CI on that exact handback HEAD;
3. post one issue #26 receipt with exact HEAD + exact CI URLs;
4. post one PR #27 receipt pointing to that issue receipt;
5. give the Owner only the short pointer/status;
6. **STOP for Architecture review**.

Do not merge.
Do not start WP-013.
Do not perform host/deployment/provider/database provisioning/port/proxy/systemd/container work.

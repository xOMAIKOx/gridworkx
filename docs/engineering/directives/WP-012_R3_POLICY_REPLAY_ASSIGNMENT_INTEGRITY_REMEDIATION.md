# GRIDWORKS WP-012 R3 — Policy, Replay and Assignment Integrity Remediation

**Status:** GO — REMEDIATION AUTHORIZED  
**Technical Authority:** GPT-5.6 Sol  
**Execution agent:** Devin  
**Worker model:** GPT-5.6 Luna  
**Worker effort:** XHigh  
**Implementation class:** D2 — bounded complex remediation  
**Reviewed R2 HEAD:** `d3309a4cdbd8393142278d7a5d25176190ac06c4`  
**Engineering branch:** `engineering/wp-012-skills-managers`

## 1. Disposition

WP-012 R2 is **REQUEST_CHANGES**.

Continue on the existing Engineering branch from exactly:

`d3309a4cdbd8393142278d7a5d25176190ac06c4`

Do **not** rebase onto `main` or onto this Architecture directive commit.

R2 made meaningful progress and its exact-head push and PR CI are green, but several R2 closure claims are not supported by the implementation or focused tests.

## 2. R2 accepted / frozen

Preserve unless strictly required by the remediation below:

- exact original WP-012 ancestry from `466faa05ad24d0dc1af35a1fa3d71a04e8920961`;
- draft PR #27 and existing Engineering branch;
- WP-012-only public API route set;
- no public XP, manager-create, recruitment, transfer or assignment mutation routes;
- Rust directly embedding the canonical skills/manager JSON;
- veteran-Gold vs fresh-Platinum invariant;
- no direct rarity power multiplier;
- bounded fatigue / no payment or stamina input;
- semantic facility IDs with no fake facility FK;
- Rust WP-004 diagnosis-confidence proof preserving symptom/evidence identity;
- semantic multi-axis trait-vector direction;
- append-close/no-delete history direction;
- WP-013 recruitment boundary;
- existing WP-001–WP-011 behavior and gates;
- no host/deployment/provider/database provisioning/port/proxy/systemd/container work.

## 3. R3-01 — shared progression policy is still duplicated / incompletely proved

R2 adds generated Go content, but the trusted Go player progression path still hard-codes policy:

- trivial multipliers are locally defined as `[10000, 5000, 2500, 0]`;
- meaningful multiplier is locally initialized as `10000`;
- the generated `SharedTrivialMultipliers` and `SharedMeaningfulMultiplier` are therefore not the controlling runtime values.

The Go drift test is also incomplete. It checks registry lengths and some anti-grind values, but does not prove equality for:

- exact player activity-class values;
- player labels / progression-curve / rules-version metadata;
- exact manager skill IDs;
- rarity potential ceilings;
- rarity trait capacities;
- trait IDs and semantic effect vectors.

The JSON schema similarly proves counts, not exact required identities/uniqueness:

- 12 player-skill entries can still contain duplicates or wrong IDs;
- 9 manager-skill entries can still contain duplicates or wrong IDs;
- the four rarity keys are not constrained to Bronze/Silver/Gold/Platinum;
- required trait IDs / uniqueness are not enforced.

### Required remediation

1. Make Go trusted progression consume the generated canonical values directly.
2. Remove duplicate tunable anti-grind values from runtime code.
3. Strengthen drift proof so every generated semantic/tunable value is checked against the canonical JSON.
4. Strengthen schema/validation so the exact required skill/rarity/trait identities and uniqueness are proven.
5. Preserve one authoritative versioned content source.

## 4. R3-02 — trusted Go progression still has fail-closed arithmetic defects

`ApplySkillEvent` has two concrete correctness defects:

1. the configured fourth/later trivial repetition produces multiplier `0`, after which the code evaluates `MaxInt64 / multiplier`, which can panic through integer division by zero;
2. derived proficiency computes `state.CumulativeXP * 10` without checked arithmetic, so a valid non-overflowing cumulative XP near `MaxInt64` can overflow before the cap is applied.

R2 also required existing state version validation, but `state.ProgressionVersion` is not rejected when it differs from the active progression version.

### Required remediation

- zero-multiplier events must deterministically award zero without panic;
- all XP/proficiency arithmetic must be overflow-safe;
- existing non-empty state progression version must match the active version before mutation;
- invalid/version/overflow failures must leave supplied state unchanged;
- add direct Go tests for:
  - repetitions 1, 2, 3, 4 and >4;
  - zero multiplier;
  - near-`MaxInt64` cumulative XP;
  - state-version mismatch;
  - no partial state mutation on error.

## 5. R3-03 — durable replay/conflict semantics remain incomplete

### Player progression

The PostgreSQL player path currently treats any existing `source_event_id` as a successful replay and returns zero award without comparing the replayed semantic payload.

Therefore:

- same source ID + changed player/skill/activity/base XP/class/repetition/time/version is not rejected as a conflict;
- repository replay results omit the canonical cumulative XP, proficiency and repetition state returned by the domain/Rust replay representation.

R2 explicitly required same source/key + changed payload to fail with a stable conflict and replay semantics to remain consistent across layers.

### Manager progression

The manager receipt digest omits `OccurrenceTime`, even though occurrence time is persisted as part of the progression event. Same source + changed occurrence time therefore does not conflict.

A successful manager replay returns only `Replay: true`, losing the original/current `TotalXP` and `Level`.

The persisted `manager_progression_events` row also omits the targeted `skill_id`, so the event history cannot identify which manager skill the trusted operation progressed.

### Required remediation

1. Give player progression durable digest-backed replay/conflict semantics.
2. Same source/key + same semantic request replays without mutation.
3. Same source/key + any changed semantic request field returns stable conflict.
4. Replay results must carry the canonical result state consistently across domain/repository layers.
5. Manager digest must cover all semantic event fields, including occurrence time.
6. Persist the targeted manager skill identity in manager progression history when a skill is progressed.
7. Make WP-012 receipt/event history immutable against UPDATE/DELETE at the DB layer where appropriate.
8. Add direct sqlmock/repository tests for same-request replay and changed-payload conflict for player and manager progression.

## 6. R3-04 — manager progression policy/validation is not fully versioned or fail-closed

Manager progression currently derives:

- level through hard-coded `total_xp / 1000 + 1`;
- skill increase directly from `AwardedXP`.

Those are progression-policy decisions but are not controlled by the canonical versioned content.

The domain also mutates total XP/level before validating the requested manager skill, so an invalid skill can leave a caller-provided state object partially mutated on error. A known manager-skill ID that is absent from the manager's state is silently accepted without progressing a skill.

### Required remediation

- express manager XP-to-level and manager skill-gain rules through the versioned progression policy/content;
- validate manager rarity and targeted skill/state before mutation;
- missing targeted manager skill must fail closed;
- any error must leave the input state unchanged;
- continue enforcing potential ceilings;
- add focused manager-domain tests for invalid rarity, unknown/missing skill, overflow, version mismatch, no-partial-mutation, level thresholds and potential cap.

## 7. R3-05 — assignment/employment lifecycle can become internally inconsistent

The DB assignment employer guard runs on assignment INSERT/UPDATE, but employment close does not prevent a manager with an active facility assignment from losing the active employer.

This creates a bad sequence:

1. manager has active employment and active assignment;
2. employment is closed successfully;
3. assignment remains active while there is no active employer;
4. later assignment close attempts UPDATE the row;
5. the assignment employer guard rejects that UPDATE because no active employer now exists.

The result is an active assignment that can become impossible to close through the accepted repository path.

### Required remediation

Choose and enforce one deterministic lifecycle invariant. Preferred bounded rule:

- employment cannot close while any active manager facility assignment exists;
- close assignment(s) first while active employment still exists;
- then employment may close.

Enforce this both:

- in the trusted repository transaction path; and
- at the PostgreSQL integrity layer so direct SQL cannot create the dangling state.

Preserve the existing assignment-company-matches-active-employer rule.

Add tests proving:

- employment close with active assignment is rejected;
- assignment close while employment active succeeds;
- employment close after assignments are closed succeeds;
- closed assignment cannot reopen/rewrite;
- no active assignment can exist without matching active employment.

## 8. R3-06 — R2 verification/handback claims must match actual proof

R2-08 required focused repository/SQL mutation tests. The R2 remediation commit does not add or modify PostgreSQL mutation tests for:

- player changed-payload replay conflict;
- manager progression replay/conflict;
- employment replay/conflict;
- assignment replay/conflict;
- employment/assignment lifecycle integrity.

The durable handback also remains stale/incomplete:

- it names the R1 directive path while citing the R2 directive commit;
- it does not record exact R2 final HEAD;
- it does not provide the exact R2 implementation commit list;
- its main verification section still lists the earlier R1 push/PR CI URLs rather than the R2 exact-head runs.

### Required remediation

Add actual repository-path tests for the repaired mutations and update the durable handback so it is self-contained and exact.

The R3 handback must include:

- exact R3 parent and final HEAD;
- exact R3 implementation commit(s);
- R3-01..R3-06 mapping;
- exact shared-source/drift proof;
- zero-multiplier and overflow tests;
- player/manager replay + changed-payload conflict tests;
- manager progression versioning/fail-closed tests;
- employment/assignment lifecycle tests;
- full repository validation/test/security results;
- final exact-head push and PR CI URLs;
- live PG18 status;
- confirmation that WP-013 and environment work remain untouched.

## 9. CI / environment

R2 exact-head push run `35962404282` and PR run `35962406694` are green.

Preserve all existing gates.

Live PostgreSQL 18 remains an honest `BLOCKED_DEPENDENCY` if no authorized native PG18 environment exists. Do not fabricate live evidence and do not weaken repository-local gates.

## 10. STOP

After R3 implementation:

1. update/commit the durable WP-012 handback;
2. post one concise issue #26 pointer;
3. post one concise PR #27 receipt;
4. give the Owner only the short pointer/status;
5. **STOP for Architecture review**.

Do not merge.
Do not start WP-013.
Do not perform host/deployment/provider/DB provisioning/port/proxy/systemd/container work.

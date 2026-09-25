# GRIDWORKS WP-012 R4 — Final Arithmetic, Receipt and Lifecycle Integrity Remediation

**Status:** GO — REMEDIATION AUTHORIZED  
**Technical Authority:** GPT-5.6 Sol  
**Execution agent:** Devin  
**Worker model:** GPT-5.6 Luna  
**Worker effort:** XHigh  
**Implementation class:** D1 — bounded exact remediation  
**Reviewed R3 HEAD:** `9cf624e1a3982eb0f2eeb63a4f60130a136f34b4`  
**Engineering branch:** `engineering/wp-012-skills-managers`

## 1. Disposition

WP-012 R3 is **REQUEST_CHANGES**.

Continue on the existing Engineering branch from exactly:

`9cf624e1a3982eb0f2eeb63a4f60130a136f34b4`

Do **not** rebase onto `main` or onto this Architecture directive commit.

R3 closes most of the R2/R3 design defects and its exact-head push/PR CI are green. This R4 is a final bounded integrity/test/handback pass only.

## 2. R3 accepted / frozen

Preserve unless strictly required by the remediation below:

- original WP-012 ancestry from `466faa05ad24d0dc1af35a1fa3d71a04e8920961`;
- PR #27 and existing Engineering branch;
- WP-012-only public route set and no public mutation expansion;
- canonical shared JSON as the progression/manager policy source;
- Go consumption of generated anti-grind and manager progression values;
- Rust consumption of the same manager progression policy;
- exact skill/manager/rarity/trait content validator direction;
- player zero-multiplier handling;
- player progression version validation and checked cumulative/proficiency arithmetic direction;
- request digests covering player and manager semantic event fields including occurrence time;
- manager progression history carrying `skill_id`;
- diagnosis-confidence / unchanged symptom-evidence proof;
- veteran-Gold vs fresh-Platinum invariant;
- semantic trait vectors;
- append-close history direction;
- assignment company ↔ active employer rule;
- no WP-013 recruitment scope;
- no host/deployment/provider/database provisioning/port/proxy/systemd/container work.

## 3. R4-01 — manager skill progression still has an arithmetic overflow hole

The Go manager progression path currently checks:

- total-XP addition; and
- `AwardedXP * SkillBPSPerXP`.

It then performs:

`target.ProficiencyBPS + int(event.AwardedXP * SharedManagerProgression.SkillBPSPerXP)`

without checking that final addition/conversion is safe.

At the accepted current policy (`skill_bps_per_xp = 1`), an event near `MaxInt64` with a non-zero existing proficiency can overflow the Go `int` addition and produce a negative/invalid skill value rather than returning `ErrOverflow`.

### Required remediation

- use checked arithmetic for the complete manager skill-gain calculation, including conversion and addition to current proficiency;
- cap only after a valid non-overflowing calculation;
- any arithmetic/version/rarity/skill error must leave the supplied manager state unchanged;
- add focused boundary tests including:
  - near-`MaxInt64` awarded XP;
  - non-zero current skill + extreme award;
  - total-XP overflow;
  - skill-gain overflow;
  - valid potential-cap behavior after safe arithmetic.

## 4. R4-02 — generated-content drift proof must be exact in both directions

R3's drift test validates canonical rarity/trait entries against generated values, but it does not reject extra generated rarity/trait keys.

The controlling requirement is that generated runtime content is an exact mechanical representation of the canonical JSON, not merely a superset containing all canonical entries.

### Required remediation

- prove exact key-set/cardinality equality for generated rarity potential, rarity trait capacity and trait effect maps;
- retain exact player/manager registries, metadata and manager-progression policy checks;
- no generated semantic/tunable value may exist without a canonical JSON source entry.

## 5. R4-03 — durable mutation receipts are still mutable

`gridworks.manager_mutation_receipts` is the durable replay/conflict authority but currently has no database guard preventing UPDATE or DELETE.

The player and manager progression event histories are append-only, but the receipt row that controls manager/employment/assignment replay can still be rewritten or deleted through direct SQL.

### Required remediation

- make WP-012 durable mutation receipts immutable after insert (UPDATE/DELETE forbidden);
- preserve deterministic replay/conflict semantics;
- extend structural/integrity proof to assert the receipt immutability guard;
- do not introduce destructive migration behavior.

## 6. R4-04 — employment/assignment invariant must be concurrency-safe

R3 added a database trigger blocking employment close while an active assignment exists. That closes the simple sequential orphaning case.

However, the trusted operations do not serialize employment-close against assignment-open using the same manager lock:

- assignment-open locks the manager row;
- employment-close currently does not;
- the DB triggers query the opposite table but do not acquire a common manager lock.

A concurrent sequence can therefore race between "close active employment" and "open assignment under the still-visible active employment" and violate the intended invariant.

The R3 directive also required the lifecycle rule to be enforced in both the trusted repository transaction path and PostgreSQL integrity layer. The repository close path currently relies on the DB trigger rather than performing the explicit trusted-domain check.

### Required remediation

Use one deterministic lock order centered on the manager identity.

At minimum:

1. employment-close resolves/locks the manager row before lifecycle checks;
2. assignment-open uses the same manager lock order;
3. employment-close explicitly checks for active assignments and returns a stable domain/API conflict before attempting the close;
4. PostgreSQL direct-write integrity uses the same manager serialization point (for example, a `FOR UPDATE` manager lock in the relevant guard functions) so direct SQL cannot race the invariant;
5. preserve assignment-close while employment is still active, then employment-close;
6. no deadlock-prone inverse lock order.

Add focused tests proving the repository lock/check ordering and stable rejection semantics. Live PG18 may remain `BLOCKED_DEPENDENCY`; do not fabricate a live concurrency test if the authorized environment is unavailable.

## 7. R4-05 — required WP-012 mutation-path tests are still absent

The R3 directive required direct repository tests for repaired WP-012 mutation paths.

At R3 HEAD, the PostgreSQL WP-012 test file still covers read paths only, while the larger `remediation_test.go` covers earlier identity/company remediation rather than WP-012 progression/employment/assignment mutations.

Add direct sqlmock/accepted repository tests for at least:

### Player progression
- first application;
- same semantic request replay;
- changed-payload same-source conflict;
- canonical replay result fields.

### Manager progression
- first application with manager + skill row locks;
- same semantic request replay;
- changed occurrence/payload conflict;
- replay total XP/level;
- skill-id history persistence.

### Employment / assignment
- open replay and changed-payload conflict;
- close replay and changed-payload conflict;
- employment close rejected while active assignment exists;
- assignment close succeeds while employment is active;
- employment close succeeds after assignment closure;
- expected manager lock ordering is exercised.

These tests must exercise the actual repository methods and transaction SQL, not only structural migration tokens.

## 8. R4-06 — durable handback must be self-contained and exact

The R3 durable handback is still not self-contained:

- its top-level `Directive` path still names the R1 directive while citing the R3 directive commit;
- its top-level exact parent is the original WP-012 parent rather than the immediate R3 Engineering parent;
- it records `f0dcc01...` as the pre-handback head rather than the final reviewed R3 head;
- its main verification section still cites the old R1 push/PR CI runs instead of the R3 exact-head runs;
- final exact-head/CI evidence is delegated to GitHub pointers rather than recorded durably in the handback.

### Required remediation

Update the durable handback so it explicitly records:

- controlling R4 directive path and commit;
- immediate R4 parent `9cf624e1a3982eb0f2eeb63a4f60130a136f34b4`;
- original WP-012 ancestry separately;
- exact R4 implementation commit list;
- exact final R4 HEAD;
- exact push CI and PR CI URLs for that final HEAD;
- R4-01..R4-06 proof mapping;
- live PG18 status;
- no WP-013 / no environment mutation statement.

The file itself must be sufficient for a fresh Architecture session without relying on a placeholder saying the exact values are only in issue comments.

## 9. CI / environment

R3 exact-head CI is verified green:

- push: `36113286323`;
- PR: `36113293980`.

Preserve all existing WP-001–WP-011 and WP-012 gates.

Live PostgreSQL 18 may remain `BLOCKED_DEPENDENCY` if no authorized native PG18 environment exists.

## 10. STOP

After R4:

1. commit the exact durable handback;
2. post one concise issue #26 handback pointer;
3. post one concise PR #27 receipt;
4. give the Owner only the short status/pointer;
5. **STOP for Architecture review**.

Do not merge.
Do not start WP-013.
Do not perform host/deployment/provider/database provisioning/port/proxy/systemd/container work.

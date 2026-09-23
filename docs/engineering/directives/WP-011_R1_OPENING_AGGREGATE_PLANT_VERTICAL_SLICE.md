# GRIDWORKS WP-011 R1 — Opening Aggregate-Plant Vertical Slice Engineering Directive

**Status:** GO  
**Technical Authority:** GPT-5.6 Sol  
**Execution agent:** Devin  
**Worker model:** GPT-5.6 Luna  
**Worker effort:** XHigh  
**Implementation class:** D2 — bounded complex  
**Exact Engineering parent:** `ca5b4f06cdba226176296dd1ee7d294ef6c7927a`

## 1. Objective

Implement the complete WP-011 opening aggregate-plant vertical slice defined by:

`docs/work-orders/WP-011_OPENING_AGGREGATE_PLANT_VERTICAL_SLICE.md`

This is the first playable GRIDWORKS recovery loop.

## 2. Repository / branch

Repository:

`xOMAIKOx/gridworkx`

Create/use Engineering branch:

`engineering/wp-011-opening-aggregate-slice`

Branch from exactly:

`ca5b4f06cdba226176296dd1ee7d294ef6c7927a`

The work-order/directive commit on `main` is not the Engineering parent.

Open one draft PR against `main` after the first implementation commit.

No rebase unless Architecture explicitly authorizes it.

## 3. Execution route

Architecture classifies WP-011 as **D2 bounded complex**:

- architecture is fixed;
- existing simulation/bridge primitives are accepted;
- implementation spans Rust + GDExtension + Godot UI + tests;
- ordinary bounded debugging is expected;
- no open architectural discovery currently justifies Luna Max.

Execute as:

**Devin + GPT-5.6 Luna XHigh**

If the implementation reveals genuine architecture ambiguity rather than ordinary debugging, STOP and return to Architecture. Do not silently escalate worker model.

## 4. Read before editing

Mandatory:

1. `AGENTS.md`
2. `AI_AGENT_COLLABORATION_PROTOCOL.md`
3. `docs/GITHUB_COLLABORATION_TRANSPORT.md`
4. WP-011 controlling work order
5. Product Philosophy first-session onboarding
6. DP2 core loop
7. WP-003/004/005 work orders and evidence
8. WP-010 work order/evidence
9. active WP-011 issue and draft PR once created

## 5. Architecture invariants

Do not violate:

- Rust `gridworks-sim` remains canonical simulation authority.
- Godot is presentation/control only.
- No local duplicate capacity, failure, production, RNG or inventory math.
- Scenario canonical state is built in Rust.
- UI wrong choices remain recoverable.
- No hidden answer flag may substitute for canonical evidence/state.
- No fake money/revenue/XP/manager/contracts.
- No network/database/persistence/auth work.
- Existing WP-010 bridge error/version/panic contracts remain intact.
- Rust 1.80.0 / Godot 4.7.2 / exact godot 0.2.4 remain pinned.

## 6. Required implementation sequence

### R1.1 canonical scenario

Implement `scenario.opening.aggregate` in Rust with:

- aggregate facility;
- aggregate fault definitions;
- aggregate material fixture;
- critical feed-conveyor motor seizure;
- non-critical output-conveyor worn belt;
- accepted graph-level bottleneck yielding about 3% capacity after critical repair.

Pin exact expected state/digests by tests.

### R1.2 scenario bridge

Add one coarse scenario creation bridge operation.

Do not add tutorial-specific simulation calculations to GDScript.

### R1.3 opening controller + UI

Build the playable Godot scene/controller:

- plant schematic;
- component details;
- evidence/diagnosis;
- interventions;
- throughput;
- input/output inventory;
- advisor;
- completion feedback.

### R1.4 canonical actions

Wire UI through canonical command batches for:

- inspect;
- diagnose;
- intervention;
- production.

### R1.5 wrong-choice recovery

Prove repairing the visible worn belt first leaves throughput at zero and the scenario remains recoverable.

### R1.6 critical recovery

Prove repairing the feed-conveyor seizure restores the pinned small partial capacity.

### R1.7 first output

Execute canonical aggregate production and prove finished aggregate enters canonical inventory.

### R1.8 deterministic restart/golden

Commit and assert a shared WP-011 deterministic golden across pure Rust and real Godot/GDExtension paths.

## 7. Acceptance criteria

- **AC-011-01:** exact parent is `ca5b4f06cdba226176296dd1ee7d294ef6c7927a`.
- **AC-011-02:** scenario state is created only by canonical Rust.
- **AC-011-03:** initial facility effective capacity is 0.
- **AC-011-04:** critical motor seizure and visible non-critical worn belt are present.
- **AC-011-05:** canonical inspect evidence is displayed; UI does not fabricate it.
- **AC-011-06:** canonical diagnosis action is executable.
- **AC-011-07:** non-critical repair first leaves facility at 0 and remains recoverable.
- **AC-011-08:** critical repair resolves the seizure and yields the pinned ~3% capacity.
- **AC-011-09:** production uses `recipe.aggregate_crush` and accepted material inventories.
- **AC-011-10:** first production is constrained by canonical facility capacity and creates canonical finished aggregate.
- **AC-011-11:** completion derives from canonical state/output.
- **AC-011-12:** no GDScript simulation math or fake fallback exists.
- **AC-011-13:** advisor explains reasoning without directly mutating outcomes.
- **AC-011-14:** portrait/mobile and wide layouts are usable.
- **AC-011-15:** bridge/unavailable/error state is explicit and non-authoritative.
- **AC-011-16:** deterministic restart returns to exact initial digest.
- **AC-011-17:** pure Rust and Godot assert one committed WP-011 golden.
- **AC-011-18:** existing WP-010 native discovery/parity remains green.
- **AC-011-19:** no fake economy/skills/managers/contracts/later-WP implementation.
- **AC-011-20:** evidence and durable handback are complete.
- **AC-011-21:** no host/deployment/provider/DB/port/container mutation occurred.

## 8. Mandatory verification

Run at minimum:

- `make validate`
- `make test`
- `make security`
- Rust workspace tests/format/clippy as current repository gates require
- existing WP-010 Godot extension build/discovery/parity
- new WP-011 Rust deterministic tests
- new WP-011 Godot headless controller/vertical-slice test
- structural checks
- visual evidence generation/validation if implemented by repository tooling

Final push and PR CI must be green.

## 9. Mandatory STOP conditions

STOP and hand back `BLOCKED` if:

- parent/baseline is not exact;
- success requires changing accepted simulation architecture;
- scenario requires new economics/contracts/skills/managers;
- success requires fake client-side simulation;
- accepted bridge semantics must be weakened;
- success requires changing Rust/Godot/godot-rust pins;
- host/deployment mutation would be required;
- a legitimate test/gate would need disabling;
- worker substitution/escalation would be required;
- repository state materially differs from the work order.

A BLOCKED handback is correct. Do not improvise around the boundary.

## 10. Durable Engineering handback

Commit:

`docs/engineering/handbacks/wp-011/WP011_R1_HANDBACK.md`

Required content:

- Work package: WP-011 R1
- Technical Authority directive path
- execution agent: Devin
- actual worker model/effort
- implementation class D2
- repository/branch
- exact parent
- implementation commits
- final HEAD
- changed files/components
- canonical scenario ID/seed
- initial/repair/final digests
- effective-capacity proof
- wrong-choice recovery proof
- first-output inventory proof
- Rust golden result
- Godot headless result
- visual evidence paths
- full verification results
- CI push/PR URLs
- AC-011-01..21 mapping
- deviations/assumptions
- blockers/risks
- worker-routing deviations
- unauthorized environment mutation status
- final state: READY FOR SOL/PRO REVIEW / BLOCKED / STOP

## 11. GitHub return path

After committing the handback:

1. post one concise WP-011 issue handback pointer;
2. post one concise draft-PR receipt pointing to that issue handback;
3. give the Owner only a concise status/pointer;
4. STOP.

Do not paste the full technical handback into chat.
Do not ask the Owner to relay it.
Do not merge.
Do not start WP-012.

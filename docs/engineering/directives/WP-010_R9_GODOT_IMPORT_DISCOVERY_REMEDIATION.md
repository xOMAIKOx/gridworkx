# GRIDWORKS WP-010 R9 — Godot Import Discovery Remediation Directive

**Status:** GO — exact remediation authorized  
**Published:** 2026-09-23  
**Owner authorization:** WP-010 repository implementation already authorized  
**Architecture state:** WP-010 ACTIVE; R1–R8 ACCEPTED; R9 ONLY  
**Collaboration ledger:** issue #22  
**Implementation surface:** existing branch `engineering/wp-010-godot-rust-extension`, draft PR #23

## 1. Authority / execution route

- **Technical Authority:** GPT-5.6 Sol
- **Execution agent:** Devin
- **Worker model:** GPT-5.6 Luna
- **Worker effort:** XHigh
- **Implementation class:** D1 — deterministic exact remediation

Devin is the Engineering/orchestration agent. Luna executes inside this fixed remediation contract.

**No worker-model/effort substitution or escalation is authorized.** If Luna XHigh is unavailable, STOP and hand back the routing blocker.

## 2. Controlling repository state

Repository:

`xOMAIKOx/gridworkx`

Existing Engineering branch:

`engineering/wp-010-godot-rust-extension`

Original required WP-010 parent:

`03c0ad7971804897a74489f8f109756f98f6e219`

Reviewed R9 blocker head:

`e48749c5567d266210888deacfcb6dac452f875c`

Current draft PR #23 head when this directive was issued:

`cf4ab864f6688da81a31d06764858e595d25bfd3`

Do not rebase, create a new branch or rewrite reviewed history.

Before editing, verify PR #23 is still on the same branch and that no later uncontrolled implementation commit has superseded this directive. If materially different, STOP and post a baseline-mismatch handback.

## 3. Controlling references

Read before mutation:

1. `AGENTS.md`
2. `AI_AGENT_COLLABORATION_PROTOCOL.md`
3. `docs/GITHUB_COLLABORATION_TRANSPORT.md`
4. `docs/work-orders/WP-010_GODOT_RUST_GDEXTENSION_INTEGRATION.md`
5. `docs/evidence/WP-010_GODOT_RUST_GDEXTENSION.md`
6. issue #22 Architecture R1–R8 acceptance / R9 review
7. issue #22 Architecture root-cause decision replacing `--quit-after` with `--import`
8. draft PR #23 current diff and CI

R1–R8 are accepted and frozen.

## 4. Exact defect

The current discovery gate uses an iteration-count forced editor exit:

```sh
--editor --quit-after 2
```

Godot 4.7.2 defines `--quit-after` as quitting after the specified number of iterations.

The failed CI reaches:

```text
[ DONE ] first_scan_filesystem
WARNING: Scan thread aborted...
Aborted (core dumped)
```

Godot 4.7.2 source shows the warning is emitted during editor teardown when an EditorFileSystem scan thread remains active.

Godot's purpose-built CLI mode for this use case is:

```text
--import
```

which starts the editor, waits for resource import/filesystem work, and then exits.

This remediation changes only the discovery lifecycle/gate.

## 5. Authorized scope

Engineering may change only what is required to replace the R9 discovery lifecycle and update its evidence:

- `.github/workflows/ci.yml`
- `tests/structural/check_wp010_integration.py`
- `docs/evidence/WP-010_GODOT_RUST_GDEXTENSION.md`
- the required R9 Engineering handback document

No production bridge code change is authorized.

## 6. Prohibited / deferred scope

Do not change:

- Rust `1.80.0`
- Godot `4.7.2`
- `godot = "=0.2.4"`
- `crates/gridworks-godot/**`
- `crates/gridworks-sim/**`
- GDExtension descriptor/staging layout
- accepted R1–R8 contracts/tests/goldens
- simulation logic
- API/network/persistence behavior
- WP-011 or later work

No host/ERIS/deployment/provider/database/port/reverse-proxy/systemd/container-runtime work is authorized.

## 7. Required implementation

Replace the discovery step with the equivalent of:

```sh
timeout 60s "$HOME/.local/bin/godot" \
  --headless \
  --path apps/game \
  --import \
  --rendering-method gl_compatibility \
  --rendering-driver opengl3 \
  --audio-driver Dummy
```

Requirements:

1. Remove `--editor` from the R9 discovery command.
2. Remove `--quit-after`.
3. Do not whitelist exit 134 or any non-zero code.
4. Require discovery/import exit status 0.
5. Do not manually create `.godot/extension_list.cfg`.
6. After exit 0, assert:
   - `apps/game/.godot/extension_list.cfg` exists;
   - exactly one nonblank entry exists;
   - the entry is `res://native/gridworks_sim.gdextension`.
7. Keep the already accepted fresh normal-runtime project parse unchanged.
8. Keep the already accepted deterministic WP-010 runtime parity scene unchanged.
9. Update the structural check to require the `--import` discovery mode and reject the old `--editor --quit-after` pattern.
10. Update WP-010 evidence to state only results actually observed.

## 8. Acceptance criteria

- **AC-R9-01:** CI discovery command uses Godot `--import`; no `--quit-after`.
- **AC-R9-02:** no non-zero discovery exit is whitelisted.
- **AC-R9-03:** discovery/import exits 0 in both push and PR CI.
- **AC-R9-04:** generated extension-list exists.
- **AC-R9-05:** extension-list contains exactly `res://native/gridworks_sim.gdextension` and no other nonblank descriptor line.
- **AC-R9-06:** normal fresh-runtime project parse remains PASS.
- **AC-R9-07:** WP-010 deterministic Godot parity remains PASS.
- **AC-R9-08:** foundation/repository/security gates remain PASS.
- **AC-R9-09:** no production bridge/R1–R8/version changes occur.
- **AC-R9-10:** evidence document records the exact command, clean exit 0, extension-list proof and final CI.
- **AC-R9-11:** no unauthorized environment/host/deployment mutation occurs.
- **AC-R9-12:** complete committed Engineering handback exists and Engineering stops.

## 9. Mandatory verification

Run the repository-established gates applicable to the branch, including at minimum:

- `make validate`
- `make test`
- `make security`
- repository formatting/static checks as currently defined
- the Godot extension build/stage path
- the `--import` discovery gate
- generated extension-list assertions
- fresh normal-runtime Godot parse
- WP-010 deterministic parity scene

Final GitHub Actions push and pull-request workflows must both be green.

Do not weaken or skip an existing legitimate gate to obtain green CI.

## 10. STOP conditions

STOP immediately and hand back `BLOCKED` if any of the following occurs:

1. current PR branch/HEAD materially differs from the controlling state;
2. Luna XHigh is unavailable and substitution would be required;
3. Godot `--import` generates the correct extension list but exits non-zero/aborts;
4. success would require another arbitrary Godot CLI combination;
5. success would require changing Rust/Godot/godot-rust versions;
6. success would require bridge/simulation changes;
7. success would require weakening/skipping a legitimate gate;
8. success would require scope expansion or host/deployment mutation.

If STOP condition 3 occurs, do **not** continue experimenting. Capture the full import log, exact exit code and extension-list contents, commit the handback, post concise GitHub pointers, and stop for Architecture.

## 11. Required durable Engineering handback

Commit exactly one canonical handback for this pass at:

`docs/engineering/handbacks/wp-010/WP010_R9_HANDBACK.md`

It must contain:

```text
ENGINEERING HANDBACK

Work package: WP-010 R9
Technical Authority directive: docs/engineering/directives/WP-010_R9_GODOT_IMPORT_DISCOVERY_REMEDIATION.md
Execution agent: Devin
Worker model: GPT-5.6 Luna
Worker effort: XHigh
Implementation class: D1

Repository / branch:
Controlling original parent:
Remediation starting HEAD:
Implementation commit(s):
Final HEAD:

Authorized scope completed:
- ...

Files/components changed:
- ...

Verification:
- make validate: PASS/FAIL
- make test: PASS/FAIL
- make security: PASS/FAIL
- Godot import discovery: PASS/FAIL + exact command/exit
- extension-list assertion: PASS/FAIL
- fresh runtime parse: PASS/FAIL
- deterministic parity: PASS/FAIL
- push CI: PASS/FAIL + URL
- PR CI: PASS/FAIL + URL

Acceptance criteria:
- AC-R9-01 ... AC-R9-12: PASS/FAIL + evidence

Deviations from directive:
- NONE / ...

Open blockers/risks:
- NONE / ...

Worker-routing deviations:
- NONE / ...

Unauthorized environment/production mutations:
- NONE

Final state:
READY FOR SOL/PRO REVIEW / BLOCKED / STOP
```

The handback must accurately state the actual worker model/effort used.

## 12. GitHub return path

After committing the handback:

1. post one concise issue #22 handback pointer containing:
   - exact final HEAD;
   - handback path;
   - worker model/effort;
   - state: `READY FOR SOL/PRO REVIEW`, `BLOCKED` or `STOP`;
   - no-work-beyond-boundary statement;
2. post one concise PR #23 receipt pointing to the issue handback;
3. give the Owner only a concise status/pointer;
4. **STOP**.

Do not paste the full handback into chat.
Do not ask the Owner to relay it to Architecture.
Do not merge PR #23.
Do not start WP-011.

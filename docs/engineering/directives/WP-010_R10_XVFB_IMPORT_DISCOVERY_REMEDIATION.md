# GRIDWORKS WP-010 R10 — Xvfb Godot Import Discovery Remediation

**Status:** GO — exact remediation authorized  
**Published:** 2026-09-23  
**Owner authorization:** WP-010 repository implementation already authorized  
**Architecture state:** WP-010 ACTIVE; R1–R8 ACCEPTED; R9 BLOCKED/STOP ACCEPTED; R10 ONLY  
**Collaboration ledger:** issue #22  
**Implementation surface:** existing branch `engineering/wp-010-godot-rust-extension`, draft PR #23

## 1. Authority / execution route

- **Technical Authority:** GPT-5.6 Sol
- **Execution agent:** Devin
- **Worker model:** GPT-5.6 Luna
- **Worker effort:** XHigh
- **Implementation class:** D1 — deterministic exact remediation

No worker-model/effort substitution or escalation is authorized. If Luna XHigh is unavailable, STOP and hand back the routing blocker.

## 2. Controlling repository state

Repository:

`xOMAIKOx/gridworkx`

Existing Engineering branch:

`engineering/wp-010-godot-rust-extension`

Original WP-010 parent:

`03c0ad7971804897a74489f8f109756f98f6e219`

Accepted R9 starting HEAD:

`cf4ab864f6688da81a31d06764858e595d25bfd3`

R9 BLOCKED/STOP exact HEAD:

`98e56649d257a402f92eeaab824a7467b420db36`

Required R9 durable handback:

`docs/engineering/handbacks/wp-010/WP010_R9_HANDBACK.md`

Do not rebase, create a new branch or rewrite reviewed history.

Before editing, verify PR #23 still points to exact HEAD `98e56649d257a402f92eeaab824a7467b420db36`. If it has moved materially, STOP and post a baseline-mismatch handback.

## 3. Controlling references

Read before mutation:

1. `AGENTS.md`
2. `AI_AGENT_COLLABORATION_PROTOCOL.md`
3. `docs/GITHUB_COLLABORATION_TRANSPORT.md`
4. `docs/work-orders/WP-010_GODOT_RUST_GDEXTENSION_INTEGRATION.md`
5. `docs/evidence/WP-010_GODOT_RUST_GDEXTENSION.md`
6. `docs/engineering/directives/WP-010_R9_GODOT_IMPORT_DISCOVERY_REMEDIATION.md`
7. R9 durable handback above
8. Godot upstream issue #111645
9. Godot upstream issue #123712
10. draft PR #23 current diff and CI

R1–R8 remain accepted and frozen.

## 4. Architecture diagnosis

R9 correctly executed the strict Godot 4.7.2 headless import gate:

```sh
godot --headless --path apps/game --import ...
```

and stopped when the process exited 134.

The final R9 CI log reached:

```text
[ DONE ] first_scan_filesystem
Initialize godot-rust ...
[ DONE ] loading_editor_layout
Aborted (core dumped)
```

Foundation/repository/security gates passed.

This failure now matches documented upstream Godot behavior closely enough to justify one bounded workaround test:

- Godot issue **#111645** is open, confirmed, tagged `topic:gdextension` + `crash`, and reproduces with:
  `rm -rf .godot; godot --headless --import`
  on projects containing GDExtensions.
- The issue discussion reports the workaround:
  **use `xvfb-run godot ...` instead of `--headless`**.
- Godot issue **#123712** separately confirms Godot **4.7.2** can complete work successfully and then crash during teardown/exit with GDExtension/addon involvement.

Architecture therefore authorizes testing the upstream-documented Xvfb import path.

This is not permission to whitelist exit 134 or treat a crash as success.

## 5. Authorized scope

Engineering may change only what is required for the R10 CI discovery workaround and its evidence:

- `.github/workflows/ci.yml`
- `tests/structural/check_wp010_integration.py`
- `docs/evidence/WP-010_GODOT_RUST_GDEXTENSION.md`
- `docs/engineering/handbacks/wp-010/WP010_R10_HANDBACK.md`

A CI-only `xvfb` package install is authorized if `xvfb-run` is not already available on the GitHub runner.

No production bridge or simulation code change is authorized.

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

### 7.1 Xvfb availability

In the Godot CI job:

1. check whether `xvfb-run` is available;
2. if not, install only the CI dependency needed to provide it, e.g. `xvfb`;
3. do not install unrelated desktop/runtime packages.

### 7.2 Discovery command

Replace the headless import discovery invocation with the equivalent of:

```sh
timeout 60s xvfb-run -a "$HOME/.local/bin/godot" \
  --path apps/game \
  --import \
  --rendering-method gl_compatibility \
  --rendering-driver opengl3 \
  --audio-driver Dummy
```

Requirements:

1. **Do not pass `--headless`** to the import/discovery process.
2. Keep `--import`.
3. Keep the CI-safe compatibility renderer unless exact evidence requires otherwise.
4. Do not whitelist exit 134 or any non-zero status.
5. Require Xvfb import/discovery exit status **0**.
6. Do not manually create `.godot/extension_list.cfg`.

### 7.3 Discovery proof

After clean exit 0, assert:

- `apps/game/.godot/extension_list.cfg` exists;
- it contains exactly one nonblank descriptor entry;
- that entry is:
  `res://native/gridworks_sim.gdextension`.

### 7.4 Existing runtime proof

After discovery succeeds:

- run the already-accepted fresh normal-runtime project parse unchanged;
- run the already-accepted WP-010 deterministic runtime parity scene unchanged;
- both must pass.

Do not convert either runtime proof to Xvfb merely to mask a failure unless Architecture separately authorizes it.

## 8. Acceptance criteria

- **AC-R10-01:** R9 BLOCKED/STOP handback remains intact and unmodified except for later cross-reference if needed.
- **AC-R10-02:** import/discovery uses `xvfb-run`.
- **AC-R10-03:** import/discovery does not use `--headless`.
- **AC-R10-04:** no non-zero discovery exit is whitelisted.
- **AC-R10-05:** Xvfb import/discovery exits 0 in both push and PR CI.
- **AC-R10-06:** generated extension-list exists.
- **AC-R10-07:** extension-list contains exactly `res://native/gridworks_sim.gdextension`.
- **AC-R10-08:** fresh normal-runtime parse remains PASS.
- **AC-R10-09:** WP-010 deterministic runtime parity remains PASS.
- **AC-R10-10:** foundation/repository/security gates remain PASS.
- **AC-R10-11:** no production bridge/R1–R8/version changes occur.
- **AC-R10-12:** evidence records upstream bug references, exact Xvfb command, clean exit 0, extension-list proof and final CI.
- **AC-R10-13:** no unauthorized environment/host/deployment mutation occurs.
- **AC-R10-14:** durable R10 Engineering handback exists and Engineering stops.

## 9. Mandatory verification

Run the repository-established gates applicable to the branch, including at minimum:

- `make validate`
- `make test`
- `make security`
- repository formatting/static checks as currently defined
- Godot extension build/stage path
- Xvfb import/discovery gate
- generated extension-list assertions
- fresh normal-runtime Godot parse
- WP-010 deterministic parity scene

Final GitHub Actions push and pull-request workflows must both be green.

Do not weaken or skip an existing legitimate gate.

## 10. STOP conditions

STOP immediately and hand back `BLOCKED` if:

1. PR #23 HEAD is not the expected R9 blocked head before mutation;
2. Luna XHigh is unavailable and substitution would be required;
3. `xvfb-run` import/discovery still exits non-zero/aborts;
4. success would require another arbitrary Godot CLI experiment;
5. success would require changing Rust/Godot/godot-rust versions;
6. success would require bridge/simulation changes;
7. success would require whitelisting exit 134 or another crash/non-zero status;
8. success would require scope expansion or host/deployment mutation.

If STOP condition 3 occurs, do not continue experimenting. Capture:
- complete Xvfb import log;
- exact exit code;
- whether `.godot/extension_list.cfg` exists after the process;
- exact extension-list contents if present;
- whether the staged `.so` existed before import;
- exact final HEAD and CI URLs.

Then commit the handback and stop.

## 11. Required durable Engineering handback

Commit exactly one canonical handback for this pass at:

`docs/engineering/handbacks/wp-010/WP010_R10_HANDBACK.md`

It must contain:

```text
ENGINEERING HANDBACK

Work package: WP-010 R10
Technical Authority directive: docs/engineering/directives/WP-010_R10_XVFB_IMPORT_DISCOVERY_REMEDIATION.md
Execution agent: Devin
Worker model: GPT-5.6 Luna
Worker effort: XHigh
Implementation class: D1

Repository / branch:
Controlling original parent:
R10 starting HEAD:
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
- xvfb import discovery: PASS/FAIL + exact command/exit
- extension-list assertion: PASS/FAIL
- fresh runtime parse: PASS/FAIL
- deterministic parity: PASS/FAIL
- push CI: PASS/FAIL + URL
- PR CI: PASS/FAIL + URL

Acceptance criteria:
- AC-R10-01 ... AC-R10-14: PASS/FAIL + evidence

Upstream references:
- godotengine/godot#111645
- godotengine/godot#123712

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

## 12. GitHub return path

After committing the handback:

1. post one concise issue #22 handback pointer with:
   - exact final HEAD;
   - handback path;
   - worker model/effort;
   - final state;
   - no-work-beyond-boundary statement;
2. post one concise PR #23 receipt pointing to issue #22;
3. give the Owner only a concise status/pointer;
4. STOP.

Do not paste the full handback into chat.
Do not ask the Owner to relay it.
Do not merge PR #23.
Do not start WP-011.

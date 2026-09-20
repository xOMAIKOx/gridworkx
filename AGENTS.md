# GRIDWORKS Agent Operating Contract

This repository uses GitHub as the authoritative collaboration channel between Architecture and Engineering.

## Roles

### Architecture
Owns:
- product/architecture decisions;
- ADRs;
- work orders;
- acceptance criteria;
- review outcomes;
- REQUEST_CHANGES/PASS/HOLD decisions.

### Engineering
Owns:
- implementation strictly within an authorized work package;
- evidence;
- tests;
- documented deviations/blockers;
- GitHub handback.

Engineering does not redefine accepted architecture in code.

## Collaboration channel

For each work package:
- the GitHub issue is the coordination/control thread;
- the draft PR is the implementation/evidence thread;
- Architecture instructions are posted in GitHub;
- Engineering questions, blockers and handbacks are posted in GitHub;
- the owner is given only a concise summary and is not used as a manual relay for long prompts.

If Architecture asks Engineering to perform work, the detailed instruction must live in GitHub.

If Engineering completes or blocks work, the detailed handback must live in GitHub.

## Owner interaction

The owner may authorize, stop, narrow or redirect work.

Do not require the owner to paste Architecture prompts into an Engineering chat or paste Engineering handbacks back to Architecture.

A concise owner-facing message should say only what happened and where the authoritative GitHub record is.

## Work-package discipline

Before starting any WP:
1. read the controlling work order;
2. read the referenced design/ADR documents;
3. verify the exact required parent commit;
4. verify scope and host authorization;
5. post/confirm any prerequisite receipt required by the WP;
6. only then implement.

If blocked by an architecture conflict:
- stop that decision path;
- post evidence in the WP issue/PR;
- request Architecture direction;
- do not invent a replacement architecture.

## Runtime policy

No Docker, Podman, Compose, Kubernetes or OCI runtime dependencies are permitted on Cortex estate hosts, including ERIS.

Native Linux/systemd is the target runtime.

Any host work must be explicitly authorized by the controlling WP.

## Engineering handback minimum

Every handback must include:
- WP identifier;
- branch;
- exact parent SHA;
- exact HEAD SHA;
- files changed;
- commands/tests executed;
- results;
- evidence paths;
- host changes made, if authorized;
- known deviations;
- unresolved risks/blockers;
- explicit scope-compliance statement.

After handback, stop unless the current WP explicitly authorizes another step.

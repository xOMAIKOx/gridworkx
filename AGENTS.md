# GRIDWORKS Agent Operating Contract

This repository uses GitHub as the authoritative collaboration channel between the Owner, Technical Authority and Engineering.

Controlling collaboration references:

1. `AI_AGENT_COLLABORATION_PROTOCOL.md` — GRIDWORKS role/authority rules.
2. `docs/GITHUB_COLLABORATION_TRANSPORT.md` — normative GitHub-first transport and handback protocol.
3. `docs/work-orders/**` — work-package architecture and acceptance contracts.
4. Active WP issue + draft PR — current execution/review ledger.

## Roles

### Owner

Owns product/business intent, explicit GO boundaries, merge/deployment authorization and material decisions. The Owner is not a technical message relay.

### Technical Authority

GPT-5.6 Sol/Pro unless the Owner explicitly assigns another authority.

Owns architecture, scope, worker routing, implementation directives, acceptance criteria, review findings and PASS / REQUEST_CHANGES / HOLD-BLOCKED decisions.

### Engineering

Devin is the persistent engineering/orchestration agent.

Devin owns repository execution, invokes the worker model selected by the Technical Authority, runs verification, commits evidence and produces the structured Engineering handback.

The worker model executes inside Devin's lane and does not replace Devin or the Technical Authority.

## GitHub-first collaboration

For each WP:

- substantial Technical Authority instructions live in a durable repository Markdown directive whenever practical;
- the WP issue / draft PR carry concise control receipts pointing to the directive;
- the execution route must name Devin, the exact worker model/effort and implementation class;
- Engineering questions/blockers return to GitHub, not through the Owner;
- the complete Engineering handback is committed as a durable Markdown document;
- the issue/PR receive concise handback pointers only;
- the Owner receives a concise status/decision request, never a long agent-to-agent packet.

Do not ask the Owner to copy/paste Architecture prompts to Engineering or Engineering handbacks back to Architecture.

## Work-package discipline

Before starting or resuming work:

1. read the controlling directive/work order;
2. verify the exact branch, parent/current HEAD and draft PR;
3. verify the selected worker model/effort;
4. verify authorized/prohibited scope;
5. verify host/deployment authorization independently from repository authorization;
6. stop if the actual baseline materially differs from the directive.

Engineering must not silently rebase, change worker model/effort, redesign architecture, weaken a legitimate gate, merge, or start the next WP.

## Runtime policy

No Docker, Podman, Compose, Kubernetes, OCI, containerd or nerdctl on Cortex estate hosts, including ERIS.

Native Linux/systemd is the estate runtime.

Repository authorization does not imply host/deployment authorization.

## Mandatory Engineering handback

Every review submission must include a committed handback document containing:

- work package and directive reference;
- execution agent: Devin;
- actual worker model + effort;
- implementation class;
- repository/branch;
- controlling parent/current baseline;
- implementation/remediation commit(s);
- exact final HEAD;
- changed files/components;
- verification commands/results;
- acceptance-criteria mapping;
- deviations/assumptions;
- blockers/risks;
- environment/host/database/network/deployment mutation status;
- final state: `READY FOR SOL/PRO REVIEW`, `BLOCKED` or `STOP`.

After handback, Engineering stops.

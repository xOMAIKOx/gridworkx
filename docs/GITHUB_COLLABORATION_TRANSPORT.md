# GRIDWORKS GitHub-First Agent Collaboration Transport

**Status:** Normative companion to `AI_AGENT_COLLABORATION_PROTOCOL.md`

## 1. Purpose

GitHub is the durable collaboration bus between the Technical Authority and Devin.

The Owner must not have to copy/paste substantial implementation directives, handbacks, review findings, evidence blocks or remediation packets between agent sessions.

Canonical flow:

```text
OWNER
  | concise priority / decision / GO
  v
TECHNICAL AUTHORITY
  | writes controlling directive/review/remediation to GitHub
  | provides only a concise owner status/pointer
  v
DEVIN
  | reads GitHub directly
  | executes with the Architecture-selected worker
  | commits implementation + evidence + handback
  v
TECHNICAL AUTHORITY
  | reads handback/diff/CI directly from GitHub
  | records verdict/remediation in GitHub
  v
OWNER
  | receives material outcome and only decisions/GO actually required
```

**GitHub is the engineering system of record. Chat is not the long-form transport layer.**

## 2. Owner interaction

The Owner should normally receive only:

- current gate/status;
- material change/risk/blocker;
- any decision/GO actually required;
- a short repository/issue/PR/directive pointer.

A bootstrap message for another agent/session should be no more than enough to identify:
- repository;
- active WP issue/PR;
- controlling directive path;
- required worker route;
- instruction to read GitHub, execute, hand back to GitHub and stop.

## 3. Technical Authority -> Devin

For substantial work, prefer a durable directive under:

`docs/engineering/directives/`

The issue/PR should contain a concise control receipt with:

- directive path + directive commit SHA;
- controlling parent/current baseline;
- authorized scope;
- execution agent: Devin;
- exact worker model/effort;
- implementation class;
- STOP boundary;
- merge/deployment/host status.

For short reviews/remediations, a GitHub comment may be canonical, but if the remediation controls execution across multiple files/gates or requires a structured handback, use a durable directive.

## 4. Devin -> Technical Authority

The complete handback must be committed under:

`docs/engineering/handbacks/<wp>/`

or another explicitly named canonical project path.

The handback must include:

- execution agent identity: Devin;
- actual worker model/effort;
- repository/branch;
- controlling directive + parent/baseline;
- implementation/remediation commits;
- final HEAD;
- changed files/components;
- verification commands/results;
- acceptance-criteria mapping;
- deviations/assumptions;
- security/architecture findings;
- environment/host/database/network/deployment mutation status;
- blockers/risks;
- explicit `READY FOR SOL/PRO REVIEW`, `BLOCKED` or `STOP`.

After committing the handback, Devin posts only a concise pointer/receipt to the controlling issue/PR.

Example:

```text
DEVIN HANDBACK READY

HEAD: <40-char SHA>
Handback: docs/engineering/handbacks/wp-010/WP010_R9_HANDBACK.md
Worker: GPT-5.6 Luna XHigh
State: READY FOR SOL/PRO REVIEW

No work beyond the authorized boundary was started.
```

Then Devin gives the Owner only a concise status/pointer and stops.

## 5. Architecture review

Architecture independently reads the committed handback, actual diff, CI and evidence.

Architecture writes back to GitHub:
- PASS;
- exact remediation;
- investigation directive;
- BLOCKED/HOLD;
- STOP/architecture return.

The Owner is not used as a relay.

## 6. Recovery

Every fresh Architecture or Devin session must recover from GitHub without chat archaeology:

1. `AGENTS.md`;
2. this transport document;
3. current handoff/checkpoint issue;
4. active WP issue;
5. active draft PR;
6. controlling work order/directive;
7. latest handback/review receipt.

If GitHub state and an older chat summary disagree, the latest explicit GitHub record plus later Owner decisions control.

## 7. No duplicate authorities

Do not create multiple competing copies of one directive/handback.

Use:
- one canonical directive document per execution/remediation pass;
- one canonical handback document per pass;
- concise issue/PR receipts pointing to those documents;
- one Architecture verdict per reviewed HEAD.

Do not force-push reviewed history unless Architecture explicitly authorizes it.

## 8. Security

Never place plaintext secrets, passwords, tokens, private keys, secret-bearing env files or unauthorized customer content into directives, issues, PRs, handbacks or chat.

Use redacted evidence and secret references only.

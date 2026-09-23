# GRIDWORKS AI Agent Collaboration Protocol

**Status:** Normative project protocol

GRIDWORKS follows the Cortex governed development model:

```text
OWNER
  -> TECHNICAL AUTHORITY (GPT-5.6 Sol/Pro)
  -> DEVIN (Engineering/orchestration)
  -> selected worker model
  -> DEVIN verification + durable handback
  -> TECHNICAL AUTHORITY independent review
  -> OWNER only when a decision/GO is required
```

The detailed transport rules are normative in:

`docs/GITHUB_COLLABORATION_TRANSPORT.md`

## 1. Authority

### Owner

Controls:
- product/business intent;
- priority;
- explicit repository GO boundaries;
- merge authorization;
- host/deployment/production GO boundaries;
- material exceptions/decisions.

### Technical Authority

Controls:
- requirements interpretation;
- architecture/security/data-integrity invariants;
- WP decomposition and scope;
- acceptance criteria;
- exact Engineering directive;
- worker-model/effort selection;
- independent review;
- PASS / REQUEST_CHANGES / HOLD-BLOCKED / STOP verdict.

### Devin

Devin is the persistent Engineering/IDE orchestration agent.

Devin:
- verifies repository/baseline state;
- uses the exact worker model/effort selected by Architecture;
- executes only the authorized scope;
- runs mandatory verification;
- commits implementation/evidence/handback;
- reports the actual worker route used;
- stops after handback.

Devin may not silently:
- substitute or escalate worker model/effort;
- redesign architecture;
- widen scope;
- infer host/deployment permission;
- self-approve or merge.

## 2. Worker routing

Default Technical Authority routing:

- solved/deterministic implementation or exact remediation -> **Devin + GPT-5.6 Luna XHigh**;
- unresolved implementation diagnosis inside fixed architecture -> **Devin + GPT-5.6 Luna Max**;
- difficult/non-convergent implementation -> explicit Architecture reroute only.

A worker change is a new Architecture decision.

## 3. Engineering directive minimum

Every substantial implementation/remediation directive must state:

- WP;
- Technical Authority;
- execution agent: Devin;
- worker model + effort;
- implementation class D1-D5;
- repository/branch/parent or reviewed HEAD;
- controlling work order/ADR/review;
- authorized scope;
- prohibited/deferred scope;
- architectural invariants;
- implementation requirements;
- numbered acceptance criteria;
- mandatory verification;
- STOP conditions;
- required handback path/content.

## 4. Mandatory STOP conditions

Engineering stops and returns evidence to Architecture when:

- baseline/parent is wrong;
- success requires architecture change or scope expansion;
- a material security/data-integrity contradiction appears;
- success requires weakening/bypassing a legitimate gate;
- host/deployment mutation would be required without explicit GO;
- selected worker model/effort is unavailable;
- further attempts would become speculative;
- actual repository state materially differs from the directive.

A BLOCKED/STOP handback is a correct outcome.

## 5. Review

Architecture reviews the actual diff, tests, CI and evidence, not Engineering confidence.

Verdicts:
- `PASS`
- `REQUEST_CHANGES — EXACT REMEDIATION`
- `REQUEST_CHANGES — INVESTIGATION REQUIRED`
- `HOLD/BLOCKED`
- `STOP / ARCHITECTURE RETURN`

PASS applies only to the exact reviewed HEAD.

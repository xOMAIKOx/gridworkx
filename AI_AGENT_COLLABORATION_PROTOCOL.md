# GRIDWORKS GitHub Collaboration Protocol

## Purpose

Keep Architecture ↔ Engineering collaboration in GitHub so the owner does not have to relay long prompts or handbacks through chat.

## Authoritative locations

For a work package such as WP-001:

1. **Work order**
   - `docs/work-orders/<WP>.md`
   - durable specification and acceptance contract.

2. **GitHub issue**
   - current coordination thread;
   - Architecture instructions/clarifications;
   - Engineering prerequisite receipts;
   - blockers;
   - short status records;
   - final handback reference.

3. **Draft pull request**
   - implementation diff;
   - test/evidence details;
   - review findings;
   - Architecture acceptance or REQUEST_CHANGES.

## Architecture → Engineering

Architecture posts detailed instructions directly in the relevant issue/PR.

Owner-facing summary should be short, for example:

> WP-001 instruction is on GitHub issue #1. Engineering should read the work order and issue, execute only the authorized scope, and hand back in GitHub.

The owner should not need to copy the full instruction.

## Engineering → Architecture

Engineering posts its detailed handback directly in the issue/PR.

Owner-facing summary should be short, for example:

> WP-001 handback is posted on issue #1 / PR #N at commit <SHA>. Architecture can review there.

The owner should not need to copy the handback.

## Blockers

If Engineering is blocked:
- post `BLOCKED` in the issue/PR;
- state exact blocker;
- include evidence/commands;
- state what decision or dependency is required;
- stop the affected work path.

Architecture replies in GitHub.

## Review states

Architecture records one of:
- `PASS`
- `REQUEST_CHANGES`
- `HOLD/BLOCKED`

A PASS applies only to the exact reviewed commit.

## Host work

Host work is permitted only when explicitly stated in the controlling work order.

For WP-001, the only currently authorized ERIS host work is the bounded prerequisite/toolchain audit and installation described in the work order.

No deployment or shared-service configuration is implied by permission to install development prerequisites.

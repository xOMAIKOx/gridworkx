# ADR-004 — Rules, Content and Localization Distribution

**Status:** ACCEPTED  
**Date:** 2026-09-19  
**Applies to:** simulation rules, resources, recipes, manager traits, challenges, localization

## Context

GRIDWORKS requires frequent balancing and global localization without rebuilding the entire mobile application for every text/content adjustment. At the same time, clients must not be able to inject arbitrary rules into the shared economy.

## Decision

Rules/content/localization are distributed as **versioned, signed content bundles**.

Bundle classes:

1. **Simulation rules bundle**
   - component definitions;
   - recipes;
   - failure modes;
   - repair/intervention definitions;
   - manager traits;
   - XP/balance parameters;
   - challenge rules.

2. **Localization bundle**
   - message catalogue;
   - terminology/glossary;
   - font/fallback metadata;
   - localized tutorial/help content.

3. **Presentation/content bundle**
   - non-authoritative descriptors and content metadata.

## Publication

Bundles are authored in the monorepo, validated in CI, assigned an immutable content/rules version, signed/checksummed, published to object storage/CDN delivery, and referenced by server-issued manifests.

The server determines which rules version is authoritative for a shared-world operation.

The client may cache multiple versions when necessary to replay historical/challenge state.

## Safety

- unsigned bundles are rejected;
- bundle hash/version is included in simulation/replay evidence;
- clients cannot select an arbitrary economy rules version;
- rollback requires an explicit published manifest;
- incompatible content changes require migration policy.

## Localization

Localization bundles may update independently from simulation rules where semantic keys remain compatible.

Normal gameplay never requires a live translation/LLM service.

## Rationale

This allows rapid balance/localization delivery while preserving authoritative, auditable shared rules.

## Consequences

- WP-001 establishes package directories and signing/checksum tooling interfaces.
- WP-021 implements locale-pack pipeline.
- WP-023 challenge evidence records rules/content version.

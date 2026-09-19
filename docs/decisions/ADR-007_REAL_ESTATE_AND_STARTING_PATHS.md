# ADR-007 — Launch Real Estate Domain and Player-Selected Starting Paths

**Status:** ACCEPTED  
**Date:** 2026-09-19  
**Applies to:** WP-001 domain scaffolding, onboarding architecture, property/business lifecycle

## Context

GRIDWORKS should appeal to players interested in mining, farming, logistics, finance, real estate and other disciplines without forcing everyone to operate businesses they do not enjoy.

Real estate must exist from the launch foundation and can expand substantially in later seasons.

## Decision

### 1. Real estate is a launch domain

WP-001 must reserve the domain, schemas/packages/directories and ownership concepts required for:

- land parcels;
- buildings;
- property ownership;
- property sale/listing;
- lease state;
- occupancy state;
- basic maintenance/condition;
- renovation jobs;
- property valuation;
- GRIDWORKS-owned inventory;
- player-owned property.

Season 1 implementation may be intentionally limited.

Initial property classes should include:
- industrial land;
- warehouses;
- factories/industrial buildings;
- offices;
- basic residential property.

Season 2+ may expand into:
- apartment blocks;
- mixed-use;
- richer residential systems;
- interior redesign;
- conversions;
- hospitality/retail;
- advanced property management;
- property companies/portfolio mechanics.

Future assets may already exist in-world at launch under owner **GRIDWORKS** with status unavailable/reserved.

### 2. Player-selected starting path

After the shared core tutorial teaches the GRIDWORKS reasoning loop, the player selects a preferred starting path.

Candidate launch paths:
- Mining / Quarrying;
- Agriculture;
- Manufacturing;
- Energy;
- Logistics;
- Real Estate;
- Finance / Markets;
- Generalist / Surprise Me.

This is a starting position, not a permanent class.

### 3. No forced unwanted businesses

Players are never required to actively operate businesses they do not want.

Supported lifecycle actions should include, where applicable:
- Operate;
- Mothball / Suspend;
- Sell;
- Lease/Contract Out later;
- Dismantle / Repurpose where physically relevant.

Mothballing reduces variable operating costs while leaving appropriate fixed holding/maintenance costs.

### 4. Finance-first path

A finance-focused player must be able to begin with capital and market/investment gameplay rather than being forced to run an industrial facility first.

### 5. Cross-domain visibility

All players should be able to browse the wider economy from early play:
- business opportunities;
- property listings;
- jobs/contracts;
- market data;
- consortium opportunities;
- finance/commodity systems as unlocked by complexity/readiness.

The starting path does not isolate the player from the rest of GRIDWORKS.

### 6. Early correction

The first-path choice must not trap the player.

A bounded early buyback/restart mechanism may allow a player to exit a starter business and redirect equivalent starter capital into another path without deleting the account.

This mechanism must be abuse-resistant and limited.

## Invariants

1. Starter path is not a class lock.
2. No player is forced to operate unwanted businesses.
3. Real estate exists in the launch architecture.
4. Property uses the same ownership, ledger, maintenance and market primitives as other assets.
5. GRIDWORKS-owned future inventory may be visible before it becomes purchasable.
6. Real estate productive assets cannot be purchased with real money.
7. Season expansion adds depth without requiring a new incompatible property model.

## Consequences

- WP-001 must scaffold real-estate/land/property domains.
- Onboarding architecture must support starter-path selection.
- Later vertical slices may differ by path while sharing the same causal rules engine.

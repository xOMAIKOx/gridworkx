# ADR-006 — GRIDWORKS-Owned Assets, Non-Political Authorities and Treasury Recirculation

**Status:** ACCEPTED  
**Date:** 2026-09-19  
**Applies to:** economy sinks/sources, property, ports/logistics, business resale, live-ops balancing

## Context

GRIDWORKS needs economic sinks such as property tax, port fees, warehouse charges and service levies, but the product must avoid political gameplay. There are no mayors, governments, elections, parties, ideological systems or political governance layer.

The platform also needs a transparent way to seed property/business inventory and optionally recirculate collected in-game value back into the economy.

## Decision

### 1. No political-government simulation

GRIDWORKS does not model elected governments or political institutions.

Allowed in-world abstractions include operational authorities such as:

- Local Authority;
- Port Authority;
- Utility Authority;
- Land/Registry Authority;
- Infrastructure Authority.

These are accounting/service abstractions only.

They do not have:
- elections;
- political parties;
- ideological positions;
- political agendas;
- player voting;
- mayor/governor roles.

### 2. GRIDWORKS-owned assets

The platform may own in-game economic assets under the visible owner label:

> **GRIDWORKS**

Examples:
- land parcels;
- warehouses;
- factories;
- offices;
- residential property;
- distressed facilities;
- businesses;
- future-season inventory.

GRIDWORKS-owned assets may be:
- fixed-price Buy Now;
- leaseable;
- reserved/unavailable;
- released through later seasons;
- occasionally offered to players under system-generated purchase offers.

GRIDWORKS-owned is platform/system inventory, not a political/state actor and not the legal operator of the software represented inside gameplay.

### 3. Player-owned asset interaction

Player-owned assets may expose:
- owner/company identity;
- Contact Owner;
- Make Offer;
- Buy Now where seller set a price;
- lease terms where supported.

### 4. Treasury / authority ledgers

Economic sinks such as:
- property tax;
- port charges;
- warehouse/local authority fees;
- registry/transaction fees;
- infrastructure usage charges;

are booked into internal non-player treasury/authority ledgers.

Players do not manage or politicize those treasuries.

Operators/admin tooling may inspect:
- revenue by authority type;
- region;
- asset class;
- period;
- reinvestment/outflow.

### 5. Recirculation

Treasury funds may be recirculated into the world through controlled sinks/sources such as:
- GRIDWORKS purchase offers for player businesses/properties;
- system procurement;
- community project funding;
- infrastructure maintenance/subsidy;
- emergency liquidity operations;
- seasonal asset releases;
- starter/recovery buyback programmes.

Treasury spending must be bounded and rules-driven so GRIDWORKS does not become an infinite guaranteed buyer.

### 6. GRIDWORKS buyback/offers

When a player lists or attempts to sell a business/property, GRIDWORKS may occasionally generate a purchase offer.

Such offers should:
- be probabilistic/conditional rather than guaranteed;
- use auditable valuation inputs;
- avoid systematically beating healthy player-market prices;
- include cooldowns/limits;
- be visible as GRIDWORKS offers;
- settle through the same ledger as player transactions.

The offer engine exists partly as:
- a liquidity backstop;
- a currency recirculation mechanism;
- a way to replenish system-owned inventory;
- a live-ops balancing tool.

## Invariants

1. No politics/elections/governance gameplay.
2. Authority labels exist only to explain economic/service charges.
3. GRIDWORKS-owned inventory is clearly distinguishable from player-owned inventory.
4. Treasury creation/spending is fully ledgered and operator-auditable.
5. No infinite guaranteed buyback floor for normal assets.
6. System offers must not secretly manipulate competitive outcomes.
7. Real money cannot purchase GRIDWORKS-owned productive property/assets directly.

## Consequences

- Economy schema must support non-player ledger principals.
- Property/business ownership must support GRIDWORKS as a system principal.
- Admin tooling must expose treasury inflow/outflow and buyback behavior.
- Business/property sale flows must allow optional GRIDWORKS offers alongside player offers.

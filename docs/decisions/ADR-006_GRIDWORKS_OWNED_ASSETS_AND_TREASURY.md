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

### 6. GRIDWORKS buyer-of-last-resort and resale model

When a player wants to sell a business/property and no acceptable player buyer exists, GRIDWORKS may act as a **buyer of last resort**.

This is a designed liquidity facility rather than a speculative competitor.

#### 6.1 Acquisition price

GRIDWORKS acquisition price should be derived from an auditable valuation model using inputs such as:

- recent comparable player transactions;
- asset/property condition;
- trailing revenue/profit/cash flow where relevant;
- inventory and included working capital;
- liabilities/encumbrances;
- regional demand;
- current market depth;
- asset liquidity;
- estimated replacement/redevelopment value.

The GRIDWORKS offer should normally be a **liquidation/convenience price below expected open-market value**. This gives the seller immediate certainty while preserving an incentive to wait for a player buyer when maximizing price matters.

The exact discount is a balance parameter, not hard-coded in this ADR.

#### 6.2 Ownership transfer

If accepted:

- ownership transfers to GRIDWORKS;
- the player receives immediate settled credits;
- the transaction is fully ledgered;
- the exact facility/business history remains attached to the asset;
- GRIDWORKS becomes the visible owner;
- the asset may be mothballed, held, leased or relisted according to system rules.

#### 6.3 GRIDWORKS resale pricing

GRIDWORKS should normally relist acquired assets.

The listing price may move with market conditions, but GRIDWORKS must not intentionally sell the asset below its own acquisition cost except through an explicitly authorized economy intervention.

Normal resale floor:

> **GRIDWORKS acquisition cost + applicable carrying/transaction costs**

The live asking price may be higher based on:

- scarcity;
- regional demand;
- comparable sales;
- asset improvement/deterioration;
- carrying period;
- sector outlook;
- lack of competing inventory.

#### 6.4 Original-owner reacquisition

A former owner may buy their old business/property back if GRIDWORKS still owns it.

This must not function as cheap temporary financing or free asset parking.

For the original seller, the reacquisition price is at least:

> **max(current GRIDWORKS listing price, GRIDWORKS acquisition price × 1.20, acquisition cost + carrying costs)**

The **20% minimum premium** is the initial design rule and remains a tunable balance parameter.

If equivalent alternative assets are unavailable, the normal scarcity/market pricing model may push the listing above that minimum.

The original owner receives no secret discount or priority price advantage merely because they previously owned the asset.

#### 6.5 Scarcity and fairness

When GRIDWORKS is the only seller of a comparable asset class in a region, scarcity may increase the asking price.

Scarcity pricing must:

- use published/consistent valuation rules;
- apply to all potential buyers;
- not target the former owner specifically;
- remain operator-auditable.

#### 6.6 Anti-abuse controls

The buyer-of-last-resort mechanism should include:

- valuation floors/ceilings;
- anti-circular-sale checks;
- cooldowns where required;
- related-account checks;
- limits on repeated sell/rebuy loops;
- full ownership/price history;
- treasury liquidity constraints where appropriate.

The purpose is to provide a reliable exit route for players who need capital, not an arbitrage engine.

The offer engine therefore serves as:

- a guaranteed liquidity backstop where the asset class is eligible;
- a currency recirculation mechanism;
- a way to replenish GRIDWORKS-owned inventory;
- a market stabilizer;
- a path for players to redeploy capital into another profession or business.


## Invariants

1. No politics/elections/governance gameplay.
2. Authority labels exist only to explain economic/service charges.
3. GRIDWORKS-owned inventory is clearly distinguishable from player-owned inventory.
4. Treasury creation/spending is fully ledgered and operator-auditable.
5. Eligible normal assets may use a guaranteed GRIDWORKS liquidity offer, but at conservative system valuation rather than guaranteed open-market value.
6. System offers must not secretly manipulate competitive outcomes.
7. Real money cannot purchase GRIDWORKS-owned productive property/assets directly.

## Consequences

- Economy schema must support non-player ledger principals.
- Property/business ownership must support GRIDWORKS as a system principal.
- Admin tooling must expose treasury inflow/outflow and buyback behavior.
- Business/property sale flows must support player offers plus an eligible GRIDWORKS buyer-of-last-resort path.
- Ownership history and original-seller reacquisition pricing must be preserved and auditable.

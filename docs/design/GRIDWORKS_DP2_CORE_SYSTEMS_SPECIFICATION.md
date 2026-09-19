# GRIDWORKS — Design Phase 2 (DP2): Core Systems Specification

**Repository:** `xOMAIKOx/gridworkx`  
**Status:** DP2 DESIGN / CONTROLLING SYSTEMS SPECIFICATION  
**Date:** 2026-09-19  
**Depends on:** `docs/design/GRIDWORKS_PRODUCT_PHILOSOPHY_AND_GAME_DESIGN.md`  
**Scope:** Define the shared simulation primitives and cross-cutting rules that all later industry, economy, social and finance modules must use.

---

## 0. DP2 purpose

DP2 converts the founding GRIDWORKS philosophy into a coherent systems architecture suitable for later decomposition into engineering work packages.

DP2 is deliberately **not** an exhaustive content catalogue. It defines the rules of the world.

The primary engineering target after DP2 is a vertical slice proving this loop:

> **broken facility → inspect → diagnose → choose intervention → recover partial production → earn → trade/help → improve skill/manager → make another economically meaningful decision**

If the shared primitives defined here do not support that loop cleanly, later modules must not bypass them with bespoke logic.

---

# PART I — CORE WORLD MODEL

## 1. Simulation principles

### 1.1 Persistent world, asynchronous play

GRIDWORKS is a persistent world in which time continues while the player is offline, subject to protective rules that respect the player.

Offline progress may include:

- construction advancing;
- production continuing while inputs, capacity and safe operating constraints permit;
- crops progressing;
- contracts approaching deadlines;
- markets moving;
- managers recovering from workload;
- maintenance due dates approaching.

Offline absence must not create punitive traps merely because the player slept, worked or did not log in.

Examples:

- storage full → affected production pauses safely;
- harvest ready → remains available or transitions to safe storage rules;
- manager fatigued → performance normalizes only through intended recovery rules, not paid energy;
- critical failure → facility may stop safely rather than silently destroy itself.

### 1.2 Simulation clock

The game needs two related notions of time:

1. **Real time** — used for construction durations, recruitment refreshes, contracts, market events and social coordination.
2. **Simulation time** — used for facility operations, production cycles, equipment operating hours, crop growth, maintenance intervals and economic accounting.

The ratio between real and simulation time may vary by context but must be deterministic and documented.

Recommended initial architecture:

- one authoritative UTC server clock;
- deterministic client simulation using server-issued timestamps;
- no trust in local device wall-clock for authoritative economic outcomes;
- offline catch-up calculated from last authoritative state + elapsed authoritative time + bounded simulation rules.

### 1.3 Tick model

Avoid forcing every object to update continuously on the server.

Use event-driven / interval-based simulation where possible.

Entities should be able to derive state from:

- last authoritative state;
- last simulation timestamp;
- elapsed time;
- active process;
- available inputs;
- current constraints;
- event history.

The client may render high-frequency visual animation locally while the underlying economy advances at coarser deterministic intervals.

### 1.4 Determinism

Given the same:

- starting state;
- rules version;
- simulation seed;
- event sequence;
- elapsed authoritative time;

the simulation should produce the same result.

Randomness must therefore be:

- seeded;
- versioned;
- auditable;
- replayable for competitive or suspicious outcomes.

---

## 2. World structure

### 2.1 World → Region → Site → Facility → System → Component

Proposed containment model:

```
World
 └── Region
      └── Site / Parcel
           └── Facility
                └── System
                     └── Component
```

Examples:

```
World
 └── Western Basin
      └── Site 14C
           └── Copper Concentrator
                └── Crushing Circuit
                     ├── Feed Conveyor
                     ├── Jaw Crusher
                     ├── Transfer Conveyor
                     └── Screen
```

### 2.2 Region

A Region contains persistent environmental/economic characteristics such as:

- climate profile;
- rainfall profile;
- temperature range;
- wind potential;
- solar potential;
- water availability;
- terrain;
- distance/friction to trade hubs;
- mineral/oil/gas resource probabilities;
- soil profiles;
- transport topology;
- environmental constraints.

Regions are intentionally imperfect.

No region should dominate all industrial categories.

### 2.3 Site / Parcel

A Site is a controllable area on which one or more facilities may exist.

Attributes may include:

- area;
- ownership/lease status;
- terrain;
- elevation;
- soil;
- water access;
- road/rail/water connections;
- grid connection;
- deposits/resources;
- zoning/usage constraints where gameplay-relevant;
- local storage;
- environmental state.

### 2.4 Facility

A Facility is an economically meaningful operating unit.

Examples:

- farm;
- greenhouse complex;
- mine;
- quarry;
- refinery;
- smelter;
- power plant;
- substation;
- warehouse;
- port terminal;
- rail yard;
- airport cargo terminal;
- water-treatment plant.

Facility state must be derived from underlying systems/components rather than represented by one arbitrary "health" number.

A summary health/availability metric may be displayed, but it is derived.

---

# PART II — FACILITY, COMPONENT AND DEPENDENCY MODEL

## 3. System graph

Each facility is represented as a directed graph of nodes and dependencies.

A node can represent:

- physical equipment;
- storage;
- process step;
- utility;
- control subsystem;
- transport connection;
- labor/staff capacity;
- environmental input.

Edges represent dependency or flow relationships such as:

- material;
- energy;
- water;
- control;
- transport;
- capacity;
- information;
- staffing.

### 3.1 Example

```
Power Supply
     ↓
Motor → Conveyor → Crusher → Screen → Stockpile
                ↑
             Feed Ore
```

If the motor cannot run:

- conveyor throughput = 0;
- crusher may be healthy but starved;
- screen may be healthy but starved;
- overall production falls.

Replacing the belt does not bypass the failed motor unless a separate solution is introduced.

### 3.2 Dependency types

Minimum shared dependency types:

- HARD: downstream cannot operate without dependency;
- CAPACITY: downstream is limited proportionally;
- QUALITY: output quality/yield degrades;
- RELIABILITY: failure risk increases;
- COST: operation continues but more expensively;
- OPTIONAL: alternative path/bonus.

### 3.3 Bottleneck rule

Facility throughput should emerge from the minimum effective capacity along active critical paths.

A new asset should not improve output where another bottleneck remains dominant.

This principle is central to GRIDWORKS.

---

## 4. Component model

Every operational component should support a shared core attribute set.

Proposed component properties:

- component_type;
- rated_capacity;
- effective_capacity;
- condition;
- wear;
- operating_hours;
- efficiency;
- reliability;
- maintenance_state;
- temperature/load/pressure or domain-relevant telemetry;
- required_inputs;
- produced_outputs;
- failure_modes;
- repair_options;
- replacement_options;
- salvage_value;
- installation_age;
- history;
- manager/skill modifiers where relevant.

Industry modules extend this base model rather than replacing it.

### 4.1 Condition versus operability

Condition and operability are distinct.

Example:

- belt condition 41%;
- belt operational = true;
- service due soon.

A component at poor condition can still function.

A component at high apparent condition may be non-operational due to a specific fault.

### 4.2 Capacity degradation

Wear may reduce:

- maximum throughput;
- efficiency;
- reliability;
- output quality;
- energy consumption;
- maintenance interval.

Effects should be component-specific.

---

# PART III — FAILURE, DIAGNOSIS, REPAIR AND MAINTENANCE

## 5. Failure model

Failures must have causal models.

Each failure mode defines:

- prerequisites;
- probability/hazard curve;
- telemetry symptoms;
- direct effect;
- downstream effects;
- repair options;
- temporary workaround options;
- consequences of continued operation;
- discoverability.

Example: motor bearing seizure

Prerequisites / influences:

- operating hours;
- load;
- lubrication quality;
- maintenance history;
- temperature.

Symptoms:

- elevated temperature;
- elevated current;
- vibration;
- reduced RPM;
- eventual stop.

### 5.1 Hidden state versus observable state

The game may know the true fault while the player sees evidence.

Player diagnosis depends on:

- available sensors;
- player skill;
- manager skill;
- inspection action;
- equipment sophistication.

The system may show confidence rather than certainty.

### 5.2 No arbitrary deception

A player should not be given materially false information merely to create difficulty.

Low skill may mean:

- incomplete evidence;
- lower confidence;
- broader list of probable causes.

It should not routinely mean fabricated facts.

---

## 6. Intervention model

Interventions should be first-class actions.

Types:

- inspect;
- test;
- clean;
- calibrate;
- lubricate;
- patch;
- repair;
- rebuild;
- replace;
- bypass;
- reroute;
- derate;
- shut down;
- outsource;
- salvage/dismantle.

Every intervention has:

- resource cost;
- labor/skill requirement;
- time;
- expected effect;
- residual risk;
- recoverable materials;
- impact on history.

### 6.1 Temporary fixes

Temporary fixes are encouraged.

Example pump:

**Patch**
- low material;
- short time;
- 50% capacity;
- short lifetime.

**Rebuild**
- medium resources;
- medium time;
- ~95% restored condition.

**Replace**
- high cost;
- longer lead time;
- new baseline condition/efficiency.

The game should reward context-sensitive decisions rather than always making "replace" optimal.

---

## 7. Maintenance

Maintenance modes may include:

- reactive;
- scheduled preventive;
- condition-based;
- predictive later.

Maintenance influences:

- reliability;
- downtime;
- cost;
- performance;
- asset life.

A facility operating at extreme throughput should often incur:

- accelerated wear;
- shorter maintenance intervals;
- higher failure risk.

Managers and player skill can improve planning/diagnosis but should not magically repeal physical constraints.

---

# PART IV — RESOURCES, INVENTORY AND PRODUCTION

## 8. Resource ontology

All tradable/consumable resources should use a common schema.

Suggested fields:

- resource_id;
- category;
- unit;
- mass/volume where needed;
- quality/grade;
- perishability;
- storage requirements;
- transport compatibility;
- hazardous/special handling flags;
- fungibility;
- base emergency-supplier policy;
- market eligibility;
- recipe usage;
- salvage/recycling pathways.

### 8.1 Resource categories

Initial categories:

- raw materials;
- agricultural products;
- fuels;
- water;
- processed materials;
- industrial components;
- construction materials;
- finished goods;
- waste/by-products;
- energy;
- services/labor;
- financial instruments later.

### 8.2 Quality/grade

Do not introduce quality tiers unless they create decisions.

Examples where grade matters:

- ore grade;
- fuel quality;
- crop quality;
- water purity;
- refined-metal purity.

---

## 9. Production recipes

A production process is a transformation with:

- required input(s);
- optional input(s);
- utility requirements;
- equipment requirements;
- labor/skill requirements;
- duration/rate;
- yield;
- quality effects;
- by-products;
- waste;
- energy/water requirements;
- capacity constraints.

Example:

```
Copper Ore
 + electricity
 + water
 → Copper Concentrate
 + tailings
```

A recipe must not assume infinite transport, storage or utilities.

### 9.1 Continuous and batch processes

Support both:

- continuous processes;
- batch processes.

Examples:

- pipeline/refinery flow: continuous;
- harvest/processing batch;
- construction batch;
- discrete manufacturing.

---

## 10. Storage

Storage is physical/economic capacity.

Properties may include:

- capacity;
- compatible resource classes;
- handling speed;
- losses;
- operating cost;
- spoilage control;
- safety constraints.

When full:

- upstream production may stop or reroute;
- material should not silently disappear;
- the player may incur opportunity cost.

---

# PART V — PLAYERS, COMPANIES, SKILLS AND MANAGERS

## 11. Player entity

Player represents the human account and persistent expertise.

Player owns/controls one or more Company entities according to future progression rules.

Player-level state includes:

- skills;
- achievements;
- reputation dimensions;
- membership state;
- cosmetics;
- challenge history;
- recruitment allowance;
- social relationships.

Economic assets belong primarily to Companies, not directly to the Player.

This separation is required for later public-company and investment systems.

---

## 12. Company entity

A Company has:

- ownership;
- cash;
- assets;
- liabilities;
- facilities;
- inventory;
- contracts;
- employees/managers;
- operating history;
- financial statements;
- reputation;
- market listings later;
- divisions/business units;
- JV/consortium interests.

A player should initially control one private company.

Multiple-company ownership may be unlocked later.

---

## 13. Player skills

Skills are learned through relevant behavior.

Possible top-level skill families:

- mechanical;
- electrical;
- process engineering;
- agriculture;
- mining;
- energy;
- water;
- logistics;
- construction;
- commerce;
- finance;
- management.

Skill XP must be tied to meaningful actions, not repetitive zero-risk spam.

### 13.1 XP sources

Examples:

- successful diagnosis;
- repair;
- facility operation;
- production;
- contract completion;
- assisting another player;
- solving a challenge;
- trading/financial analysis where appropriate.

### 13.2 Anti-grind principle

Repeated trivial actions should face diminishing XP or eligibility thresholds.

A player should improve by operating systems, not by tapping one free action 10,000 times.

---

## 14. Manager entity

Manager fields:

- manager_id;
- name;
- rarity;
- current employer;
- employment history;
- facility specialization;
- skills;
- traits;
- experience;
- workload;
- fatigue;
- morale;
- mentor/deputy links;
- contract state;
- market availability.

### 14.1 Rarity

Rarity bands:

- Bronze;
- Silver;
- Gold;
- Platinum.

Rarity controls:

- potential skill ceilings;
- number/rarity of traits;
- adaptability;
- unusual combinations.

Rarity does **not** directly encode a universal production multiplier.

### 14.2 Daily recruitment

Initial design target:

- free: 5 candidate interviews/day;
- member: 8 candidate interviews/day;
- same underlying rarity odds;
- no cash purchase of extra rolls;
- Platinum available to both.

Candidate options:

- recruit;
- reject;
- shortlist;
- later referral/trade.

Recruitment refresh should use authoritative server time.

### 14.3 Manager trading

Narrative model: contract/employment transfer.

Potential transaction modes:

- permanent employment transfer;
- temporary contract;
- consortium assignment;
- apprenticeship/deputy placement.

Manager history remains immutable/auditable.

---

# PART VI — CONSTRUCTION, TIME AND ASSISTANCE

## 15. Construction jobs

Construction is represented as a Job.

Job fields:

- target;
- prerequisites;
- required resources;
- committed resources;
- base duration;
- skill requirements;
- current progress;
- assistance capacity;
- state;
- start/end timestamps;
- cancellation/salvage policy.

### 15.1 Queue model

Membership may improve planning/automation but must not provide unlimited parallel productive throughput.

Design target:

- all players have equal fundamental construction capacity for equivalent assets;
- members may queue future jobs;
- members may gain administrative convenience;
- additional actual capacity must be earned in-game and available to free players.

---

## 16. Player assistance

Assistance jobs allow other players to contribute to:

- construction;
- repair;
- commissioning;
- harvesting;
- logistics;
- inspection;
- maintenance.

Assistance produces:

- time reduction or quality improvement;
- helper skill XP;
- payment/resources;
- reputation.

Assistance cap is job-specific.

Initial balancing range for time reduction: **25–40% maximum external reduction**.

No number in DP2 is final balance; the invariant is that assistance matters but cannot trivialize project time.

---

# PART VII — ECONOMY, TRADE AND CONTRACTS

## 17. Base currency

GRIDWORKS needs a non-real-money operating currency.

Working generic term: **Credits**.

Credits are earned in-game and must not be directly convertible to/from real money.

Credits are used for:

- player trades;
- wages/manager contracts;
- services;
- emergency supplier;
- construction services;
- financial markets later.

Premium cosmetic entitlements, if any, must remain separate from tradeable economic currency.

---

## 18. Sources and sinks

A stable economy requires designed currency/resource sinks.

Potential credit sources:

- NPC/bootstrap demand;
- player contracts;
- starter recovery missions;
- system/community projects;
- challenge rewards;
- economic activity.

Potential sinks:

- maintenance;
- energy/fuel inputs;
- labor/manager wages;
- transport;
- market fees;
- taxes/fees where appropriate to simulation;
- construction;
- emergency supplier premium;
- listing fees;
- repair;
- insurance later;
- research/services later.

The economy must not depend on endless currency creation without sinks.

---

## 19. Player market

Core market functions:

- buy orders;
- sell orders;
- market price history;
- quantity;
- quality/grade;
- location/region;
- delivery terms where needed.

Market settlement must be server-authoritative.

### 19.1 Geography

Price can vary by region because transportation is not free.

This creates arbitrage opportunities for logistics/trading players.

### 19.2 Emergency supplier

System supplier is a recovery floor, not the default source.

Rules:

- price materially above healthy player-market price;
- quantity limits;
- essential/basic resource focus;
- no rare/strategic endgame bypass.

---

## 20. Contracts

Contract entity supports:

- issuer;
- counterparty or open offer;
- contract type;
- resource/service;
- quantity;
- quality;
- schedule;
- delivery location;
- price;
- collateral/penalty later;
- deadline;
- status;
- completion record;
- subcontract permission.

Initial contract types:

- resource supply;
- transport;
- maintenance;
- construction;
- engineering/inspection;
- assistance.

Later:

- facility management;
- financing;
- power purchase/offtake;
- long-term supply.

### 20.1 Recovery path

A stuck player can browse Work Exchange contracts and turn time/skill into credits.

This is a foundational anti-frustration loop.

---

# PART VIII — REPUTATION AND TRUST

## 21. Reputation dimensions

Do not reduce reputation to one universal star rating.

Potential dimensions:

- contract reliability;
- delivery punctuality;
- quality compliance;
- assistance/community contribution;
- financial default history;
- trading integrity;
- consortium participation.

Reputation influences:

- partner willingness;
- contract terms;
- financing cost later;
- manager recruitment attractiveness later;
- consortium admission;
- listing eligibility.

---

# PART IX — LEADERBOARDS AND COMPETITIVE FAIRNESS

## 22. Leaderboard families

At minimum:

- Overall;
- Agriculture;
- Mining;
- Energy;
- Oil & Gas;
- Manufacturing;
- Logistics;
- Trade;
- Water;
- Circular Economy;
- Engineering;
- Contracts;
- Finance;
- Consortium.

### 22.1 Time windows

- weekly challenge;
- monthly;
- season;
- all-time.

### 22.2 Scale classes

Working scale:

- Local;
- Regional;
- National;
- Continental;
- Global.

Exact thresholds deferred to balancing.

---

## 23. Metrics

Raw output is secondary.

First-class measures include:

- net profit;
- ROI;
- ROC/ROIC;
- ROE where applicable;
- margin;
- capital efficiency;
- cost per unit;
- reliability;
- utilization;
- yield;
- water/energy efficiency;
- recovery rate;
- downtime;
- contract performance;
- risk-adjusted return later.

### 23.1 Composite scores

Composite scores must be transparent.

Every score must expose:

- component metrics;
- weightings;
- measured period;
- eligibility criteria.

No opaque "power score."

### 23.2 Reward invariant

Leaderboard rewards may confer prestige and bounded economic value.

They must not create a permanent compounding advantage that materially increases the chance of winning the next season.

---

# PART X — STANDARDIZED CHALLENGES

## 24. Challenge sandbox

A standardized challenge is isolated from persistent-company wealth.

Server provides:

- challenge version;
- deterministic seed;
- starting map;
- starting resources;
- allowed technology;
- normalized manager/staff state;
- objective;
- scoring formula;
- start/end window.

Client simulates locally under deterministic rules and submits:

- final state;
- event/replay digest;
- score inputs;
- integrity metadata.

Server validates before leaderboard acceptance.

---

# PART XI — CONSORTIUMS AND JOINT VENTURES

## 25. Consortium

A Consortium is a social/economic organization of players/companies.

Functions may include:

- shared projects;
- contribution ledger;
- contracts;
- pooled objectives;
- shared infrastructure;
- rankings;
- governance later.

### 25.1 Joint venture

A JV is a specific jointly owned economic entity/project.

Fields:

- participating companies;
- contribution amounts/types;
- ownership percentages;
- governance rules;
- revenue distribution;
- asset ownership;
- exit rules.

DP2 requires the domain separation between Consortium and JV even if both are implemented later.

---

# PART XII — FINANCIAL FOUNDATION

## 26. Accounting model

All Companies should use a real double-entry-inspired accounting model from the start, even before the stock exchange is built.

Minimum statements/ledgers:

- cash ledger;
- revenue;
- operating expenses;
- assets;
- liabilities;
- inventory;
- capital expenditure;
- depreciation model later if useful;
- profit/loss;
- cash flow;
- owner equity.

This avoids retrofitting finance onto an economy that never tracked capital properly.

### 26.1 Metrics

The accounting model must support correct calculation of:

- ROI for projects/investments;
- ROC/ROIC;
- ROE later;
- operating margin;
- net margin;
- asset turnover;
- free cash flow;
- debt ratios later.

Definitions must be fixed/versioned before competitive finance leaderboards launch.

---

## 27. Stock exchange — domain reservation

DP2 does not fully specify exchange microstructure, but the core domain must not prevent it.

Required future capabilities:

- companies remain distinct legal/economic entities;
- share classes/ownership ledger;
- issuance;
- dilution;
- dividends;
- public financial reporting;
- historical financial statements;
- trading ledger;
- corporate actions;
- acquisitions;
- insolvency/restructuring.

No real-money conversion or cash-out.

---

## 28. Commodities — domain reservation

Resources marked market-eligible may become commodity markets.

Future systems may include:

- spot market;
- regional basis;
- forwards/futures;
- hedging;
- producer/consumer exposure;
- speculative trading;
- position/risk controls.

Again, no real-money conversion.

---

# PART XIII — FIRST-SESSION VERTICAL SLICE

## 29. Recommended opening scenario

DP2 selects **a small aggregate-processing plant** as the preferred first vertical-slice facility unless later UX testing demonstrates that a farm is materially clearer.

Why aggregate plant:

- visually legible;
- simple material flow;
- clear dependency chain;
- immediate conveyor/motor example;
- easy to show partial output;
- introduces power, material handling, storage and market without agricultural season complexity.

### 29.1 Starting state

Facility:

- quarry feed stockpile;
- feed conveyor;
- electric motor;
- crusher;
- screen;
- output conveyor;
- finished aggregate stockpile;
- simple grid connection.

Faults/wear:

- feed conveyor motor seized — critical;
- belt worn — non-critical;
- screen partly blocked — throughput penalty;
- crusher healthy;
- output conveyor healthy;
- limited spare materials;
- no current production.

### 29.2 Advisor script intent

Step 1 — Inspect.

Advisor explains:

> Do not spend yet. Find what prevents flow.

Step 2 — Evidence.

Player sees worn belt first but belt remains operational.

Step 3 — Motor diagnosis.

Motor telemetry indicates seizure.

Step 4 — Choice.

Available actions might include:

- repair motor;
- replace motor;
- replace belt;
- clear screen;
- rent temporary haul truck.

Several are valid, but only some restart production efficiently.

Step 5 — Partial recovery.

Cheap repair restores motor to reduced capacity.

Plant starts at approximately 2–5% throughput.

Step 6 — First economic outcome.

First aggregate reaches stockpile.

Player can sell a small batch into bootstrap demand.

Step 7 — Second decision.

Player chooses between:

- belt repair;
- screen clearing;
- savings;
- buying material;
- doing a Work Exchange job.

Advisor becomes optional.

### 29.3 Failure-safe design

If the player spends initial resources badly:

- dismantle recovery;
- emergency supplier;
- starter Work Exchange;
- salvage;
- temporary haulage workaround.

No real-money prompt may be the required recovery route.

---

# PART XIV — PROGRESSION

## 30. Progression dimensions

GRIDWORKS should not rely on one universal player level.

Progression consists of:

- player skill;
- manager skill/history;
- company capital;
- facilities/assets;
- reputation;
- regional access;
- technology/tool access;
- market sophistication;
- organizational scale;
- achievements/challenge prestige.

A display "level" may exist for onboarding/UI, but it must not become the sole gating model.

### 30.1 Unlock philosophy

Unlocks should primarily represent complexity readiness and capability, not arbitrary grind.

Example:

Player does not access futures markets on day one because they lack the accounting/market context, not because "Level 38" is sacred.

---

# PART XV — MONETIZATION BOUNDARIES

## 31. Membership

Current working commercial baseline:

- free game;
- annual membership around **$35/year**;
- free recruitment interviews: 5/day;
- member interviews: ~8/day;
- identical rarity odds.

Possible member QoL:

- advanced analytics;
- extended history;
- additional watchlists;
- saved blueprints/layouts;
- queue planning/automation;
- extra cosmetic/company customization;
- early/additional content where fair;
- larger shortlists/administrative limits.

### 31.1 Forbidden monetization paths

Never allow:

real money → tradeable credits → commodities/shares/managers/bonds.

Never allow:

real money → guaranteed competitive production advantage.

No paid loot-box manager draws.

No cash-only Platinum.

---

# PART XVI — CLIENT/SERVER AUTHORITY

## 32. Client responsibilities

Where practical:

- rendering;
- local UI;
- high-frequency animation;
- deterministic facility simulation;
- offline cache;
- local planning;
- challenge execution;
- diagnostic presentation.

### 32.1 Server responsibilities

Authoritative for:

- accounts;
- ownership;
- company cash/ledger;
- markets;
- trade settlement;
- contracts;
- manager recruitment results;
- manager transfers;
- world/region persistent state;
- challenge seeds/results;
- leaderboards;
- membership entitlements;
- consortium/JV ownership;
- public financial markets later;
- authoritative time;
- anti-cheat/risk signals.

### 32.2 Sync policy

Client submits intents and/or deterministic state transitions.

Server validates:

- elapsed time;
- resources;
- ownership;
- prerequisites;
- rule version;
- transaction integrity.

Critical economy state must never be accepted solely because the client says it happened.

---

# PART XVII — EVENT MODEL

## 33. Domain events

Design systems around durable events where useful.

Examples:

- FacilityCreated;
- ComponentInstalled;
- InspectionCompleted;
- FaultDetected;
- RepairStarted;
- RepairCompleted;
- ProductionStarted;
- ProductionStopped;
- ResourceProduced;
- ResourceConsumed;
- InventoryChanged;
- ContractIssued;
- ContractAccepted;
- ContractCompleted;
- TradeSettled;
- AssistanceProvided;
- SkillAdvanced;
- ManagerRecruited;
- ManagerTransferred;
- ConstructionStarted;
- ConstructionCompleted;
- DividendDeclared later.

Not every implementation must be pure event sourcing, but these event boundaries should be explicit for analytics, replay, anti-cheat and debugging.

---

# PART XVIII — VERSIONING AND BALANCE

## 34. Rules versioning

Every simulation result should be attributable to a rules version.

Version:

- recipes;
- failure curves;
- manager traits;
- scoring formulas;
- market fees;
- XP curves;
- challenge rules.

This is essential when balancing a live persistent economy.

### 34.1 Balance migration

Rule changes must avoid silently destroying player investments.

Where possible:

- grandfather state safely;
- compensate material changes;
- announce major economy adjustments;
- avoid retroactive invalidation of earned assets.

---

# PART XIX — ECONOMIC STABILITY

## 35. Inflation/deflation controls

The economy requires monitoring and bounded controls.

Track:

- money supply;
- velocity;
- resource production;
- resource destruction/consumption;
- market depth;
- price volatility;
- concentration;
- hoarding;
- emergency supplier usage.

Possible controls:

- maintenance/resource sinks;
- market fees;
- construction demand;
- system procurement;
- emergency supply pricing;
- community projects;
- NPC/bootstrap demand;
- storage/carry costs where justified.

Do not solve inflation by arbitrarily deleting player wealth.

---

# PART XX — SECURITY AND ABUSE MODEL

## 36. Threat classes

At minimum:

- modified client;
- local clock tampering;
- replay attacks;
- duplicate transaction submission;
- resource duplication;
- multi-account farming;
- wash trading;
- collusion;
- automation/bots;
- market manipulation;
- leaderboard replay forgery.

### 36.1 Security principle

Any state that can materially affect another player must be server-verifiable.

A fully offline player may plan and simulate locally, but synchronization into the shared economy requires authoritative validation.

---

# PART XXI — ANALYTICS AND OBSERVABILITY

## 37. Product telemetry

GRIDWORKS should measure whether the philosophy works.

Key product signals:

- first facility restart rate;
- time to first meaningful diagnosis;
- percentage choosing wrong-but-recoverable intervention;
- recovery success without real-money spend;
- Work Exchange use;
- trade participation;
- specialization diversity;
- member/free competitive distribution;
- leaderboard representation by membership tier;
- manager rarity distribution;
- economic concentration;
- churn around timers/maintenance;
- session frequency;
- optional advisor usage.

### 37.1 Fairness telemetry

Explicitly monitor:

- whether members statistically dominate economy/leaderboards beyond skill/engagement expectations;
- whether any paid QoL mechanic becomes de facto power;
- whether Platinum distribution is materially distorting markets;
- whether free players hit practical progression walls.

If data shows a paid feature has become power, design must correct it.

---

# PART XXII — DP2 ACCEPTANCE CRITERIA

DP2 is considered design-complete when the following are true:

1. Every initial industry can be represented using shared Region/Site/Facility/System/Component primitives.
2. Facilities can express hard dependencies, capacity constraints and alternative paths.
3. Failures are causal and diagnosable.
4. Repairs/replacements/bypasses are represented consistently.
5. Production recipes consume/produce shared resource entities.
6. Storage/logistics constraints can limit production.
7. Skills and managers modify diagnosis/operation without becoming universal power multipliers.
8. A stuck player has a non-cash recovery route through salvage, trade, emergency supply and Work Exchange.
9. Company accounting can calculate profit, ROI and ROC/ROIC.
10. Shared economy state is server-authoritative.
11. Standardized challenges can run deterministically.
12. Membership cannot directly create tradeable economic power.
13. The opening aggregate-plant vertical slice can be implemented without inventing special-case architecture.
14. Later stock/commodity systems can be added without collapsing Player and Company identity.
15. All competitive scoring can be versioned and explained.
16. Economic/anti-cheat telemetry exists in the design.

---

# PART XXIII — DEFERRED DP3+ DECISIONS

The following are intentionally not fixed in DP2:

- exact client game engine/framework;
- exact backend language/framework;
- exact database technology;
- exact art style;
- exact world-map topology;
- exact resource counts/recipe numbers;
- exact manager rarity probabilities;
- exact skill XP curves;
- exact timer durations;
- exact market fees;
- exact credit denomination;
- exact leaderboard formulas/weights;
- exact stock-exchange order-book mechanics;
- exact bond/futures implementation;
- exact company insolvency rules;
- exact consortium governance model;
- exact advertising implementation;
- exact membership price by store/region;
- exact server topology.

These belong in architecture decisions, balance specifications or implementation work packages after the core domain is accepted.

---

# PART XXIV — DP3 HANDOFF TARGET

DP3 should turn this systems specification into an **engineering-ready architecture and work-package ledger**.

Recommended DP3 outputs:

1. domain/entity schema;
2. simulation service boundaries;
3. authoritative-state matrix;
4. API/event contracts;
5. persistence model;
6. deterministic simulation contract;
7. vertical-slice UX/state flow;
8. test strategy;
9. security/anti-cheat architecture;
10. repository/runtime standards;
11. WP-001 onward engineering decomposition;
12. model routing per work package;
13. evidence/handback template;
14. acceptance gate sequence.

DP3 must preserve the controlling invariants from the Product Philosophy baseline and DP2.

---

## DP2 controlling summary

GRIDWORKS is built from a small set of reusable causal primitives:

> **Resources flow through facilities composed of dependent systems and components. Components wear, fail, are diagnosed and repaired. Players and managers learn by operating them. Companies own assets, earn and spend credits, trade with one another, contract for work and accumulate financial histories. Geography creates scarcity. Scarcity creates specialization. Specialization creates trade. Trade creates an economy. Cooperation and transparent performance metrics create competition without destruction.**

All later industries and financial systems are expressions of this shared model, not independent minigames.


---

# PART XXV — BUSINESS LIFECYCLE AND TURNAROUND MODEL

## 38. Business lifecycle

The core lifecycle is:

```
Opportunity
→ Due diligence
→ Acquire / Develop
→ Operate
→ Diagnose / Improve
→ Expand / Optimize
→ Hold / List / JV / Sell
→ Reinvest
```

A player may repeatedly re-enter the lifecycle in different industries without replaying the same starting puzzle.

### 38.1 Opportunity types

- **Greenfield** — undeveloped site/resource opportunity;
- **Going concern** — functioning operation with known history/cash flow;
- **Turnaround** — distressed or barely operational asset/business;
- **Player listing** — business/facility offered by another player;
- **Restructuring/tender** — later-stage insolvency or forced-sale opportunity.

### 38.2 Procedural business generation

Generated businesses must begin from a generated **history/state narrative** and derive present condition using the same rules engine that governs normal operations.

Possible history variables:

- facility age;
- operating strategy;
- maintenance quality;
- utilization history;
- investment history;
- management quality;
- staffing;
- recent incidents;
- utility reliability;
- logistics constraints;
- debt/capital structure later;
- market/contract situation.

The generator must not directly select arbitrary visible faults as the primary mechanism.

### 38.3 Due diligence

Before purchase, the player receives bounded information such as:

- estimated mechanical/electrical/process condition;
- known maintenance backlog;
- estimated production capability;
- historical revenue/profit where available;
- major contracts;
- major liabilities;
- asking price;
- confidence/uncertainty.

Relevant player/manager skills may improve precision or expose additional risks.

### 38.4 Persistence after sale

A sold facility/business remains a persistent world object unless the new owner later closes/dismantles it.

Preserve:

- ownership lineage;
- facility layout;
- installed components;
- maintenance history;
- manager/employment history where transferred;
- applicable contracts;
- operating statistics;
- major incidents;
- historical financials.

---

# PART XXVI — SOCIAL, IDENTITY AND MESSAGING

## 39. Identity hierarchy

```
Account
 └── PlayerProfile
      └── Ownership/Control
           └── CompanyGroup
                ├── Company
                ├── Company
                └── JV interests
```

Account identity is not an economic asset.

PlayerProfile is not a Company.

Company ownership can change without transferring the human account.

### 39.1 Player profile fields

Minimum domain fields:

- player_id;
- unique_handle;
- display_name;
- avatar_asset_id;
- bio;
- locale;
- timezone preference;
- reputation summary;
- achievements;
- privacy settings;
- DM permissions;
- notification preferences;
- block list;
- moderation state.

### 39.2 Company/group identity

Fields:

- company_id;
- legal/game name;
- parent_group_id;
- logo/emblem asset;
- short description;
- headquarters region;
- owner/share ledger;
- public/private state;
- ticker later;
- reputation;
- financial summary;
- consortium memberships.

Unique naming and impersonation rules are enforced server-side.

---

## 40. Messaging domain

Messaging is durable and server-owned.

Core entities:

- Conversation;
- ConversationParticipant;
- Message;
- MessageAttachment;
- ReadReceipt / last-read cursor;
- MessageReport;
- ChannelRole/Permission;
- ContextLink.

### 40.1 DM inbox

DMs are asynchronous inbox-style conversations.

Required behavior:

- conversation list;
- unread counts;
- archive;
- pin;
- mute;
- block;
- report;
- search;
- contextual initiation from another player/company/contract;
- durable history subject to moderation/privacy retention policy.

### 40.2 Consortium chat

Structured channels at launch:

- General;
- Projects;
- Trade;
- Management;
- Announcements.

Permissions depend on consortium role.

System-generated consortium/project events are visually differentiated from human messages.

### 40.3 Job/project/contract/JV threads

Economic objects may own dedicated discussion threads.

When the object closes, the thread becomes historical/read-only according to retention policy rather than disappearing.

### 40.4 Global chat

Not required for launch.

If introduced later, it requires separate moderation and anti-spam review.

---

## 41. Notification domain

Notification categories:

- Operations;
- Markets;
- Contracts;
- Messages;
- Consortium;
- Managers;
- Finance;
- Security.

Delivery channels:

- in-app;
- push;
- silent/background;
- none.

Per-category user preferences are mandatory.

---

# PART XXVII — INTERNATIONALIZATION AND LOCALIZATION

## 42. Semantic-content architecture

Simulation/domain objects store semantic IDs, not translated player-facing prose.

Example:

```
event_type = facility.production_stopped
reason_id  = fault.motor.bearing_seizure
component  = component.conveyor_motor
```

Localization renders these IDs into the selected language.

### 42.1 Message formatting

Support:

- plural categories;
- grammatical variants;
- parameter ordering;
- locale-aware date/time;
- locale-aware numeric/financial formatting;
- RTL;
- bidirectional text.

Use ICU MessageFormat semantics or equivalent.

### 42.2 Units

Simulation stores canonical SI values.

Presentation can render user/locale preferences, e.g.:

- tonnes vs short tons;
- km vs miles;
- °C vs °F;
- litres vs gallons.

Conversion happens at presentation boundaries, never by altering simulation state.

### 42.3 Language packs

Language packs may contain:

- strings;
- localized tutorials;
- terminology catalogue;
- locale fonts/fallback config;
- localized media where applicable.

They are:

- versioned;
- signed;
- cacheable;
- independently updateable;
- downloadable where size warrants.

### 42.4 Localization pipeline

Required stages:

1. canonical source string/content authored;
2. glossary/terminology resolution;
3. AI-assisted translation where appropriate;
4. automated structural validation;
5. human/native QA for critical content;
6. package/version publication;
7. client fallback validation.

CI must fail appropriately for malformed placeholders and required-locale omissions.

### 42.5 Chat translation

Optional message translation service:

- never replaces stored original;
- original always viewable;
- can be invoked manually or via preference;
- failure does not block chat;
- translation is not authoritative evidence of what sender wrote.

---

# PART XXVIII — DP3 STACK CONSTRAINTS NOW ACCEPTED

## 43. Technical choices carried into DP3

The following DP2-deferred choices are now accepted as DP3 working architecture:

| Layer | Decision |
|---|---|
| Game client | Godot 4.x |
| Shared deterministic rules engine | Rust |
| Backend application/services | Go |
| Primary transactional DB | PostgreSQL 18 |
| Local/offline store | SQLite |
| Messaging/event transport | NATS + JetStream |
| Cache | Valkey only when justified |
| Object storage | S3-compatible |
| Admin/back-office | React + TypeScript |
| External API | HTTPS/JSON initially |
| Realtime | WebSocket |
| Runtime | Native Linux/systemd |
| Repository | Monorepo |
| Search | PostgreSQL initially |
| Identity UX | Guest-first, account-link/protect later |
| Content definitions | Data-driven/versioned/signed |

### 43.1 One-rules-engine invariant

The Rust rules engine is the sole canonical implementation of simulation mathematics and deterministic world rules.

Consumers include:

- Godot client;
- server validator;
- challenge runner;
- test harness;
- balance simulator/tools.

Backend Go services orchestrate ownership, transactions, markets, social systems, persistence and authoritative state. They must not duplicate core simulation logic.

### 43.2 Ledger-first economy

Economic value transitions must be auditable from the first implementation.

Do not scatter untraceable `coins += x` mutations throughout services.

Use durable transaction/journal structures suitable for later:

- profit/loss;
- ROI;
- ROC/ROIC;
- inventory valuation;
- dividends;
- stock exchange;
- fraud investigation.

---

# PART XXIX — UPDATED DP2 ACCEPTANCE ADDITIONS

DP2/DP3 design continuity additionally requires:

17. Player, Account and Company remain distinct domains.
18. DM inbox, Consortium chat and object-linked project/job threads have durable domain models.
19. Notifications are user-configurable and not manipulative.
20. Businesses support greenfield, going-concern and turnaround acquisition paths.
21. Procedurally generated distressed assets derive from causal operating histories.
22. Sold businesses/facilities preserve ownership and operating history.
23. Localization is semantic-key based and language-neutral at the simulation layer.
24. RTL, Unicode, locale formatting and canonical-unit conversion are architectural requirements.
25. Global launch localization does not depend on an online LLM.
26. The accepted stack preserves the single Rust rules engine across client/server validation.
27. The economic ledger is auditable from initial implementation.


---

# PART XXX — REAL ESTATE, SYSTEM OWNERSHIP AND ECONOMIC AUTHORITIES

## 44. Real-estate domain foundation

Real estate is represented through the same shared asset/economic primitives as other GRIDWORKS systems.

Initial domain entities should include:

- LandParcel;
- Building;
- PropertyOwnership;
- PropertyListing;
- Lease;
- Occupancy;
- RenovationJob;
- PropertyCondition;
- PropertyValuation.

Initial property classes:

- industrial land;
- warehouse;
- factory/industrial building;
- office;
- basic residential property.

Later seasons may add richer residential, mixed-use, interior-design, conversion, hospitality/retail and advanced portfolio systems without changing the core ownership model.

### 44.1 GRIDWORKS system principal

Ownership must support a non-player system principal displayed to players as **GRIDWORKS**.

GRIDWORKS-owned property/business assets may have:

- fixed purchase price;
- lease terms;
- unavailable/reserved state;
- future-season release state.

Player-owned assets support owner identity, contact and negotiated offers.

### 44.2 Non-political authorities

Taxes/fees may be attributed to service abstractions such as Local Authority or Port Authority, but these entities have no political gameplay or governance mechanics.

Authority income is ledgered to non-player treasury principals.

### 44.3 Treasury recirculation

Treasury balances may fund bounded system actions, including:

- GRIDWORKS asset purchases;
- business/property buyback offers;
- system procurement;
- community projects;
- infrastructure funding;
- liquidity/recovery mechanisms.

All such flows must use the same economic ledger and be operator-auditable.

---

# PART XXXI — STARTING PATHS AND BUSINESS SUSPENSION

## 45. Starting-path selection

The common introductory lesson teaches the GRIDWORKS reasoning model, after which a player selects a preferred starting path.

The selected path controls the first economic scenario but does not lock later progression.

Candidate paths:

- Mining / Quarrying;
- Agriculture;
- Manufacturing;
- Energy;
- Logistics;
- Real Estate;
- Finance / Markets;
- Generalist / Surprise Me.

### 45.1 No forced operation

Any owned business may support lifecycle states including:

- Active;
- Mothballed/Suspended;
- Listed for Sale;
- Sold;
- Lease/Managed by Other later;
- Dismantled/Repurposed where applicable.

Mothballing stops normal production and materially reduces variable operating costs while preserving appropriate fixed holding/maintenance costs.

### 45.2 Early correction path

A bounded starter-business buyback/switch mechanism may let a new player change direction without account reset.

It must be limited and abuse-resistant.


### 44.4 GRIDWORKS liquidity acquisition and resale

Eligible player-owned businesses/properties may be sold directly to GRIDWORKS as a buyer of last resort when no acceptable player purchaser exists.

Required domain state:

- SystemPurchaseOffer;
- SystemAcquisition;
- SystemInventoryAsset;
- SystemResaleListing;
- CarryingCostLedger;
- PriorOwnerReference;
- ValuationSnapshot;
- ReacquisitionEligibility;
- ReacquisitionPriceFloor.

The GRIDWORKS acquisition offer is derived from a server-authoritative valuation model and settles immediately through the normal economic ledger if accepted.

The acquired asset preserves:

- ownership lineage;
- condition;
- installed equipment;
- operating history;
- maintenance history;
- financial history;
- contracts where legally/game-mechanically transferred;
- manager/employment relationships where included in the transaction.

GRIDWORKS normally relists the asset.

Normal resale floor:

> acquisition cost + applicable carrying/transaction costs

Former-owner reacquisition floor:

> max(current GRIDWORKS listing price, acquisition price × 1.20, acquisition price + carrying/transaction costs)

The 20% premium is an initial balance value and must be data-driven/versioned rather than hard-coded into business logic.

Scarcity may increase GRIDWORKS listing prices when equivalent inventory is unavailable. Scarcity rules must apply consistently to all buyers and be operator-auditable.

### 44.5 Anti-arbitrage requirements

The system must prevent repeated sale/reacquisition loops from acting as free financing or guaranteed arbitrage.

Controls may include:

- cooldowns;
- repeated-cycle limits;
- related-account checks;
- valuation haircuts;
- carrying costs;
- treasury liquidity constraints;
- full sale/reacquisition history;
- anomaly detection.

The design goal is **reliable exit liquidity, not risk-free temporary financing**.


---

# PART XXXII — ACCEPTED TIME AND SEASON MODEL

## 46. Time baseline

Initial operational time compression:

> **1 game day = 1 real hour**

All authoritative elapsed-time calculations use server-issued UTC timestamps.

The time ratio is versioned content/rules configuration and may not be silently changed by the client.

### 46.1 Clock domains

Maintain separate concepts for:

- authoritative real time;
- compressed operational simulation time;
- regional climate calendar.

### 46.2 Regional climate calendar

Climate follows regional real-calendar cycles rather than the compressed operational clock.

Season profiles may affect:

- resource demand;
- crop/biological processes;
- heating/cooling load;
- generation availability;
- water availability;
- logistics risk;
- construction productivity;
- property costs.

### 46.3 Offline-safe settlement

Offline catch-up must apply the same deterministic rules as online play while enforcing protective behavior such as safe shutdown, storage limits and bounded degradation.

No destructive outcome should occur solely because the user failed to log in during an artificially narrow window.

### 46.4 Timed jobs

Timed jobs must support data-driven duration and legitimate in-game modifiers.

Paid entitlements may improve planning/queue administration, but must not create an instant productive completion path unavailable through gameplay.

### 46.5 Release terminology

Product expansions are called Releases. Climate/world seasons remain separate simulation concepts.

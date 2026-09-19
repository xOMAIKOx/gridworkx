# GRIDWORKS — Product Philosophy & Game Design Baseline

**Repository:** `xOMAIKOx/gridworkx`  
**Status:** DESIGN BASELINE / CONTROLLING PRODUCT INTENT  
**Date:** 2026-09-19  
**Purpose:** Preserve the founding game philosophy and design decisions so they can be expanded into an implementation-grade application/game specification and then decomposed into engineering work orders.

---

## 1. Product thesis

GRIDWORKS is a persistent cooperative economic and systems-simulation game for mobile, designed around **building, repairing, producing, trading, financing, optimizing and collaborating rather than destroying other players**.

The game should feel like a living industrial/economic world. Players may specialize in agriculture, mining, energy, oil and gas, manufacturing, water, logistics, trade, engineering, finance or other interconnected disciplines. They can build businesses, help other players, trade resources and expertise, form joint ventures/consortia, invest in public companies, participate in commodity and financial markets, and compete through efficiency, profitability and engineering skill.

The core design ambition is:

> **Free players can play the entire game indefinitely and can ultimately obtain everything that affects gameplay. Paying players get more game, more convenience, more customization and faster access — but not superior power.**

This principle is not marketing copy. It is a controlling game-design constraint.

GRIDWORKS must deliberately reject the dominant mobile pattern of selling overwhelming competitive power, predatory loot boxes, artificial resource starvation, or $50–$100 packs required for meaningful progression.

The desired commercial relationship is:

> **Paying should remove friction from managing the game, enrich the experience, and fund continuing development — not remove the challenge of playing it.**

---

## 2. Foundational design laws

### 2.1 The world is the puzzle

GRIDWORKS should not be a collection of disconnected puzzle levels sitting beside a city builder.

The player's world itself is the puzzle.

The fundamental loop is:

**Observe → Diagnose → Decide → Act → Observe consequences → Adapt**

The same loop must work at different scales:

- first-session conveyor fault;
- a farm suffering poor yield;
- a mine with a processing bottleneck;
- a refinery constrained by utilities;
- an electrical grid with insufficient transmission;
- a national-scale logistics network;
- a public company making a capital-allocation decision.

A beginner may ask:

> Which component do I repair first?

A veteran consortium may ask:

> Which intervention restores the most economic output before a major contract deadline?

It is the same game mechanic at different scale.

### 2.2 Every system must behave like a system

Different game resources cannot simply be different-colored rocks.

Every material, facility, industry and logistics mode should create different dependencies, constraints, failure modes and optimization problems.

Complexity is justified only where it produces a decision.

> **If a distinction creates no meaningful decision, do not add it merely as another icon or inventory item.**

### 2.3 Every failure has a cause

Nothing should break merely because a hidden arbitrary timer says "broken."

Failures may be probabilistic, but they must emerge from understandable causes such as:

- wear;
- operating hours;
- excessive load;
- poor maintenance;
- bad environmental conditions;
- incorrect process settings;
- over-capacity operation;
- material properties;
- equipment age;
- insufficient staffing;
- upstream/downstream failures.

### 2.4 Every problem should permit multiple viable responses

Problems should rarely have one scripted answer.

A failed conveyor may be addressed by:

- repairing the motor;
- replacing the motor;
- repairing the belt;
- bypassing the line with trucks;
- rerouting material;
- outsourcing production;
- temporarily operating at reduced capacity;
- rebuilding the line.

The correct choice depends on the player's resources, time, contracts, skill, goals and risk appetite.

### 2.5 Mistakes matter but must remain recoverable

Players should be allowed to make bad decisions.

If a player spends scarce steel repairing a conveyor belt when the real immediate failure is a seized motor, the repair should not magically solve the problem.

However, the player must not be permanently trapped or manipulated into spending real money.

Recovery mechanisms include:

- dismantling/recovery of a percentage of materials;
- resale;
- player-to-player trade;
- small work contracts;
- salvage;
- cheaper temporary workarounds;
- emergency system-supplier purchases using earned game currency;
- assistance from other players.

### 2.6 Information exists before the decision

GRIDWORKS should reward diagnosis, not hidden gotchas.

A player should be able to inspect relevant evidence before acting.

Example:

**Motor**
- temperature: 112°C;
- rotation: 0 RPM;
- current draw: 184% rated;
- condition: critical.

**Belt**
- condition: 41%;
- operational: yes;
- service recommended: soon.

The game does not need to say "repair the motor first." The evidence should allow the player to infer it.

Optional advisory assistance may explain the reasoning, especially during onboarding, but skilled players must be able to disable or ignore hints.

### 2.7 Partial operation is valuable

Players should learn that restoring 2–5% useful throughput can be strategically superior to attempting a perfect 100% rebuild immediately.

GRIDWORKS should celebrate the moment a wrecked plant produces its first output.

Example:

- plant condition: 31%;
- available equipment: 43%;
- throughput: 0%.

After one intelligent intervention:

- plant condition: 34%;
- available equipment: 47%;
- throughput: 3%.

The plant is still a disaster, but production has started, revenue exists, and the player now has options.

---

## 3. First-session onboarding

The first login should include a tightly authored advisor-guided recovery task. Its purpose is not merely UI training. It must teach the player **how GRIDWORKS thinks**.

### 3.1 Opening premise

The player inherits a barely functioning industrial/agricultural facility.

Possible opening facility:

- aggregate plant;
- grain-processing operation;
- simple farm + pump + storage system.

The facility contains multiple damaged or suboptimal components, but only a small subset need intervention to restart useful output.

### 3.2 Advisor behavior

The advisor should guide one focused recovery task, not repair the whole facility.

Desired teaching sequence:

1. inspect before spending;
2. distinguish critical failure from non-critical wear;
3. reason about dependencies;
4. choose a minimal intervention;
5. restore 2–5% output;
6. observe revenue/material flow;
7. allow the player to decide what comes next.

Example advisor philosophy:

> "Before spending anything, inspect the line. We only need enough production to start earning."

When the player notices a worn belt:

> "The belt looks bad. But does it actually prevent the line from running?"

After diagnosing a seized motor:

> "The motor is the immediate failure. Repairing the belt first would improve reliability, but it would not restart production."

After startup:

> "It's ugly. It's inefficient. But it works. Now we have options."

The advisor may remain accessible later, but should stop nagging after the initial guided task.

Possible in-world role names:

- Foreman;
- Operator;
- Chief Engineer;
- Dispatcher.

---

## 4. Progressive complexity

GRIDWORKS can become a deep world simulation, but it must not expose all complexity in the first ten minutes.

The player should discover systems gradually through consequences.

Early gameplay might begin with:

- connect pump to field;
- restore conveyor;
- add power;
- store output;
- sell a small batch;
- fulfill a simple contract.

Later the same player may manage:

- soil chemistry;
- hydroponics;
- multiple power sources;
- rail logistics;
- commodity markets;
- balance sheets;
- public companies;
- futures;
- consortium infrastructure.

The design requirement is **progressive complexity**, not simplification of the endgame.

---

## 5. Core economic pillars

The current design baseline includes at least the following interoperable pillars.

### 5.1 Agriculture

Potential systems:

- field crops;
- horticulture;
- livestock where appropriate;
- greenhouse production;
- vertical farming;
- hydroponics;
- aquaponics;
- irrigation;
- soil quality;
- crop rotation;
- temperature and season;
- drainage;
- nutrients;
- pollination;
- storage;
- processing.

Real-world relationships should create gameplay.

Example: vertical farming may offer excellent yield per area but high energy demand. The player should discover this by operating the system rather than reading a lecture.

Circular relationships should be possible:

food waste → anaerobic digestion → biogas + fertilizer → energy + agriculture.

### 5.2 Mining and extraction

Potential resources include:

- iron ore;
- bauxite;
- copper;
- limestone;
- phosphate;
- coal;
- gold;
- silver;
- platinum-group metals;
- diamonds;
- uranium;
- lithium;
- nickel;
- cobalt;
- graphite;
- rare-earth elements such as neodymium, dysprosium and praseodymium;
- quarry products;
- forestry where appropriate.

Each resource family must produce different economic/engineering challenges.

Examples:

- iron ore: scale and bulk logistics;
- gold: low-volume/high-value ore processing and recovery;
- diamonds: geology, sorting and recovery;
- rare earths: separation complexity and waste-management challenge;
- lithium: hard-rock vs brine routes;
- coal: energy and metallurgical use cases.

Mining problems may include:

- drilling availability;
- blast cycle;
- haulage;
- crusher bottlenecks;
- conveyor capacity;
- ore grade;
- recovery rate;
- processing capacity;
- energy;
- water;
- tailings;
- rehabilitation.

### 5.3 Oil and gas

Oil and gas should be a complete industry, not a single resource node.

Potential chain:

exploration → seismic work → test well → field development → production → separation/processing → pipeline/LNG/refinery/petrochemicals.

Outputs may feed:

- transport fuels;
- aviation;
- bitumen;
- lubricants;
- petrochemical feedstocks;
- fertilizer/chemical chains;
- manufacturing.

Failure and optimization systems may include:

- reservoir pressure;
- pump performance;
- water cut;
- separators;
- compressors;
- pipeline capacity;
- processing bottlenecks.

### 5.4 Energy

Potential generation/storage systems:

- coal;
- oil;
- natural gas;
- nuclear;
- hydro;
- solar;
- wind;
- geothermal;
- biomass;
- biogas;
- batteries;
- pumped hydro;
- hydrogen later if justified.

Energy types should have real characteristics and trade-offs rather than moralized "good/bad" labels.

Examples:

- solar/wind are intermittent;
- gas may respond quickly;
- nuclear is capital-intensive but high-output;
- hydro is geography-dependent;
- batteries shift energy rather than produce it;
- transmission can constrain otherwise adequate generation.

### 5.5 Water

Potential systems:

- extraction;
- pumping;
- treatment;
- storage;
- desalination;
- irrigation;
- industrial water;
- wastewater;
- recycling;
- groundwater impacts.

Water should be a meaningful production constraint in agriculture, mining, industry and urban systems.

### 5.6 Industry and manufacturing

Potential chains:

- crushing;
- concentrating;
- smelting;
- refining;
- steel;
- chemicals;
- food processing;
- electronics;
- motors;
- industrial equipment;
- construction materials;
- advanced components.

The design should allow supply chains to emerge across many players rather than requiring every player to vertically integrate.

### 5.7 Logistics

Core modes:

- road;
- rail;
- water/sea;
- air;
- pipeline;
- conveyor/material handling.

Each mode must be economically different.

**Road:** flexible, excellent last-mile, higher long-distance bulk cost, congestion possible.

**Rail:** expensive infrastructure, strong for bulk and containers, efficient over distance.

**Water:** huge capacity, low cost per tonne, slow, geography/port dependent.

**Air:** fast, expensive, low bulk suitability, good for urgent/high-value cargo.

**Pipeline:** high capital cost, excellent continuous flow for oil/gas/water/chemicals.

**Conveyor/material handling:** short-distance, extremely high throughput for industrial/mining complexes.

A player should be able to specialize primarily in logistics.

### 5.8 Circular systems and environmental constraints

Sustainability should be represented as engineering/economic consequences, not moral points.

Potential systems:

- waste recovery;
- recycling;
- waste heat;
- biogas;
- mine rehabilitation;
- tailings management;
- water treatment;
- emissions;
- groundwater;
- land impacts;
- resource substitution.

Players should be able to discover profitable closed loops.

---

## 6. Geography and regional specialization

Player regions should differ in meaningful ways.

Example region A:

- excellent farmland;
- high rainfall;
- poor mineral resources;
- river transport;
- moderate wind.

Example region B:

- arid;
- excellent solar;
- copper;
- oil and gas;
- limited freshwater;
- long distance to port.

Example region C:

- mountainous;
- hydro potential;
- gold;
- forestry;
- expensive roads.

There should be no perfect starting region.

Geography creates:

- opportunity;
- constraint;
- trade;
- specialization;
- infrastructure decisions.

The game should allow emergent regional identities rather than hard-coding "farming zone" or "mining zone" roles.

---

## 7. Player identity and specialization

A mature player should not necessarily be "Level 97 with everything."

Possible identities include:

- agricultural producer;
- mine operator;
- power producer;
- refinery operator;
- manufacturer;
- water operator;
- logistics company;
- port operator;
- engineer/maintenance specialist;
- commodity trader;
- investor;
- financier;
- distressed-asset specialist;
- consortium organizer;
- challenge competitor.

GRIDWORKS should support multiple legitimate definitions of success.

---

## 8. Production, trade and recovery

Players should be able to obtain needed resources primarily through three routes:

1. **Produce it** — own/operate the infrastructure.
2. **Trade for it** — buy from other players or contracts.
3. **Work for it** — help another player and earn credits/resources.

Real-money purchase should not be required as a fourth economic route.

### 8.1 Player market

Player trade is fundamental.

Players should be able to specialize and rely on others.

Example:

- copper mine sells concentrate/ore;
- another player smelts/refines;
- manufacturer buys copper and produces motors;
- wind-turbine company buys motors and steel;
- logistics players move the goods.

### 8.2 Emergency/system supplier

If a system resource store exists, it should act as a safety net and price stabilizer, not replace the player economy.

Example:

- player-market steel: ~92–108 credits;
- emergency supplier: 140 credits;
- daily limits.

This prevents deadlock while preserving incentives to trade.

---

## 9. Work Exchange and player assistance

A player who made mistakes or lacks resources should be able to recover by helping others.

Possible jobs:

- clear maintenance backlog;
- help construction;
- inspect facility;
- harvest assistance;
- freight loading;
- delivery legs;
- plant commissioning;
- consortium project contribution.

The helper contributes **time/attention/skill**, then earns:

- credits;
- materials;
- reputation;
- skill progression.

### 9.1 Construction assistance

Player assistance can reduce another player's build/repair time within strict caps.

Example:

Base build time: 8h

- player A: -45m;
- player B: -30m;
- player C: -45m.

Remaining: 6h.

The design should cap assistance so mass participation cannot collapse an 8-hour build into seconds. A maximum reduction around 25–40% may be explored and balanced later.

### 9.2 Skills matter

Assistance value should depend on relevant player skill.

Examples:

- electrical;
- mechanical;
- agriculture;
- mining;
- logistics;
- civil works;
- process engineering.

This allows a player with modest assets to become economically valuable through expertise.

---

## 10. Player skills

Players themselves gain persistent skills through doing.

Examples:

- mining;
- agriculture;
- engineering;
- maintenance;
- logistics;
- commerce;
- energy;
- water;
- finance.

Skill should improve through relevant activity:

- helping;
- producing;
- diagnosing;
- repairing;
- operating facilities;
- trading;
- completing contracts.

Player skill and manager skill are separate systems.

---

## 11. Facility managers

Facility managers are a major progression and identity system.

Possible manager types:

- mine manager;
- farm manager;
- refinery manager;
- power-station manager;
- logistics manager;
- port manager;
- manufacturing manager;
- water manager;
- oilfield manager.

### 11.1 Rarity and potential

Working rarity model:

- Bronze;
- Silver;
- Gold;
- Platinum.

Rarity should represent potential, trait combinations and specialization — not a simple production multiplier.

A veteran Gold manager must be able to outperform a newly recruited Platinum manager in relevant circumstances.

Example:

**Gold Manager — Level 42**
- Mining Operations 86
- Maintenance 71
- Safety 64
- Logistics 53

**Platinum Manager — Level 1**
- Mining Operations 37
- Maintenance 31
- Safety 28
- Logistics 42
- Potential: Exceptional

### 11.2 Manager development

Managers gain experience by operating facilities.

Potential dimensions:

- technical skill;
- maintenance;
- safety;
- leadership;
- logistics;
- energy efficiency;
- crisis response;
- mentoring.

Managers should develop histories and become meaningful world entities.

### 11.3 Traits with trade-offs

Avoid universally dominant managers.

Possible traits:

**Aggressive Operator**
- improved throughput when healthy;
- shorter maintenance interval / higher wear.

**Maintenance First**
- reliability improvement;
- slightly lower peak throughput.

**Cost Cutter**
- lower operating expense;
- higher deferred-maintenance risk.

**Mentor**
- junior managers gain experience faster;
- less direct operating benefit.

**Crisis Specialist**
- excellent failure recovery;
- limited normal-operation bonus.

### 11.4 Managers improve information, not only numbers

A low-skill manager might report:

> Conveyor problem likely caused by belt wear — confidence 63%.

A highly experienced manager might report:

> Belt wear is secondary. Motor current suggests bearing seizure — confidence 92%.

This directly supports the central Observe → Diagnose loop.

### 11.5 Fatigue/workload

Managers may have:

- workload;
- morale;
- fatigue;
- specialization;
- experience.

Fatigue must not become an energy/stamina paywall.

Recovery occurs through management:

- rest;
- deputies;
- rotation;
- staffing;
- automation;
- improved processes.

### 11.6 Daily recruitment

Current baseline:

- free player: **5 recruitment opportunities per day**;
- paid member: **more daily opportunities**, initial working target around **8/day**;
- probabilities remain the same for free and paid players;
- Platinum is available to both;
- recruitment draws cannot be purchased separately with real money.

A recruitment opportunity should reveal a candidate rather than automatically dump permanent staff into an infinite card inventory.

Players may:

- recruit;
- reject;
- shortlist;
- potentially refer/trade.

### 11.7 Platinum

Platinum must not be directly purchasable.

Two acquisition paths are desirable:

1. luck through normal/earned recruitment;
2. long-term mastery/reputation achievements that improve access to elite recruitment without making Platinum a cash purchase.

### 11.8 Manager market

Manager trading is desirable, but should be framed as employment/contract transfer rather than literal ownership of people.

Potential systems:

- permanent recruitment transfer;
- temporary management contract/lease;
- consortium assignment;
- apprenticeship/deputy arrangement.

A manager's employment history should travel with them and create world stories.

---

## 12. Time and construction

Construction time should matter.

Illustrative scale:

- small repair: seconds/minutes;
- warehouse: minutes;
- factory: tens of minutes/hours;
- major power plant: hours;
- megaproject: days.

The game must not weaponize timers.

Examples of forbidden design:

- 4d 17h timer intended to force a $39.99 skip;
- sleep interruption because a harvest will otherwise disappear;
- punishment for failing to log in at exact times.

When storage fills, production may pause rather than destroy the player's effort.

### 12.1 Acceleration

Earned speedups may exist.

Paid membership may provide convenience, but acceleration must be bounded so it cannot create overwhelming compounding economic power.

Player assistance can also reduce build time within a cap.

Real money must never create effectively unlimited construction throughput.

---

## 13. Consortiums, cooperatives and joint ventures

The social structure should emphasize cooperation rather than clans built around combat.

Working terminology:

- Consortium;
- Cooperative;
- Venture;
- Enterprise;
- Alliance.

"Consortium" currently fits best.

Players should be able to contribute different resources/capabilities to shared projects.

Example:

**Regional Hydro Project**

Requires:

- concrete;
- steel;
- turbines/components;
- logistics;
- construction effort;
- capital.

Different players contribute according to specialization.

Ownership/reward can be proportional.

### 13.1 Joint ventures

Players may create joint ventures around:

- mines;
- farms;
- power projects;
- rail;
- ports;
- refineries;
- major factories.

A JV might record proportional ownership and distribute revenue accordingly.

### 13.2 Global/community megaprojects

Examples:

- transcontinental rail restoration;
- ocean cleanup;
- desert reclamation;
- regional hydro;
- large water systems;
- grid interconnects.

Community effort should visibly alter the shared world and may improve logistics/economics for many players.

---

## 14. Competition without destruction

Another player must not be able to burn down, raid or destroy a player's work while they are offline.

Competition should occur through:

- efficiency;
- profitability;
- engineering;
- contracts;
- reputation;
- market performance;
- standardized challenges;
- leaderboards.

### 14.1 Standardized challenges

Weekend/seasonal challenges should remove accumulated wealth advantage.

All players receive identical:

- map;
- capital;
- technology;
- staff/managers or standardized equivalents;
- starting resources;
- time window.

Examples:

- maximum sustainable food output with a 5 MW cap;
- 50 MW generation in minimum land area;
- supply a town with zero fossil generation;
- maintain farm output under 40% water reduction.

This is skill-based PvP through engineering.

---

## 15. Leaderboards and scoring

GRIDWORKS should have multiple leaderboards rather than one global "biggest wins" board.

Possible categories:

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
- Consortium;
- Finance.

Possible time horizons:

- all-time;
- season;
- monthly;
- weekly challenge.

### 15.1 Raw output is not enough

Leaderboard design must avoid "most tonnes wins" as the dominant criterion.

Measures should include combinations of:

- net profit;
- ROI;
- ROC/ROIC;
- ROE where applicable;
- profit margin;
- capital efficiency;
- yield;
- cost per unit;
- energy efficiency;
- water efficiency;
- recovery rate;
- asset utilization;
- reliability;
- downtime;
- contract fulfilment;
- sustainable operation;
- risk-adjusted return.

Examples:

**Agriculture**
- net profit;
- ROI/ROIC;
- yield per hectare;
- profit per hectare;
- water efficiency;
- energy efficiency;
- contract reliability.

**Mining**
- net profit;
- ROI/ROIC;
- recovery rate;
- cost per tonne;
- utilization;
- energy per tonne;
- downtime;
- rehabilitation performance.

**Energy**
- profit;
- ROI/ROIC;
- cost per MWh;
- availability;
- reliability;
- capacity factor;
- storage efficiency.

**Logistics**
- profit per tonne-km;
- on-time delivery;
- asset utilization;
- empty-running percentage;
- cost per delivery;
- network efficiency.

### 15.2 Transparent scoring

Composite indices may exist, but must be auditable.

Example:

Agriculture Operating Index:

- Profitability: 91
- Yield efficiency: 84
- Water efficiency: 96
- Energy efficiency: 72
- Contract reliability: 99
- Asset utilization: 88
- Overall: 88.6

The player must understand how to improve.

Avoid opaque "Power Score 8,734,222" systems.

### 15.3 Scale and efficiency must both be recognized

Have separate rankings for:

- absolute net profit;
- return on capital;
- total generation;
- profit per MW;
- total mined tonnes;
- cost per recovered tonne.

Large companies can be recognized for scale; brilliant small operators can be recognized for efficiency.

### 15.4 Divisions

Players should compete against peers as well as globally.

Possible business-scale divisions:

- Local;
- Regional;
- National;
- Continental;
- Global.

This keeps leaderboards accessible to new players.

### 15.5 Rewards

Top 5/10/100 rewards should emphasize prestige rather than permanent productive advantage.

Possible rewards:

- seasonal company emblems;
- animated facility banners;
- trophies;
- profile/company badges;
- cosmetic facility skins;
- titles;
- limited earned currency;
- elite recruitment opportunities.

Do not award permanent superior equipment that makes previous winners more likely to win again.

---

## 16. Reputation

Reputation must remain separate from wealth and technical skill.

Reputation can represent:

- contract completion;
- delivery reliability;
- fair trading;
- assistance history;
- project participation;
- default history.

A modest player with 99.8% contract completion may be a more attractive partner than a very wealthy player with poor reliability.

Reputation should carry genuine economic value in:

- contracts;
- hiring;
- JVs;
- lending;
- investment;
- consortium selection.

---

## 17. Player contracts

Beyond spot-market trading, players should be able to issue contracts.

Example:

**Request for Supply**
- 10,000 steel;
- 1,500/day delivery;
- seven-day term;
- 130 credits/unit.

A player may accept and potentially subcontract.

Possible contract types:

- supply;
- logistics;
- maintenance;
- construction;
- engineering;
- facility management;
- financing later.

Contracts create a more realistic economy and allow specialization.

---

## 18. Companies and corporate structure

The player should operate through companies rather than merely an abstract account.

Possible progression:

- one initial private company;
- multiple divisions/operations later;
- possibly multiple companies;
- consortium/JV ownership.

Examples:

- Maik Farms;
- Maik Energy;
- Maik Logistics.

The player remains distinct from the company.

This supports investment, public listings, acquisitions, debt and corporate history.

---

## 19. Finance module

Finance is not a side minigame. It can become an entire profession within GRIDWORKS.

A financially minded player should be able to enjoy the game with little or no direct industrial operation.

Possible roles:

- equity investor;
- bond investor;
- commodity trader;
- analyst;
- project financier;
- distressed-asset investor;
- capital allocator.

### 19.1 Financial metrics

Track meaningful financial performance, potentially including:

- revenue;
- operating profit;
- net profit;
- free cash flow;
- ROI;
- ROC/ROIC;
- ROE;
- margins;
- asset turnover;
- debt-to-equity;
- dividend payout;
- dividend yield;
- revenue/profit growth.

### 19.2 Public companies and stock exchange

Private player companies may later qualify to list on the in-game GRIDWORKS Exchange.

Possible listing requirements:

- operating history;
- minimum size/value;
- financial statements;
- reputation/default thresholds;
- listing standards.

Example IPO:

**Maik Resources PLC**
- valuation: 80M credits;
- shares offered: 20%;
- IPO price: 8 credits/share;
- capital raised: 16M credits.

The owner trades dilution for growth capital.

### 19.3 Shareholders and dividends

Other players can buy shares in listed player companies.

They may receive dividends when declared and benefit/lose from company performance.

Investment decisions should be based on visible company information such as:

- revenue;
- profit;
- ROIC;
- debt;
- free cash flow;
- reserves;
- management performance;
- contract reliability;
- commodity exposure.

### 19.4 Companies can fail

Financial gameplay requires real risk.

Poor decisions may cause:

- losses;
- dividend cuts;
- share-price declines;
- restructuring;
- asset sales;
- takeover;
- insolvency/reorganization mechanisms.

Failure must affect the business without deleting the player's GRIDWORKS account.

### 19.5 Commodities

Physical commodity markets should reflect player supply/demand where feasible.

Potential markets:

- copper;
- iron ore;
- gold;
- oil;
- gas;
- coal;
- wheat;
- corn;
- fertilizer;
- electricity;
- lithium;
- rare-earth products;
- steel.

Price changes should propagate through supply chains.

Example: rapid solar construction increases demand for copper/materials, raising input costs and stimulating new production.

### 19.6 Futures and hedging

Advanced players may eventually use futures/forward-style contracts.

Example:

- farmer locks wheat delivery price 60 days ahead;
- buyer hedges supply risk;
- speculator takes the other side.

This can teach hedging and price-risk management through gameplay.

It must be introduced only at advanced stages.

### 19.7 Debt and bonds

Companies may finance projects with debt rather than only equity.

Potential player-issued corporate bonds:

- principal;
- maturity;
- coupon;
- credit/risk assessment.

Strong companies can borrow more cheaply than highly leveraged speculative operators.

### 19.8 Acquisitions

Later-stage systems may support:

- full acquisitions;
- majority stakes;
- minority stakes;
- distressed purchases;
- JVs as alternatives.

### 19.9 Fictional crypto/digital assets

A fictional crypto/digital-asset module may exist later if it creates meaningful economic gameplay.

It should not be added merely because "crypto is cool."

Hard boundary:

> **No conversion to/from real money and no real-money cash-out.**

Conventional corporate/commodity markets should be established before any crypto system.

---

## 20. Fair-play monetization constitution

This section is controlling.

### 20.1 What real money must never buy

Real money must never directly purchase:

- superior competitive equipment;
- permanent production multipliers;
- tradable commodities/resources;
- shares;
- bonds;
- manager contracts that convey economic power;
- game currency used to acquire tradable productive assets;
- paid random manager loot boxes;
- unlimited construction speed;
- guaranteed Platinum managers.

This prevents real-money advantage from leaking into the player economy.

### 20.2 Free-play guarantee

A free player must ultimately be able to:

- access the complete core game indefinitely;
- reach maximum gameplay skill;
- obtain all gameplay-affecting equipment;
- recruit Platinum managers;
- trade;
- join Consortia;
- compete in challenges;
- participate in markets;
- become economically successful;
- appear at the top of skill/efficiency leaderboards.

### 20.3 Membership

Working target: approximately **$35/year**.

Membership should be deliberately good value.

Possible benefits:

- additional daily manager recruitment opportunities (e.g. 8 vs 5);
- same rarity probabilities as free players;
- expanded blueprint storage;
- advanced analytics/history;
- more market watchlists;
- more alerts;
- larger manager shortlist;
- cosmetic/company branding;
- additional content/early access;
- planning/automation convenience.

Avoid giving members multiplicative economic throughput.

Example:

Members may be allowed to queue future actions for automatic start, but not run three factories/construction crews in parallel solely because they paid.

### 20.4 Manager recruitment and money

Recruitment opportunities cannot be purchased separately.

No:

- 10 recruitment packs for $9.99;
- 2% Platinum paid rolls;
- cash-only managers.

Membership may provide more daily opportunities, but the odds remain the same.

### 20.5 Advertising

If used:

- no forced advertising;
- optional rewarded ads only;
- possibly a low-cost permanent ad-removal purchase.

Advertising must not become the game loop.

### 20.6 Small paid packs/content

Small purchases may exist, likely within roughly $0.99–$9.99, when they provide:

- cosmetics;
- content;
- campaigns;
- expansion modules;
- convenience that does not create durable competitive superiority.

The game should not depend on whale economics.

---

## 21. Device/server architecture philosophy

The game should perform the majority of simulation and rendering on the player's device where practical.

Device responsibilities may include:

- simulation;
- rendering;
- local physics/pathing;
- local saves/cache;
- offline play;
- puzzle generation from deterministic seeds.

Server responsibilities may include:

- identity/accounts;
- authoritative player state;
- cloud sync;
- purchases;
- markets;
- contracts;
- leaderboards;
- Consortiums;
- challenge seeds;
- anti-cheat validation;
- social systems;
- financial exchange state.

A standardized challenge can distribute a deterministic seed; the device simulates locally and submits compact replay/result evidence for validation.

This reduces server compute while retaining a shared persistent economy.

Client install size is expected to be driven largely by art/audio/engine assets, not simulation logic. Asset packs/regions may be downloadable.

---

## 22. Anti-cheat and market integrity

Because GRIDWORKS contains persistent player markets, abuse must be treated as a product-system problem from the beginning.

Potential abuse:

- wash trading;
- multi-account manipulation;
- fake transactions;
- collusion;
- pump-and-dump behavior;
- cornering thin markets;
- exploit-driven resource creation;
- leaderboard cheating.

Possible safeguards:

- transaction history;
- server-side authoritative market settlement;
- major-holder disclosure;
- trading fees;
- listing standards;
- circuit breakers;
- anti-wash-trade detection;
- related-account analysis;
- position limits where justified;
- replay/seed validation for competitive challenges;
- economic anomaly detection.

No real-money cash-out is part of the risk boundary.

---

## 23. Learning through consequences

GRIDWORKS should teach without becoming educational software.

Examples of knowledge players may absorb naturally:

- why vertical farming has high energy demand;
- why storage cannot replace generation;
- why transmission limits matter;
- why rail/water transport suits bulk goods;
- why air freight is poor for iron ore;
- why bottlenecks dominate throughput;
- why preventive maintenance can outperform reactive repair;
- why capital efficiency matters;
- why diversification and debt create trade-offs;
- why commodity hedging exists;
- why closed-loop systems can be economically useful.

The player should learn because the world responds coherently to their decisions.

---

## 24. Tone and emotional experience

GRIDWORKS should feel:

- intelligent;
- satisfying;
- constructive;
- curious;
- deep but learnable;
- respectful of the player's time;
- socially cooperative;
- economically alive.

It should not feel:

- childish;
- manipulative;
- hyper-casual;
- "tap +500 coins";
- combat-centric;
- whale-driven;
- punitive toward players who sleep/work/live offline.

The emotional reward is:

> **"I made this system work."**

---

## 25. What GRIDWORKS is not

GRIDWORKS is not:

- a pay-to-win PvP game;
- a raid/war/clan game;
- a loot-box economy;
- a simple idle clicker;
- a reskinned resource-generator collection;
- a city builder where every player follows the same optimal path;
- a financial service involving real-money investment or cash-out;
- an educational lecture disguised as a game.

---

## 26. Current high-level module map

1. Core Simulation / Time / World
2. Player Identity / Company
3. Facilities / Components / Maintenance
4. Skills / Experience
5. Managers / Recruitment / Contracts
6. Agriculture
7. Mining / Extraction
8. Oil & Gas
9. Energy
10. Water
11. Manufacturing / Processing
12. Logistics
13. Circular Economy / Environmental Systems
14. Trading / Marketplace
15. Contracts / Work Exchange
16. Consortiums / JVs
17. Community Megaprojects
18. Leaderboards / Challenges / Seasons
19. Reputation
20. Finance / Accounting
21. Stock Exchange
22. Commodities
23. Debt / Bonds
24. Advanced Financial Markets
25. Social / Notifications
26. Monetization / Membership / Cosmetics
27. Anti-Cheat / Market Integrity
28. Backend Sync / Persistence / Authoritative Services
29. Offline / Deterministic Simulation
30. Analytics / Telemetry / Balance Operations

This is a product-domain map, not yet an engineering decomposition.

---

## 27. Design governance

This file captures **product intent and philosophy**. It is not yet the final engineering specification.

Subsequent design work should turn this baseline into implementation-grade documents covering at least:

- core loop;
- simulation model;
- world model;
- facility/component model;
- resource ontology;
- production recipes and dependencies;
- failure/maintenance model;
- skill progression;
- manager system;
- economy and resource sinks;
- market mechanics;
- contract mechanics;
- Consortium/JV mechanics;
- finance/accounting;
- public-company/exchange mechanics;
- leaderboards/scoring;
- membership/monetization;
- anti-cheat;
- offline/server authority boundaries;
- UX/onboarding;
- live-ops and seasonal challenge model;
- data model;
- APIs/events;
- observability;
- security;
- testing;
- balancing methodology.

### 27.1 Architecture/design vs engineering

The design/architecture role should retain control of:

- product philosophy;
- invariants;
- system boundaries;
- work-package decomposition;
- acceptance criteria;
- cross-module consistency;
- architecture review;
- final acceptance/rejection of engineering handbacks.

Engineering should receive bounded, detailed work orders and return evidence through GitHub.

### 27.2 Model routing baseline

The expected engineering workflow may use different reasoning tiers depending on task complexity.

Working model:

- **Architecture / specification / difficult review:** GPT-5.6 Sol at high/extra-high reasoning, or the strongest available architecture/review model.
- **Well-specified implementation:** GPT-5.6 Luna XHigh where the problem is already solved in the work order.
- **Investigative/non-convergent implementation:** Luna Max or another stronger autonomous implementation model.
- **Review and acceptance:** return to Sol-class architecture/review reasoning.

Model choice is subordinate to task characteristics. Do not use a more expensive/reasoning-heavy model merely by habit, and do not use a cheaper worker when unresolved architecture is being delegated accidentally.

### 27.3 GitHub as collaboration channel

GRIDWORKS should use the repository as the durable collaboration record between Design/Architecture and Engineering.

Expected pattern:

1. Design/Architecture records controlling specification and decisions in GitHub.
2. Architecture issues a bounded work order.
3. Engineering implements on a dedicated branch/PR.
4. Engineering posts a structured handback with:
   - commit SHA;
   - branch;
   - files changed;
   - tests/evidence;
   - known deviations;
   - unresolved questions/blockers.
5. Architecture reviews the exact commit/PR and records PASS / REQUEST_CHANGES / HOLD.
6. The owner receives only a concise status/decision summary rather than acting as a manual message courier.

This avoids losing design context in long chat transcripts.

---

## 28. Design invariants to protect during future specification

The following must not be diluted casually:

1. Free players can ultimately obtain all gameplay-affecting capability.
2. Paying improves experience, convenience and content, not fundamental power.
3. No direct real-money purchase of tradable economic power.
4. No paid random manager loot boxes.
5. No player destruction/raiding as a core mechanic.
6. Competition is through skill, efficiency, profitability, markets, reputation and standardized challenges.
7. Every system has meaningful causal behavior.
8. Problems permit multiple viable solutions.
9. Mistakes cost something but remain recoverable.
10. Information exists before meaningful decisions.
11. Partial recovery/operation is strategically valid.
12. Player trade and specialization are foundational.
13. Helping other players is a legitimate progression path.
14. Skills improve through doing/helping/operating.
15. Managers develop through experience and are not simple rarity multipliers.
16. Raw output alone should rarely determine leaderboard success.
17. ROI and ROC/ROIC are first-class economic metrics.
18. Finance can be a complete player profession.
19. No real-money cash-out from the simulated financial economy.
20. Complexity is progressively revealed.
21. The device should carry simulation work where safe/practical; servers preserve shared authority.
22. The world should reward understanding, not arbitrary hidden rules.
23. GRIDWORKS should respect the player's time.
24. The core emotional reward is making interconnected systems work.

---

## 29. Immediate next design phase

The next phase should not jump directly to implementation.

Architecture/design should produce the detailed specification in a controlled order, likely beginning with:

1. **World + simulation clock + persistence model**
2. **Core facility/component/dependency/failure model**
3. **Resource ontology and production-chain rules**
4. **Player/company/skill/manager domain model**
5. **Economy, trade, contracts and anti-inflation mechanisms**
6. **First-session onboarding vertical slice**
7. **Industry modules using shared primitives**
8. **Consortium/JV/social systems**
9. **Leaderboard/scoring framework**
10. **Finance and exchange module**
11. **Monetization/membership implementation rules**
12. **Client/server/offline authority model**
13. **Engineering decomposition and acceptance gates**

A vertical slice should prove the philosophy before the full world is built:

> broken facility → inspect → diagnose → partial repair → production → earn → trade/help → improve skill/manager → make a second meaningful decision.

If that loop is not satisfying, adding thirty industries will not save the game.

---

## 30. Working product statement

> **GRIDWORKS is a persistent cooperative world economy where players restore and build interconnected systems by diagnosing problems, producing resources, trading, investing, transporting, engineering and collaborating. The world rewards understanding and efficient operation rather than destructive PvP or spending power. Players can specialize in physical industries or finance, compete through transparent operational and financial metrics, and participate fully without paying for superior power.**

This statement, together with the design invariants above, is the current controlling product intent.


---

## 31. Business lifecycle, acquisition and turnaround gameplay

GRIDWORKS must support business creation, improvement, sale, acquisition and persistent ownership history as first-class gameplay.

A player who has restored and optimized one facility must not be forced into a rinse-and-repeat copy of the same starting puzzle. After selling a successful operation, the player may choose among:

- **Greenfield development** — build a new operation from an undeveloped site;
- **Going-concern acquisition** — buy an operating business with known cash flow and identifiable optimization opportunities;
- **Turnaround acquisition** — acquire a distressed or barely functioning business with uncertain condition and potentially high upside.

Turnaround gameplay is a legitimate specialization in its own right:

> **Buy broken businesses → diagnose them → repair/restructure them → improve profitability and capital efficiency → retain or sell.**

Procedurally generated opportunities must be generated from a coherent operating history, not by randomly coloring components red. Example histories may include deferred maintenance, over-aggressive throughput, inadequate capital investment, utility constraints, poor logistics design, excessive leverage, weak management or a recent incident. The resulting facility state must be derivable from the same causal rules used for player-operated assets.

Acquisition should include imperfect but fair due diligence. Skill, managers and specialist expertise may improve what the player can infer before purchase.

Businesses sold by players should persist in the world with their facility configuration, asset history, maintenance record, applicable contracts, managers if included, ownership lineage and operating history.

Later transaction structures may include:

- asset sale;
- operating-business sale;
- full company acquisition;
- majority/minority stake;
- JV;
- restructuring/insolvency sale.

The Opportunity Board should combine system-generated and player-listed opportunities. Infinite free rerolling of generated opportunities is prohibited because it would allow players to cherry-pick favorable seeds without economic cost.

---

## 32. Player identity, company identity and social communication

GRIDWORKS must distinguish:

1. **Account** — authentication/security identity;
2. **Player Profile** — the human's persistent public game identity;
3. **Company / Group** — the economic entities the player owns or controls.

A Player Profile should support:

- immutable internal player ID;
- globally unique handle;
- non-unique display name;
- avatar/profile image or approved game-generated avatar;
- biography;
- language/locale;
- reputation summaries;
- achievements;
- privacy controls;
- block/mute/report controls;
- notification preferences.

Company/group identities should support unique names, logos/emblems, profile pages and later public-market identity/tickers. Company identity is separate from player identity because companies can be sold, shared, invested in or acquired.

### 32.1 Messaging model

Launch architecture should include three persistent communication modes:

- **Private DM inbox** — asynchronous conversation threads with unread state, pin/archive/mute/block/report/search;
- **Consortium chat** — persistent role-aware channels such as General, Projects, Trade, Management and Announcements;
- **Object-linked job/project/contract/JV threads** — contextual discussion attached to the actual economic object and archived with its history.

Important system events may appear in contextual threads as visually distinct system messages.

Unrestricted global chat is not a launch requirement. Communication should initially be relationship/context driven to reduce spam, scams and moderation load.

### 32.2 Notifications

A unified notification center should support categories such as:

- operations;
- markets;
- contracts;
- messages;
- consortium activity;
- managers;
- finance;
- security/account.

Players control whether each category is:

- in-game only;
- push;
- silent;
- disabled.

Push notifications must not use manipulative guilt or manufactured urgency.

### 32.3 Social safety

Minimum launch controls:

- block;
- mute;
- report;
- DM permissions;
- spam/rate limits;
- moderation audit trail;
- banned-name/impersonation controls;
- safe attachment policy;
- restricted links initially;
- admin review tooling.

---

## 33. Global language and localization philosophy

Global-language support is a **launch architecture requirement**, not a post-launch translation task.

The authoritative simulation and economic state must be language-neutral. The rules engine stores stable semantic IDs such as:

- `resource.copper_ore`;
- `fault.motor_bearing_seizure`;
- `contract.supply`;
- `event.production_stopped`.

The client renders those concepts using the player's locale.

### 33.1 Localization invariants

1. No authoritative gameplay state stores translated prose as the source of truth.
2. All system/player-facing game text uses stable localization keys.
3. Pluralization and grammar must use ICU-grade message semantics or equivalent.
4. Numbers, dates, percentages and financial values must be locale aware.
5. Simulation uses canonical units; presentation may convert to locale/user-preferred units.
6. RTL layout capability is designed from the beginning.
7. Unicode names/search/moderation are first-class requirements.
8. Fonts use locale-appropriate fallback chains and downloadable packs where useful.
9. Language packs are versioned independently from simulation rules.
10. CI detects missing keys, broken placeholders and malformed translations.
11. A controlled terminology/glossary system governs specialist industrial and financial vocabulary.
12. AI-assisted translation may accelerate production, but critical onboarding, financial, legal and monetization content receives human/native QA.
13. Chat translation is optional and must always preserve access to the original message.
14. English fallback always exists.
15. Localization cannot depend on a live LLM/API for normal game operation.

Potential launch locale set should target approximately 15–20 high-value languages with high quality rather than dozens of poor translations. Working candidates include English, Spanish, Portuguese (Brazil), French, German, Italian, Polish, Turkish, Bulgarian, Arabic, Japanese, Korean, Simplified Chinese, Traditional Chinese, Indonesian, Vietnamese and Thai, with additional languages selected from market validation.

---

## 34. Technical direction accepted for DP3

The following architecture direction is accepted as the working DP3 baseline unless a concrete blocker is found:

- **Game client:** Godot 4.x;
- **Single deterministic simulation/rules engine:** Rust;
- **Backend services:** Go;
- **Primary database:** PostgreSQL 18;
- **Local/offline database:** SQLite;
- **Messaging/event backbone:** NATS + JetStream where durability is required;
- **Cache/ephemeral state:** Valkey only when justified by a real use case;
- **Object storage:** S3-compatible;
- **Admin/back-office:** React + TypeScript;
- **External/mobile API:** HTTPS/JSON initially;
- **Realtime delivery:** WebSocket only where realtime interaction is actually required;
- **Runtime:** native Linux/systemd on the estate; no container dependency;
- **Repository:** monorepo;
- **Identity UX:** guest-first mobile entry with later account protection/linking;
- **Content/configuration:** data-driven, versioned, signed and server-controlled;
- **Search:** PostgreSQL initially; introduce dedicated search infrastructure only when justified.

A central invariant is:

> **One rules engine.**

The Rust simulation package is shared by the mobile client, server-side validator, challenge runner, test harness and balance tooling so that core facility/economy rules are not reimplemented inconsistently across languages.

The backend may be modular without being prematurely fragmented into dozens of microservices.


---

## 35. Non-political authorities and GRIDWORKS treasury recirculation

GRIDWORKS deliberately excludes political gameplay.

There are no mayors, elected governments, political parties, ideological factions or player-controlled public offices.

Operational charges may still exist and be explained through neutral service/economic abstractions such as:

- Local Authority;
- Port Authority;
- Utility Authority;
- Infrastructure Authority;
- Registry Authority.

These are not political actors. They exist only to explain and account for charges such as property tax, port fees, warehouse/local charges, registry fees or infrastructure usage.

Internally, those charges may accumulate in auditable non-player treasury ledgers. GRIDWORKS operators may use that in-game value under controlled rules to recirculate money back into the player economy through:

- system procurement;
- community project funding;
- infrastructure funding;
- starter/recovery buybacks;
- occasional GRIDWORKS purchase offers for player-owned businesses or property;
- release/acquisition of GRIDWORKS-owned inventory.

GRIDWORKS buyback offers must never become an unlimited guaranteed purchaser at above-market value. Offers should be conditional, bounded, auditable and priced from transparent valuation inputs.

The purpose is to create a closed-loop economic sink/source system rather than simply destroying all tax/fee revenue.

## 36. GRIDWORKS-owned assets

GRIDWORKS may own in-game land, property, businesses and facilities as platform/system inventory.

A player clicking such an asset may see, for example:

**Owner: GRIDWORKS**  
**Price: 420,000 Cr**  
**[Purchase]**

A player-owned asset may instead expose:

**Owner: Maik Industrial Properties**  
**[View Company] [Contact Owner] [Make Offer]**

GRIDWORKS-owned assets may be:

- fixed-price purchasable;
- leaseable;
- reserved/unavailable;
- held for a future season;
- used to seed a new market;
- repurchased from players under bounded system offers.

Future Season 2+ inventory may already exist visibly in the world from launch under GRIDWORKS ownership but remain unavailable until the relevant release.

GRIDWORKS-owned productive assets are never directly purchasable with real money.

## 37. Real estate and player-selected starts

Real estate is part of the launch-world foundation, even if Season 1 implements only a constrained subset.

Core launch concepts include:

- land parcels;
- industrial buildings;
- warehouses;
- offices;
- basic residential property;
- ownership;
- sale/listing;
- leasing foundations;
- vacancy/occupancy;
- maintenance;
- renovation;
- valuation;
- GRIDWORKS-owned and player-owned inventory.

New players should not be forced into one industry or given a burdensome portfolio of every business.

After the shared introductory reasoning tutorial, players choose a starting path such as:

- Mining/Quarrying;
- Agriculture;
- Manufacturing;
- Energy;
- Logistics;
- Real Estate;
- Finance/Markets;
- Generalist / Surprise Me.

This is a starting position, not a class.

Players may later diversify freely.

A business the player does not want to operate may be:

- operated;
- mothballed/suspended;
- sold;
- leased/contracted out later;
- dismantled/repurposed where appropriate.

The game should make it clear that the starting business is only the player's first foothold in the world, not their permanent identity.

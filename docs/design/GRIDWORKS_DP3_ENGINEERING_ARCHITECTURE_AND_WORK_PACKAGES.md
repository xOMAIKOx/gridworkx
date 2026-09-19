# GRIDWORKS — Design Phase 3 (DP3): Engineering Architecture & Work-Package Design

**Repository:** `xOMAIKOx/gridworkx`  
**Status:** DP3 ARCHITECTURE / CONTROLLING ENGINEERING DESIGN  
**Date:** 2026-09-19  
**Depends on:**  
- `docs/design/GRIDWORKS_PRODUCT_PHILOSOPHY_AND_GAME_DESIGN.md`
- `docs/design/GRIDWORKS_DP2_CORE_SYSTEMS_SPECIFICATION.md`

---

## 0. Purpose

DP3 converts the accepted product philosophy and DP2 core systems model into an engineering architecture suitable for implementation through bounded GitHub work packages.

The primary objectives are:

1. preserve the **single deterministic rules engine** across client and server validation;
2. define domain boundaries without premature microservice fragmentation;
3. establish authoritative-state ownership;
4. define persistence, event, API and sync boundaries;
5. include social/messaging, business lifecycle and localization as launch architecture;
6. create an auditable economy from the first commit;
7. prepare a work-package sequence suitable for Luna XHigh / Max implementation and Sol-class review.

---

# PART I — ACCEPTED TECHNOLOGY STACK

## 1. Stack

| Layer | Decision |
|---|---|
| Game client | Godot 4.x |
| Shared deterministic simulation/rules engine | Rust |
| Backend services | Go |
| Primary database | PostgreSQL 18 |
| Local/offline database | SQLite |
| Messaging/event transport | NATS + JetStream |
| Cache | Valkey only when justified |
| Object storage | S3-compatible |
| Admin/back-office | React + TypeScript |
| External/mobile API | HTTPS/JSON initially |
| Realtime delivery | WebSocket |
| Runtime | Native Linux/systemd |
| Repository model | Monorepo |
| Search | PostgreSQL initially |
| Identity UX | Guest-first, account-link/protect later |
| Content/rules data | Versioned, signed, data-driven |

No container dependency is part of the target runtime architecture.

---

# PART II — MONOREPO ARCHITECTURE

## 2. Proposed repository layout

```
gridworkx/
├── apps/
│   ├── game/                  # Godot client
│   ├── admin/                 # React/TypeScript back-office
│   └── tools/                 # balance/content/dev utilities
├── crates/
│   └── gridworks-sim/         # canonical Rust simulation/rules engine
├── services/
│   ├── api/                   # Go public/mobile API
│   ├── worker/                # background jobs
│   ├── realtime/              # websocket/session delivery
│   └── market/                # may begin in API module; split only when justified
├── packages/
│   ├── schemas/               # shared contracts/schema generation
│   ├── content/               # versioned game definitions
│   └── localization/          # locale packs/glossary/source catalogues
├── db/
│   ├── migrations/
│   └── seeds/
├── docs/
│   ├── design/
│   ├── architecture/
│   ├── decisions/
│   └── work-orders/
├── ops/
│   ├── systemd/
│   ├── config/
│   └── scripts/
└── tests/
    ├── integration/
    ├── economy/
    ├── deterministic/
    └── security/
```

The repository is modular but remains one engineering coordination unit unless scale creates a concrete reason to split.

---

# PART III — ONE RULES ENGINE

## 3. Canonical simulation package

`crates/gridworks-sim` is the only canonical implementation of:

- facility dependency graphs;
- component state;
- production transformations;
- resource flows;
- bottlenecks;
- equipment wear;
- failure generation;
- repair effects;
- manager operational modifiers where simulation-related;
- player-skill operational modifiers where simulation-related;
- deterministic offline advancement;
- standardized challenge execution;
- calculation of operational KPIs.

The simulation package is consumed by:

- Godot mobile client;
- server validation process;
- challenge validator/runner;
- deterministic replay tests;
- balance simulation tools.

Go services must not reimplement simulation mathematics.

## 3.1 FFI boundary

The Rust engine exposes a stable boundary to Godot.

Preferred properties:

- explicit versioned command/input schema;
- deterministic state snapshots;
- no direct database access from Rust simulation core;
- no network access from simulation core;
- no dependence on client wall clock;
- serializable replay/event inputs;
- narrow FFI surface.

Exact Godot integration mechanism is implementation-deferred, but must preserve these constraints.

---

# PART IV — DOMAIN BOUNDARIES

## 4. Core domains

Initial logical domains:

1. Identity & Account
2. Player Profile
3. Company & Ownership
4. World / Region / Site
5. Facility / Components / Simulation
6. Resources / Inventory
7. Construction / Maintenance
8. Skills
9. Managers / Recruitment
10. Trade / Marketplace
11. Contracts / Work Exchange
12. Reputation
13. Consortium / JV
14. Social / Messaging
15. Notifications
16. Leaderboards / Challenges
17. Accounting / Ledger
18. Business Lifecycle / Acquisition / Disposal
19. Localization / Content
20. Admin / Moderation / Live Ops
21. Finance / Exchange — reserved for later implementation

These are **domain boundaries**, not necessarily one process per domain.

---

# PART V — AUTHORITATIVE STATE MATRIX

## 5. State authority

| State | Primary authority |
|---|---|
| Account/authentication | Server |
| Player profile | Server |
| Company ownership | Server |
| Cash / ledger | Server |
| Inventory quantities | Server |
| Market orders/trades | Server |
| Contracts | Server |
| Manager recruitment result | Server |
| Manager employment/transfer | Server |
| Consortium/JV ownership | Server |
| Business sale/acquisition | Server |
| Message history | Server |
| Notification preferences | Server |
| Localization preference | Server + local cache |
| Facility simulation rendering | Client |
| Facility deterministic advancement | Shared Rust engine; server-verifiable |
| Challenge simulation | Client Rust engine + server validation |
| UI layout/preferences | Client |
| Local planning drafts | Client |
| Offline cache | Client |
| Economic settlement from offline actions | Server validated |

Any state affecting another player must be server-verifiable.

---

# PART VI — PERSISTENCE

## 6. PostgreSQL model

PostgreSQL is authoritative for transactional domain state.

Initial table families should include:

- accounts;
- player_profiles;
- player_settings;
- companies;
- company_ownership;
- company_groups;
- regions;
- sites;
- facilities;
- facility_snapshots;
- inventories;
- resources;
- construction_jobs;
- skills;
- managers;
- manager_employment_history;
- recruitment_candidates/results;
- market_orders;
- trades;
- contracts;
- contract_events;
- reputation_events;
- consortiums;
- consortium_members;
- jv_entities;
- jv_ownership;
- conversations;
- conversation_participants;
- messages;
- message_reports;
- notifications;
- economic_transactions;
- ledger_entries;
- business_listings;
- ownership_history;
- localization_preferences;
- rules_versions.

Exact schemas are a dedicated work package.

---

## 7. Ledger-first economy

All material economic mutations must be auditable.

Examples:

- purchase;
- sale;
- contract payment;
- manager compensation;
- construction spend;
- repair spend;
- market fee;
- resource valuation event where required;
- dividend later;
- debt later.

A transaction/journal layer must exist from the first vertical slice.

No uncontrolled service-level mutations such as `coins += x`.

Economic writes should have:

- transaction ID;
- source actor;
- target actor;
- reason/type;
- related object;
- timestamp;
- amount/value entries;
- idempotency key;
- audit metadata.

---

# PART VII — EVENT ARCHITECTURE

## 8. NATS usage

NATS provides cross-module/event delivery.

JetStream is used only where durable replay is required.

Example event subjects:

```
facility.production.started
facility.production.stopped
facility.fault.detected
repair.completed
construction.completed
market.trade.executed
contract.accepted
contract.completed
manager.recruited
manager.transferred
company.acquired
company.sold
message.created
consortium.project.updated
challenge.result.accepted
```

Do not publish high-frequency visual simulation ticks.

Use domain events, not animation events.

---

# PART VIII — API ARCHITECTURE

## 9. Mobile API

Initial public API:

- HTTPS;
- JSON;
- explicit versioning;
- idempotency for economic writes;
- optimistic concurrency where appropriate;
- server-issued timestamps;
- stable semantic IDs.

Primary API groups:

- auth/session;
- profile;
- company;
- world/facility;
- inventory;
- construction/repair;
- manager/recruitment;
- market;
- contracts;
- consortium/JV;
- business opportunities/sales;
- messages;
- notifications;
- challenge/leaderboards;
- localization/content manifests.

## 9.1 Realtime

WebSocket carries:

- incoming DMs;
- Consortium messages;
- object-thread messages;
- contract changes;
- market fills;
- important facility events;
- notification fan-out;
- leaderboard/challenge updates where useful.

It must not carry continuous local simulation.

---

# PART IX — IDENTITY AND PROFILE

## 10. Guest-first identity

First launch should create a guest identity quickly.

The player can begin play before completing a conventional account-registration flow.

Later prompt:

> Protect your company/account

Link options may include:

- Apple;
- Google;
- email/OIDC;
- estate identity integration where appropriate.

Guest-to-linked migration must preserve:

- player ID;
- company ownership;
- progress;
- messages;
- purchases;
- reputation.

---

## 11. Player profiles

Server fields:

- immutable player ID;
- unique handle;
- display name;
- avatar;
- bio;
- locale;
- public achievements;
- reputation summary;
- privacy options;
- DM policy;
- notification policy.

Unique handle rules must be Unicode/confusable aware.

---

# PART X — COMPANY AND BUSINESS IDENTITY

## 12. Company hierarchy

Supported structure:

```
Player
 └── Group / Holding Company
      ├── Operating Company
      ├── Operating Company
      └── JV interests
```

Company identity persists independently from the player.

Company names and later public tickers require uniqueness rules.

A company profile may expose:

- name/logo;
- industry;
- HQ/region;
- ownership;
- reputation;
- operating metrics;
- financial metrics;
- consortium memberships;
- public/private state;
- messaging/trade/contact actions.

---

# PART XI — BUSINESS LIFECYCLE

## 13. Opportunity engine

Business opportunity classes:

- Greenfield;
- Going Concern;
- Turnaround;
- Player Listing;
- Restructuring/Tender later.

The opportunity engine generates a **causal history first**, then derives present facility/business condition.

Example generated history inputs:

- age;
- maintenance strategy;
- utilization;
- recent incident;
- capital investment;
- management quality;
- staffing;
- utility quality;
- logistics constraints;
- contracts;
- liabilities.

The same Rust rules used by player operations must derive resulting technical state where applicable.

## 13.1 Opportunity Board

The board is finite and time-varying.

It must not provide infinite free rerolls.

Member advantage may include:

- saved searches;
- filters;
- alerts;
- watchlists;

but not privileged access to superior generated businesses.

## 13.2 Due diligence

Due diligence returns imperfect but fair information.

Skill/manager expertise may improve:

- confidence;
- fault visibility;
- financial interpretation;
- reserve/asset estimate quality.

## 13.3 Sale structures

Early implementation may begin with:

- facility/business sale.

Architecture must reserve later support for:

- asset sale;
- company sale;
- stake sale;
- acquisition;
- JV;
- restructuring.

Ownership history is permanent/auditable.

---

# PART XII — SOCIAL AND MESSAGING

## 14. Messaging model

Durable entities:

- Conversation;
- ConversationParticipant;
- Message;
- MessageAttachment;
- ReadCursor;
- MessageReport;
- ContextLink;
- ChannelPermission.

## 14.1 DMs

Private messages use an inbox model:

- unread;
- archive;
- pin;
- mute;
- block;
- report;
- search;
- persistent history;
- contextual launch from company/profile/contract.

## 14.2 Consortium chat

Initial channels:

- General;
- Projects;
- Trade;
- Management;
- Announcements.

Role permissions control access.

## 14.3 Object-linked chat

Contracts, projects, jobs and JVs may own their own thread.

Closed objects preserve thread history as read-only or archived.

## 14.4 Attachments

Attachments live in S3-compatible object storage.

Initial restrictions:

- approved image/document types only;
- size caps;
- malware/moderation scanning;
- audit metadata;
- rate limits.

---

# PART XIII — NOTIFICATIONS

## 15. Notification service

Categories:

- Operations;
- Markets;
- Contracts;
- Messages;
- Consortium;
- Managers;
- Finance;
- Security.

Delivery:

- in-app;
- push;
- silent/background;
- disabled.

All categories are user configurable.

---

# PART XIV — LOCALIZATION / INTERNATIONALIZATION

## 16. Semantic localization model

All player-facing system text is rendered from stable semantic keys.

Authoritative domain data never relies on translated prose.

Example:

```
resource.copper_ore
fault.motor_bearing_seizure
event.production_stopped
contract.supply.completed
```

## 16.1 Formatting requirements

Localization runtime must support:

- plural categories;
- grammar variants;
- parameter ordering;
- locale date/time;
- number separators;
- percentages;
- financial formatting;
- RTL;
- bidi text;
- unit preferences.

Use ICU-equivalent semantics.

## 16.2 Units

Simulation uses canonical SI.

Presentation converts to player preference.

## 16.3 Locale packs

Packages contain:

- strings;
- tutorial text;
- glossary;
- font/fallback metadata;
- localized media if required.

Locale packs are:

- versioned;
- signed;
- cacheable;
- downloadable;
- independently deployable.

## 16.4 Translation workflow

1. canonical English/source content;
2. glossary resolution;
3. AI-assisted translation;
4. placeholder/schema validation;
5. human review for critical surfaces;
6. signed package publication;
7. fallback testing.

Critical surfaces requiring human/native review:

- onboarding;
- monetization;
- legal/privacy;
- finance;
- industrial terminology;
- safety/moderation;
- store listings.

## 16.5 Chat translation

Optional service.

Rules:

- original message remains canonical;
- translated version is clearly marked;
- service failure never blocks messaging;
- no hidden replacement of sender text.

---

# PART XV — ADMIN / BACK-OFFICE

## 17. Admin application

React + TypeScript.

Initial functional areas:

- player/account support;
- moderation;
- company/profile/name moderation;
- message/report review;
- economy dashboard;
- market anomaly review;
- rules/content publication;
- localization review;
- opportunity generation inspection;
- challenge authoring;
- manager/resource/recipe configuration;
- membership/support;
- audit logs.

The admin application is not part of the Godot client.

---

# PART XVI — CONTENT AND CONFIGURATION

## 18. Data-driven definitions

Simulation/content definitions should be versioned data, not scattered constants.

Examples:

- resource definitions;
- component definitions;
- failure modes;
- repair options;
- recipes;
- manager traits;
- recruitment tables;
- skill progression;
- challenge definitions;
- leaderboard formulas;
- emergency supplier policies.

Content publication must include:

- schema validation;
- semantic version/rules version;
- signatures/checksum;
- compatibility metadata.

---

# PART XVII — OFFLINE MODEL

## 19. SQLite client cache

Local SQLite stores:

- cached profile/company state;
- cached facility snapshots;
- content manifests;
- localization packs;
- local planning;
- queued non-authoritative intents;
- UI settings;
- last sync metadata.

SQLite does not become authority for shared economic state.

## 19.1 Offline catch-up

Server-issued authoritative timestamps define elapsed time.

Rust engine computes deterministic progression from last accepted snapshot.

Server validates permissible resulting transitions before settlement.

---

# PART XVIII — SECURITY

## 20. Security requirements

Threats:

- modified client;
- clock tampering;
- replay;
- duplicate writes;
- resource duplication;
- fraudulent market actions;
- multi-account farming;
- botting;
- message spam;
- impersonation;
- malicious attachments.

Required patterns:

- idempotency keys;
- authoritative timestamps;
- server validation;
- economic audit ledger;
- signed content/rules;
- replay validation;
- rate limiting;
- moderation/report flows;
- confusable-name detection;
- session/device risk signals.

---

# PART XIX — TEST STRATEGY

## 21. Required test layers

### 21.1 Rust simulation

- deterministic replay tests;
- property tests;
- failure-model tests;
- production/bottleneck tests;
- repair-path tests;
- offline advancement tests;
- rules-version compatibility tests.

### 21.2 Go backend

- domain unit tests;
- API contract tests;
- idempotency tests;
- transaction integrity;
- authorization/ownership tests;
- market settlement;
- messaging permissions;
- sale/acquisition accounting.

### 21.3 Integration

- Godot ↔ Rust;
- client ↔ API;
- API ↔ PostgreSQL;
- event delivery;
- offline sync;
- social realtime;
- localization package updates.

### 21.4 Economy

- currency conservation/injection tests;
- exploit simulations;
- market manipulation tests;
- member/free fairness analytics;
- business-sale transfer invariants.

---

# PART XX — INITIAL WORK-PACKAGE LEDGER

## 22. Work-package sequence

### WP-001 — Repository and native runtime baseline
**Model:** Luna XHigh  
**Review:** Sol  
**Scope:** monorepo structure, build/test skeleton, native Linux/systemd conventions, CI foundation.

### WP-002 — Rust simulation kernel
**Model:** Luna XHigh if DP3 contracts are sufficient; Luna Max if determinism/FFI investigation is required  
**Review:** Sol  
**Scope:** simulation state model, rules-version framework, deterministic command execution, serialization.

### WP-003 — Core facility/component/dependency graph
**Model:** Luna XHigh  
**Review:** Sol  
**Scope:** Facility/System/Component primitives, dependency types, bottleneck calculation.

### WP-004 — Failure / diagnosis / intervention model
**Model:** Luna XHigh  
**Review:** Sol  
**Scope:** causal faults, observability, confidence, repair/bypass/rebuild actions.

### WP-005 — Resources / recipes / inventory
**Model:** Luna XHigh  
**Review:** Sol

### WP-006 — PostgreSQL foundation and ledger
**Model:** Luna XHigh  
**Review:** Sol  
**Scope:** company/economic ledger, migrations, transaction/idempotency primitives.

### WP-007 — Identity / guest-first / player profile
**Model:** Luna XHigh  
**Review:** Sol

### WP-008 — Company / ownership / business identity
**Model:** Luna XHigh  
**Review:** Sol

### WP-009 — Go API foundation
**Model:** Luna XHigh  
**Review:** Sol

### WP-010 — Godot ↔ Rust integration
**Model:** Luna Max if integration unknowns remain  
**Review:** Sol

### WP-011 — Opening aggregate-plant vertical slice
**Model:** Luna XHigh after WP-002–010 converge  
**Review:** Sol

### WP-012 — Skills and manager domain
**Model:** Luna XHigh  
**Review:** Sol

### WP-013 — Recruitment system
**Model:** Luna XHigh  
**Review:** Sol

### WP-014 — Contracts / Work Exchange
**Model:** Luna XHigh  
**Review:** Sol

### WP-015 — Marketplace / trade settlement
**Model:** Luna Max if matching/settlement complexity remains unresolved  
**Review:** Sol

### WP-016 — Business lifecycle / opportunity generator / due diligence
**Model:** Luna Max  
**Review:** Sol  
**Reason:** causal procedural generation requires stronger investigation and property testing.

### WP-017 — Sale/acquisition transfer
**Model:** Luna XHigh  
**Review:** Sol

### WP-018 — Social messaging / DM inbox
**Model:** Luna XHigh  
**Review:** Sol

### WP-019 — Consortium and contextual threads
**Model:** Luna XHigh  
**Review:** Sol

### WP-020 — Notifications / push abstraction
**Model:** Luna XHigh  
**Review:** Sol

### WP-021 — Localization engine / language-pack pipeline
**Model:** Luna Max initially  
**Review:** Sol  
**Reason:** RTL, ICU semantics, font fallback and hot-update packaging need architectural care.

### WP-022 — Admin/back-office foundation
**Model:** Luna XHigh  
**Review:** Sol

### WP-023 — Challenge sandbox / replay validation
**Model:** Luna Max  
**Review:** Sol

### WP-024 — Leaderboards / scoring
**Model:** Luna XHigh  
**Review:** Sol

### WP-025 — Economy telemetry / fairness monitoring
**Model:** Luna XHigh  
**Review:** Sol

Finance/exchange implementation begins only after the operational economy is proven stable.

---

# PART XXI — WORK-PACKAGE HANDOFF STANDARD

## 23. Engineering handback

Every engineering work package must return through GitHub with:

- WP identifier;
- branch;
- exact commit SHA;
- parent SHA;
- files changed;
- architecture/spec references;
- commands/tests executed;
- test results;
- evidence paths;
- migrations;
- API/schema changes;
- known deviations;
- unresolved risks;
- security implications;
- performance notes;
- explicit statement that no out-of-scope host/infra work was performed.

Architecture reviews the exact commit and records:

- PASS;
- REQUEST_CHANGES;
- HOLD/BLOCKED.

The owner receives a concise summary only.

---

# PART XXII — DP3 INVARIANTS

## 24. Non-negotiable architecture invariants

1. One canonical Rust rules engine.
2. Go backend does not duplicate simulation mathematics.
3. PostgreSQL is the transactional/economic authority.
4. Economic mutations are ledgered/auditable.
5. Client wall clock is never authoritative.
6. Any shared economic state is server-verifiable.
7. Local simulation is allowed; shared settlement is validated.
8. Player, Account and Company are separate.
9. Company ownership can change without account transfer.
10. Business history persists after sale.
11. Generated distressed businesses derive from causal histories.
12. Messaging is durable and context-aware.
13. DM inbox, Consortium chat and job/project threads are launch architecture.
14. Global chat is not required for launch.
15. Notifications are opt-configurable and non-manipulative.
16. Localization is semantic-key driven and simulation-neutral.
17. RTL, Unicode and locale formatting are designed in from the beginning.
18. Normal game operation never depends on a live LLM.
19. Paid membership does not create tradeable economic power.
20. Repository remains monorepo unless a concrete scaling/ownership reason requires a split.
21. Native Linux/systemd is the target runtime.
22. Infrastructure is added only for a demonstrated need; no ornamental cache/search/microservice layers.

---

# PART XXIII — DP3 EXIT CRITERIA

DP3 is architecture-complete when:

1. accepted stack and repository structure are committed;
2. Rust simulation boundary is documented;
3. authoritative-state matrix is accepted;
4. ledger-first persistence model is accepted;
5. identity/player/company separation is accepted;
6. business lifecycle is represented;
7. messaging/social domain is represented;
8. localization/global-launch domain is represented;
9. first vertical slice dependencies are known;
10. work-package ledger exists with model routing;
11. engineering handback/review standard exists;
12. unresolved items are isolated into ADRs/work packages rather than hidden in implementation.

DP3 does not authorize engineering implementation by itself. Each work package requires explicit issuance/authorization under the repository collaboration process.


---

# PART XXIV — ACCEPTED ADR REGISTER

## 25. Accepted architecture decisions

The following ADRs are controlling for DP3 and all engineering work packages:

1. `docs/decisions/ADR-001_GODOT_RUST_INTEGRATION.md`
   - Godot integrates the canonical Rust simulation through GDExtension;
   - coarse-grained, versioned FFI;
   - no simulation duplication in GDScript/Go.

2. `docs/decisions/ADR-002_IDENTITY_AND_ACCOUNT_LINKING.md`
   - server-issued guest-first identity;
   - later account protection/linking;
   - immutable player continuity through account linking.

3. `docs/decisions/ADR-003_SCHEMA_AND_SERIALIZATION.md`
   - repository-owned canonical schema catalogue;
   - JSON external API;
   - generated cross-language contracts;
   - compact binary encoding only where justified;
   - versioned event/snapshot/replay envelopes.

4. `docs/decisions/ADR-004_CONTENT_RULES_DISTRIBUTION.md`
   - versioned signed rules/content/localization bundles;
   - server-issued authoritative manifests;
   - independently deployable localization where semantic compatibility permits.

5. `docs/decisions/ADR-005_INITIAL_SERVER_TOPOLOGY.md`
   - initial native process set: API, worker, realtime, simulation validator boundary;
   - PostgreSQL + NATS + S3-compatible storage;
   - no Valkey/search cluster until measured need;
   - no container runtime on the target estate.

Engineering may not replace these decisions inside a work package. If implementation reveals a concrete blocker, Engineering must stop that decision path, document evidence, and request an ADR amendment from Architecture.

---

# PART XXV — DP3 GATE SEQUENCE

## 26. Architecture-to-engineering gates

### G0 — Product/design baseline
PASS when:
- Product Philosophy baseline exists;
- DP2 core systems specification exists;
- DP3 architecture exists.

**Current status:** PASS.

### G1 — Architecture decision closure
PASS when:
- client/simulation integration is decided;
- identity/account-linking is decided;
- schemas/serialization are decided;
- content/rules distribution is decided;
- initial server topology is decided.

**Current status:** PASS via ADR-001..ADR-005.

### G2 — WP-001 issue readiness
PASS when Architecture publishes a bounded WP-001 engineering work order containing:
- exact parent SHA;
- allowed paths;
- prohibited actions;
- required repository layout;
- build/test gates;
- native-runtime constraints;
- evidence/handback format.

**Current status:** NOT YET ISSUED.

### G3 — Repository foundation acceptance
WP-001 must PASS Architecture review before WP-002+ implementation branches are authorized.

### G4 — Simulation kernel acceptance
WP-002 must prove:
- deterministic command execution;
- versioned snapshot/command schema;
- seeded RNG;
- no wall-clock/network/database dependency;
- repeatable digest tests.

### G5 — Vertical-slice foundation
WP-003..WP-010 converge before WP-011 is authorized as an integrated vertical slice.

### G6 — Product-loop proof
WP-011 must demonstrate:

> inspect → diagnose → intervene → partial production → economic settlement → second decision

without violating the fair-play or authority invariants.

No broad industry expansion is authorized before this loop is accepted.

---

# PART XXVI — WP ROUTING REFINEMENT

## 27. Model routing rules

Use model capability according to uncertainty rather than work-package importance alone.

### Luna XHigh
Default implementation worker when:
- architecture is already decided;
- acceptance criteria are explicit;
- work is bounded;
- failure mode is normal engineering/debugging.

### Luna Max
Use when:
- investigation is genuinely required;
- multiple implementation approaches must be tested;
- FFI/platform behavior is uncertain;
- procedural-generation properties require exploration;
- localization/RTL/font behavior needs empirical validation;
- challenge replay/determinism behavior is non-convergent.

### Sol / Sol Pro / highest available architecture reviewer
Use for:
- architecture;
- ADR changes;
- cross-domain design;
- security/economic invariant review;
- work-order authoring;
- acceptance/rejection of engineering handbacks.

Engineering must not silently escalate unresolved architecture into code because a stronger implementation model is available.

---

# PART XXVII — NEXT AUTHORIZED DESIGN ACTION

## 28. Next Architecture task

The next design action after DP3/ADR closure is to author **WP-001 — Repository and Native Runtime Baseline** as a detailed GitHub work order.

WP-001 is repository-only unless the owner separately authorizes host work.

It should create the monorepo skeleton and engineering conventions necessary for later work, but it must not implement substantive game systems prematurely.

DP3/ADR completion itself does **not** authorize Engineering to start WP-001.


---

# PART XXVIII — ADR-006 / ADR-007 INTEGRATION

## 29. Additional accepted architecture decisions

6. `docs/decisions/ADR-006_GRIDWORKS_OWNED_ASSETS_AND_TREASURY.md`
   - no political-government gameplay;
   - neutral operational authorities only;
   - GRIDWORKS-owned system inventory;
   - treasury/authority ledgers;
   - bounded value recirculation and GRIDWORKS buyback offers.

7. `docs/decisions/ADR-007_REAL_ESTATE_AND_STARTING_PATHS.md`
   - real estate is a launch architecture domain;
   - WP-001 must scaffold land/property concepts;
   - future-season GRIDWORKS-owned property may be visible but unavailable;
   - starter-path selection is required;
   - players are never forced to operate unwanted businesses.

## 29.1 Domain additions

Add to the DP3 domain map:

22. Real Estate / Land / Property
23. Non-Player Treasury / Authority Ledgers
24. Starter Path / Onboarding Selection

## 29.2 Persistence additions

PostgreSQL foundation must reserve/support table families for:

- land_parcels;
- buildings;
- property_ownership;
- property_listings;
- leases;
- occupancy;
- renovation_jobs;
- authority_principals;
- treasury_accounts;
- treasury_transactions;
- system_purchase_offers;
- starter_path_selection;
- business_operating_state.

Exact schemas remain WP-owned, but WP-001 scaffolding must acknowledge these domains.

## 29.3 Work-package ledger amendment

WP-001 scope now explicitly includes repository/domain scaffolding for:

- real estate / land / property;
- GRIDWORKS system-principal ownership;
- non-player treasury ledgers;
- starter-path selection;
- business operating-state model (active/mothballed/listed/etc.).

WP-001 does **not** implement the full real-estate simulation.

A later work package should implement the Season-1 real-estate baseline after the core vertical slice is proven, while preserving the already-defined domain contracts.

## 29.4 Gate update

G2 WP-001 issue readiness additionally requires the work order to name ADR-006 and ADR-007 as controlling references.

Engineering may not omit these domains from the repository foundation simply because full gameplay arrives later.


---

# PART XXIX — GRIDWORKS LIQUIDITY FACILITY

## 30. Buyer-of-last-resort architecture

ADR-006 now defines GRIDWORKS as a buyer of last resort for eligible player-owned businesses and properties.

Engineering architecture must therefore support:

- system valuation of eligible assets;
- immediate GRIDWORKS purchase settlement;
- transfer to GRIDWORKS system ownership;
- preservation of full asset/business history;
- dynamic GRIDWORKS resale pricing;
- acquisition-cost floor accounting;
- carrying-cost accumulation;
- original-owner reacquisition detection;
- minimum original-owner reacquisition premium (initially 20% over GRIDWORKS acquisition cost, subject to the higher live listing price/carrying-cost floor);
- scarcity-aware listing price changes;
- anti-circular-sale and anti-arbitrage controls;
- treasury liquidity/accounting;
- complete price/ownership audit history.

### 30.1 Pricing invariants

Normal GRIDWORKS resale price must not intentionally fall below GRIDWORKS acquisition cost plus applicable carrying/transaction costs except through an explicitly authorized economy intervention.

Former owners do not receive a preferential reacquisition discount.

If comparable inventory is unavailable, scarcity may increase the listing price for all buyers.

### 30.2 Work-package impact

WP-001 must scaffold the domain boundaries and schema locations for:

- valuation policy;
- system purchase offer;
- GRIDWORKS inventory ownership;
- resale listing;
- carrying cost;
- reacquisition rule;
- treasury settlement.

WP-006 ledger work must support GRIDWORKS/system-principal purchases and resale accounting.

WP-017 sale/acquisition transfer must implement the actual transfer and reacquisition invariants.

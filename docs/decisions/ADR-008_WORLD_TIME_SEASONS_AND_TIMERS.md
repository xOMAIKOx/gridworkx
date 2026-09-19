# ADR-008 — World Time, Regional Seasons and Timer Philosophy

**Status:** ACCEPTED  
**Date:** 2026-09-19  
**Applies to:** simulation clock, construction, production, maintenance, agriculture, contracts, markets, offline progression

## Decision

### 1. Core simulation time ratio

The initial GRIDWORKS timing baseline is:

> **1 in-game day = 1 real hour**

Therefore:
- 24 in-game days = 24 real hours;
- 7 in-game days = 7 real hours;
- 30 in-game days = 30 real hours.

The ratio is a versioned balance parameter but should be treated as a major world invariant, not changed casually after launch.

### 2. Separate clocks

GRIDWORKS distinguishes:

- **Operational simulation time** — compressed at the 1 game day / 1 real hour baseline;
- **Authoritative real time** — UTC server time used for deadlines, sync, refreshes and settlement;
- **World climate calendar** — linked to real-world calendar progression and regional climate rules.

World climate seasons do not cycle every few compressed game days.

### 3. Regional seasons

Climate seasonality is regional rather than globally uniform.

Examples:
- northern temperate regions: spring/summer/autumn/winter aligned to calendar months;
- southern temperate regions: opposite seasonal cycle;
- tropical regions: wet/dry cycles;
- arid regions: hot/cool/rainfall cycles.

Regional seasonality may affect:
- crop yields and planting;
- heating/cooling demand;
- electricity demand;
- oil/gas demand;
- water demand;
- hydro availability;
- solar output;
- logistics reliability;
- construction productivity;
- property operating costs.

### 4. Timer classes

Not all actions use the same waiting scale.

Typical categories:
- seconds/minutes: inspection, configuration, minor intervention;
- minutes/hours: repair, maintenance, component replacement;
- hours/days: construction, redevelopment, large infrastructure.

Exact values are data-driven and balanced later.

### 5. Timer philosophy

Timers exist to create planning and operational choices, not artificial frustration.

Time may be reduced through in-game mechanisms such as:
- player assistance;
- manager/project-management skill;
- prefabrication;
- logistics improvement;
- better equipment;
- earned construction capacity.

There is no real-money instant-finish path that creates productive power.

### 6. Offline-safe behavior

Offline play must not punish normal life.

Examples:
- production continues while safe and supplied;
- storage full causes safe pause/reroute rather than silent loss;
- critical faults may trigger safe shutdown;
- crops do not require an unrealistically narrow login window to avoid destruction;
- manager recovery proceeds under normal rules;
- contracts use reasonable real-time deadlines;
- the player returns to consequences and decisions, not arbitrary catastrophe caused solely by absence.

### 7. Release terminology

Feature/content expansions use **Release** terminology:
- Release 1;
- Release 2;
- Release 3.

The word **season** is reserved primarily for world/climate seasonality and challenge/competition periods only where explicitly qualified.

## Invariants

1. Client wall clock is never authoritative.
2. Simulation-time compression and climate seasonality are separate concepts.
3. Regional climate drives economic consequences.
4. Real money cannot bypass productive timers.
5. Offline absence alone must not create punitive failure.
6. Time configuration is versioned and server-controlled.

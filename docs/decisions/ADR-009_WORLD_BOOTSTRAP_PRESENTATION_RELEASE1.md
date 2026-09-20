# ADR-009 — World Structure, Bootstrap Economy, Presentation and Release-1 Posture

**Status:** ACCEPTED  
**Date:** 2026-09-20  
**Applies to:** world topology, trade/logistics, presentation, launch completeness, age/social posture, bootstrap economy

## Decision

### 1. One logical global world

GRIDWORKS presents one logical persistent global economy divided into economically distinct regions.

Regions differ by:
- climate;
- geography;
- resources;
- transport links;
- land/property markets;
- power/water availability;
- local demand;
- industrial concentration.

Player-visible shards are not part of the product model unless later scale makes them technically unavoidable.

### 2. Multi-layer world presentation

GRIDWORKS uses four principal presentation layers:

1. **World Map** — macro economy, trade, player/company presence, logistics routes, regional opportunities.
2. **Player Portfolio Map** — all businesses/assets owned or controlled by the player, with lightweight animated activity.
3. **Site / Business View** — industry-specific 2.5D representation of a mine, refinery, farm, warehouse, property, etc.
4. **Operational Detail View** — diagnostics, components, failures, repairs, production, telemetry and management.

### 3. Visual architecture

Initial visual architecture:
- 2.5D isometric/oblique presentation;
- stylized low-poly 3D assets rendered in Godot;
- 2D information overlays and panels;
- clean semi-realistic industrial-diorama aesthetic;
- limited/fixed site-view rotation, initially four orientations;
- no free-camera full-3D requirement for Release 1.

The asset system should preserve the option for richer future reskins without changing authoritative gameplay state.

### 4. Functional animation

Animation is functional and readable rather than cinematic.

Examples:
- trucks moving on roads;
- trains moving on rail;
- ships using ports/routes;
- conveyors running;
- turbines rotating;
- pumps/motors showing active state;
- construction cranes/activity;
- smoke/steam/lights and other operational cues;
- lightweight ambient people/vehicle activity where useful.

Detailed character simulation is not a Release-1 requirement.

### 5. Transport model

Core transport modes:
- road;
- rail;
- ship/water;
- air.

Each mode has distinct:
- cost;
- travel time;
- capacity;
- route/infrastructure dependencies;
- cargo compatibility;
- risk;
- handling requirements.

Hazardous cargo may be supported with higher cost, restrictions, insurance/risk and possible loss/disruption events. Deep disaster/crash simulation is not required for Release 1.

### 6. Bootstrap economy invariant

GRIDWORKS must remain economically functional with only one active player.

The GRIDWORKS system principal may provide bounded:
- buy orders/demand;
- sell orders/supply;
- contracts;
- transport opportunities;
- properties/businesses for sale;
- emergency supply;
- buyer-of-last-resort liquidity;
- market-seeding inventory.

System participation should recede as healthy player liquidity/depth develops.

Operator intervention may be available through admin/live-ops controls, but must be bounded, auditable and ledgered.

### 7. Starting experience

All players begin with the shared playable aggregate-plant introduction that teaches:

> inspect → diagnose → intervene → restore partial production → settle economically → make a second decision

The player retains the introductory facility after onboarding.

After the shared introduction, the player selects a preferred economic path such as:
- mining/quarrying;
- agriculture;
- manufacturing;
- energy;
- logistics;
- real estate;
- finance/markets;
- generalist/surprise me.

The selection is not a permanent class lock.

### 8. Release-1 posture

Release 1 must be commercially complete and compelling, not a thin shell.

Core domains such as industry, trade, property, finance, contracts, social/messaging, managers, skills, localization, markets and GRIDWORKS liquidity must exist in usable form.

Release 1 may defer depth that depends on mature systems, including examples such as:
- advanced derivatives;
- deep mergers/restructuring;
- rich interior/property decoration;
- detailed disaster simulation;
- highly advanced demographic simulation;
- unrestricted global chat;
- broad arbitrary user-uploaded media.

The default question for Release 1 planning is:

> **What can reasonably be deferred without making the launched product feel incomplete?**

not:

> What is the smallest thing we can ship?

### 9. Age/social posture

Design the online/social architecture for approximately **12+/13+** suitability.

No formal identity/age-verification requirement is part of the normal player flow.

Release-1 social controls must include:
- block;
- mute;
- report;
- rate limits;
- DM permissions;
- name moderation;
- impersonation/confusable controls;
- audit/moderation tooling.

Prefer generated/customizable avatars/logos before broad arbitrary user-uploaded media.

### 10. Asset and icon pipeline

Art production should be modular and data-driven.

Required foundations:
- style guide;
- camera/grid scale;
- material/LOD conventions;
- reusable industrial/transport/property/terrain kits;
- vector-first icon/symbol system where practical;
- stable semantic IDs for icons and visual assets;
- localization-safe UI graphics.

Examples:
- `resource.copper_ore`;
- `transport.rail`;
- `industry.real_estate`;
- `fault.motor_bearing_seizure`.

## Invariants

1. One logical world, many differentiated regions.
2. The economy must function with one active player.
3. GRIDWORKS system participation supports markets but does not permanently replace players.
4. Presentation does not become authoritative simulation logic.
5. Release 1 is a full commercial product, not an MVP shell.
6. No free-camera full-3D dependency for Release 1.
7. Transport modes have meaningful economic differences.
8. Social architecture is moderation-ready from launch.

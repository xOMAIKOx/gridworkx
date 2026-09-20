# Native Go service boundaries

WP-001 uses one Go workspace with four explicit process boundaries:

- `api`: public/application HTTP boundary;
- `worker`: scheduled and background work boundary;
- `realtime`: WebSocket/session delivery boundary;
- `sim-validator`: explicit server-side simulation validation boundary.

The shared runtime package provides only health/version smoke contracts. Domain orchestration, persistence, messaging and simulation validation remain separately authorized work-package concerns. No service starts automatically and no service unit is installed by WP-001.

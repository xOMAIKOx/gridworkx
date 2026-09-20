# GRIDWORKS Network Port Allocation

**Status:** ARCHITECTURE BASELINE  
**Date:** 2026-09-20  
**Purpose:** Define the GRIDWORKS application-owned port block and distinguish it from shared estate infrastructure ports.

## 1. GRIDWORKS application-owned block

Reserve the following contiguous TCP block in the Cortex estate port register:

> **18080-18089/TCP — GRIDWORKS**

### Active assignments

| Port | Protocol | Service | Purpose | Exposure |
|---|---|---|---|---|
| 18080 | TCP/HTTP | `gridworks-api` | Application HTTP API, health/version endpoint, future authenticated API traffic | Origin/internal only; reverse proxy may connect over approved private/WireGuard path |
| 18081 | TCP/HTTP | `gridworks-worker` | Worker health/version and future bounded operational/control endpoint if retained | Loopback/internal only; never public |
| 18082 | TCP/HTTP + WebSocket upgrade | `gridworks-realtime` | Realtime client sessions, messaging/notifications/market update fan-out | Origin/internal only; reverse proxy may connect over approved private/WireGuard path |
| 18083 | TCP/HTTP | `gridworks-sim-validator` | Simulation-validator health/version and future internal validation endpoint | Backend/internal only; never public |

### Reserved GRIDWORKS expansion ports

| Port | Status | Purpose |
|---|---|---|
| 18084 | RESERVED | Future GRIDWORKS process/split point if an accepted ADR requires one |
| 18085 | RESERVED | Future GRIDWORKS process/split point |
| 18086 | RESERVED | Future GRIDWORKS process/split point |
| 18087 | RESERVED | Future GRIDWORKS process/split point |
| 18088 | RESERVED | Future GRIDWORKS process/split point |
| 18089 | RESERVED | Future GRIDWORKS process/split point |

Reserved ports must not be assigned to another Cortex application without Architecture explicitly releasing them.

A reserved port does not authorize creation of a new service.

## 2. Public ingress

GRIDWORKS does not require a unique public TCP port.

Public application traffic should use the estate reverse-proxy layer on the normal shared web ports:

| Port | Protocol | Purpose | Ownership |
|---|---|---|---|
| 443 | TCP/HTTPS/WSS | Public API, game HTTPS traffic, WebSocket upgrade, admin web access as applicable | Shared estate reverse proxy |
| 80 | TCP/HTTP | Optional HTTP-to-HTTPS redirect only | Shared estate reverse proxy |

Ports 80/443 are shared estate infrastructure ports and are **not** part of the GRIDWORKS-owned 18080-18089 block.

The origin services should not be directly Internet-exposed.

## 3. Shared infrastructure dependency ports

These are dependencies GRIDWORKS may connect to. They are not application-owned ports and should not be reallocated as part of the GRIDWORKS block.

| Port | Protocol | Dependency | Direction from GRIDWORKS | Notes |
|---|---|---|---|---|
| 5432 | TCP | PostgreSQL | outbound/internal | Standard PostgreSQL service port unless the estate database deployment assigns a different registered port |
| 4222 | TCP | NATS / JetStream client | outbound/internal | Application client connection to NATS |
| 443 | TCP/HTTPS | S3-compatible object storage, identity/OIDC, external APIs/content endpoints | outbound | Actual endpoint may be internal or external; HTTPS preferred |
| 8222 | TCP/HTTP | NATS monitoring | operations/admin only if deployed | GRIDWORKS application code should not require this port for normal runtime |

If the selected S3-compatible estate service, PostgreSQL instance or NATS deployment uses a non-standard estate port, that port remains owned by that shared service and should be referenced from configuration rather than copied into the GRIDWORKS application block.

## 4. Client/device networking

Godot/mobile clients require no inbound listening port.

Clients initiate outbound connections using:
- HTTPS on 443;
- WSS/WebSocket upgrade on 443;
- normal platform networking/DNS supplied by the device.

Clients must not connect directly to 18080-18089.

## 5. Admin application

The React admin application does not currently require a dedicated production application port. It may be built as static content and served through the approved web/reverse-proxy path.

A development server port such as Vite's default must not be treated as a production estate allocation.

If a dedicated admin origin process is later accepted, it should consume one of the reserved GRIDWORKS ports rather than claiming an unrelated port.

## 6. Binding policy

Default bindings remain:

```text
GRIDWORKS_API_BIND=127.0.0.1:18080
GRIDWORKS_WORKER_BIND=127.0.0.1:18081
GRIDWORKS_REALTIME_BIND=127.0.0.1:18082
GRIDWORKS_SIM_VALIDATOR_BIND=127.0.0.1:18083
```

A future deployment ADR may bind API/realtime to a private/WireGuard interface instead of loopback where the reverse proxy is on another host.

Any such change requires:
- exact source/destination path;
- firewall rule;
- reverse-proxy/origin authorization;
- port-register update if needed.

It must never imply direct public exposure of 18080-18089.

## 7. Port-register request

Requested registration:

```text
Application: GRIDWORKS / GridWorx
Owner: CortexSSG
TCP block: 18080-18089
18080 gridworks-api
18081 gridworks-worker
18082 gridworks-realtime
18083 gridworks-sim-validator
18084-18089 reserved for GRIDWORKS future accepted service split/expansion
Public ingress: shared 443 (80 redirect optional) via estate reverse proxy
Shared dependencies: PostgreSQL 5432, NATS 4222, HTTPS 443; not owned by GRIDWORKS
```

## 8. Constraints

- No UDP application port is currently required.
- No Docker/Podman/OCI/container networking is permitted.
- No new port implies authorization to deploy or start a service.
- Any new process requiring a port must first consume an available reserved GRIDWORKS port or obtain Architecture approval for a revised allocation.

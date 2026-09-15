# Architectural Decisions

## ADR-001: One reverse proxy per service, not a shared gateway

**Status:** Decided

**Context:** Where ignis terminates TLS itself, it needs a proxy in front of it. Other services (buem and more) need the same.

**Decision:** Each service gets its own proxy, deployed with it. No shared gateway.

**Reason:** A shared gateway ties every service's uptime to one process, and reconfiguring it for one service risks every other service behind it. A per-service proxy limits that to one service, and the same Caddy setup drops straight onto the next service.

**Scope:** This applies to the proxy a service deploys for itself. A platform ingress that terminates TLS for everything behind it is a different layer, and `environment/http` exists for that case (ADR-006).

**Rejected:** A single shared API gateway. Reasonable at a different scale, but the wrong fit for per-service ownership.

## ADR-002: Proxy and app in separate containers

**Status:** Decided

**Context:** Caddy and the app could run as two processes in one container.

**Decision:** They run as separate containers.

**Reason:** Docker expects one process per container. Two processes under a wrapper script lose supervision: if the app crashed, Docker wouldn't notice, since the script keeps running. Separate containers also let the proxy be reloaded or recreated without touching the app, which we relied on repeatedly during this work.

**Rejected:** One container with a process supervisor (e.g. `supervisord`). Extra complexity for nothing this deployment needs.

## ADR-003: Reuse one CA across container recreation

**Status:** Decided for local development; production open

**Context:** A fresh Caddy container generates a new, untrusted CA, breaking trust the browser already had.

**Decision:** Mount the host's Caddy CA storage into the container so it reuses the same CA.

**Reason:** Trust attaches to the CA, not the individual certificate (see Crosscutting Concepts). Reuse one CA and trust survives recreation, reloads, and adding more services signed by the same CA.

**Rejected:** Re-running `caddy trust` after every recreation. Defeats the point and doesn't scale past one machine.

**Open for production:** Mounting one developer's local CA is local-only. Production needs public HTTPS (Let's Encrypt) or a shared internal CA. Which one depends on where ignis sits relative to the internet, not yet decided.

## ADR-004: No credential on ignis; the network is the trust boundary

**Status:** Decided

**Context:** ignis should be reachable only by a trusted caller. An earlier revision gated Caddy on a static, shared `X-Api-Key`. That key was visible in the page source of its only caller, Building Configurator, which runs in the browser, so it stopped casual access and kept no real secret.

**Decision:** Remove the key. ignis checks no credential. It is deployed on an internal network reachable only over VPN, behind a platform that has already authenticated the end user, and reachability is decided by the port mapping.

**Reason:** The gate was not load-bearing where ignis actually runs, and it cost every caller a header, every environment an extra secret to rotate, and the local-development path a configuration step for no security it did not already have.

**Rejected:** Keeping the static key as a speed bump. It implied a guarantee the deployment did not provide, which is worse than stating plainly that there is none.

**Open:** A real caller credential, should ignis ever be deployed somewhere the network is not the trust boundary. A short-lived token from the orchestration layer and mutual TLS are both candidates; neither is chosen.

## ADR-005: Group containers into per-concern namespaces, named <service>-<role>

**Status:** Decided

**Context:** The energy-modelling system is several services (ignis, buem, and more), each with a proxy, app, and sometimes a database.

**Decision:** Group containers by concern into a Docker Compose namespace (project name), here `building-simulation`, shared by all building-modelling services. Name each container `<service>-<role>`: `ignis-app`, `ignis-db`, `ignis-reverse-proxy`, and later `buem-app`, `buem-reverse-proxy`, and so on.

**Reason:** A container name then says both which service it belongs to and what role it plays, and the namespace groups related services together. The host port is the only thing that can clash between them, so it is set per service (`HOST_PORT` or `HOST_HTTPS_PORT`); the internal port stays fixed, since container isolation keeps it from clashing.

**Rejected:** Directory-derived project names and ad-hoc container names. They don't convey role or concern and don't group cleanly as more services are added.

**Open:** When buem joins, two repos each declaring the same namespace will coexist but Compose may warn about "orphan" containers. A single top-level compose (`include:`) or a shared external network resolves it. Decide when buem is wired in.

## ADR-006: Two Compose environments, http and https, rather than one with a toggle

**Status:** Decided

**Context:** Running ignis locally meant standing up Caddy: a `caddy` CLI install on the host, a one-off `caddy trust`, a `CADDY_DATA_DIR` that fails validation when wrong, and a key in `env/proxy.env`. None of that is wanted by a developer who needs the API up to test against. In production TLS is often already terminated upstream, by a firewall or the orchestrator's own reverse proxy.

**Decision:** Split `environment/` into `http/` and `https/`, each with its own compose and env files and each usable with no further configuration. The developer picks a directory. `environment/http` runs the app and the database only; `environment/https` adds Caddy for TLS.

**Reason:** TLS is a property of a deployment, not of the service. Two directories say that plainly, and neither is the default, so nothing has to argue for one.

**Rejected:** A single environment with a profile or an override file selecting the proxy. Compose merges `volumes:` by appending, and a service left with both `build:` and `image:` can still fall back to building locally, so overrides here fail quietly rather than loudly.

**Shared:** The two dockerfiles stay at `environment/`. The published image is identical for both environments, and the publishing workflow builds from that one path; per-directory copies would give CI two identical files with no correct one to pick, and they would drift.

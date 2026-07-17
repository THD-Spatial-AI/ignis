# Architectural Decisions

## ADR-001: One reverse proxy per service, not a shared gateway

**Status:** Decided

**Context:** ignis needs TLS and an access check in front of it. Other services (buem and more) will need the same.

**Decision:** Each service gets its own proxy, deployed with it. No shared gateway.

**Reason:** A shared gateway ties every service's uptime to one process. Reconfiguring the proxy for one service (as we did to fix its CORS handling) would risk every other service behind it. A per-service proxy limits that to one service, and the same Caddy setup drops straight onto the next service.

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

## ADR-004: Static API key as a prototype, not production auth

**Status:** Open

**Context:** ignis should be reachable only by a trusted caller. The real caller, an orchestration layer, doesn't exist yet. Building Configurator stands in for it.

**Decision (provisional):** Caddy gates on a static, shared `X-Api-Key`.

**Reason:** Enough to prove the proxy can enforce "only a known caller gets in." Not enough for production: Building Configurator runs in the browser, so the key is visible in its page source.

**Rejected for now, not yet chosen between:** a short-lived token from the real orchestration layer, or mutual TLS giving each service its own identity. Both wait until that layer exists. Revisit then.

## ADR-005: Group containers into per-concern namespaces, named <service>-<role>

**Status:** Decided

**Context:** The energy-modelling system is several services (ignis, buem, and more), each with a proxy, app, and sometimes a database.

**Decision:** Group containers by concern into a Docker Compose namespace (project name), here `building-simulation`, shared by all building-modelling services. Name each container `<service>-<role>`: `ignis-app`, `ignis-db`, `ignis-reverse-proxy`, and later `buem-app`, `buem-reverse-proxy`, and so on.

**Reason:** A container name then says both which service it belongs to and what role it plays, and the namespace groups related services together. The host port is the only thing that can clash between them, so it is set per service via `HOST_HTTPS_PORT`; the internal port stays fixed, since container isolation keeps it from clashing.

**Rejected:** Directory-derived project names and ad-hoc container names. They don't convey role or concern and don't group cleanly as more services are added.

**Open:** When buem joins, two repos each declaring the same namespace will coexist but Compose may warn about "orphan" containers. A single top-level compose (`include:`) or a shared external network resolves it. Decide when buem is wired in.

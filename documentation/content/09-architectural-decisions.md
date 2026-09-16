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

## ADR-005: One Compose project per repository and transport, with a shared-network alias as the cross-repository name

**Status:** Decided

**Context:** The energy-modelling system is several services (ignis, buem, city2tabula, weather and more), each with an app, sometimes a proxy and sometimes a database, expected to run side by side on one developer machine. Both of ignis's own environments must be able to run at the same time.

**Decision:** Name the Compose project after the repository and the transport: `ignis-http` and `ignis-https`. Name services by role, `db`, `build-db` and `proxy`, except the one that joins the shared network, which takes the repository name, `ignis`. Declare no `container_name`, so Compose generates a name per project. Let each project hold its own database volume. Attach only the HTTP environment to `tentacron-net`, with the explicit alias `ignis-http`.

**Reason:** Four separate mechanisms, each chosen for a measured failure.

The project name is what Compose takes destructive action against, not a label for related things. `docker compose down --remove-orphans` removes every container carrying the project label, so a project shared with another repository lets a command run there destroy this stack, silently and without naming what it took.

Fixed container names are what prevented the two environments running together, since a container name is unique across the host whatever the project is called. Generated names are unique by construction. Nothing outside the repository resolves them, because the cross-repository name is the alias.

Separate database volumes are forced by running both at once: two Postgres containers must never share one data directory. The cost is seeding each environment once, and the data is rebuilt from the workbook rather than authored, so a second copy costs disk rather than provenance.

Only the HTTP environment attaches to the shared network. Callers inside that network carry no TLS trust configuration, so an HTTPS alias fails at handshake and cannot be fixed in their config, only in their code. HTTPS is therefore host-facing and browser-facing.

A service name is registered as a DNS alias on every network the service joins, and an explicit `aliases:` entry adds to it rather than replacing it. Two services sharing a name on one shared network therefore round-robin between different repositories' stacks, with nothing reporting a fault. That is why the service on the shared network carries the repository name and the others, which never join it, stay short.

**Rejected:** One namespace per concern (`building-simulation`), shared by every building-modelling service. Grouping related containers is worth wanting, but a Compose project is the wrong place to put it. Splitting by transport alone (`building-simulation-http`) leaves the cross-repository case, which is the one that destroys containers.

Pinning the database volume to one name across both environments, which this decision reverses. It was correct while only one environment could run, and is corruption once both can.

`external: true` for the shared network, which fails to start when the network is absent and so breaks a clean checkout. The minimal `name:` declaration creates it on first use instead. It must be written with the same key in every repository: Compose labels the network with the key and refuses to start against a mismatch, while the `name:` values agree and both files look correct. The key is therefore spelled identically to the network name, so that a field which is really a cross-repository contract does not read as a local label.

A shared network declared with a driver, subnet or ipam block. A joining project whose options differ attempts to delete and recreate the network, which is a cross-repository destructive action of the kind this decision exists to remove.

**Open:** None. Cross-repository callers resolve the `ignis-http` alias on `tentacron-net`, which is the interface and survives any rename of project, service or container. The earlier answer, a caller attaching to `building-simulation_default` and resolving `ignis-app` by container name, failed because a project's default network takes its name from the project, so it is an implementation detail and a caller relying on it gets no signal when it changes: the service starts normally and fails at its first outbound call.

## ADR-006: Two Compose environments, http and https, rather than one with a toggle

**Status:** Decided

**Context:** Running ignis locally meant standing up Caddy: a `caddy` CLI install on the host, a one-off `caddy trust`, a `CADDY_DATA_DIR` that fails validation when wrong, and a key in `env/proxy.env`. None of that is wanted by a developer who needs the API up to test against. In production TLS is often already terminated upstream, by a firewall or the orchestrator's own reverse proxy.

**Decision:** Split `environment/` into `http/` and `https/`, each with its own compose and env files and each usable with no further configuration. The developer picks a directory. `environment/http` runs the app and the database only; `environment/https` adds Caddy for TLS.

**Reason:** TLS is a property of a deployment, not of the service. Two directories say that plainly, and neither is the default, so nothing has to argue for one.

**Rejected:** A single environment with a profile or an override file selecting the proxy. Compose merges `volumes:` by appending, and a service left with both `build:` and `image:` can still fall back to building locally, so overrides here fail quietly rather than loudly.

**Shared:** The two dockerfiles stay at `environment/`. The published image is identical for both environments, and the publishing workflow builds from that one path; per-directory copies would give CI two identical files with no correct one to pick, and they would drift.

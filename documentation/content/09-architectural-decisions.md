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

## ADR-006: Multi-stage builds — one binary, one image, one job

**Status:** Decided

**Context:** The original `environment/Dockerfile` was single-stage: it installed the full Go toolchain, `build-essential`, `git`, and `wget`, then compiled all three of `cmd/`'s binaries (`ignis`, `build_db`, `validate`) into one image, shipping the compiler and full source tree in the final container. Its own default `CMD` was an interactive `bash` shell rather than the server, because the image doubled as a shared toolbox for whichever of the three binaries someone needed at the moment — not a single-purpose deployable service.

**Decision:** Split into purpose-specific images. `environment/Dockerfile` is now multi-stage (`golang:1.26.0-bookworm AS builder` compiles, `debian:bookworm-slim` final stage holds only the compiled `ignis` binary + `ca-certificates`/`wget`), runs as non-root, and its `CMD` runs the server directly — no shell, no override needed in Compose. `environment/Dockerfile.build_db` is a second, separate multi-stage image holding only the `build_db` binary, wired into `docker-compose.yml` as `ignis-build-db`, gated behind Compose's `profiles: ["seed"]` so it never runs on a normal `docker compose up`. `validate` is dropped from production images entirely — it is a standalone scientific-correctness CLI tool (`cmd/validate`, checks calculated output against expected values within a tolerance), unrelated to the live request path, dev/CI-only.

**Reason:** `build_db` prints its own warning on every run — "This will DROP and recreate all country tables!" — so it must never be reachable from a container's automatic startup path; the old single-image design had no mechanism preventing that. Splitting by purpose also gives each image the same single-responsibility discipline already applied to this project's code: one job per image, named for that job. Measured effect: `ignis-app` went from ~1.39GB/304MB compressed (single-stage, full toolchain + source shipped) to 232MB/66MB (~4-5x smaller); `build_db` is a separate 175MB/47MB image. Neither ships a compiler, `git`, or the source tree, shrinking the attack surface available to anyone who ever gets a shell in a running container.

**Rejected:** Keeping one shared toolbox image and relying on a documented convention ("don't run build_db by accident") to prevent it running automatically — rejected because it depends on every future operator remembering a warning rather than the system preventing the mistake structurally.

## ADR-007: Bake static reference data into the image; don't mount it

**Status:** Decided

**Context:** `build_db` reads exactly one sheet, `Calc.Set.Building`, from a TABULA Webtool workbook (`internal/db/table_constructor.go`) to seed Postgres. The original `tabula-calculator.xlsx` (28MB) was bind-mounted into the `ignis-build-db` container at runtime, on the assumption that the data might change independently of the image.

**Decision:** That assumption was wrong for this data specifically — it's static source-of-truth reference data, read once to seed the database, with all real modification happening in Postgres afterward, never in the spreadsheet. A trimmed derivative, `data/tabula-calculator-lite.xlsx` (11MB — `Calc.Set.Building` plus the original `Info`/attribution sheet, every other sheet removed), is now `COPY`'d directly into `environment/Dockerfile.build_db` at build time. The bind mount is removed from `docker-compose.yml` entirely.

**Reason:** Baking in static data makes a given `build_db` image tag fully reproducible — it always seeds the exact same building-variant dataset, with no dependency on a particular file existing on whatever host happens to run it. General rule going forward: **does this data change independently of a new image build?** Yes → mount it (bind mount/named volume/env var, per ADR-002's separation-of-concerns logic). No → `COPY` it into the image. The trimmed file also drops sheets `build_db` never reads (`Calc.Set.System` alone, unused, was larger uncompressed than the sheet actually used), and remains under the original workbook's CC BY 4.0 license (see `data/README.md`, `ATTRIBUTIONS.md`).

**Rejected:** Keeping the bind mount and just documenting "don't change this file casually" — rejected for the same reason as ADR-006's rejected alternative: a convention that depends on remembering it, instead of the system enforcing it by construction. Also rejected: trimming via direct zip/XML surgery on the `.xlsx` — Excel's shared-strings/style/calc-chain indices cross-reference each other, and a manual edit risks corrupting those references; trimmed instead via Excel/LibreOffice's own "delete sheet" + save.

## ADR-008: Publish ignis-app to GHCR on version tags; deploy via a standalone compose file

**Status:** Decided

**Context:** Today, deploying `ignis-app` anywhere requires the full repo cloned and the Go toolchain present on that machine, because `docker-compose.yml`'s `ignis-app` entry uses `build: context: ./.. dockerfile: environment/Dockerfile`. `go.yml` (the existing CI) only compiles and tests; nothing publishes a deployable artifact.

**Decision:** Added `.github/workflows/docker-publish.yml`: on a pushed `vX.Y.Z` tag (or manual dispatch), it builds `environment/Dockerfile` and pushes `ghcr.io/thd-spatial-ai/ignis:vX.Y.Z` + `:latest` to GitHub Container Registry, using the built-in `GITHUB_TOKEN` (no extra secret). Added `environment/docker-compose.prod.yml`: a **standalone** deploy-time file (not a partial override merged with `docker-compose.yml`) that references `image: ghcr.io/thd-spatial-ai/ignis:${IGNIS_IMAGE_TAG:-latest}` for `ignis-app` instead of `build:`, alongside unchanged `ignis-db`/`ignis-reverse-proxy` entries (both already pull public images and needed no change). `ignis-build-db` is intentionally not included in the prod file — it's a rare, manual admin operation, not part of the standing deploy path.

**Reason:** This is the "Package" step of the CI pipeline chain (Checkout → Compile → Test → Build → Package → E2E → Deployable System) that `go.yml` alone never reached. A deploy server running off `docker-compose.prod.yml` needs only Docker itself plus small config files (`docker-compose.prod.yml`, `docker.env`, `Caddyfile`) — no source, no compiler. Tagging with the version (not just `latest`) makes rollback possible, the same principle as the `kubectl rollout history`/`rollback` already used elsewhere. The prod file is standalone rather than a merged override specifically because Compose's merge rules can leave a service with both `build:` and `image:` set, in which case Compose can still fall back to building locally if no matching image is cached — silently defeating the point on a source-less deploy target.

**Rejected:** A partial override file (`docker-compose.override.yml`-style) containing only the changed `image:` key, merged via `-f docker-compose.yml -f docker-compose.prod.yml` — rejected due to the build/image merge ambiguity above. Also rejected: triggering the publish job on every push to `main` — rejected in favor of tag-triggered publishing, so pushing an image to the registry is a deliberate release action, not an automatic side effect of every commit.

**Open:** The real production Caddyfile still needs its site block changed from `localhost` (a local-development-only address relying on Caddy's own non-public CA, per ADR-003) to the actual public domain, at which point Caddy obtains a real Let's Encrypt certificate automatically. Not yet done — noted here so it isn't lost.

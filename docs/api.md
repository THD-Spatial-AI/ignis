# API reference

The interactive reference below is generated from the OpenAPI spec
([`openapi.yaml`](openapi.yaml)), the machine-readable source of truth for
every endpoint, schema, and error. Download it to generate a client, load it
into Postman, or import it into another tool.

## How to consume the API

ignis has no authentication of its own. It runs behind a reverse proxy, and
that proxy is the only way in: any backend service that wants to use ignis
sends a valid `X-Api-Key` header, and the proxy rejects anything else with
`403` before ignis ever sees the request. ignis is meant to be called by a
trusted server-side caller, never directly by a browser or an end user's
client, since the key must not be visible outside that caller. If the system
in front of ignis has its own user-facing login (EnerPlanET, for example,
uses Keycloak), that login authenticates the user **to that system**, not to
ignis; that system's own backend then calls ignis on the user's behalf, using
the API key. The full model is in the spec's description.

!!! note "Base URL"
    All paths are served through the reverse proxy. In local development that is
    `https://localhost`; in a deployment it is whatever host the proxy is
    published on.

## Testing it yourself

The Swagger UI below can call a locally running ignis directly.

1. Start the stack, from `environment/`: `docker compose up -d`. On a first
   run, load the TABULA data once: `docker compose exec ignis-app ./bin/build_db`.
2. Serve these docs locally with `mkdocs serve`. The reverse proxy already
   allows requests from `http://localhost:8000` (its default port).
3. If your browser has never trusted the local proxy's certificate, open
   `https://localhost` directly once and accept it (or run `caddy trust`).
4. Click **Authorize** below and enter the API key checked by the reverse
   proxy (`X-Api-Key`; the prototype default is set in
   `environment/Caddyfile`). It applies to every **Try it out** call from
   then on.
5. Expand an endpoint, click **Try it out**, fill in the parameters, and
   **Execute**.

<swagger-ui src="openapi.yaml"/>

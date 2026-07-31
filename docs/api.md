# API reference

## What you can do

- **Look up a building's typical energy profile.** Give ignis a country, a
  building type (e.g. single-family house), and a construction period, and
  `/api/v1/variants/{country}/match` returns the matching TABULA archetypes —
  the "existing state" and its refurbishment levels (medium, advanced) — each
  with a code you can pass to the next two endpoints.
- **See the full physical parameter set behind a result.** `/api/v1/data/{code}`
  returns every input the calculation used for that archetype: envelope areas,
  U-values, climate data, solar gains, and more (~200 fields). `/api/v1/fields`
  explains what each one means, in plain language.
- **Calculate the annual heating demand.** `/api/v1/calculate/{code}` runs the
  ISO 13790 pipeline and returns `q_h_nd` in kWh/(m²·a).
- **Override any input to model a specific scenario**, instead of accepting
  the archetype's textbook defaults — the same calculate endpoint accepts an
  optional JSON body:

    | You want to change | Field(s) |
    |---|---|
    | The building's actual floor area | `A_ref` |
    | How cold/mild the local climate is | `HeatingDays`, `Theta_e`, `theta_i` |
    | How much sun the site gets | `I_Sol_South`, `I_Sol_East`, `I_Sol_West`, `I_Sol_North`, `I_Sol_Hor` |
    | Extra heat loss at wall/roof/window junctions | `delta_U_ThermalBridging_Original`, `delta_U_ThermalBridging_Refurbished` |

    Every field is independent and optional — send just the one you want to
    change, and every other input still comes from the archetype's TABULA
    defaults. See `CalculateRequest` in the reference below for exact types
    and validation rules (e.g. `HeatingDays` and the solar-irradiance fields
    must not be negative).

!!! example "A milder winter than the archetype assumes"
    ```bash
    curl -sk https://localhost/api/v1/calculate/DE.N.SFH.01.Gen.ReEx.001.001 \
      -H "X-Api-Key: ..." -H "Content-Type: application/json" \
      -d '{"HeatingDays": 150}'
    ```
    Returns the same archetype's `q_h_nd`, recalculated for a 150-heating-day
    winter instead of the TABULA default — every other input unchanged.

## Interactive reference

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

**Step 1:** Start the stack, from `environment/`: `docker compose -f docker-compose.quickstart.yml up -d`.
   On a first run, load the TABULA data once: `docker compose -f docker-compose.quickstart.yml --profile seed run --rm ignis-build-db`.

**Step 2:** Serve these docs locally with `mkdocs serve`. The reverse proxy already
   allows requests from `http://localhost:8000` (its default port).

**Step 3:** If your browser has never trusted the local proxy's certificate, open
   `https://localhost` directly once and accept it (or run `caddy trust`).

**Step 4:** Click **Authorize** below and enter the API key checked by the reverse
   proxy (`X-Api-Key`; the prototype default is set in
   `environment/env/proxy.env`). It applies to every **Try it out** call from
   then on.

**Step 5:** Expand an endpoint, click **Try it out**, fill in the parameters, and
   **Execute**.

!!! bug "Swagger UI Bug"
    The Swagger UI might not load correctly, just reload the browser page and it should load correctly.

<swagger-ui src="openapi.yaml"/>

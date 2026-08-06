# API reference

The interactive reference below is generated from the OpenAPI spec ([`openapi.yaml`](openapi.yaml)). Download it to generate a client or import it into Postman.

## Authentication

ignis has no authentication of its own. The reverse proxy in front of it is the only way in: callers send a valid `X-Api-Key` header, and the proxy rejects anything else with `403` before the request reaches ignis.

!!! warning "Call it server-side only"
    The API key must not be visible outside the calling service. Call ignis from a trusted server-side caller, never directly from a browser or an end user's client. A user-facing login in front of your own system (EnerPlanET uses Keycloak) authenticates the user to *that system*, not to ignis. That system's backend then calls ignis on the user's behalf using the API key.

!!! note "Base URL"
    All paths are served through the reverse proxy. In local development that is `https://localhost`; in a deployment it is whatever host the proxy is published on.

## Call sequence

| Step | Endpoint | Returns |
|---|---|---|
| 1 | `GET /api/v1/variants/{country}/match?type=...&period=...` | The matching TABULA archetypes for that country, building type and construction period: the existing state and its refurbishment levels (medium, advanced). Each carries a `code` used by the next two endpoints. |
| 2 | `GET /api/v1/data/{code}` | Every physical input behind that archetype (roughly 200 fields): envelope areas, U-values, climate data, solar gains. `GET /api/v1/fields` gives a plain-language description of each field. |
| 3 | `POST /api/v1/calculate/{code}` | Annual heating demand `q_h_nd` in kWh/(m²·a), from the ISO 13790 pipeline. |

Step 2 is optional. Call it when you need to inspect or report the inputs; `calculate` uses them either way.

## Overriding archetype inputs

`calculate` accepts an optional JSON body. Every field is independent and optional: send only what you want to change, and every other input still comes from the archetype's TABULA defaults.

| You want to change | Field(s) |
|---|---|
| The building's actual floor area | `A_ref` |
| How cold or mild the local climate is | `HeatingDays`, `Theta_e`, `theta_i` |
| How much sun the site gets | `I_Sol_South`, `I_Sol_East`, `I_Sol_West`, `I_Sol_North`, `I_Sol_Hor` |
| Extra heat loss at wall, roof and window junctions | `delta_U_ThermalBridging_Original`, `delta_U_ThermalBridging_Refurbished` |
| The building's actual room height or storey count | `h_room`, `n_Storey` |
| How airtight the building is, or how much it is ventilated | `n_air_infiltration`, `n_air_use` |
| How much heat the building's structure can store | `c_m` |

See `CalculateRequest` in the reference below for exact types and validation rules; `HeatingDays` and the solar-irradiance fields must not be negative.

!!! example "Overriding a single field"
    ```bash
    curl -sk https://localhost/api/v1/calculate/DE.N.SFH.01.Gen.ReEx.001.001 -H "X-Api-Key: ..." -H "Content-Type: application/json" -d '{"HeatingDays": 150}'
    ```
    Returns the same archetype's `q_h_nd` recalculated for a 150-heating-day winter, every other input unchanged.

## Describing real surfaces

The body also accepts an optional `surfaces` list, one entry per physical element, in place of the archetype's generic wall and window slots.

- Each entry takes an `id`, a `type` (`wall`, `window`, `roof`, `floor`), an `area` in m², a `u_value` in W/(m²·K), and an `azimuth`.
- ignis merges the entries in each category into the single area and area-weighted U-value the calculation needs.
- A category you do not list keeps the archetype's default for that category.
- `azimuth` is degrees clockwise from North (0 North, 90 East, 180 South, 270 West). For windows it also selects which direction's solar irradiance the window counts toward, rounded to the nearest of North, East, South or West. A window with no `azimuth` is assumed to face South.

??? example "A building with five listed surfaces"
    ```json
    {
      "surfaces": [
        {"id": "wall-north", "type": "wall",   "area": 45.0, "u_value": 0.85, "azimuth": 0},
        {"id": "wall-south", "type": "wall",   "area": 45.0, "u_value": 0.85, "azimuth": 180},
        {"id": "win-south",  "type": "window", "area": 8.0,  "u_value": 1.2,  "azimuth": 180},
        {"id": "roof-1",     "type": "roof",   "area": 90.0, "u_value": 0.4,  "azimuth": -1},
        {"id": "floor-1",    "type": "floor",  "area": 90.0, "u_value": 0.5,  "azimuth": -1}
      ]
    }
    ```

    List as many windows as the building has, each with its own area, U-value and orientation, rather than fitting them into two slots.

## Testing it yourself

The Swagger UI below can call a locally running ignis directly.

**Step 1:** Start the stack, from `environment/`: `docker compose -f docker-compose.quickstart.yml up -d`. On a first run, load the TABULA data once: `docker compose -f docker-compose.quickstart.yml --profile seed run --rm ignis-build-db`.

**Step 2:** Serve these docs locally with `mkdocs serve`. The reverse proxy already allows requests from `http://localhost:8000` (its default port).

**Step 3:** If your browser has never trusted the local proxy's certificate, open `https://localhost` directly once and accept it, or run `caddy trust`.

**Step 4:** Click **Authorize** below and enter the API key checked by the reverse proxy (`X-Api-Key`; the prototype default is set in `environment/env/proxy.env`). It applies to every **Try it out** call from then on.

**Step 5:** Expand an endpoint, click **Try it out**, fill in the parameters, then **Execute**.

!!! bug "Swagger UI bug"
    The Swagger UI sometimes fails to load. Reload the browser page and it should come up correctly.

<swagger-ui src="openapi.yaml"/>
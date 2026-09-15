# API reference

The interactive reference lives in its own standalone page, [`openapi/index.html`](openapi/index.html), so it can be opened directly without running `mkdocs serve`. It renders [`openapi/openapi.yaml`](openapi/openapi.yaml); download that file to generate a client or import it into Postman.

## Authentication

ignis carries no credential and asks for none. Every endpoint is reachable by anything that can open a connection to it, so who that is gets decided at the network layer: ignis is published on an internal network reachable over VPN, behind a platform that has already authenticated the end user.

!!! warning "Do not publish ignis on a public interface"
    There is nothing in ignis to stop an unauthenticated caller. Keep it on an internal network, or behind something that authenticates for it. A user-facing login in front of your own system (EnerPlanET uses Keycloak) authenticates the user to *that system*, not to ignis; that system's backend then calls ignis on the user's behalf.

!!! note "Base URL"
    Depends on which environment you started. `environment/http` publishes the app directly, so the base URL is `http://localhost:8080`. `environment/https` puts Caddy in front of it for TLS, giving `https://localhost`. In a deployment it is whatever host ignis is published on. Every path below is the same either way.

## Call sequence

| Step | Endpoint | Returns |
|---|---|---|
| 1 | `GET /api/v1/variants/{country}/match?type=...&period=...` | The matching TABULA archetypes for that country, building type and construction period: the existing state and its refurbishment levels (medium, advanced). Each carries a `code` used by the next two endpoints. |
| 2 | `GET /api/v1/data/{code}` | Every physical input behind that archetype (roughly 200 fields): envelope areas, U-values, climate data, solar gains. `GET /api/v1/fields` gives a plain-language description of each field. |
| 3 | `POST /api/v1/calculate/{code}` | Annual heating demand `q_h_nd` in kWh/(m²·a), from the ISO 13790 pipeline. |

Step 2 is optional. Call it when you need to inspect or report the inputs; `calculate` uses them either way.

If you have a construction year rather than a period code, send `year=` in place of `period=` at step 1 and ignis resolves it. Exactly one of the two is required. `year` must be a non-negative integer no later than the current calendar year. A year before the earliest or after the latest band defined for that country and type resolves to the oldest or newest period. `GET /api/v1/periods/{country}` lists a country's bands as `{ period, year_from, year_to }`, oldest first, with `year_from` 0 and `year_to` 9999 marking the open-ended bands. Step 2's response also carries the resolved variant's own band, at `tabula_data.BasicParameters.BuildingAppearance.Year1_Building` / `Year2_Building`.

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
    curl -s http://localhost:8080/api/v1/calculate/DE.N.SFH.01.Gen.ReEx.001.001 -H "Content-Type: application/json" -d '{"HeatingDays": 150}'
    ```
    Returns the same archetype's `q_h_nd` recalculated for a 150-heating-day winter, every other input unchanged.

## Describing real surfaces

The body also accepts an optional `surfaces` list, one entry per physical element, in place of the archetype's generic wall and window slots.

- Each entry takes an `id`, a `type` (`wall`, `window`, `roof`, `floor`, `door`), an `area` in m², a `u_value` in W/(m²·K), and, for windows, an `azimuth` and a `tilt`.
- ignis merges the entries in each category into the single area and area-weighted U-value the calculation needs.
- A category you do not list keeps the archetype's default for that category.
- You own the accuracy of the geometry you send. ignis uses the areas, U-values and orientations as given and does not check them against the archetype.
- `azimuth` is degrees clockwise from North (0 North, 90 East, 180 South, 270 West). For windows it also selects which direction's solar irradiance the window counts toward, rounded to the nearest of North, East, South or West. A window with no `azimuth` is assumed to face South.
- `tilt` is degrees from horizontal: 0 is a flat skylight, 90 a vertical window. A window within 5 degrees of horizontal counts toward horizontal irradiance rather than a compass direction. city2tabula measures tilt the other way round (0 for a vertical wall, 90 for a flat roof), so convert its values with `ignis_tilt = 90 - c2t_tilt` before sending them.

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

[Open the API reference](openapi/index.html), which can call a locally running ignis directly, no `mkdocs serve` needed to view it.

**Step 1:** Start the stack, from `environment/http/`: `docker compose -f docker-compose.prod.yml up -d`. On a first run, load the TABULA data once: `docker compose -f docker-compose.prod.yml --profile seed run --rm ignis-build-db`.

**Step 2:** Serve `docs/openapi/` on `http://localhost:8000` (`python -m http.server 8000` from that directory works), since `ALLOWED_ORIGINS` allows that origin already. Opening the file directly (`file://`) works for reading the reference, but **Try it out** needs an allowed origin.

**Step 3:** Pick `http://localhost:8080` from the **Servers** dropdown, then expand an endpoint, click **Try it out**, fill in the parameters, and **Execute**.

!!! info "Using the HTTPS environment instead"
    `environment/https` serves the same API on `https://localhost`. Select that server in the dropdown, and on the quickstart file open `https://localhost` in a tab once and accept the certificate warning first: browser JavaScript cannot click through it the way a manual page load can.
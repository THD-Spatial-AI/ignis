---
audience: developer
---

# Validation

The pipeline is validated against the TABULA reference values across 2,147 buildings in 20 European countries, with a tolerance of ±2.5% on `q_h_nd`. Last run: November 2025.

---

## Results by country

| Country | Buildings | Pass rate |
|---|---|---|
| Austria | 165 | 100% |
| Belgium | 99 | 100% |
| Bulgaria | 78 | 100% |
| Cyprus | 36 | 100% |
| Czech Republic | 84 | 100% |
| Denmark | 91 | 100% |
| France | 120 | 100% |
| Germany | 232 | 100% |
| Greece | 144 | 100% |
| Hungary | 45 | 100% |
| Ireland | 118 | 100% |
| Italy | 106 | 100% |
| Netherlands | 135 | 100% |
| Norway | 69 | 100% |
| Poland | 78 | 100% |
| Serbia | 111 | 100% |
| Slovenia | 112 | 100% |
| Sweden | 171 | 100% |
| United Kingdom | 81 | 100% |
| Spain | 72 | 22.2% |

**Overall: 97.4% (2,091 / 2,147 buildings passing)**

---

## Spain: known issue

The failures are systematic: calculated values deviate from the TABULA reference in one direction, which points to a Spain-specific parameter rather than a general pipeline bug.

**Status:** under investigation.

### Measured 2026-09-23

A re-measurement of all 72 Spanish variants against TABULA defaults, taken after the table above. It was run through `POST /api/v1/calculate/{code}` with an empty request body, which reaches the same pipeline as `./bin/validate`:

| Measure | Value |
|---|---|
| Passing | 48 of 72 |
| Signed mean error | -2.32% |
| Rows computing below the reference | 71 of 72, none above |
| Failure range | 2.50% to 3.57% |

The offset is uniform rather than scattered: -2.32% for references below 5 kWh/(m²·a), -2.32% from 5 to 20, and -2.27% above 20. Germany, Italy and France show no such bias in a 25-row sample each.

!!! warning "Spain's pass count is unstable by construction"
    A country-wide offset of about 2.3% sits just inside the ±2.5% gate, so small differences move rows across the line. Stored references are rounded to one decimal, which at the lowest magnitudes here (references as low as 1.7) is worth up to 2.9% on its own. Any pass count published for Spain needs its measurement date attached.

---

## Validation methodology

- **Tolerance:** ±2.5% on `q_h_nd` (annual heating energy demand, kWh/(m²·a))
- **Reference:** TABULA Excel workbook (`data/tabula-calculator-lite.xlsx`)
- **Pipeline:** 17-level cascading calculation: geometry, envelope, U-values, climate, solar gains, thermal bridges, heat transfer coefficients

!!! note "Scope"
    This runs the pipeline directly against each database row's unmodified
    TABULA defaults. It never goes through the `/api/v1/calculate/{code}`
    HTTP handler, so it does not exercise that endpoint's optional overrides
    (`A_ref`, `surfaces` and the rest). Those are additive: a request with no
    body reaches the same code path this validation covers.

---

## Running validation

```bash
go build -buildvcs=false -o bin/ ./cmd/...
./bin/validate
```

Requires a populated database (`./bin/build_db` must have been run first) and a valid `.env` file.

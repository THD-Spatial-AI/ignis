# API reference

Base URL: `http://localhost:8080/api/v1`

All responses are JSON. All endpoints are read-only — ignis does not modify the database.

---

## Health check

```
GET /health
```

Returns `200 OK` when the server is running. Does not check the database connection.

**Response**

```json
{ "status": "OK" }
```

---

## List variants

```
GET /api/v1/variants/{country_iso2}
```

Returns all available building variant codes for a country.

**Path parameters**

| Parameter | Description | Example |
|---|---|---|
| `country_iso2` | ISO 3166-1 alpha-2 country code | `DE` |

**Response**

```json
{
  "country": "germany",
  "data": ["DE.N.SFH.01.Gen", "DE.N.SFH.01.ReEx", ...]
}
```

---

## Match variants

```
GET /api/v1/variants/{country_iso2}/match?type={type}&period={period}
```

Returns all refurbishment levels for a specific building type and construction period, ordered from existing state to most-refurbished.

**Path parameters**

| Parameter | Description | Example |
|---|---|---|
| `country_iso2` | ISO 3166-1 alpha-2 country code | `DE` |

**Query parameters**

| Parameter | Description | Example |
|---|---|---|
| `type` | Building type code | `SFH`, `MFH`, `TH`, `AB` |
| `period` | Construction period code | `01`, `02`, `03` |

**Response**

```json
{
  "country": "germany",
  "prefix": "DE.N.SFH.01",
  "data": [
    { "code": "DE.N.SFH.01.Gen", "label": "Existing state" },
    { "code": "DE.N.SFH.01.ReEx", "label": "Medium refurbishment" },
    { "code": "DE.N.SFH.01.ReAd", "label": "Advanced refurbishment" }
  ]
}
```

---

## Get variant data

```
GET /api/v1/data/{code}
```

Returns the raw TABULA parameters for a building variant. Used by client applications to populate their own forms or models.

**Path parameters**

| Parameter | Description | Example |
|---|---|---|
| `code` | TABULA variant code | `DE.N.SFH.01.Gen` |

**Response**

```json
{
  "country": "germany",
  "variant_code": "DE.N.SFH.01.Gen",
  "tabula_data": { ... },
  "expected_q_h_nd": 123.45
}
```

!!! note
    `tabula_data` contains the full set of ~200 TABULA parameters. `expected_q_h_nd` is the reference value from the workbook, used for validation.

---

## Calculate heat demand

```
POST /api/v1/calculate/{code}
```

Runs the 17-level ISO 13790 pipeline for the specified building variant and returns the annual heating energy demand.

**Path parameters**

| Parameter | Description | Example |
|---|---|---|
| `code` | TABULA variant code | `DE.N.SFH.01.Gen` |

**Request body** (optional)

```json
{ "A_ref": 150.0 }
```

`A_ref` overrides the reference floor area stored in the TABULA record. Omit the body to use the TABULA default.

**Response**

```json
{
  "variant_code": "DE.N.SFH.01.Gen",
  "q_h_nd": 123.45,
  "unit": "kWh/(m2.a)"
}
```

**Error responses**

| Status | Condition |
|---|---|
| `400` | Invalid variant code format, unknown country, invalid `A_ref` |
| `404` | Variant code not found in database |
| `500` | Pipeline execution failed or returned a non-finite result |

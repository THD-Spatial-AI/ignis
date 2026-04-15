# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Upstream data source

The building typology data and the 17-level calculation methodology both come from the **TABULA** project (Typology Approach for Building Stock Energy Assessment), coordinated by Institut Wohnen und Umwelt (IWU) under the EU Intelligent Energy Europe Programme (CC BY 4.0).

- WebTool: <https://episcope.eu/building-typology/tabula-webtool/>
- Source data download: <https://episcope.eu/welcome/>
- Full attribution: [ATTRIBUTIONS.md](ATTRIBUTIONS.md)

The Excel workbook (`data/tabula-calculator.xlsx`) is the authoritative reference for every formula in `internal/calc/`. When a calculation result diverges from TABULA, compare against the corresponding cell in the `Calc.Set.Building` sheet of that workbook.

## Overview

HDCP Go is a Go implementation of the **Heat Demand Calculation Pipeline (HDCP)** for estimating building energy performance across European countries. It processes TABULA building datasets (from an Excel workbook) into PostgreSQL, then runs a 17-level cascading calculation to produce `q_h_nd` — annual heating energy demand in kWh/(m²·a).

## Commands

### Build

```bash
go build -o bin/app        cmd/app/main.go       # HTTP API server
go build -o bin/build_db   cmd/build_db/main.go  # Load Excel → PostgreSQL
go build -o bin/validate   cmd/validate/main.go  # Run pipeline against all DB rows
```

### Run

```bash
# 1. Load Excel workbook into PostgreSQL (destructive — drops and recreates all country tables)
./bin/build_db

# 2. Validate pipeline accuracy against stored expected values
./bin/validate

# 3. Start HTTP API
./bin/app
```

### Lint / vet

```bash
go vet ./...
```

There are no automated test files in this repository. Validation is done via `cmd/validate`, which runs the pipeline against every building in the database and checks that calculated `q_h_nd` is within ±2.5% of the TABULA reference value.

## Configuration

Copy `.env.example` to `.env`. Key variables:

| Variable | Default | Purpose |
|---|---|---|
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_NAME` | `hdcp` | Database name |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | — | Database password |
| `DB_SSL_MODE` | `disable` | SSL mode |

The Excel file is auto-discovered from `data/*.xlsx`; default is `data/tabula-calculator.xlsx`.

## Architecture

### Three executables

| Binary | Package | Purpose |
|---|---|---|
| `cmd/app` | API server | Gin HTTP server exposing calculation and data endpoints |
| `cmd/build_db` | DB importer | Reads `Calc.Set.Building` sheet in the Excel workbook and bulk-loads one PostgreSQL table per country into the `tabula` schema |
| `cmd/validate` | Validation tool | Queries every row across all country tables and runs the pipeline in parallel goroutines, reporting pass/fail at ±2.5% tolerance |

### Calculation pipeline (`internal/hdcp` + `internal/calc`)

The core of the project. `Pipeline.Run()` in `internal/hdcp/pipeline.go` executes 17 levels sequentially. Each level is a struct (`CalcLevel1`…`CalcLevel17`) in `internal/calc/calc_level_NN.go`.

- **Level 0** (`models.TabulaBuildingParameters`) — raw TABULA data loaded from the database; this is the input to the pipeline, not a `CalcLevel` struct.
- **Levels 1–17** — each level takes specific prior levels as constructor arguments (not all levels depend on their immediate predecessor). See `pipeline.go` for the exact dependency graph.
- **Final output** — `CalcLevel17.QHNd` = `q_ht - eta_h_gn * (q_sol + q_int)` in kWh/(m²·a).

The formula in every `calcXxx()` method mirrors the corresponding Excel cell formula from `data/tabula-calculator.xlsx`. Keep this correspondence intact when fixing or extending calculations.

### Data model (`internal/models/tabula.go`)

`TabulaBuildingParameters` is the single input type for the pipeline. It has two top-level sub-trees:

- `BasicParameters` — `BuildingThematic` (codes, storey count, room height) + `Envelope` (all surface areas and volumes)
- `AdvancedParameters` — 16 nested structs covering U-values, insulation, solar gains, thermal bridges, climate conditions, etc.

JSON struct tags map directly to PostgreSQL column names (and Excel column headers). The reflection-based `populateStructFromMap` in `cmd/validate/main.go` and the equivalent `populateStructFromMap` in `internal/db/repository/tabula_repository.go` both use these tags to hydrate the struct from a `map[string]interface{}` row.

### Database layout

All TABULA tables live in the `tabula` schema. Each country gets its own table named after the lowercase country name with spaces replaced by `_` (e.g., `tabula.germany`, `tabula.united_kingdom`). The table name is derived from the ISO 3166-1 alpha-2 prefix of the building variant code (e.g., `DE` → `germany`) via `TabulaCountryHelper` in `internal/utils/tabula_country_helper.go`.

### HTTP API (`internal/api`)

Built with **Gin**. Routes (all under `/api/v1`):

| Method | Path | Description |
|---|---|---|
| `GET` | `/health` | Health check |
| `GET` | `/api/v1/variants/:country_iso2` | List all variant codes for a country |
| `GET` | `/api/v1/data/:code` | Return raw TABULA record for a variant |
| `POST` | `/api/v1/calculate/:code` | Run the pipeline; optional body `{"A_ref": 150.0}` overrides the reference floor area |

The `Handler` struct holds a `TabulaRepository`; the `HDCPService` is instantiated per-request (stateless). The variant code's first two characters are the ISO2 country code used to resolve the correct table.

### Key conventions

- `Code_BuildingVariant` format: `<ISO2>.<ResidentialType>.<BuildingType>.<Index>.<Variant>` (e.g., `DE.N.SFH.01.Gen`).
- Country table names use snake_case (`united_kingdom`, not `United Kingdom`).
- All physical quantities carry SI units documented in struct field comments (m², W/m²K, kWh/(m²·a), etc.).
- Panic recovery is used in `Pipeline.handleError` so a single level failure degrades gracefully rather than crashing the server.

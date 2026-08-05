# Data

This directory contains the TABULA Webtool workbook used to seed the heat demand database.

## tabula-calculator.xlsx

| Field | Detail |
|-------|--------|
| **File** | `tabula-calculator.xlsx` |
| **Source** | TABULA Webtool — [building-typology.eu](https://webtool.building-typology.eu/) |
| **Author** | Institut Wohnen und Umwelt (IWU), Darmstadt, Germany |
| **Project** | Intelligent Energy Europe, IEE/09/739/SI2.558245 |
| **License** | [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/) |

The workbook contains per-country building typology data (U-values, areas, infiltration rates, climate parameters) and reference heating demand values (`q_h_nd`) used to validate the calculation pipeline within ±2.5 %.

**Citation:**

> Loga, T., Stein, B., Diefenbach, N., Born, R. (2016): *Deutsche Wohngebäudetypologie. Beispielhafte Maßnahmen zur Verbesserung der Energieeffizienz von typischen Wohngebäuden.* 2nd edition. Institut Wohnen und Umwelt, Darmstadt.

This file is distributed under the same CC BY 4.0 terms as the original TABULA dataset. See [`ATTRIBUTIONS.md`](../ATTRIBUTIONS.md) for the full attribution statement.

## tabula-calculator-lite.xlsx

A trimmed derivative of `tabula-calculator.xlsx`, created for this project and published alongside it on GitHub. `build_db` only ever reads the `Calc.Set.Building` sheet (see `internal/db/table_constructor.go`) — every other sheet in the full TABULA Webtool workbook (country-specific system/demo calculations, charts, auxiliary climate tables, etc.) is dead weight for this tool's purposes, and made up the majority of the original file's size.

Sheets kept:
- **`Info`** — the original workbook's own attribution/info sheet, kept as-is so provenance travels with the file itself, not just this README.
- **`Calc.Set.Building`** — the only sheet `build_db` actually reads.
- **`BlankSheet`** (hidden) — a template artifact from the original workbook; harmless, left in rather than risking broken internal references by removing it.

Result: **11 MB**, down from the original **28 MB** — everything else (roughly two-thirds of the file, including a single sheet, `Calc.Set.System`, that alone exceeded the size of the sheet actually used) was unused by this tool.

This is a derivative of TABULA Webtool data and remains under the same **CC BY 4.0** license as the original — see attribution above and in [`ATTRIBUTIONS.md`](../ATTRIBUTIONS.md). It is what `environment/ignis-db.dockerfile` bakes into the `ignis-build-db` image (see that file's comments): the data is static source-of-truth reference data, read once to seed Postgres and never modified independently afterward, so shipping it inside the image — rather than mounting it in at runtime — makes a given image tag fully reproducible.

## Usage

The `build_db` binary reads whichever `.xlsx` it finds in this directory (`filepath.Glob("data/*.xlsx")`, first match alphabetically — currently `tabula-calculator-lite.xlsx`, since both files satisfy the same `Calc.Set.Building`-only requirement):

```bash
go build -o bin/build_db cmd/build_db/main.go
./bin/build_db
```

Inside Docker (`environment/ignis-db.dockerfile`), only `tabula-calculator-lite.xlsx` is present in the image at all, so there's no ambiguity there.

!!! warning "Destructive operation"
    Running `build_db` drops and recreates all TABULA country tables. Do not run against a database that holds production data without a backup.

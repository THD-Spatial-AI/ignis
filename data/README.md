# Data

This directory contains the TABULA Webtool workbook used to seed the heat demand database.

## tabula-calculator.xlsx

| Field | Detail |
|-------|--------|
| **File** | `tabula-calculator.xlsx` |
| **Source** | TABULA Webtool, [building-typology.eu](https://webtool.building-typology.eu/) |
| **Author** | Institut Wohnen und Umwelt (IWU), Darmstadt, Germany |
| **Project** | Intelligent Energy Europe, IEE/09/739/SI2.558245 |
| **Credit** | IEE Projects TABULA + EPISCOPE (www.episcope.eu) |
| **Terms of use** | [TABULA usage rules](https://episcope.eu/communication/download/): non-exclusive use, with the credit above visibly mentioned as the source |

The workbook contains per-country building typology data (U-values, areas, infiltration rates, climate parameters) and reference heating demand values (`q_h_nd`) used to validate the calculation pipeline within ±2.5%.

**Citation:**

> Loga, T., Stein, B., Diefenbach, N., Born, R. (2016): *Deutsche Wohngebäudetypologie. Beispielhafte Maßnahmen zur Verbesserung der Energieeffizienz von typischen Wohngebäuden.* 2nd edition. Institut Wohnen und Umwelt, Darmstadt.

This file is distributed under the TABULA usage rules, not under the MIT licence of this repository. See [`ATTRIBUTIONS.md`](../ATTRIBUTIONS.md) for the full attribution statement.

## tabula-calculator-lite.xlsx

A trimmed derivative of `tabula-calculator.xlsx`, created for this project and published alongside it on GitHub. `build_db` reads only the `Calc.Set.Building` sheet (see `internal/db/table_constructor.go`). Every other sheet in the full TABULA Webtool workbook (country-specific system and demo calculations, charts, auxiliary climate tables) is unused here and made up most of the original file's size.

Sheets kept:
- **`Info`**: the original workbook's attribution sheet, kept so provenance travels with the file rather than only in this README.
- **`Calc.Set.Building`**: the only sheet `build_db` reads.
- **`BlankSheet`** (hidden): a template artefact from the original workbook, left in place because removing it risks breaking internal references.

The result is **11 MB**, down from **28 MB**. The removed sheets were roughly two-thirds of the file, including `Calc.Set.System`, which alone was larger than the sheet in use.

This is a derivative of TABULA Webtool data and is distributed under the same TABULA usage rules as the original; see the attribution above and in [`ATTRIBUTIONS.md`](../ATTRIBUTIONS.md). `environment/ignis-db.dockerfile` bakes it into the `ignis-build-db` image. The data is static reference data, read once to seed Postgres and never modified afterwards, so shipping it inside the image rather than mounting it at runtime makes a given image tag reproducible.

## Usage

The `build_db` binary reads whichever `.xlsx` it finds in this directory (`filepath.Glob("data/*.xlsx")`, first match alphabetically, currently `tabula-calculator-lite.xlsx`, since both files satisfy the same `Calc.Set.Building` requirement):

```bash
go build -o bin/build_db cmd/build_db/main.go
./bin/build_db
```

Inside Docker (`environment/ignis-db.dockerfile`) only `tabula-calculator-lite.xlsx` is present, so there is no ambiguity.

!!! warning "Destructive operation"
    Running `build_db` drops and recreates all TABULA country tables. Do not run against a database that holds production data without a backup.

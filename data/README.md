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

The workbook contains per-country building typology data (U-values, areas, infiltration rates, climate parameters) and reference heating demand values (`q_h_nd`) used to validate the HDCP pipeline within ±2.5 %.

**Citation:**

> Loga, T., Stein, B., Diefenbach, N., Born, R. (2016): *Deutsche Wohngebäudetypologie. Beispielhafte Maßnahmen zur Verbesserung der Energieeffizienz von typischen Wohngebäuden.* 2nd edition. Institut Wohnen und Umwelt, Darmstadt.

This file is distributed under the same CC BY 4.0 terms as the original TABULA dataset. See [`ATTRIBUTIONS.md`](../ATTRIBUTIONS.md) for the full attribution statement.

## Usage

The `build_db` binary reads this file and loads it into PostgreSQL:

```bash
go build -o bin/build_db cmd/build_db/main.go
./bin/build_db
```

!!! warning "Destructive operation"
    Running `build_db` drops and recreates all TABULA country tables. Do not run against a database that holds production data without a backup.

![Ignis logo](docs/assets/logo/ignis-logo-dark.svg#gh-dark-mode-only)
![Ignis logo](docs/assets/logo/ignis-logo-light.svg#gh-light-mode-only)

[![CI](https://github.com/thd-spatial-ai/ignis/actions/workflows/ci.yml/badge.svg)](https://github.com/thd-spatial-ai/ignis/actions/workflows/ci.yml)&nbsp;&nbsp;&nbsp;[![MkDocs](https://github.com/thd-spatial-ai/ignis/actions/workflows/docs.yml/badge.svg)](https://thd-spatial-ai.github.io/ignis)&nbsp;&nbsp;&nbsp;[![Go](https://github.com/THD-Spatial-AI/ignis/actions/workflows/go.yml/badge.svg)](https://github.com/THD-Spatial-AI/ignis/actions/workflows/go.yml)&nbsp;&nbsp;&nbsp;[![codecov](https://codecov.io/gh/THD-Spatial-AI/ignis/graph/badge.svg?token=CTUZED1ELJ)](https://codecov.io/gh/THD-Spatial-AI/ignis)&nbsp;&nbsp;&nbsp;[![GitHub release](https://img.shields.io/github/v/release/thd-spatial-ai/ignis?include_prereleases&label=release&logo=github)](https://github.com/thd-spatial-ai/ignis/releases)

Go microservice implementing the **EN ISO 13790** annual heating energy demand calculation pipeline derived from [tabula-calculator.xlsx](https://episcope.eu/welcome/) *(Accessed on: 26.06.26)*. The calculation method has been documented in [TABULA CommonCalculationMethod](https://episcope.eu/fileadmin/tabula/public/docs/report/TABULA_CommonCalculationMethod.pdf) *(Accessed on: 26.06.2026)*. The tool covers all European building typologies across 20 countries defined by **TABULA & EPISCOPE (IEE Projects)**.

Results are validated against the workbook's own output: **19 of 20 countries at 100% accuracy, 2,091 of 2,147 buildings passing.** See the [validation report](docs/validation.md).

---

## Compatibility

| Dependency | Version |
| ---------- | ------- |
| Go | 1.26+ |
| PostgreSQL | 15 to 17 |

---

## Quick start

| Step | Command | Description |
| ---- | ------- | ----------- |
| 1 | `cp .env.example .env` | Configure DB connection, ALLOWED_ORIGINS, APP_PORT |
| 2 | `make build` | Compile all binaries into bin/ |
| 3 | `make create-db` | Create the PostgreSQL database named in .env |
| 4 | `./bin/build_db` | Load TABULA workbook inside PostgreSQL |
| 5 | `make run` | Start API on APP_PORT (default 8080) |

Full setup and API documentation: [thd-spatial-ai.github.io/ignis](https://thd-spatial-ai.github.io/ignis)

Architecture documentation (arc42): under development, not yet published.

---

## Testing

![test coverage](https://codecov.io/github/THD-Spatial-AI/ignis/graphs/icicle.svg?token=CTUZED1ELJ)

In the graphic above, the top section is the whole project, the sections below it are folders, and the smallest are individual files. Slice size is the number of statements and slice colour is the coverage.

> [!NOTE]
> **AI usage disclaimer:** Some tests in this repository were written with the assistance of AI coding tools, then reviewed and validated by a maintainer before merging.

```bash
go test ./...
go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -html=coverage.out
```

---

## Local docs

```bash
python -m venv .venv
.venv/bin/pip install -r docs/requirements.txt
.venv/bin/mkdocs serve
```

---

## License

MIT License. Copyright 2026 BigGeoData & Spatial AI, Technische Hochschule Deggendorf. See [LICENSE](LICENSE) for the full text. The licence covers the code. The TABULA workbooks in `data/` are distributed under the TABULA usage rules; see [ATTRIBUTIONS.md](ATTRIBUTIONS.md).

Found a security issue? See [SECURITY.md](SECURITY.md) for how to report it privately.

## Acknowledgements

Developed in the context of the RENvolveIT research project (<https://projekte.ffg.at/projekt/5127011>), funded by CETPartnership under the 2023 joint call for research proposals, co-funded by the European Commission (GA N°101069750).

<img src="docs/assets/sponsors/CETP-logo.svg" alt="CETPartnership" width="144" height="72">&nbsp;&nbsp;&nbsp;<img src="docs/assets/sponsors/EN_Co-fundedbytheEU_RGB_POS.png" alt="EU" width="180" height="40">

Building-characteristic data source: **IEE Projects TABULA + EPISCOPE (www.episcope.eu)**, used under the [TABULA usage rules](https://episcope.eu/communication/download/) (accessed 08.07.2026).

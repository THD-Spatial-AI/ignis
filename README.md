![Ignis logo](docs/assets/logo/ignis-logo-dark.svg#gh-dark-mode-only)
![Ignis logo](docs/assets/logo/ignis-logo-light.svg#gh-light-mode-only)

[![CI](https://github.com/thd-spatial-ai/ignis/actions/workflows/ci.yml/badge.svg)](https://github.com/thd-spatial-ai/ignis/actions/workflows/ci.yml)&nbsp;&nbsp;&nbsp;[![MkDocs](https://github.com/thd-spatial-ai/ignis/actions/workflows/docs.yml/badge.svg)](https://thd-spatial-ai.github.io/ignis)&nbsp;&nbsp;&nbsp;[![codecov](https://codecov.io/gh/THD-Spatial-AI/ignis/graph/badge.svg?token=CTUZED1ELJ)](https://codecov.io/gh/THD-Spatial-AI/ignis)&nbsp;&nbsp;&nbsp;[![GitHub release](https://img.shields.io/github/v/release/thd-spatial-ai/ignis?include_prereleases&label=release&logo=github)](https://github.com/thd-spatial-ai/ignis/releases)

Go microservice implementing the **EN ISO 13790** annual heating energy demand calculation pipeline derived from [tabula-calculator.xlsx](https://episcope.eu/welcome/) *(Accessed on: 26.06.26)*. The calculation method has been documented in [TABULA CommonCalculationMethod](https://episcope.eu/fileadmin/tabula/public/docs/report/TABULA_CommonCalculationMethod.pdf) *(Accessed on: 26.06.2026)*. The tool covers all European building typologies across 20 countries defined by **TABULA & EPISCOPE (IEE Projects)**.

The results have been validated against the Excel Workbook output. So far **19/20 countries at 100% accuracy in total 2,091 / 2,147 buildings validated.** See the [validation report](docs/validation.md).

---

## Compatibility

| Dependency | Version |
| ---------- | ------- |
| Go | 1.26+ |
| PostgreSQL | >= 15 |

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

MIT License — Copyright 2026 BigGeoData & Spatial AI, Technische Hochschule Deggendorf. See [LICENSE](LICENSE) for the full text.

Found a security issue? See [SECURITY.md](SECURITY.md) for how to report it privately.

## Acknowledgements

Developed in the context of the RENvolveIT research project (<https://projekte.ffg.at/projekt/5127011>), funded by CETPartnership under the 2023 joint call for research proposals, co-funded by the European Commission (GA N°101069750).

<img src="docs/assets/sponsors/CETP-logo.svg" alt="CETPartnership" width="144" height="72">&nbsp;&nbsp;&nbsp;<img src="docs/assets/sponsors/EN_Co-fundedbytheEU_RGB_POS.png" alt="EU" width="180" height="40">

**TABULA & EPISCOPE (IEE Projects):** building-characteristic data ([episcope.eu](https://episcope.eu/iee-project/tabula/), accessed 08.07.2026)

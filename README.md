![Ignis logo](docs/assets/logo/ignis-logo-dark.svg#gh-dark-mode-only)
![Ignis logo](docs/assets/logo/ignis-logo-light.svg#gh-light-mode-only)

[![CI](https://github.com/thd-spatial-ai/ignis/actions/workflows/ci.yml/badge.svg)](https://github.com/thd-spatial-ai/ignis/actions/workflows/ci.yml) [![MkDocs](https://github.com/thd-spatial-ai/ignis/actions/workflows/docs.yml/badge.svg)](https://thd-spatial-ai.github.io/ignis) [![GitHub release](https://img.shields.io/github/v/release/thd-spatial-ai/ignis?include_prereleases&label=release&logo=github)](https://github.com/thd-spatial-ai/ignis/releases)

Go microservice implementing the **EN ISO 13790** annual heating energy demand calculation pipeline derived from [tabula-calculator.xlsx](https://episcope.eu/welcome/) *(Accessed on: 26.06.26)*. The calculation method has been documented in [TABULA CommonCalculationMethod](https://episcope.eu/fileadmin/tabula/public/docs/report/TABULA_CommonCalculationMethod.pdf) *(Accessed on: 26.06.2026)*. The tool covers all European building typologies across 20 countries defined by **TABULA & EPISCOPE (IEE Projects)**.

The results have been validated against the Excel Workbook output. So far **19/20 countries at 100% accuracy — 2,091 / 2,147 buildings validated.** See the [validation report](docs/validation.md).

---

## Compatibility

| Dependency | Version |
|---|---|
| Go | 1.26+ |
| PostgreSQL | 15 – 17 |

---

## Quick start

```bash
cp .env.example .env          # configure DB connection and ALLOWED_ORIGINS
go build -o bin/build_db cmd/build_db/main.go
./bin/build_db                # load TABULA workbook → PostgreSQL
go build -o bin/app cmd/app/main.go
./bin/app                     # start API on :8080
```

Full setup and API documentation: [thd-spatial-ai.github.io/ignis](https://thd-spatial-ai.github.io/ignis)

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

## Acknowledgements

Developed in the context of the RENvolveIT research project (<https://projekte.ffg.at/projekt/5127011>), funded by CETPartnership under the 2023 joint call for research proposals, co-funded by the European Commission (GA N°101069750).

<img src="docs/assets/sponsors/CETP-logo.svg" alt="CETPartnership" width="144" height="72">&nbsp;&nbsp;&nbsp;<img src="docs/assets/sponsors/EN_Co-fundedbytheEU_RGB_POS.png" alt="EU" width="180" height="40">

**TABULA & EPISCOPE (IEE Projects):** building-characteristic data ([episcope.eu](https://episcope.eu/iee-project/tabula/), accessed 13.11.2025)

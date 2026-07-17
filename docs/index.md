# ignis

Go microservice implementing the ISO 13790 annual heating energy demand calculation pipeline, validated against the [TABULA](https://episcope.eu/building-typology/tabula-webtool/) European building typology database.

ignis is designed as an internal microservice. Given a TABULA building variant code, it returns the annual heating energy demand in kWh/(m²·a) via a simple REST API.

**19/20 countries at 100% accuracy: 2,091 / 2,147 buildings validated.** See the [validation report](validation.md).

---

## Documentation

| Section | Description |
|---|---|
| [Getting started](getting-started.md) | Installation, configuration, building, running |
| [API reference](api.md) | All endpoints with request/response examples |
| [Validation](validation.md) | Per-country accuracy results |
| [Architecture (PDF)](documentation.pdf) | Full arc42 architecture documentation |

---

## Repository

[github.com/thd-spatial-ai/ignis](https://github.com/thd-spatial-ai/ignis) · [Issue tracker](https://github.com/thd-spatial-ai/ignis/issues) · [Contributing](https://github.com/thd-spatial-ai/ignis/blob/main/CONTRIBUTING.md)

# HDCP Go

A Go implementation of the Heat Demand Calculation Pipeline (HDCP) for estimating and validating building energy performance across European countries.

**19/20 countries at 100% validation accuracy — 2,091/2,147 buildings passing validation.**

See [MULTI_COUNTRY_VALIDATION_REPORT.md](MULTI_COUNTRY_VALIDATION_REPORT.md) for detailed results.

---

## Project Structure

```
hdcp-go/
├── cmd/
│   ├── app/            # HTTP API server
│   ├── build_db/       # Database rebuild tool
│   └── validate/       # Validation tool
├── data/
│   ├── tabula_models/  # Exported JSON models (auto-generated, gitignored)
│   └── tabula-calculator.xlsx
├── examples/
│   └── batch_by_code.json
├── internal/
│   ├── api/            # HTTP handlers
│   ├── calc/           # 17-level calculation pipeline
│   ├── config/         # Configuration and environment
│   ├── db/             # Database access
│   ├── hdcp/           # Pipeline orchestration
│   ├── models/         # Data models
│   ├── service/        # Business logic
│   └── utils/          # Helpers
└── go.mod
```

---

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Excel workbook with Tabula building data (`data/tabula-calculator.xlsx`)

---

## Configuration

Copy `.env.example` to `.env` and fill in your values:

```bash
cp .env.example .env
```

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=hdcp
DB_SSL_MODE=disable

# Country for validation
COUNTRY=germany

# Data files
EXCEL_FILE=data/tabula-calculator.xlsx
```

---

## Building

```bash
go build -o bin/validate   cmd/validate/main.go
go build -o bin/build_db   cmd/build_db/main.go
go build -o bin/app        cmd/app/main.go
```

---

## Usage

### 1. Build the database

Load building data from the Excel workbook into PostgreSQL:

```bash
./bin/build_db
```

### 2. Validate calculations

```bash
# Single country
./bin/validate -country germany

# All countries
./test_all_countries.sh
```

### 3. Run the API server

```bash
./bin/app
```

---

## Validation Methodology

- **Tolerance:** ±2% on final `q_h_nd` (annual heating energy demand)
- **Metric:** kWh/(m²·a)
- **Pipeline:** 17-level cascading calculation covering building geometry, envelope, thermal properties, climate conditions, solar gains, thermal bridges, heat transfer coefficients, and final energy demand

---

## Supported Countries

Austria, Belgium, Bulgaria, Cyprus, Czech Republic, Denmark, France, Germany, Greece, Hungary, Ireland, Italy, Netherlands, Norway, Poland, Serbia, Slovenia, Sweden, United Kingdom

> Spain is not yet producing valid results (under investigation).

---

## Contributing

Bug reports, feature requests, and pull requests are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md) and the [Code of Conduct](CODE_OF_CONDUCT.md) before getting started.

---

## License

MIT License — Copyright (c) 2026 Technische Hochschule Deggendorf (THD-Spatial-AI). See [LICENSE](LICENSE) for the full text.

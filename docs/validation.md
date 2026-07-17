# Validation

Pipeline validated against the TABULA reference values (±2% tolerance on `q_h_nd`) across 2,147 buildings in 20 European countries. Last run: November 2025.

---

## Results by country

| Country | Buildings | Pass rate |
|---|---|---|
| 🇦🇹 Austria | 165 | 100% ✓ |
| 🇧🇪 Belgium | 99 | 100% ✓ |
| 🇧🇬 Bulgaria | 78 | 100% ✓ |
| 🇨🇾 Cyprus | 36 | 100% ✓ |
| 🇨🇿 Czech Republic | 84 | 100% ✓ |
| 🇩🇰 Denmark | 91 | 100% ✓ |
| 🇫🇷 France | 120 | 100% ✓ |
| 🇩🇪 Germany | 232 | 100% ✓ |
| 🇬🇷 Greece | 144 | 100% ✓ |
| 🇭🇺 Hungary | 45 | 100% ✓ |
| 🇮🇪 Ireland | 118 | 100% ✓ |
| 🇮🇹 Italy | 106 | 100% ✓ |
| 🇳🇱 Netherlands | 135 | 100% ✓ |
| 🇳🇴 Norway | 69 | 100% ✓ |
| 🇵🇱 Poland | 78 | 100% ✓ |
| 🇷🇸 Serbia | 111 | 100% ✓ |
| 🇸🇮 Slovenia | 112 | 100% ✓ |
| 🇸🇪 Sweden | 171 | 100% ✓ |
| 🇬🇧 United Kingdom | 81 | 100% ✓ |
| 🇪🇸 Spain | 72 | 22.2% ⚠ |

**Overall: 97.4% (2,091 / 2,147 buildings passing)**

---

## Spain: known issue

56 of 72 buildings fail (77.8%). The failures are systematic (calculated values consistently deviate from the TABULA reference), which points to a Spain-specific parameter rather than a general pipeline bug.

**Status:** under investigation.

### Investigation checklist

- [ ] Compare climate zone parameters against Mediterranean-specific values
- [ ] Check construction type and attic/cellar condition code handling for Spanish archetypes
- [ ] Validate thermal bridge calculations against Spanish building standards
- [ ] Field-by-field comparison of a failing building against the TABULA Excel workbook
- [ ] Cross-reference with the Spanish national energy performance methodology (CTE DB-HE)

---

## Validation methodology

- **Tolerance:** ±2% on `q_h_nd` (annual heating energy demand, kWh/(m²·a))
- **Reference:** TABULA Excel workbook (`data/tabula-calculator.xlsx`)
- **Pipeline:** 17-level cascading calculation: geometry, envelope, U-values, climate, solar gains, thermal bridges, heat transfer coefficients

---

## Running validation

```bash
go build -o bin/validate cmd/validate/main.go
./bin/validate
```

Requires a populated database (`./bin/build_db` must have been run first) and a valid `.env` file.

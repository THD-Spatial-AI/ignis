package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/thd-spatial-ai/ignis/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrVariantNotFound indicates that a building variant was not found in the requested table.
var ErrVariantNotFound = errors.New("tabula variant not found")

// TabulaRepository provides read access to TABULA datasets.
type TabulaRepository struct {
	pool   *pgxpool.Pool
	schema string
}

// NewTabulaRepository constructs a new repository instance.
func NewTabulaRepository(pool *pgxpool.Pool, schema string) *TabulaRepository {
	return &TabulaRepository{pool: pool, schema: schema}
}

// ListVariants returns the available building variant codes for a given country table.
func (r *TabulaRepository) ListVariants(ctx context.Context, tableName string) ([]string, error) {
	query := fmt.Sprintf(`SELECT "Code_BuildingVariant" FROM %s ORDER BY "Code_BuildingVariant"`, r.qualifyTable(tableName))

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query variants: %w", err)
	}
	defer rows.Close()

	var variants []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, fmt.Errorf("failed to scan variant code: %w", err)
		}
		variants = append(variants, code)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate variant rows: %w", err)
	}

	return variants, nil
}

// MatchVariants returns all variant codes whose Code_BuildingVariant starts with the given prefix.
// The prefix should be in the form "CC.N.TYPE.PERIOD" (e.g. "DE.N.SFH.01") — a trailing "."
// is appended automatically so that only variants of that exact building+period are returned.
func (r *TabulaRepository) MatchVariants(ctx context.Context, tableName, prefix string) ([]string, error) {
	pattern := prefix + ".%"
	query := fmt.Sprintf(
		`SELECT "Code_BuildingVariant" FROM %s WHERE "Code_BuildingVariant" LIKE $1 ORDER BY "Code_BuildingVariant"`,
		r.qualifyTable(tableName),
	)

	rows, err := r.pool.Query(ctx, query, pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to match variants: %w", err)
	}
	defer rows.Close()

	var codes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, fmt.Errorf("failed to scan variant code: %w", err)
		}
		codes = append(codes, code)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate matched variants: %w", err)
	}

	return codes, nil
}

// ConstructionPeriod is one TABULA construction-year band for a country.
type ConstructionPeriod struct {
	Period   string `json:"period"`
	YearFrom int    `json:"year_from"`
	YearTo   int    `json:"year_to"`
}

// periodFromYearClass turns a Code_ConstructionYearClass ("AT.01") into the bare
// period index ("01") used by Code_BuildingVariant and the match endpoint.
func periodFromYearClass(yearClass string) string {
	if i := strings.LastIndex(yearClass, "."); i >= 0 {
		return yearClass[i+1:]
	}
	return yearClass
}

// ResolvePeriodByYear returns the period index whose construction-year band
// contains year, scoped to one building type. typePrefix is "CC.N.TYPE".
// Year1_Building 0 and Year2_Building 9999 are the open-ended sentinels some
// countries use; others simply stop at their oldest or newest recorded
// period, so a year before/after every band for that type clamps to the
// nearest edge (see clampToNearestPeriod) rather than returning no match.
// Returns "" only when the type has no bands at all, or the year falls in a
// gap between two defined bands.
func (r *TabulaRepository) ResolvePeriodByYear(ctx context.Context, tableName, typePrefix string, year int) (string, error) {
	query := fmt.Sprintf(
		`SELECT "Code_ConstructionYearClass" FROM %s
		 WHERE "Code_BuildingVariant" LIKE $1
		   AND "Year1_Building" <= $2 AND "Year2_Building" >= $2
		 LIMIT 1`,
		r.qualifyTable(tableName),
	)

	var yearClass string
	err := r.pool.QueryRow(ctx, query, typePrefix+".%", year).Scan(&yearClass)
	if err == nil {
		return periodFromYearClass(yearClass), nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("failed to resolve period for %s year %d: %w", typePrefix, year, err)
	}

	yearClass, err = r.clampToNearestPeriod(ctx, tableName, typePrefix, year)
	if err != nil {
		return "", err
	}
	if yearClass == "" {
		return "", nil
	}
	return periodFromYearClass(yearClass), nil
}

// clampToNearestPeriod handles a year that matched no band for typePrefix. If
// year precedes every defined band it returns the oldest one; if it follows
// every band it returns the newest. A year that falls inside the type's
// overall range but between two defined bands (a gap) is left as "" — which
// neighbour to prefer there is undefined, so it stays a genuine no-match.
func (r *TabulaRepository) clampToNearestPeriod(ctx context.Context, tableName, typePrefix string, year int) (string, error) {
	boundsQuery := fmt.Sprintf(
		`SELECT MIN("Year1_Building"), MAX("Year2_Building") FROM %s WHERE "Code_BuildingVariant" LIKE $1`,
		r.qualifyTable(tableName),
	)

	var minYear1, maxYear2 *int
	if err := r.pool.QueryRow(ctx, boundsQuery, typePrefix+".%").Scan(&minYear1, &maxYear2); err != nil {
		return "", fmt.Errorf("failed to load year bounds for %s: %w", typePrefix, err)
	}
	if minYear1 == nil || maxYear2 == nil {
		return "", nil // no bands at all for this type
	}

	var edgeQuery string
	switch {
	case year < *minYear1:
		edgeQuery = fmt.Sprintf(
			`SELECT "Code_ConstructionYearClass" FROM %s
			 WHERE "Code_BuildingVariant" LIKE $1
			 ORDER BY "Year1_Building" ASC LIMIT 1`,
			r.qualifyTable(tableName),
		)
	case year > *maxYear2:
		edgeQuery = fmt.Sprintf(
			`SELECT "Code_ConstructionYearClass" FROM %s
			 WHERE "Code_BuildingVariant" LIKE $1
			 ORDER BY "Year2_Building" DESC LIMIT 1`,
			r.qualifyTable(tableName),
		)
	default:
		return "", nil // inside the overall range but in a gap between bands
	}

	var yearClass string
	if err := r.pool.QueryRow(ctx, edgeQuery, typePrefix+".%").Scan(&yearClass); err != nil {
		return "", fmt.Errorf("failed to clamp period for %s year %d: %w", typePrefix, year, err)
	}
	return yearClass, nil
}

// ListPeriods returns the country's distinct construction-year bands, oldest first.
func (r *TabulaRepository) ListPeriods(ctx context.Context, tableName string) ([]ConstructionPeriod, error) {
	query := fmt.Sprintf(
		`SELECT DISTINCT "Code_ConstructionYearClass", "Year1_Building", "Year2_Building"
		 FROM %s
		 WHERE "Code_ConstructionYearClass" IS NOT NULL
		   AND "Year1_Building" IS NOT NULL AND "Year2_Building" IS NOT NULL
		 ORDER BY "Year1_Building"`,
		r.qualifyTable(tableName),
	)

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query construction periods: %w", err)
	}
	defer rows.Close()

	var periods []ConstructionPeriod
	for rows.Next() {
		var yearClass string
		var from, to int
		if err := rows.Scan(&yearClass, &from, &to); err != nil {
			return nil, fmt.Errorf("failed to scan construction period: %w", err)
		}
		periods = append(periods, ConstructionPeriod{
			Period:   periodFromYearClass(yearClass),
			YearFrom: from,
			YearTo:   to,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate construction periods: %w", err)
	}

	return periods, nil
}

// GetVariant loads the full TABULA record and key metadata for a specific building variant.
func (r *TabulaRepository) GetVariant(ctx context.Context, tableName, buildingCode string) (*models.TabulaBuildingParameters, string, float64, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE "Code_BuildingVariant" = $1 LIMIT 1`, r.qualifyTable(tableName))

	rows, err := r.pool.Query(ctx, query, buildingCode)
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to query building data: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, "", 0, fmt.Errorf("failed to iterate building rows: %w", err)
		}
		return nil, "", 0, ErrVariantNotFound
	}

	dataMap, err := rowsToDataMap(rows)
	if err != nil {
		return nil, "", 0, err
	}

	tabulaData := initializeTabulaData()
	populateStructFromMap(tabulaData, dataMap)

	buildingID := fmt.Sprintf("%v", dataMap["Code_BuildingVariant"])
	expectedQHND := toFloat64(dataMap["q_h_nd"])

	return tabulaData, buildingID, expectedQHND, nil
}

func (r *TabulaRepository) qualifyTable(tableName string) string {
	if r.schema == "" {
		return pgx.Identifier{tableName}.Sanitize()
	}
	return pgx.Identifier{r.schema, tableName}.Sanitize()
}

func rowsToDataMap(rows pgx.Rows) (map[string]interface{}, error) {
	values, err := rows.Values()
	if err != nil {
		return nil, fmt.Errorf("failed to read row values: %w", err)
	}

	descriptions := rows.FieldDescriptions()
	dataMap := make(map[string]interface{}, len(descriptions))

	for i, fd := range descriptions {
		if i >= len(values) {
			continue
		}

		name := string(fd.Name)
		dataMap[name] = normalizeValue(values[i])
	}

	return dataMap, nil
}

func normalizeValue(value interface{}) interface{} {
	switch v := value.(type) {
	case []byte:
		return string(v)
	default:
		return v
	}
}

func toFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return f
		}
	}
	return 0
}

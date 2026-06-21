package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/THD-Spatial-AI/hdcp-go/internal/models"

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

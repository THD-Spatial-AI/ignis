package repository

import (
	"context"

	"github.com/thd-spatial-ai/ignis/internal/models"
)

// TabulaReader is the read interface used by API handlers.
// Using an interface instead of the concrete type makes handlers testable without a live database.
type TabulaReader interface {
	// ListVariants returns all building variant codes for a given country table.
	ListVariants(ctx context.Context, tableName string) ([]string, error)

	// MatchVariants returns variant codes whose Code_BuildingVariant starts with prefix + ".".
	// prefix should be "CC.N.TYPE.PERIOD" (e.g. "DE.N.SFH.01").
	MatchVariants(ctx context.Context, tableName, prefix string) ([]string, error)

	// ResolvePeriodByYear returns the TABULA period index (e.g. "03") whose
	// construction-year band contains year, scoped to one building type.
	// typePrefix is "CC.N.TYPE" (e.g. "DE.N.SFH"). Returns "" if no band for
	// that type contains the year.
	ResolvePeriodByYear(ctx context.Context, tableName, typePrefix string, year int) (string, error)

	// ListPeriods returns the country's construction-year bands, oldest first.
	// year_from 0 and year_to 9999 are passed through as the open-ended sentinels.
	ListPeriods(ctx context.Context, tableName string) ([]ConstructionPeriod, error)

	// GetVariant loads the full TABULA record for a specific building variant code.
	// Returns the building parameters, the variant code string, the reference q_h_nd, and any error.
	GetVariant(ctx context.Context, tableName, variantCode string) (*models.TabulaBuildingParameters, string, float64, error)
}

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

	// GetVariant loads the full TABULA record for a specific building variant code.
	// Returns the building parameters, the variant code string, the reference q_h_nd, and any error.
	GetVariant(ctx context.Context, tableName, variantCode string) (*models.TabulaBuildingParameters, string, float64, error)
}

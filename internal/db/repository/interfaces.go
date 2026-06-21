package repository

import (
	"context"

	"github.com/THD-Spatial-AI/hdcp-go/internal/models"
)

// TabulaReader is the read interface used by API handlers.
// Using an interface instead of the concrete type makes handlers testable without a live database.
type TabulaReader interface {
	// ListVariants returns all building variant codes for a given country table.
	ListVariants(ctx context.Context, tableName string) ([]string, error)

	// GetVariant loads the full TABULA record for a specific building variant code.
	// Returns the building parameters, the variant code string, the reference q_h_nd, and any error.
	GetVariant(ctx context.Context, tableName, variantCode string) (*models.TabulaBuildingParameters, string, float64, error)
}

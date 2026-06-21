package service

import (
	"context"
	"fmt"
	"github.com/THD-Spatial-AI/hdcp-go/internal/db/repository"
	"github.com/THD-Spatial-AI/hdcp-go/internal/hdcp"
	"github.com/THD-Spatial-AI/hdcp-go/internal/models"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Re-export types for easier use in handlers
type TabulaBuildingParameters = models.TabulaBuildingParameters

// HDCPService provides business logic for HDCP calculations
type HDCPService struct {
	logger     *hdcp.Logger
	repository *repository.BuildingRepository
}

// NewHDCPService creates a new HDCP service instance without database
func NewHDCPService() *HDCPService {
	return &HDCPService{
		logger: hdcp.NewLogger(log.New(os.Stdout, "", 0)),
	}
}

// NewHDCPServiceWithDB creates a new HDCP service with database support.
// schema is the PostgreSQL schema name (e.g. "tabula").
func NewHDCPServiceWithDB(pool *pgxpool.Pool, schema string) *HDCPService {
	return &HDCPService{
		logger:     hdcp.NewLogger(log.New(os.Stdout, "", 0)),
		repository: repository.NewBuildingRepository(pool, schema),
	}
}

// NewHDCPServiceWithLogger creates a new HDCP service with custom logger
func NewHDCPServiceWithLogger(logger *hdcp.Logger) *HDCPService {
	return &HDCPService{
		logger: logger,
	}
}

// GetBuildingByCode retrieves building parameters from database by building code
func (s *HDCPService) GetBuildingByCode(ctx context.Context, country, buildingCode string) (*models.TabulaBuildingParameters, error) {
	if s.repository == nil {
		return nil, fmt.Errorf("database repository not initialized")
	}

	// Normalize country name for table name
	tableName := normalizeCountryName(country)

	return s.repository.GetByBuildingCode(ctx, tableName, buildingCode)
}

// normalizeCountryName normalizes country names for table names
func normalizeCountryName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	return name
}

// CalculateHeatingDemand executes the HDCP calculation pipeline.
// Returns the calculated q_h_nd (annual heating energy demand in kWh/(m²·a)).
func (s *HDCPService) CalculateHeatingDemand(buildingParams *models.TabulaBuildingParameters) (float64, error) {
	pipeline := hdcp.NewPipeline(buildingParams, s.logger)
	return pipeline.Run()
}

// CalculateHeatingDemandWithDetails executes the HDCP calculation pipeline
// and returns the fully populated Pipeline struct for inspection of intermediate levels.
func (s *HDCPService) CalculateHeatingDemandWithDetails(buildingParams *models.TabulaBuildingParameters) (*hdcp.Pipeline, error) {
	pipeline := hdcp.NewPipeline(buildingParams, s.logger)
	if _, err := pipeline.Run(); err != nil {
		return nil, err
	}
	return pipeline, nil
}

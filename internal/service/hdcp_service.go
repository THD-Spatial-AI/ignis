package service

import (
	"log"
	"os"

	"github.com/thd-spatial-ai/ignis/internal/hdcp"
	"github.com/thd-spatial-ai/ignis/internal/models"
)

// Re-export types for easier use in handlers
type TabulaBuildingParameters = models.TabulaBuildingParameters

// IgnisService provides business logic for HDCP calculations
type IgnisService struct {
	logger *hdcp.Logger
}

// NewIgnisService creates a new HDCP service instance
func NewIgnisService() *IgnisService {
	return &IgnisService{
		logger: hdcp.NewLogger(log.New(os.Stdout, "", 0)),
	}
}

// NewIgnisServiceWithLogger creates a new HDCP service with custom logger
func NewIgnisServiceWithLogger(logger *hdcp.Logger) *IgnisService {
	return &IgnisService{
		logger: logger,
	}
}

// CalculateHeatingDemand executes the HDCP calculation pipeline.
// Returns the calculated q_h_nd (annual heating energy demand in kWh/(m²·a)).
func (s *IgnisService) CalculateHeatingDemand(buildingParams *models.TabulaBuildingParameters) (float64, error) {
	pipeline := hdcp.NewPipeline(buildingParams, s.logger)
	return pipeline.Run()
}

// CalculateHeatingDemandWithDetails executes the HDCP calculation pipeline
// and returns the fully populated Pipeline struct for inspection of intermediate levels.
func (s *IgnisService) CalculateHeatingDemandWithDetails(buildingParams *models.TabulaBuildingParameters) (*hdcp.Pipeline, error) {
	pipeline := hdcp.NewPipeline(buildingParams, s.logger)
	if _, err := pipeline.Run(); err != nil {
		return nil, err
	}
	return pipeline, nil
}

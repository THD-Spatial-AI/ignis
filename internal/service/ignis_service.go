package service

import (
	"log"
	"os"

	"github.com/thd-spatial-ai/ignis/internal/models"
	"github.com/thd-spatial-ai/ignis/internal/pipeline"
)

// Re-export types for easier use in handlers
type TabulaBuildingParameters = models.TabulaBuildingParameters

// IgnisService provides business logic for heating demand calculations
type IgnisService struct {
	logger *pipeline.Logger
}

// NewIgnisService creates a new IgnisService instance
func NewIgnisService() *IgnisService {
	return &IgnisService{
		logger: pipeline.NewLogger(log.New(os.Stdout, "", 0)),
	}
}

// NewIgnisServiceWithLogger creates a new IgnisService with custom logger
func NewIgnisServiceWithLogger(logger *pipeline.Logger) *IgnisService {
	return &IgnisService{
		logger: logger,
	}
}

// CalculateHeatingDemand executes the heating demand calculation pipeline.
// Returns the calculated q_h_nd (annual heating energy demand in kWh/(m²·a)).
func (s *IgnisService) CalculateHeatingDemand(buildingParams *models.TabulaBuildingParameters) (float64, error) {
	p := pipeline.NewPipeline(buildingParams, s.logger)
	return p.Run()
}

// CalculateHeatingDemandWithDetails executes the heating demand calculation pipeline
// and returns the fully populated Pipeline struct for inspection of intermediate levels.
func (s *IgnisService) CalculateHeatingDemandWithDetails(buildingParams *models.TabulaBuildingParameters) (*pipeline.Pipeline, error) {
	p := pipeline.NewPipeline(buildingParams, s.logger)
	if _, err := p.Run(); err != nil {
		return nil, err
	}
	return p, nil
}

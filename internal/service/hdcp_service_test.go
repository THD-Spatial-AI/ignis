package service

import (
	"log"
	"os"
	"testing"

	"github.com/thd-spatial-ai/ignis/internal/hdcp"
	"github.com/thd-spatial-ai/ignis/internal/models"
)

// newTestParams returns a fully initialised TabulaBuildingParameters with all
// nested pointers non-nil, mirroring internal/calc's own test fixture.
func newTestParams() *models.TabulaBuildingParameters {
	return &models.TabulaBuildingParameters{
		BasicParameters: &models.BasicParameters{
			BuildingAppearance: &models.BuildingThematic{},
			Envelope:           &models.Envelope{},
		},
		AdvancedParameters: &models.AdvancedParameters{
			AirInfiltration:       &models.AirInfiltration{},
			ClimateConditions:     &models.ClimateConditions{},
			Uvalues:               &models.Uvalues{},
			Insulation:            &models.InsulationThicknesses{},
			SolarGains:            &models.SolarGains{},
			ThermalBridges:        &models.ThermalBridgeParameters{},
			HeatLosses:            &models.TransmissionHeatLoss{},
			ThermalResistances:    &models.ThermalResistances{},
			InsulationMeasures:    &models.InsulationPredefinedMeasures{},
			ActualInsulation:      &models.ActualInsulationThicknesses{},
			HeatTransfer:          &models.HeatTransferCoefficients{},
			PredefinedCodes:       &models.PredefinedCodes{F_Corr_CeilingHeight: 1.0},
			MeasureTypes:          &models.MeasureTypeCodes{},
			SolarTransmittance:    &models.SolarEnergyTransmittance{},
			MeasureFractions:      &models.MeasureAreaFractions{},
			AdditionalResistances: &models.AdditionalThermalResistance{},
		},
	}
}

func TestNewIgnisService(t *testing.T) {
	s := NewIgnisService()
	if s.logger == nil {
		t.Error("expected non-nil logger")
	}
}

func TestNewIgnisServiceWithLogger(t *testing.T) {
	logger := hdcp.NewLogger(log.New(os.Stdout, "", 0))
	s := NewIgnisServiceWithLogger(logger)
	if s.logger != logger {
		t.Error("expected the provided logger to be used")
	}
}

func TestCalculateHeatingDemand_success(t *testing.T) {
	s := NewIgnisService()
	// An all-zero building is a division-by-zero case in the pipeline's own
	// arithmetic (see internal/hdcp's own zero-input test) and legitimately
	// yields NaN, not an error - the point here is that it completes without panicking.
	_, err := s.CalculateHeatingDemand(newTestParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCalculateHeatingDemand_pipelineError(t *testing.T) {
	s := NewIgnisService()
	// nil params trigger a nil-pointer panic inside the pipeline, which
	// handleError converts into a returned error rather than a process crash.
	_, err := s.CalculateHeatingDemand(nil)
	if err == nil {
		t.Fatal("expected error from pipeline when params are nil")
	}
}

func TestCalculateHeatingDemandWithDetails_success(t *testing.T) {
	s := NewIgnisService()
	pipeline, err := s.CalculateHeatingDemandWithDetails(newTestParams())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pipeline == nil || pipeline.Lvl17 == nil {
		t.Error("expected fully populated pipeline with Lvl17 set")
	}
}

func TestCalculateHeatingDemandWithDetails_pipelineError(t *testing.T) {
	s := NewIgnisService()
	pipeline, err := s.CalculateHeatingDemandWithDetails(nil)
	if err == nil {
		t.Fatal("expected error from pipeline when params are nil")
	}
	if pipeline != nil {
		t.Error("expected nil pipeline on error")
	}
}

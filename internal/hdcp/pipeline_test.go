package hdcp

import (
	"io"
	"log"
	"testing"

	"github.com/thd-spatial-ai/ignis/internal/models"
)

// newTestParams returns a fully initialised TabulaBuildingParameters with zero field values.
// The pipeline must not panic on zero input — if it does, that is a pipeline bug.
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

func newDiscardLogger() *Logger {
	return NewLogger(log.New(io.Discard, "", 0))
}

// TestPipelineRun_zeroInput verifies that the pipeline completes without error
// when all building parameters are zero-valued. This is a smoke test for nil-safety.
func TestPipelineRun_zeroInput_noError(t *testing.T) {
	p := NewPipeline(newTestParams(), newDiscardLogger())
	result, err := p.Run()
	if err != nil {
		t.Fatalf("unexpected pipeline error: %v", err)
	}
	if result < 0 {
		t.Errorf("expected non-negative result, got %.4f", result)
	}
}

// TestPipelineRun_nilParams verifies that a nil TabulaBuildingParameters causes a
// recoverable error rather than a process crash.
func TestPipelineRun_nilParams_returnsError(t *testing.T) {
	p := NewPipeline(nil, newDiscardLogger())
	_, err := p.Run()
	if err == nil {
		t.Error("expected error for nil params, got nil")
	}
}

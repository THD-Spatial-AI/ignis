package service

import (
	"context"
	"log"
	"os"
	"strings"
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
	if s.repository != nil {
		t.Error("expected nil repository when constructed without a DB")
	}
}

func TestNewIgnisServiceWithDB(t *testing.T) {
	// NewIgnisServiceWithDB only stores the pool; it never dials, so nil is safe here.
	s := NewIgnisServiceWithDB(nil, "tabula")
	if s.repository == nil {
		t.Error("expected non-nil repository when constructed with a DB")
	}
}

func TestNewIgnisServiceWithLogger(t *testing.T) {
	logger := hdcp.NewLogger(log.New(os.Stdout, "", 0))
	s := NewIgnisServiceWithLogger(logger)
	if s.logger != logger {
		t.Error("expected the provided logger to be used")
	}
}

func TestGetBuildingByCode_noRepository_returnsError(t *testing.T) {
	s := NewIgnisService()
	_, err := s.GetBuildingByCode(context.Background(), "germany", "DE.N.SFH.01.Gen")
	if err == nil {
		t.Fatal("expected error when repository is not initialized")
	}
	if !strings.Contains(err.Error(), "repository not initialized") {
		t.Errorf("error = %q, want mention of uninitialized repository", err.Error())
	}
}

func TestNormalizeCountryName(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Germany", "germany"},
		{" United Kingdom ", "united_kingdom"},
		{"Czech-Republic", "czech_republic"},
	}
	for _, tc := range cases {
		if got := normalizeCountryName(tc.in); got != tc.want {
			t.Errorf("normalizeCountryName(%q) = %q, want %q", tc.in, got, tc.want)
		}
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

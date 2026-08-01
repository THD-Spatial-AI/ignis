package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/thd-spatial-ai/ignis/internal/db/repository"
	"github.com/thd-spatial-ai/ignis/internal/models"

	"github.com/gin-gonic/gin"
)

// mockRepo implements repository.TabulaReader for use in handler tests.
type mockRepo struct {
	listVariants  func(ctx context.Context, tableName string) ([]string, error)
	matchVariants func(ctx context.Context, tableName, prefix string) ([]string, error)
	getVariant    func(ctx context.Context, tableName, code string) (*models.TabulaBuildingParameters, string, float64, error)
}

func (m *mockRepo) ListVariants(ctx context.Context, tableName string) ([]string, error) {
	return m.listVariants(ctx, tableName)
}

func (m *mockRepo) MatchVariants(ctx context.Context, tableName, prefix string) ([]string, error) {
	if m.matchVariants != nil {
		return m.matchVariants(ctx, tableName, prefix)
	}
	return nil, nil
}

func (m *mockRepo) GetVariant(ctx context.Context, tableName, code string) (*models.TabulaBuildingParameters, string, float64, error) {
	return m.getVariant(ctx, tableName, code)
}

func newTestHandler(repo repository.TabulaReader) *Handler {
	return &Handler{repo: repo}
}

// serve builds a minimal Gin router with a single route and executes the request.
// Using ServeHTTP ensures the full Gin middleware chain runs and responses are flushed.
func serve(method, path, routePattern string, handler func(*gin.Context), body []byte) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(gin.Recovery())
	switch method {
	case http.MethodGet:
		r.GET(routePattern, handler)
	case http.MethodPost:
		r.POST(routePattern, handler)
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

// minimalBuilding returns a zero-valued building. The pipeline will produce NaN
// for zero input (divisions by zero), so the handler must return 500.
func minimalBuilding() *models.TabulaBuildingParameters {
	return &models.TabulaBuildingParameters{
		BasicParameters: &models.BasicParameters{
			BuildingAppearance: &models.BuildingThematic{},
			Envelope:           &models.Envelope{A_C_Ref_Input: 100},
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

// realisticBuilding returns a building with enough non-zero physical values for the
// pipeline to produce a finite heat-demand result. Values are representative but
// not taken from a specific TABULA record.
func realisticBuilding() *models.TabulaBuildingParameters {
	return &models.TabulaBuildingParameters{
		BasicParameters: &models.BasicParameters{
			BuildingAppearance: &models.BuildingThematic{N_Storey: 2},
			Envelope: &models.Envelope{
				A_C_Ref_Input:  150,
				A_C_IntDim:     150,
				V_C:            375,
				A_Roof_1:       75,
				A_Wall_1:       120,
				A_Floor_1:      75,
				A_Window_1:     20,
				A_Window_South: 12,
				A_Door_1:       2,
			},
		},
		AdvancedParameters: &models.AdvancedParameters{
			AirInfiltration:   &models.AirInfiltration{N_air_infiltration: 0.1, N_air_use: 0.4},
			ClimateConditions: &models.ClimateConditions{HeatingDays: 185, Theta_e: 0, Theta_i: 20},
			Uvalues: &models.Uvalues{
				U_Roof_1: 0.5, U_Wall_1: 0.8, U_Floor_1: 0.6, U_Window_1: 2.0, U_Door_1: 2.0,
			},
			Insulation:         &models.InsulationThicknesses{},
			SolarGains:         &models.SolarGains{I_Sol_South: 500, I_Sol_Horizontal: 800},
			ThermalBridges:     &models.ThermalBridgeParameters{},
			HeatLosses:         &models.TransmissionHeatLoss{},
			ThermalResistances: &models.ThermalResistances{},
			InsulationMeasures: &models.InsulationPredefinedMeasures{},
			ActualInsulation:   &models.ActualInsulationThicknesses{},
			HeatTransfer: &models.HeatTransferCoefficients{
				Phi_int: 4.0, F_sh_hor: 0.9, F_sh_vert: 0.9, F_f: 0.7, F_w: 0.9, C_m: 165000,
			},
			PredefinedCodes: &models.PredefinedCodes{F_Corr_CeilingHeight: 1.0},
			MeasureTypes:    &models.MeasureTypeCodes{},
			SolarTransmittance: &models.SolarEnergyTransmittance{
				G_gl_n_Window_1: 0.6,
			},
			MeasureFractions:      &models.MeasureAreaFractions{},
			AdditionalResistances: &models.AdditionalThermalResistance{},
		},
	}
}

// --- CalculateHeatDemand tests ---

func TestCalculateHeatDemand_variantNotFound_returns404(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return nil, "", 0, repository.ErrVariantNotFound
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestCalculateHeatDemand_malformedCode_returns400(t *testing.T) {
	h := newTestHandler(&mockRepo{})
	// Gin route param cannot be empty, so only test non-empty malformed codes here.
	cases := []struct{ path, code string }{
		{"/calculate/DE", "DE"},
		{"/calculate/1E.SFH.01", "1E.SFH.01"},
	}
	for _, tc := range cases {
		w := serve(http.MethodPost, tc.path, "/calculate/:code", h.CalculateHeatDemand, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("code %q: expected 400, got %d", tc.code, w.Code)
		}
	}
}

func TestCalculateHeatDemand_negativeARef_returns400(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return minimalBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)
	body, _ := json.Marshal(map[string]float64{"A_ref": -50})
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for negative A_ref, got %d", w.Code)
	}
}

func TestCalculateHeatDemand_zeroARef_returns400(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return minimalBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)
	body, _ := json.Marshal(map[string]float64{"A_ref": 0})
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for zero A_ref, got %d", w.Code)
	}
}

func TestCalculateHeatDemand_negativeHeatingDays_returns400(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)
	body, _ := json.Marshal(map[string]int{"HeatingDays": -1})
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for negative HeatingDays, got %d", w.Code)
	}
}

func TestCalculateHeatDemand_negativeISolSouth_returns400(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)
	body, _ := json.Marshal(map[string]float64{"I_Sol_South": -1})
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for negative I_Sol_South, got %d", w.Code)
	}
}

func TestCalculateHeatDemand_negativeThermalBridging_returns400(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)
	body, _ := json.Marshal(map[string]float64{"delta_U_ThermalBridging_Original": -0.1})
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for negative delta_U_ThermalBridging_Original, got %d", w.Code)
	}
}

func TestCalculateHeatDemand_zeroHRoom_returns400(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)
	body, _ := json.Marshal(map[string]float64{"h_room": 0})
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for zero h_room, got %d", w.Code)
	}
}

func TestCalculateHeatDemand_zeroNStorey_returns400(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)
	body, _ := json.Marshal(map[string]int{"n_Storey": 0})
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for zero n_Storey, got %d", w.Code)
	}
}

// TestCalculateHeatDemand_hRoomOverride_changesResult proves h_room reaches the
// pipeline: realisticBuilding() leaves H_room at its zero value (no ventilation
// heat loss contribution — see calc_level_01.go's H_Ventilation formula), so a
// positive override must produce a different q_h_nd, not an identical one.
func TestCalculateHeatDemand_hRoomOverride_changesResult(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)

	baseline := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, []byte("{}"))
	if baseline.Code != http.StatusOK {
		t.Fatalf("baseline: expected 200, got %d — body: %s", baseline.Code, baseline.Body.String())
	}
	overridden := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand,
		[]byte(`{"h_room": 2.8}`))
	if overridden.Code != http.StatusOK {
		t.Fatalf("overridden: expected 200, got %d — body: %s", overridden.Code, overridden.Body.String())
	}

	var baseResp, overrideResp map[string]any
	if err := json.Unmarshal(baseline.Body.Bytes(), &baseResp); err != nil {
		t.Fatalf("baseline response is not valid JSON: %v", err)
	}
	if err := json.Unmarshal(overridden.Body.Bytes(), &overrideResp); err != nil {
		t.Fatalf("overridden response is not valid JSON: %v", err)
	}
	if baseResp["q_h_nd"] == overrideResp["q_h_nd"] {
		t.Errorf("expected h_room override to change q_h_nd, got identical value %v for both", baseResp["q_h_nd"])
	}
}

func TestCalculateHeatDemand_negativeNAirInfiltration_returns400(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)
	body, _ := json.Marshal(map[string]float64{"n_air_infiltration": -0.1})
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for negative n_air_infiltration, got %d", w.Code)
	}
}

func TestCalculateHeatDemand_negativeNAirUse_returns400(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)
	body, _ := json.Marshal(map[string]float64{"n_air_use": -0.1})
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for negative n_air_use, got %d", w.Code)
	}
}

func TestCalculateHeatDemand_zeroCM_returns400(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)
	body, _ := json.Marshal(map[string]float64{"c_m": 0})
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for zero c_m, got %d", w.Code)
	}
}

// TestCalculateHeatDemand_cMOverride_changesResult proves c_m reaches the pipeline
// (calc_level_12.go's time constant tau = C_m / (H_tr + H_ve), which feeds the gain
// utilization factor).
func TestCalculateHeatDemand_cMOverride_changesResult(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)

	baseline := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, []byte("{}"))
	if baseline.Code != http.StatusOK {
		t.Fatalf("baseline: expected 200, got %d — body: %s", baseline.Code, baseline.Body.String())
	}
	overridden := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand,
		[]byte(`{"c_m": 45}`))
	if overridden.Code != http.StatusOK {
		t.Fatalf("overridden: expected 200, got %d — body: %s", overridden.Code, overridden.Body.String())
	}

	var baseResp, overrideResp map[string]any
	if err := json.Unmarshal(baseline.Body.Bytes(), &baseResp); err != nil {
		t.Fatalf("baseline response is not valid JSON: %v", err)
	}
	if err := json.Unmarshal(overridden.Body.Bytes(), &overrideResp); err != nil {
		t.Fatalf("overridden response is not valid JSON: %v", err)
	}
	if baseResp["q_h_nd"] == overrideResp["q_h_nd"] {
		t.Errorf("expected c_m override to change q_h_nd, got identical value %v for both", baseResp["q_h_nd"])
	}
}

// TestCalculateHeatDemand_heatingDaysOverride_changesResult proves the override actually
// reaches the pipeline: zero heating days must produce a lower q_h_nd than the building's
// default 185 heating days (realisticBuilding), not an identical result.
func TestCalculateHeatDemand_heatingDaysOverride_changesResult(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)

	baseline := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, []byte("{}"))
	if baseline.Code != http.StatusOK {
		t.Fatalf("baseline: expected 200, got %d — body: %s", baseline.Code, baseline.Body.String())
	}
	overridden := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand,
		[]byte(`{"HeatingDays": 0}`))
	if overridden.Code != http.StatusOK {
		t.Fatalf("overridden: expected 200, got %d — body: %s", overridden.Code, overridden.Body.String())
	}

	var baseResp, overrideResp map[string]any
	if err := json.Unmarshal(baseline.Body.Bytes(), &baseResp); err != nil {
		t.Fatalf("baseline response is not valid JSON: %v", err)
	}
	if err := json.Unmarshal(overridden.Body.Bytes(), &overrideResp); err != nil {
		t.Fatalf("overridden response is not valid JSON: %v", err)
	}
	if baseResp["q_h_nd"] == overrideResp["q_h_nd"] {
		t.Errorf("expected HeatingDays override to change q_h_nd, got identical value %v for both", baseResp["q_h_nd"])
	}
}

// TestCalculateHeatDemand_zeroBuildingProducesNaN verifies that the handler returns 500
// when the pipeline produces a NaN result (all-zero building inputs cause division-by-zero).
func TestCalculateHeatDemand_zeroBuildingProducesNaN_returns500(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return minimalBuilding(), "DE.N.SFH.01.Gen", 0, nil
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, []byte("{}"))
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 for NaN result, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCalculateHeatDemand_repoError_returns500(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return nil, "", 0, errors.New("connection refused")
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestCalculateHeatDemand_malformedJSONBody_returns400(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return minimalBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, []byte("not json"))
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for malformed JSON body, got %d", w.Code)
	}
}

func TestCalculateHeatDemand_pipelinePanics_returns500(t *testing.T) {
	mock := &mockRepo{
		// A nil building triggers a nil-pointer panic inside the pipeline,
		// which the pipeline itself converts into an error - exercises the
		// handler's "Pipeline execution failed" branch.
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return nil, "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, []byte("{}"))
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d — body: %s", w.Code, w.Body.String())
	}
}

func TestCalculateHeatDemand_success_returns200withExpectedShape(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 100.0, nil
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodPost, "/calculate/DE.N.SFH.01.Gen", "/calculate/:code", h.CalculateHeatDemand, []byte("{}"))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	for _, key := range []string{"variant_code", "q_h_nd", "unit"} {
		if _, ok := resp[key]; !ok {
			t.Errorf("response missing field %q", key)
		}
	}
}

// --- GetVariants tests ---

func TestGetVariants_unknownCountry_returns400(t *testing.T) {
	h := newTestHandler(&mockRepo{})
	w := serve(http.MethodGet, "/variants/ZZ", "/variants/:country_iso2", h.GetVariants, nil)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for unknown country, got %d", w.Code)
	}
}

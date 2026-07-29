package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/thd-spatial-ai/ignis/internal/db/repository"
	"github.com/thd-spatial-ai/ignis/internal/models"
)

// --- GetVariants ---

func TestGetVariants_success_returns200withVariants(t *testing.T) {
	mock := &mockRepo{
		listVariants: func(_ context.Context, _ string) ([]string, error) {
			return []string{"DE.N.SFH.01.Gen", "DE.N.SFH.02.Gen"}, nil
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/variants/DE", "/variants/:country_iso2", h.GetVariants, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Country string   `json:"country"`
		Data    []string `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp.Country != "germany" {
		t.Errorf("country = %q, want %q", resp.Country, "germany")
	}
	if len(resp.Data) != 2 {
		t.Errorf("expected 2 variants, got %d", len(resp.Data))
	}
}

func TestGetVariants_repoError_returns500(t *testing.T) {
	mock := &mockRepo{
		listVariants: func(_ context.Context, _ string) ([]string, error) {
			return nil, errors.New("connection refused")
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/variants/DE", "/variants/:country_iso2", h.GetVariants, nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// --- GetVariantData ---

func TestGetVariantData_malformedCode_returns400(t *testing.T) {
	h := newTestHandler(&mockRepo{})
	w := serve(http.MethodGet, "/data/DE", "/data/:code", h.GetVariantData, nil)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetVariantData_unknownCountry_returns400(t *testing.T) {
	h := newTestHandler(&mockRepo{})
	w := serve(http.MethodGet, "/data/ZZ.N.SFH.01.Gen", "/data/:code", h.GetVariantData, nil)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetVariantData_variantNotFound_returns404(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return nil, "", 0, repository.ErrVariantNotFound
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/data/DE.N.SFH.01.Gen", "/data/:code", h.GetVariantData, nil)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetVariantData_repoError_returns500(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return nil, "", 0, errors.New("connection refused")
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/data/DE.N.SFH.01.Gen", "/data/:code", h.GetVariantData, nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestGetVariantData_success_returns200(t *testing.T) {
	mock := &mockRepo{
		getVariant: func(_ context.Context, _, _ string) (*models.TabulaBuildingParameters, string, float64, error) {
			return realisticBuilding(), "DE.N.SFH.01.Gen", 123.4, nil
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/data/DE.N.SFH.01.Gen", "/data/:code", h.GetVariantData, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	for _, key := range []string{"country", "variant_code", "tabula_data", "expected_q_h_nd"} {
		if _, ok := resp[key]; !ok {
			t.Errorf("response missing field %q", key)
		}
	}
}

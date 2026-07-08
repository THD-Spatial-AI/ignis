package handler

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/thd-spatial-ai/ignis/internal/models"
)

func TestGetFieldMetadata_returnsAllFieldsWithNonEmptyDescriptions(t *testing.T) {
	h := newTestHandler(nil) // GetFieldMetadata does not touch the repo
	w := serve(http.MethodGet, "/api/v1/fields", "/api/v1/fields", h.GetFieldMetadata, nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var body struct {
		Data []models.FieldMetadata `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	if len(body.Data) != len(models.AllFieldMetadata) {
		t.Fatalf("expected %d fields, got %d", len(models.AllFieldMetadata), len(body.Data))
	}

	seen := make(map[string]bool, len(body.Data))
	for _, f := range body.Data {
		seen[f.Key] = true
		if f.Label == "" || f.SimpleDescription == "" || f.ExpertDescription == "" || f.Path == "" {
			t.Errorf("field %q has an empty label/description/path: %+v", f.Key, f)
		}
	}

	// Spot-check the fields that motivated this endpoint (the ones that were
	// silently returning 0 in Building Configurator due to the nesting bug).
	for _, key := range []string{"HeatingDays", "Theta_e", "theta_i"} {
		if !seen[key] {
			t.Errorf("expected field metadata for %q, not found", key)
		}
	}
}

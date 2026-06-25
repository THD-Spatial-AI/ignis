package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

// --- isoFromVariantCode ---

func TestIsoFromVariantCode_valid(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"DE.N.SFH.01.Gen", "DE"},
		{"AT.N.SFH.01.Gen", "AT"},
		{"de.N.SFH.01.Gen", "DE"}, // lowercase prefix is uppercased
	}
	for _, tc := range cases {
		got, err := isoFromVariantCode(tc.input)
		if err != nil {
			t.Errorf("isoFromVariantCode(%q) unexpected error: %v", tc.input, err)
			continue
		}
		if got != tc.want {
			t.Errorf("isoFromVariantCode(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestIsoFromVariantCode_invalid(t *testing.T) {
	cases := []string{
		"",          // empty
		"DE",        // no dot
		"D.SFH",     // prefix shorter than 2 chars + dot
		"1E.SFH.01", // non-alpha prefix
		"D2.SFH.01", // non-alpha prefix
	}
	for _, code := range cases {
		_, err := isoFromVariantCode(code)
		if err == nil {
			t.Errorf("isoFromVariantCode(%q): expected error, got nil", code)
		}
	}
}

// --- tableNameFromISO ---

func TestTableNameFromISO_knownCountries(t *testing.T) {
	cases := []struct {
		iso  string
		want string
	}{
		{"DE", "germany"},
		{"AT", "austria"},
		{"FR", "france"},
	}
	for _, tc := range cases {
		got, err := tableNameFromISO(tc.iso)
		if err != nil {
			t.Errorf("tableNameFromISO(%q) unexpected error: %v", tc.iso, err)
			continue
		}
		if got != tc.want {
			t.Errorf("tableNameFromISO(%q) = %q, want %q", tc.iso, got, tc.want)
		}
	}
}

func TestTableNameFromISO_invalidLength(t *testing.T) {
	cases := []string{"D", "", "DEU"}
	for _, iso := range cases {
		_, err := tableNameFromISO(iso)
		if err == nil {
			t.Errorf("tableNameFromISO(%q): expected error for invalid length", iso)
		}
	}
}

func TestTableNameFromISO_unknownCountry(t *testing.T) {
	// "ZZ" is not a real ISO2 code; the country helper returns a fallback string.
	// tableNameFromISO must reject it (the fallback is not a valid TABULA table).
	_, err := tableNameFromISO("ZZ")
	if err == nil {
		t.Error("tableNameFromISO(\"ZZ\"): expected error for unknown country, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "ZZ") {
		t.Errorf("tableNameFromISO(\"ZZ\"): error should mention the code, got: %v", err)
	}
}

// --- refurbishmentLabel ---

func TestRefurbishmentLabel_knownPositions(t *testing.T) {
	cases := []struct {
		index int
		want  string
	}{
		{0, "Existing state"},
		{1, "Medium refurbishment"},
		{2, "Advanced refurbishment"},
	}
	for _, tc := range cases {
		got := refurbishmentLabel(tc.index)
		if got != tc.want {
			t.Errorf("refurbishmentLabel(%d) = %q, want %q", tc.index, got, tc.want)
		}
	}
}

func TestRefurbishmentLabel_beyondKnownPositions(t *testing.T) {
	got := refurbishmentLabel(5)
	if got == "" {
		t.Error("refurbishmentLabel(5): expected non-empty fallback label")
	}
}

// --- MatchVariants handler ---

func TestMatchVariants_missingParams_returns400(t *testing.T) {
	h := newTestHandler(&mockRepo{})
	cases := []struct{ path, route string }{
		{"/variants/DE/match", "/variants/:country_iso2/match"},             // no type or period
		{"/variants/DE/match?type=SFH", "/variants/:country_iso2/match"},   // missing period
		{"/variants/DE/match?period=01", "/variants/:country_iso2/match"},  // missing type
	}
	for _, tc := range cases {
		w := serve(http.MethodGet, tc.path, tc.route, h.MatchVariants, nil)
		if w.Code != http.StatusBadRequest {
			t.Errorf("path %q: expected 400, got %d", tc.path, w.Code)
		}
	}
}

func TestMatchVariants_unknownCountry_returns400(t *testing.T) {
	h := newTestHandler(&mockRepo{})
	w := serve(http.MethodGet, "/variants/ZZ/match?type=SFH&period=01", "/variants/:country_iso2/match", h.MatchVariants, nil)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for unknown country, got %d", w.Code)
	}
}

func TestMatchVariants_returnsLabelledVariants(t *testing.T) {
	mock := &mockRepo{
		matchVariants: func(_ context.Context, _, _ string) ([]string, error) {
			return []string{"DE.N.SFH.01.Gen", "DE.N.SFH.01.ReEx", "DE.N.SFH.01.Add"}, nil
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/variants/DE/match?type=SFH&period=01", "/variants/:country_iso2/match", h.MatchVariants, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Data []struct {
			Code  string `json:"code"`
			Label string `json:"label"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(resp.Data) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(resp.Data))
	}
	wantLabels := []string{"Existing state", "Medium refurbishment", "Advanced refurbishment"}
	for i, entry := range resp.Data {
		if entry.Label != wantLabels[i] {
			t.Errorf("entry[%d].label = %q, want %q", i, entry.Label, wantLabels[i])
		}
	}
}

func TestMatchVariants_emptyResult_returns200withEmptyList(t *testing.T) {
	mock := &mockRepo{
		matchVariants: func(_ context.Context, _, _ string) ([]string, error) {
			return nil, nil
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/variants/DE/match?type=SFH&period=99", "/variants/:country_iso2/match", h.MatchVariants, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Data []any `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(resp.Data) != 0 {
		t.Errorf("expected empty data list, got %d entries", len(resp.Data))
	}
}

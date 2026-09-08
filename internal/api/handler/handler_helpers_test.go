package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/thd-spatial-ai/ignis/internal/db/repository"
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
		{"/variants/DE/match", "/variants/:country_iso2/match"},           // no type or period
		{"/variants/DE/match?type=SFH", "/variants/:country_iso2/match"},  // missing period
		{"/variants/DE/match?period=01", "/variants/:country_iso2/match"}, // missing type
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

func TestMatchVariants_repoError_returns500(t *testing.T) {
	mock := &mockRepo{
		matchVariants: func(_ context.Context, _, _ string) ([]string, error) {
			return nil, errors.New("connection refused")
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/variants/DE/match?type=SFH&period=01", "/variants/:country_iso2/match", h.MatchVariants, nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
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

func TestNew_constructsHandlerWithRepo(t *testing.T) {
	h := New(nil, "tabula")
	if h == nil || h.repo == nil {
		t.Fatal("expected New to return a Handler with a non-nil repo")
	}
}

// --- MatchVariants: year resolution ---

func TestMatchVariants_byYear_resolvesPeriodAndReturnsVariants(t *testing.T) {
	var gotPrefix string
	var gotYear int
	mock := &mockRepo{
		resolvePeriodByYear: func(_ context.Context, _, typePrefix string, year int) (string, error) {
			gotPrefix, gotYear = typePrefix, year
			return "03", nil
		},
		matchVariants: func(_ context.Context, _, prefix string) ([]string, error) {
			if prefix != "DE.N.SFH.03" {
				t.Errorf("MatchVariants prefix = %q, want DE.N.SFH.03", prefix)
			}
			return []string{"DE.N.SFH.03.Gen", "DE.N.SFH.03.ReEx"}, nil
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/variants/DE/match?type=SFH&year=1975", "/variants/:country_iso2/match", h.MatchVariants, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
	if gotPrefix != "DE.N.SFH" || gotYear != 1975 {
		t.Errorf("ResolvePeriodByYear called with (%q, %d), want (\"DE.N.SFH\", 1975)", gotPrefix, gotYear)
	}
	var resp struct {
		Data []struct {
			Code string `json:"code"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(resp.Data))
	}
}

func TestMatchVariants_bothYearAndPeriod_returns400(t *testing.T) {
	h := newTestHandler(&mockRepo{})
	w := serve(http.MethodGet, "/variants/DE/match?type=SFH&period=03&year=1975", "/variants/:country_iso2/match", h.MatchVariants, nil)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 when both period and year given, got %d", w.Code)
	}
}

func TestMatchVariants_nonIntegerYear_returns400(t *testing.T) {
	h := newTestHandler(&mockRepo{})
	w := serve(http.MethodGet, "/variants/DE/match?type=SFH&year=nineteen", "/variants/:country_iso2/match", h.MatchVariants, nil)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for non-integer year, got %d", w.Code)
	}
}

func TestMatchVariants_negativeYear_returns400(t *testing.T) {
	mock := &mockRepo{
		resolvePeriodByYear: func(_ context.Context, _, _ string, _ int) (string, error) {
			t.Error("ResolvePeriodByYear must not be called for a negative year")
			return "", nil
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/variants/DE/match?type=SFH&year=-5", "/variants/:country_iso2/match", h.MatchVariants, nil)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for negative year, got %d", w.Code)
	}
}

func TestMatchVariants_futureYear_returns400(t *testing.T) {
	mock := &mockRepo{
		resolvePeriodByYear: func(_ context.Context, _, _ string, _ int) (string, error) {
			t.Error("ResolvePeriodByYear must not be called for a future year")
			return "", nil
		},
	}
	h := newTestHandler(mock)
	futureYear := time.Now().Year() + 1
	w := serve(http.MethodGet, fmt.Sprintf("/variants/DE/match?type=SFH&year=%d", futureYear), "/variants/:country_iso2/match", h.MatchVariants, nil)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for a year in the future, got %d", w.Code)
	}
}

func TestMatchVariants_yearWithNoMatchingArchetype_returns200EmptyList(t *testing.T) {
	mock := &mockRepo{
		resolvePeriodByYear: func(_ context.Context, _, _ string, _ int) (string, error) {
			return "", nil // no band for this type contains the year
		},
		matchVariants: func(_ context.Context, _, _ string) ([]string, error) {
			t.Error("MatchVariants must not be called when no period resolves")
			return nil, nil
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/variants/DE/match?type=SFH&year=1750", "/variants/:country_iso2/match", h.MatchVariants, nil)
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

func TestMatchVariants_resolvePeriodError_returns500(t *testing.T) {
	mock := &mockRepo{
		resolvePeriodByYear: func(_ context.Context, _, _ string, _ int) (string, error) {
			return "", errors.New("connection refused")
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/variants/DE/match?type=SFH&year=1975", "/variants/:country_iso2/match", h.MatchVariants, nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// --- ListPeriods handler ---

func TestListPeriods_returnsBands(t *testing.T) {
	mock := &mockRepo{
		listPeriods: func(_ context.Context, _ string) ([]repository.ConstructionPeriod, error) {
			return []repository.ConstructionPeriod{
				{Period: "01", YearFrom: 0, YearTo: 1859},
				{Period: "02", YearFrom: 1860, YearTo: 1918},
				{Period: "12", YearFrom: 2016, YearTo: 9999},
			}, nil
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/periods/DE", "/periods/:country_iso2", h.ListPeriods, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Data []repository.ConstructionPeriod `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(resp.Data) != 3 {
		t.Fatalf("expected 3 bands, got %d", len(resp.Data))
	}
	if resp.Data[0].YearFrom != 0 || resp.Data[2].YearTo != 9999 {
		t.Errorf("open-ended sentinels not preserved: %+v", resp.Data)
	}
}

func TestListPeriods_unknownCountry_returns400(t *testing.T) {
	h := newTestHandler(&mockRepo{})
	w := serve(http.MethodGet, "/periods/ZZ", "/periods/:country_iso2", h.ListPeriods, nil)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for unknown country, got %d", w.Code)
	}
}

func TestListPeriods_repoError_returns500(t *testing.T) {
	mock := &mockRepo{
		listPeriods: func(_ context.Context, _ string) ([]repository.ConstructionPeriod, error) {
			return nil, errors.New("connection refused")
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/periods/DE", "/periods/:country_iso2", h.ListPeriods, nil)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

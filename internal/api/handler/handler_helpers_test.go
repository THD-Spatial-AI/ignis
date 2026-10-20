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

func TestRefurbishmentLabel_withoutDescription_usesTabulaName(t *testing.T) {
	cases := []struct {
		number int
		want   string
	}{
		{1, "Existing state"},
		{2, "Usual refurbishment"},
		{3, "Advanced refurbishment"},
	}
	for _, tc := range cases {
		got := refurbishmentLabel(repository.RefurbishmentVariant{Number: tc.number})
		if got != tc.want {
			t.Errorf("refurbishmentLabel(number %d) = %q, want %q", tc.number, got, tc.want)
		}
	}
}

func TestRefurbishmentLabel_withDescription_usesDescription(t *testing.T) {
	got := refurbishmentLabel(repository.RefurbishmentVariant{Number: 3, Description: "passive house"})
	if got != "passive house" {
		t.Errorf("refurbishmentLabel = %q, want %q", got, "passive house")
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
		matchVariants: func(_ context.Context, _, _ string) ([]repository.RefurbishmentVariant, error) {
			// Variant 2 absent, as in the newest BE and NL periods.
			return []repository.RefurbishmentVariant{
				{Code: "NL.N.SFH.06.Gen.ReEx.001.001", Number: 1},
				{Code: "NL.N.SFH.06.Gen.ReEx.001.003", Number: 3, Description: "nZEB"},
			}, nil
		},
	}
	h := newTestHandler(mock)
	w := serve(http.MethodGet, "/variants/NL/match?type=SFH&period=06", "/variants/:country_iso2/match", h.MatchVariants, nil)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d — body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Data []struct {
			Code        string `json:"code"`
			Number      int    `json:"number"`
			Label       string `json:"label"`
			Description string `json:"description"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(resp.Data))
	}
	if got := resp.Data[0]; got.Number != 1 || got.Label != "Existing state" || got.Description != "" {
		t.Errorf("entry[0] = %+v, want number 1, label %q, no description", got, "Existing state")
	}
	if got := resp.Data[1]; got.Number != 3 || got.Label != "nZEB" || got.Description != "nZEB" {
		t.Errorf("entry[1] = %+v, want number 3, label and description %q", got, "nZEB")
	}
}

func TestMatchVariants_subtype_scopesPrefixAndIsEchoed(t *testing.T) {
	cases := []struct{ query, wantSubtype string }{
		{"", "Gen"},
		{"&subtype=", "Gen"},
		{"&subtype=Semi", "Semi"},
	}
	for _, tc := range cases {
		var gotPrefix string
		mock := &mockRepo{
			matchVariants: func(_ context.Context, _, prefix string) ([]repository.RefurbishmentVariant, error) {
				gotPrefix = prefix
				return nil, nil
			},
		}
		h := newTestHandler(mock)
		w := serve(http.MethodGet, "/variants/BE/match?type=TH&period=01"+tc.query, "/variants/:country_iso2/match", h.MatchVariants, nil)
		if w.Code != http.StatusOK {
			t.Fatalf("query %q: expected 200, got %d, body: %s", tc.query, w.Code, w.Body.String())
		}
		var resp struct {
			Prefix  string `json:"prefix"`
			Subtype string `json:"subtype"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("invalid JSON: %v", err)
		}
		wantPrefix := "BE.N.TH.01." + tc.wantSubtype
		if gotPrefix != wantPrefix || resp.Prefix != wantPrefix || resp.Subtype != tc.wantSubtype {
			t.Errorf("query %q: repo prefix %q, response prefix %q subtype %q; want prefix %q subtype %q",
				tc.query, gotPrefix, resp.Prefix, resp.Subtype, wantPrefix, tc.wantSubtype)
		}
	}
}

func TestMatchVariants_repoError_returns500(t *testing.T) {
	mock := &mockRepo{
		matchVariants: func(_ context.Context, _, _ string) ([]repository.RefurbishmentVariant, error) {
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
		matchVariants: func(_ context.Context, _, _ string) ([]repository.RefurbishmentVariant, error) {
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
		matchVariants: func(_ context.Context, _, prefix string) ([]repository.RefurbishmentVariant, error) {
			if prefix != "DE.N.SFH.03.Gen" {
				t.Errorf("MatchVariants prefix = %q, want DE.N.SFH.03.Gen", prefix)
			}
			return []repository.RefurbishmentVariant{
				{Code: "DE.N.SFH.03.Gen.ReEx.001.001", Number: 1},
				{Code: "DE.N.SFH.03.Gen.ReEx.001.002", Number: 2},
			}, nil
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
		matchVariants: func(_ context.Context, _, _ string) ([]repository.RefurbishmentVariant, error) {
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

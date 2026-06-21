package handler

import (
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

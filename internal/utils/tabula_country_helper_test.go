package utils

import "testing"

func TestCodeToCountry_knownCodes(t *testing.T) {
	h := NewTabulaCountryHelper()
	cases := []struct {
		code string
		want string
	}{
		{"DE", "germany"},
		{"AT", "austria"},
		{"FR", "france"},
		{"GB", "united_kingdom"},
		{"PL", "poland"},
	}
	for _, tc := range cases {
		got := h.CodeToCountry(tc.code)
		if got != tc.want {
			t.Errorf("CodeToCountry(%q) = %q, want %q", tc.code, got, tc.want)
		}
	}
}

func TestCodeToCountry_unknownCode_returnsLowercaseFallback(t *testing.T) {
	h := NewTabulaCountryHelper()
	// NOTE: the docstring says "empty string if not found" but the implementation
	// returns strings.ToLower(code) as a fallback. Tests document actual behaviour.
	got := h.CodeToCountry("XX")
	if got != "xx" {
		t.Errorf("CodeToCountry(\"XX\") = %q, want \"xx\" (lowercase fallback)", got)
	}
}

func TestCodeToCountry_lowercaseInput_normalised(t *testing.T) {
	h := NewTabulaCountryHelper()
	// Input is uppercased internally, so "de" resolves the same as "DE".
	got := h.CodeToCountry("de")
	if got != "germany" {
		t.Errorf("CodeToCountry(\"de\") = %q, want \"germany\"", got)
	}
}

func TestCountryToCode_knownCountries(t *testing.T) {
	h := NewTabulaCountryHelper()
	cases := []struct {
		country string
		want    string
	}{
		{"germany", "DE"},
		{"austria", "AT"},
		{"france", "FR"},
	}
	for _, tc := range cases {
		got := h.CountryToCode(tc.country)
		if got != tc.want {
			t.Errorf("CountryToCode(%q) = %q, want %q", tc.country, got, tc.want)
		}
	}
}

func TestCountryToCode_unknownCountry_returnsEmpty(t *testing.T) {
	h := NewTabulaCountryHelper()
	got := h.CountryToCode("narnia")
	if got != "" {
		t.Errorf("CountryToCode(\"narnia\") = %q, want empty string", got)
	}
}

func TestGetAllCountryCodes(t *testing.T) {
	h := NewTabulaCountryHelper()
	codes := h.GetAllCountryCodes()
	if len(codes) != len(countryMapping) {
		t.Errorf("GetAllCountryCodes returned %d codes, want %d", len(codes), len(countryMapping))
	}
	found := false
	for _, c := range codes {
		if c == "DE" {
			found = true
		}
	}
	if !found {
		t.Error("GetAllCountryCodes: expected \"DE\" in result")
	}
}

func TestGetAllCountryNames(t *testing.T) {
	h := NewTabulaCountryHelper()
	names := h.GetAllCountryNames()
	if len(names) != len(countryMapping) {
		t.Errorf("GetAllCountryNames returned %d names, want %d", len(names), len(countryMapping))
	}
	found := false
	for _, n := range names {
		if n == "germany" {
			found = true
		}
	}
	if !found {
		t.Error("GetAllCountryNames: expected \"germany\" in result")
	}
}

func TestIsValidCountryCode(t *testing.T) {
	h := NewTabulaCountryHelper()
	if !h.IsValidCountryCode("de") {
		t.Error("IsValidCountryCode(\"de\") = false, want true (case-insensitive)")
	}
	if h.IsValidCountryCode("XX") {
		t.Error("IsValidCountryCode(\"XX\") = true, want false")
	}
}

func TestIsValidCountryName(t *testing.T) {
	h := NewTabulaCountryHelper()
	if !h.IsValidCountryName("United Kingdom") {
		t.Error("IsValidCountryName(\"United Kingdom\") = false, want true (space/case normalised)")
	}
	if h.IsValidCountryName("narnia") {
		t.Error("IsValidCountryName(\"narnia\") = true, want false")
	}
}

func TestRoundTrip_codeToCountryToCode(t *testing.T) {
	h := NewTabulaCountryHelper()
	// Spot-check that code → country → code round-trips correctly.
	codes := []string{"DE", "AT", "FR", "PL", "NL", "BE", "IT", "ES"}
	for _, code := range codes {
		country := h.CodeToCountry(code)
		if country == "" {
			t.Errorf("CodeToCountry(%q) returned empty", code)
			continue
		}
		back := h.CountryToCode(country)
		if back != code {
			t.Errorf("round-trip %q → %q → %q", code, country, back)
		}
	}
}

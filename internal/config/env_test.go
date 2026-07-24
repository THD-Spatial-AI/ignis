package config

import "testing"

func TestLoadEnv_doesNotPanicWithoutDotenvFile(t *testing.T) {
	// No .env file exists in internal/config/ - godotenv.Load() returns an
	// error that LoadEnv deliberately swallows so the app can run purely off
	// real environment variables (e.g. in Docker).
	LoadEnv()
}

func TestGetEnv_returnsValueWhenSet(t *testing.T) {
	t.Setenv("IGNIS_TEST_KEY", "value")
	if got := GetEnv("IGNIS_TEST_KEY", "fallback"); got != "value" {
		t.Errorf("GetEnv = %q, want %q", got, "value")
	}
}

func TestGetEnv_returnsFallbackWhenUnset(t *testing.T) {
	t.Setenv("IGNIS_TEST_KEY", "")
	if got := GetEnv("IGNIS_TEST_KEY", "fallback"); got != "fallback" {
		t.Errorf("GetEnv = %q, want %q", got, "fallback")
	}
}

func TestGetEnvAsInt_returnsParsedValue(t *testing.T) {
	t.Setenv("IGNIS_TEST_INT", "42")
	if got := GetEnvAsInt("IGNIS_TEST_INT", 7); got != 42 {
		t.Errorf("GetEnvAsInt = %d, want %d", got, 42)
	}
}

func TestGetEnvAsInt_returnsFallbackWhenUnset(t *testing.T) {
	t.Setenv("IGNIS_TEST_INT", "")
	if got := GetEnvAsInt("IGNIS_TEST_INT", 7); got != 7 {
		t.Errorf("GetEnvAsInt = %d, want %d", got, 7)
	}
}

func TestGetEnvAsInt_returnsFallbackWhenUnparseable(t *testing.T) {
	t.Setenv("IGNIS_TEST_INT", "not-a-number")
	if got := GetEnvAsInt("IGNIS_TEST_INT", 7); got != 7 {
		t.Errorf("GetEnvAsInt = %d, want %d", got, 7)
	}
}

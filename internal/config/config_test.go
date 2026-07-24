package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAppConfig_default(t *testing.T) {
	t.Setenv("APP_PORT", "")
	cfg := loadAppConfig()
	if cfg.Port != defaultPort {
		t.Errorf("Port = %q, want default %q", cfg.Port, defaultPort)
	}
}

func TestLoadAppConfig_fromEnv(t *testing.T) {
	t.Setenv("APP_PORT", "9090")
	cfg := loadAppConfig()
	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9090")
	}
}

func TestLoadDBConfig_missingPassword_panics(t *testing.T) {
	t.Setenv("DB_PASSWORD", "")
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when DB_PASSWORD is unset, got none")
		}
	}()
	loadDBConfig()
}

func TestLoadDBConfig_populatesFromEnv(t *testing.T) {
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_HOST", "db.internal")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_NAME", "ignis_test")
	t.Setenv("DB_USER", "tester")
	t.Setenv("DB_SSL_MODE", "disable")

	cfg := loadDBConfig()

	if cfg.Host != "db.internal" || cfg.Port != "5433" || cfg.Name != "ignis_test" ||
		cfg.User != "tester" || cfg.Password != "secret" || cfg.SSLMode != "disable" {
		t.Errorf("unexpected DBConfig: %+v", cfg)
	}
	if cfg.Schemas == nil || cfg.Schemas.Tabula != tabulaSchema {
		t.Errorf("Schemas.Tabula = %+v, want %q", cfg.Schemas, tabulaSchema)
	}
}

func TestLoadDBConfig_defaults(t *testing.T) {
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_SSL_MODE", "")

	cfg := loadDBConfig()

	if cfg.Host != "localhost" || cfg.Port != "5432" || cfg.Name != "ignis" ||
		cfg.User != "postgres" || cfg.SSLMode != "require" {
		t.Errorf("unexpected defaults: %+v", cfg)
	}
}

func TestLoadSchemas(t *testing.T) {
	s := loadSchemas()
	if s.Tabula != tabulaSchema {
		t.Errorf("Tabula = %q, want %q", s.Tabula, tabulaSchema)
	}
}

func TestLoadDataPaths_fallbackWhenNoXlsx(t *testing.T) {
	dir := t.TempDir()
	restore := chdir(t, dir)
	defer restore()

	paths := loadDataPaths()
	want := dataDir + "tabula-calculator.xlsx"
	if paths.ExcelFile != want {
		t.Errorf("ExcelFile = %q, want fallback %q", paths.ExcelFile, want)
	}
}

func TestLoadDataPaths_findsXlsxFile(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "data"), 0755); err != nil {
		t.Fatal(err)
	}
	xlsxPath := filepath.Join(dir, "data", "my-workbook.xlsx")
	if err := os.WriteFile(xlsxPath, []byte("not a real workbook"), 0644); err != nil {
		t.Fatal(err)
	}

	restore := chdir(t, dir)
	defer restore()

	paths := loadDataPaths()
	want := dataDir + "my-workbook.xlsx"
	if paths.ExcelFile != want {
		t.Errorf("ExcelFile = %q, want %q", paths.ExcelFile, want)
	}
}

func TestLoadConfig(t *testing.T) {
	t.Setenv("DB_PASSWORD", "secret")
	cfg := LoadConfig()
	if cfg.App == nil || cfg.DB == nil || cfg.Data == nil {
		t.Errorf("LoadConfig returned incomplete config: %+v", cfg)
	}
}

// chdir switches the working directory for the duration of a test and
// returns a function that restores the original directory.
func chdir(t *testing.T, dir string) func() {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	return func() {
		if err := os.Chdir(orig); err != nil {
			t.Fatal(err)
		}
	}
}

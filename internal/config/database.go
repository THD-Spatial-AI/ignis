package config

import (
	"path/filepath"
)

const (
	tabulaSchema = "tabula"
	dataDir      = "data/"
)

// loadDBConfig loads database configuration from environment
func loadDBConfig() *DBConfig {
	password := GetEnv("DB_PASSWORD", "")
	if password == "" {
		panic("DB_PASSWORD is not set — refusing to start with an empty database password")
	}
	return &DBConfig{
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     GetEnv("DB_PORT", "5432"),
		Name:     GetEnv("DB_NAME", "ignis"),
		User:     GetEnv("DB_USER", "postgres"),
		Password: password,
		SSLMode:  GetEnv("DB_SSL_MODE", "require"),
		Schemas:  loadSchemas(),
	}
}

// loadSchemas loads schema configuration
func loadSchemas() *Schemas {
	return &Schemas{
		Tabula: tabulaSchema,
	}
}

// loadDataPaths loads data file paths
func loadDataPaths() *DataPaths {
	// Find the first .xlsx file in the data directory
	files, err := filepath.Glob(dataDir + "*.xlsx")
	if err != nil || len(files) == 0 {
		// Fallback to default if no file found
		return &DataPaths{
			ExcelFile: dataDir + "tabula-calculator.xlsx",
		}
	}
	return &DataPaths{
		ExcelFile: files[0],
	}
}

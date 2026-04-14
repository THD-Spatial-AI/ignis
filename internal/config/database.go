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
	return &DBConfig{
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     GetEnv("DB_PORT", "5432"),
		Name:     GetEnv("DB_NAME", "hdcp"),
		User:     GetEnv("DB_USER", "postgres"),
		Password: GetEnv("DB_PASSWORD", ""),
		SSLMode:  GetEnv("DB_SSL_MODE", "disable"),
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

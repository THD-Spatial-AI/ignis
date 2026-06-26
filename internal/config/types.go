package config

// Schemas holds schema configurations
type Schemas struct {
	Tabula string
}

// DBConfig holds database connection settings
type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string

	// Database structure
	Schemas *Schemas
}

// AppConfig holds HTTP server settings
type AppConfig struct {
	Port string // TCP port the server listens on, e.g. "8080"
}

// DataPaths holds file system paths for input data
type DataPaths struct {
	ExcelFile string
}

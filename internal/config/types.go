package config

// Schemas holds schema configurations
type Schemas struct {
	Tabula string
}

// Database configuration
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

// Data paths
type DataPaths struct {
	ExcelFile string
}

package config

// Config holds the application configuration
type Config struct {
	App  *AppConfig // HTTP server settings
	DB   *DBConfig  // Database connection settings
	Data *DataPaths // Data file paths
}

// LoadConfig loads configuration from environment
func LoadConfig() Config {
	LoadEnv()

	return Config{
		App:  loadAppConfig(),
		DB:   loadDBConfig(),
		Data: loadDataPaths(),
	}
}

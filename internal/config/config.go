package config

// Config holds the application configuration
type Config struct {
	DB   *DBConfig  // Database connection settings
	Data *DataPaths // Data file paths
}

// LoadConfig loads configuration from environment
func LoadConfig() Config {
	LoadEnv()

	return Config{
		DB:   loadDBConfig(),
		Data: loadDataPaths(),
	}
}

package config

import "os"

type Config struct {
	Port               string
	DBConnectionString string
}

func Load() (*Config, error) {
	return &Config{
		Port:               getEnv("PORT", "8080"),
		DBConnectionString: getEnv("DB_CONNECTION_STRING", ""),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

package config

import "os"

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	WhisperURL  string
	OllamaURL   string
	OllamaModel string
}

func Load() *Config {
	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "recruitment"),
		DBPassword:  getEnv("DB_PASSWORD", "recruitment"),
		DBName:      getEnv("DB_NAME", "recruitment"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-dev-key-change-in-prod"),
		WhisperURL:  getEnv("WHISPER_URL", "http://localhost:9000"),
		OllamaURL:   getEnv("OLLAMA_URL", "http://localhost:11434"),
		OllamaModel: getEnv("OLLAMA_MODEL", "qwen3:4b"),
	}
}

func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=disable"
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

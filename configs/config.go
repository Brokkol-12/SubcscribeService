package configs

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Db  DbConfig
	App AppConfig
}

type AppConfig struct {
	Port        string
	LogLevel    string
	ReadTimeout time.Duration
}

type DbConfig struct {
	Dsn string
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	cfg := &Config{
		App: AppConfig{
			Port:        getEnv("PORT", "8081"),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
			ReadTimeout: 10 * time.Second,
		},
		Db: DbConfig{
			Dsn: getEnv("DSN", ""),
		},
	}
	if cfg.Db.Dsn == "" {
		log.Fatal("DNS is required")
	}
	return cfg

}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

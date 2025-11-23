package config

import (
	"fmt"
	"log/slog"
	"os"
)

// const (
// 	ENV_PATH = "../.env"
// )

// func init() {
// 	if err := godotenv.Load(ENV_PATH); err != nil {
// 		panic(".env file is not found")
// 	}
// }

type DBConfig struct {
	PORT   string
	USER   string
	PASS   string
	HOST   string
	DBName string
}

type Config struct {
	DBConfig DBConfig
	port     string
	logLevel slog.Leveler
}

func (c *Config) LogLevel() slog.Leveler {
	return c.logLevel
}

func (c *Config) Port() string {
	return c.port
}

func MustLoad() *Config {

	return &Config{
		port:     getEnv("SERVE_PORT"),
		logLevel: slog.LevelDebug,
		DBConfig: DBConfig{
			PORT:   getEnv("POSTGRES_PORT"),
			USER:   getEnv("POSTGRES_USER"),
			PASS:   getEnv("POSTGRES_PASS"),
			HOST:   getEnv("POSTGRES_HOST"),
			DBName: getEnv("POSTGRES_DB_NAME"),
		},
	}
}

func getEnv(key string) string {
	if value, isExist := os.LookupEnv(key); isExist {
		return value
	}

	panic(fmt.Sprintf("key: %s is not found", key))
}

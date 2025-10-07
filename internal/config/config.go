package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env         string
	ServerPort  string
	DatabaseURL string
	SecretKey   string
}

// Load загружает конфигурацию из .env и переменных окружения
func Load(envFile string) (*Config, error) {
	_ = godotenv.Load(envFile) // если нет .env — не страшно, используем системные env

	cfg := &Config{}

	// Env
	cfg.Env = getEnv("ENV", "dev")

	// Server Port
	portStr := getEnv("SERVER_PORT", "8080") // дефолт

	if _, err := strconv.Atoi(portStr); err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	cfg.ServerPort = portStr

	// Database
	host := getEnv("DB_HOST", "algoplatform-db")
	portDB := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "user")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "algoplatform_db")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Проверка обязательных
	required := map[string]string{
		"DB_HOST":     host,
		"DB_USER":     user,
		"DB_PASSWORD": password,
		"DB_NAME":     dbname,
	}
	for key, val := range required {
		if val == "" {
			return nil, fmt.Errorf("missing required environment variable: %s", key)
		}
	}

	cfg.DatabaseURL = fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, portDB, dbname, sslmode,
	)

	//Secret Key
	cfg.SecretKey = getEnv("SECRET_KEY", "my_secret_key_algo")

	return cfg, nil
}

// getEnv возвращает переменную окружения или значение по умолчанию
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ConfigModel struct {
	ListenAddr string
	BackendURL string
	SessionTTL int
	Difficulty int

	RedisHost string
	RedisPort int

	RateLimitEnabled bool
	RateLimit        int
}

var Config ConfigModel

func getEnvString(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvInt(key string, fallback int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid value for %s, using default: %v\n", key, fallback)
		return fallback
	}
	return value
}

func getEnvBool(key string, fallback bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}
	value := valueStr == "true"
	return value
}

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env probably not found, using env")
	}

	Config = ConfigModel{
		ListenAddr: getEnvString("LISTEN_ADDR", ":8080"),
		BackendURL: getEnvString("BACKEND_URL", "http://localhost:3000"),

		SessionTTL: getEnvInt("SESSION_TTL", 3600),
		Difficulty: getEnvInt("DIFFICULTY", 16),

		RedisHost: getEnvString("REDIS_HOST", "127.0.0.1"),
		RedisPort: getEnvInt("REDIS_PORT", 6379),

		RateLimitEnabled: getEnvBool("RATE_LIMIT_ENABLED", false),
		RateLimit:        getEnvInt("RATE_LIMIT", 10),
	}
}

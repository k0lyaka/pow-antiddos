package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ConfigModel struct {
	ListenAddr string `env:"LISTEN_ADDR"`
	BackendURL string `env:"BACKEND_URL"`
	SessionTTL int    `env:"SESSION_TTL"`

	RedisHost string `env:"REDIS_HOST"`
	RedisPort int    `env:"REDIS_PORT"`
}

var Config ConfigModel

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redis_host := os.Getenv("REDIS_HOST")
	if len(redis_host) == 0 {
		redis_host = "127.0.0.1"
	}

	redis_port := os.Getenv("REDIS_PORT")
	if len(redis_port) == 0 {
		redis_port = "6379"
	}
	redis_port_parsed, err := strconv.Atoi(redis_port)

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sessionTTL, _ := strconv.Atoi(os.Getenv("SESSION_TTL"))
	Config = ConfigModel{
		ListenAddr: os.Getenv("LISTEN_ADDR"),
		BackendURL: os.Getenv("BACKEND_URL"),
		SessionTTL: sessionTTL,
		RedisHost:  redis_host,
		RedisPort:  redis_port_parsed,
	}
}

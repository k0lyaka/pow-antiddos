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
}

var Config ConfigModel

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sessionTTL, _ := strconv.Atoi(os.Getenv("SESSION_TTL"))
	Config = ConfigModel{
		ListenAddr: os.Getenv("LISTEN_ADDR"),
		BackendURL: os.Getenv("BACKEND_URL"),
		SessionTTL: sessionTTL,
	}
}

package config

import (
	"errors"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	SSLMode    string
	JwtSecret  string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
		return nil, err
	}

	cfg := &Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBName:     os.Getenv("DB_NAME"),
		SSLMode:    os.Getenv("DB_SSLMODE"),
		JwtSecret:  os.Getenv("JWT_SECRET"),
	}

	if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBName == "" || cfg.SSLMode == "" || cfg.JwtSecret == "" {
		log.Println("Error in the configuration data and check the config!")
		return nil, errors.New("Error in the configuration data and check the config")
	}

	return cfg, nil
}

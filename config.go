package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser  string
	DBPass  string
	DBHost  string
	DBPort  string
	DBName  string
	DBTable string
	Poll    int
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	c := &Config{
		DBUser:  os.Getenv("DB_USER"),
		DBPass:  os.Getenv("DB_PASS"),
		DBHost:  os.Getenv("DB_HOST"),
		DBPort:  os.Getenv("DB_PORT"),
		DBName:  os.Getenv("DB_NAME"),
		DBTable: os.Getenv("DB_TABLE"),
		Poll:    30,
	}

	if c.DBUser == "" || c.DBHost == "" || c.DBName == "" || c.DBTable == "" {
		return nil, fmt.Errorf("missing required DB config in environment")
	}

	if p := os.Getenv("POLL_SECONDS"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			c.Poll = v
		}
	}

	if c.DBPort == "" {
		c.DBPort = "3306"
	}

	return c, nil
}

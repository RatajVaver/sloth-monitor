package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser       string
	DBPass       string
	DBHost       string
	DBPort       string
	DBName       string
	DBTable      string
	Poll         int
	AvgRatio     float64
	QueryRetries int
	QueryTimeout int
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	c := &Config{
		DBUser:       os.Getenv("DB_USER"),
		DBPass:       os.Getenv("DB_PASS"),
		DBHost:       os.Getenv("DB_HOST"),
		DBPort:       os.Getenv("DB_PORT"),
		DBName:       os.Getenv("DB_NAME"),
		DBTable:      os.Getenv("DB_TABLE"),
		Poll:         30,
		QueryRetries: 1,
		QueryTimeout: 3,
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

	if r := os.Getenv("QUERY_RETRIES"); r != "" {
		if v, err := strconv.Atoi(r); err == nil && v > 0 {
			c.QueryRetries = v
		}
	}

	if t := os.Getenv("QUERY_TIMEOUT"); t != "" {
		if v, err := strconv.Atoi(t); err == nil && v > 0 {
			c.QueryTimeout = v
		}
	}

	if ar := os.Getenv("AVG_RATIO"); ar != "" {
		if v, err := strconv.ParseFloat(ar, 64); err == nil {
			if v >= 0 && v <= 1 {
				c.AvgRatio = v
			}
		}
	}

	return c, nil
}

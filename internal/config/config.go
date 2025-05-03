package config

import (
	"os"
	"strconv"
)

type Config struct {
	RateLimitIP        int
	BlockDurationIP    int
	RateLimitToken     int
	BlockDurationToken int
	RedisAddr          string
	RedisPassword      string
	RedisDB            int
}

func Load() *Config {
	return &Config{
		RateLimitIP:        getEnvAsInt("RATE_LIMIT_IP", 5),
		BlockDurationIP:    getEnvAsInt("BLOCK_DURATION_IP", 300),
		RateLimitToken:     getEnvAsInt("RATE_LIMIT_TOKEN", 10),
		BlockDurationToken: getEnvAsInt("BLOCK_DURATION_TOKEN", 300),
		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		RedisDB:            getEnvAsInt("REDIS_DB", 0),
	}
}

func getEnv(key string, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	if valStr := os.Getenv(name); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}

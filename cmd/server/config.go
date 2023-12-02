package main

import (
	"log"
	"os"
	"strconv"
)

type RedisConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDb       string
}

type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}
type Config struct {
	ServicePort     string
	StoreBackend    string
	DebugMode       bool
	BaseUrl         string
	ProjectTitle    string
	ProjectSubtitle string
	Redis           RedisConfig
	TLS             TLSConfig
}

func NewConfig() Config {
	return Config{
		ServicePort:     getEnv("SERVICE_PORT", ":8080"),
		BaseUrl:         getEnv("BASE_URL", "localhost:8080"),
		StoreBackend:    getEnv("STORE_BACKEND", "in-memory"),
		ProjectTitle:    getEnv("TITLE", "Secure Share"),
		ProjectSubtitle: getEnv("SUBTITLE", "Share short-lived secret that can be accessed only once."),
		DebugMode:       getEnvAsBool("DEBUG_MODE", false),
		Redis: RedisConfig{
			RedisAddr:     getEnv("REDIS_ADDR", "redis:6379"),
			RedisDb:       getEnv("REDIS_DB", "0"),
			RedisPassword: getEnv("REDIS_PASSWORD", ""),
		},
		TLS: TLSConfig{
			Enabled:  getEnvAsBool("TLS_ENABLED", false),
			CertFile: getEnv("CERT_FILE", ""),
			KeyFile:  getEnv("KEY_FILE", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Printf("Error parsing %s as bool: %s. Using default value: %v", key, err, defaultValue)
		return defaultValue
	}

	return value
}

package main

import (
	"log/slog"
	"os"
	"strconv"
)

type RedisConfig struct {
	Address  string
	Password string
	Db       string
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
			Address:  getEnv("REDIS_ADDR", "redis:6379"),
			Db:       getEnv("REDIS_DB", "0"),
			Password: getEnv("REDIS_PASSWORD", ""),
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		logger.Error("error parsing environment variable as bool. Using default value.", "environment_variable", valueStr, "default_value", defaultValue, "error", err)
		return defaultValue
	}

	return value
}

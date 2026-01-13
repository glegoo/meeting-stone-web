package httpserver

import (
	"os"
	"strconv"
)

type Config struct {
	Host      string
	Port      int
	AppEnv    string
	WebOrigin string
}

func LoadConfigFromEnv() Config {
	port := 8080
	if v := os.Getenv("PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			port = p
		}
	}

	return Config{
		Host:      envOrDefault("HOST", "0.0.0.0"),
		Port:      port,
		AppEnv:    envOrDefault("APP_ENV", "dev"),
		WebOrigin: os.Getenv("WEB_ORIGIN"),
	}
}

func (c Config) Addr() string {
	return c.Host + ":" + strconv.Itoa(c.Port)
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

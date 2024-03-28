package config

import (
	"os"

	"github.com/spf13/cast"
)

// Config ...
type Config struct {
	Environment string // develop, staging, production

	UserServiceHost string
	UserServicePort int

	ProductServiceHost string
	ProductServicePort int

	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string

	// context timeout in seconds
	CtxTimeout int
	RedisHost  string
	RedisPort  int

	LogLevel           string
	HTTPPort           string
	CasbinConfigPath   string
	SigningKey         string
	AccessTokenTimeout int
	AuthCSVPath        string
}

// Load loads environment vars and inflates Config
func Load() Config {
	c := Config{}

	c.Environment = cast.ToString(getOrReturnDefault("ENVIRONMENT", "develop"))
	c.PostgresHost = cast.ToString(getOrReturnDefault("POSTGRES_HOST", "db"))
	c.PostgresPort = cast.ToInt(getOrReturnDefault("POSTGRES_PORT", 5432))
	c.PostgresDatabase = cast.ToString(getOrReturnDefault("POSTGRES_DB", "db"))
	c.PostgresUser = cast.ToString(getOrReturnDefault("POSTGRES_USER", "postgres"))
	c.PostgresPassword = cast.ToString(getOrReturnDefault("POSTGRES_PASSWORD", "123"))
	c.LogLevel = cast.ToString(getOrReturnDefault("LOG_LEVEL", "debug"))
	c.HTTPPort = cast.ToString(getOrReturnDefault("HTTP_PORT", ":8080"))
	c.UserServiceHost = cast.ToString(getOrReturnDefault("USER_SERVICE_HOST", "user_service"))
	c.UserServicePort = cast.ToInt(getOrReturnDefault("USER_SERVICE_PORT", 5000))
	c.ProductServiceHost = cast.ToString(getOrReturnDefault("PRODUCT_SERVICE_HOST", "product_service"))
	c.ProductServicePort = cast.ToInt(getOrReturnDefault("PRODUCT_SERVICE_PORT", 6000))
	c.RedisHost = cast.ToString(getOrReturnDefault("REDIS_HOST", "redis"))
	c.RedisPort = cast.ToInt(getOrReturnDefault("REDIS_PORT", 6379))
	c.CasbinConfigPath = cast.ToString(getOrReturnDefault("CASBIN_CONFIG_PATH", "./config/rbac_model.conf"))
	c.AuthCSVPath = cast.ToString(getOrReturnDefault("AUTH_CSV_PATH", "./config/auth.csv"))
	c.SigningKey = cast.ToString(getOrReturnDefault("SIGNING_KEY", "nodirbek"))
	c.AccessTokenTimeout = cast.ToInt(getOrReturnDefault("ACCESS_TOKEN_TIMEOUT", 3600))

	c.CtxTimeout = cast.ToInt(getOrReturnDefault("CTX_TIMEOUT", 7))

	return c
}

func getOrReturnDefault(key string, defaultValue interface{}) interface{} {
	_, exists := os.LookupEnv(key)
	if exists {
		return os.Getenv(key)
	}

	return defaultValue
}

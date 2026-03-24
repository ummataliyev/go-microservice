package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Environment  string `mapstructure:"environment"`
	Debug        bool   `mapstructure:"debug"`
	TrustedHosts string `mapstructure:"trusted_hosts"`
	AppName      string `mapstructure:"app_name"`
	AppVersion   string `mapstructure:"app_version"`
	APIPrefix    string `mapstructure:"api_prefix"`
	DBProvider   string `mapstructure:"db_provider"`
}

type JWTConfig struct {
	SecretKey          string        `mapstructure:"secret_key"`
	Algorithm          string        `mapstructure:"algorithm"`
	AccessTokenExpiry  time.Duration `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry"`
}

type AuthConfig struct {
	MaxAttempts    int `mapstructure:"max_attempts"`
	WindowSeconds  int `mapstructure:"window_seconds"`
	LockoutSeconds int `mapstructure:"lockout_seconds"`
}

type PostgresConfig struct {
	DSN             string        `mapstructure:"dsn"`
	PoolSize        int           `mapstructure:"pool_size"`
	PoolMaxIdle     int           `mapstructure:"pool_max_idle"`
	PoolMaxLifetime time.Duration `mapstructure:"pool_max_lifetime"`
}

type MySQLConfig struct {
	DSN             string        `mapstructure:"dsn"`
	PoolSize        int           `mapstructure:"pool_size"`
	PoolMaxIdle     int           `mapstructure:"pool_max_idle"`
	PoolMaxLifetime time.Duration `mapstructure:"pool_max_lifetime"`
}

type MongoConfig struct {
	URI      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db"`
	Timeout  time.Duration `mapstructure:"timeout"`
}

type RateLimitConfig struct {
	Enabled  bool          `mapstructure:"enabled"`
	LimitGet int           `mapstructure:"limit_get"`
	TimeGet  time.Duration `mapstructure:"time_get"`
	LimitPPD int           `mapstructure:"limit_ppd"`
	TimePPD  time.Duration `mapstructure:"time_ppd"`
}

type CORSConfig struct {
	AllowedOrigins string `mapstructure:"allowed_origins"`
}

type PaginationConfig struct {
	MaxPerPage int `mapstructure:"max_per_page"`
}

type LoggingConfig struct {
	Level                   string `mapstructure:"level"`
	SlowRequestThresholdMS  int    `mapstructure:"slow_request_threshold_ms"`
}

type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Postgres   PostgresConfig   `mapstructure:"postgres"`
	MySQL      MySQLConfig      `mapstructure:"mysql"`
	Mongo      MongoConfig      `mapstructure:"mongo"`
	Redis      RedisConfig      `mapstructure:"redis"`
	RateLimit  RateLimitConfig  `mapstructure:"rate_limit"`
	CORS       CORSConfig       `mapstructure:"cors"`
	Pagination PaginationConfig `mapstructure:"pagination"`
	Logging    LoggingConfig    `mapstructure:"logging"`
}

func Load() (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.environment", "development")
	v.SetDefault("server.debug", false)
	v.SetDefault("server.trusted_hosts", "*")
	v.SetDefault("server.app_name", "go-microservice")
	v.SetDefault("server.app_version", "0.1.0")
	v.SetDefault("server.api_prefix", "/api/v1")
	v.SetDefault("server.db_provider", "postgres")

	v.SetDefault("jwt.algorithm", "HS256")
	v.SetDefault("jwt.access_token_expiry", 15*time.Minute)
	v.SetDefault("jwt.refresh_token_expiry", 7*24*time.Hour)

	v.SetDefault("auth.max_attempts", 5)
	v.SetDefault("auth.window_seconds", 900)
	v.SetDefault("auth.lockout_seconds", 1800)

	v.SetDefault("postgres.pool_size", 10)
	v.SetDefault("postgres.pool_max_idle", 5)
	v.SetDefault("postgres.pool_max_lifetime", 30*time.Minute)

	v.SetDefault("mysql.pool_size", 10)
	v.SetDefault("mysql.pool_max_idle", 5)
	v.SetDefault("mysql.pool_max_lifetime", 30*time.Minute)

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.db", 0)
	v.SetDefault("redis.timeout", 5*time.Second)

	v.SetDefault("rate_limit.enabled", true)
	v.SetDefault("rate_limit.limit_get", 60)
	v.SetDefault("rate_limit.time_get", time.Minute)
	v.SetDefault("rate_limit.limit_ppd", 30)
	v.SetDefault("rate_limit.time_ppd", time.Minute)

	v.SetDefault("cors.allowed_origins", "*")

	v.SetDefault("pagination.max_per_page", 100)

	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.slow_request_threshold_ms", 500)

	// Read .env file
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	_ = v.ReadInConfig() // ignore error if .env doesn't exist

	// Environment variables override
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	env := strings.ToLower(c.Server.Environment)

	// Reject wildcard trusted_hosts in production
	if env == "production" && strings.TrimSpace(c.Server.TrustedHosts) == "*" {
		return fmt.Errorf("wildcard trusted_hosts is not allowed in production")
	}

	// Require JWT secret in non-development environments
	if env != "development" && c.JWT.SecretKey == "" {
		return fmt.Errorf("jwt secret_key is required in %s environment", env)
	}

	// Validate db_provider
	provider := strings.ToLower(c.Server.DBProvider)
	if provider != "postgres" && provider != "mysql" && provider != "mongo" {
		return fmt.Errorf("unsupported db_provider: %s (must be postgres, mysql, or mongo)", c.Server.DBProvider)
	}

	return nil
}

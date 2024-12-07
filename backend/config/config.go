package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// Config -.
type Config struct {
	BindAddr         string `toml:"bind_addr"`
	LogLevel         string `toml:"log_level"`
	MongoUrl         string `toml:"mongo_url"`
	MongoDBName      string `toml:"mongo_db_name"`
	SessionKey       string `toml:"session_key"`
	WriteTimeout     int    `toml:"write_timeout"`
	ReadTimeout      int    `toml:"read_timeout"`
	IdleTimeout      int    `toml:"idle_timeout"`
	GracefulTimeout  int    `toml:"graceful_timeout"`
	JWTSigningKey    string `toml:"jwt_signing_key"`
	JWTSessionLength int    `toml:"jwt_session_length"`
	JWTCookieDomain  string `toml:"jwt_cookie_domain"`
	JWTSecureCookie  bool   `toml:"jwt_secure_cookie"`
	AllowedOrigin    string `toml:"allowed_origin"`
}

// NewConfig returns app config.
func NewConfig() (Config, error) {
	cfg := Config{}

	_, err := toml.DecodeFile("./config/config.toml", &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

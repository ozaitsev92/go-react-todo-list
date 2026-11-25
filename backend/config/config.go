// Package config loads configuration from defaults, an optional TOML file
// (CONFIG_PATH) and environment variables, in increasing precedence.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

const (
	// DefaultPath is used when CONFIG_PATH is not set.
	DefaultPath = "./config/config.toml"

	// MinJWTKeyLength is the minimum accepted HMAC key length in bytes.
	MinJWTKeyLength = 32

	// StorageMongo selects the MongoDB repositories.
	StorageMongo = "mongo"
	// StorageMemory selects the in-memory repositories.
	StorageMemory = "memory"
)

// Config is the application configuration. Integer durations are seconds,
// except JWTSessionLength (minutes).
type Config struct {
	BindAddr          string `toml:"bind_addr"`
	LogLevel          string `toml:"log_level"`
	Storage           string `toml:"storage"`
	MongoURL          string `toml:"mongo_url"`
	MongoDBName       string `toml:"mongo_db_name"`
	MongoTimeout      int    `toml:"mongo_timeout"`
	ReadHeaderTimeout int    `toml:"read_header_timeout"`
	ReadTimeout       int    `toml:"read_timeout"`
	WriteTimeout      int    `toml:"write_timeout"`
	IdleTimeout       int    `toml:"idle_timeout"`
	GracefulTimeout   int    `toml:"graceful_timeout"`
	MaxBodyBytes      int64  `toml:"max_body_bytes"`
	PasswordHashCost  int    `toml:"password_hash_cost"`
	JWTSigningKey     string `toml:"jwt_signing_key"`
	JWTIssuer         string `toml:"jwt_issuer"`
	JWTSessionLength  int    `toml:"jwt_session_length"`
	JWTCookieDomain   string `toml:"jwt_cookie_domain"`
	JWTSecureCookie   bool   `toml:"jwt_secure_cookie"`
}

// Default returns the defaults. It has no JWT key: one must be supplied.
func Default() Config {
	return Config{
		BindAddr:          "8080",
		LogLevel:          "info",
		Storage:           StorageMongo,
		MongoURL:          "mongodb://localhost:27017",
		MongoDBName:       "todo",
		MongoTimeout:      10,
		ReadHeaderTimeout: 5,
		ReadTimeout:       15,
		WriteTimeout:      15,
		IdleTimeout:       60,
		GracefulTimeout:   15,
		MaxBodyBytes:      64 << 10,
		PasswordHashCost:  12,
		JWTIssuer:         "todo-api",
		JWTSessionLength:  30,
		JWTCookieDomain:   "",
		JWTSecureCookie:   true,
	}
}

// Load reads, merges and validates the configuration.
func Load() (Config, error) {
	return load(os.LookupEnv)
}

func load(lookup func(string) (string, bool)) (Config, error) {
	cfg := Default()

	path := DefaultPath
	if p, ok := lookup("CONFIG_PATH"); ok && p != "" {
		path = p
	}

	if _, err := toml.DecodeFile(path, &cfg); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return Config{}, fmt.Errorf("config: read %s: %w", path, err)
	}

	if err := applyEnv(&cfg, lookup); err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func applyEnv(cfg *Config, lookup func(string) (string, bool)) error {
	str := map[string]*string{
		"BIND_ADDR":         &cfg.BindAddr,
		"LOG_LEVEL":         &cfg.LogLevel,
		"STORAGE":           &cfg.Storage,
		"MONGO_URL":         &cfg.MongoURL,
		"MONGO_DB_NAME":     &cfg.MongoDBName,
		"JWT_SIGNING_KEY":   &cfg.JWTSigningKey,
		"JWT_ISSUER":        &cfg.JWTIssuer,
		"JWT_COOKIE_DOMAIN": &cfg.JWTCookieDomain,
	}
	for name, dst := range str {
		if v, ok := lookup(name); ok {
			*dst = v
		}
	}

	ints := map[string]*int{
		"MONGO_TIMEOUT":       &cfg.MongoTimeout,
		"READ_HEADER_TIMEOUT": &cfg.ReadHeaderTimeout,
		"READ_TIMEOUT":        &cfg.ReadTimeout,
		"WRITE_TIMEOUT":       &cfg.WriteTimeout,
		"IDLE_TIMEOUT":        &cfg.IdleTimeout,
		"GRACEFUL_TIMEOUT":    &cfg.GracefulTimeout,
		"PASSWORD_HASH_COST":  &cfg.PasswordHashCost,
		"JWT_SESSION_LENGTH":  &cfg.JWTSessionLength,
	}
	for name, dst := range ints {
		if v, ok := lookup(name); ok {
			n, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("config: %s must be an integer: %w", name, err)
			}
			*dst = n
		}
	}

	if v, ok := lookup("MAX_BODY_BYTES"); ok {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return fmt.Errorf("config: MAX_BODY_BYTES must be an integer: %w", err)
		}
		cfg.MaxBodyBytes = n
	}

	if v, ok := lookup("JWT_SECURE_COOKIE"); ok {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return fmt.Errorf("config: JWT_SECURE_COOKIE must be a boolean: %w", err)
		}
		cfg.JWTSecureCookie = b
	}

	return nil
}

// Validate reports all invalid settings at once.
func (c Config) Validate() error {
	var errs []error

	if c.BindAddr == "" {
		errs = append(errs, errors.New("bind_addr is required"))
	}
	switch strings.ToLower(c.LogLevel) {
	case "debug", "info", "warn", "error":
	default:
		errs = append(errs, fmt.Errorf("log_level %q must be one of debug, info, warn, error", c.LogLevel))
	}
	switch c.Storage {
	case StorageMongo:
		if c.MongoURL == "" || c.MongoDBName == "" {
			errs = append(errs, errors.New("mongo_url and mongo_db_name are required when storage is mongo"))
		}
	case StorageMemory:
	default:
		errs = append(errs, fmt.Errorf("storage %q must be %q or %q", c.Storage, StorageMongo, StorageMemory))
	}
	if len(c.JWTSigningKey) < MinJWTKeyLength {
		errs = append(errs, fmt.Errorf("jwt_signing_key (JWT_SIGNING_KEY) must be at least %d bytes", MinJWTKeyLength))
	}
	if c.JWTSessionLength <= 0 {
		errs = append(errs, errors.New("jwt_session_length must be positive"))
	}
	if c.PasswordHashCost < 10 || c.PasswordHashCost > 16 {
		errs = append(errs, errors.New("password_hash_cost must be between 10 and 16"))
	}
	if c.MaxBodyBytes <= 0 {
		errs = append(errs, errors.New("max_body_bytes must be positive"))
	}
	for name, v := range map[string]int{
		"mongo_timeout": c.MongoTimeout, "read_header_timeout": c.ReadHeaderTimeout,
		"read_timeout": c.ReadTimeout, "write_timeout": c.WriteTimeout,
		"idle_timeout": c.IdleTimeout, "graceful_timeout": c.GracefulTimeout,
	} {
		if v <= 0 {
			errs = append(errs, fmt.Errorf("%s must be positive", name))
		}
	}

	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("config: invalid configuration: %w", err)
	}

	return nil
}

// Seconds converts an integer number of seconds into a time.Duration.
func Seconds(n int) time.Duration { return time.Duration(n) * time.Second }

// JWTSession returns the session length as a duration.
func (c Config) JWTSession() time.Duration { return time.Duration(c.JWTSessionLength) * time.Minute }

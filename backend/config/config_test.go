package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const goodKey = "0123456789abcdef0123456789abcdef"

func envOf(m map[string]string) func(string) (string, bool) {
	return func(k string) (string, bool) {
		v, ok := m[k]
		return v, ok
	}
}

func writeFile(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoad_DefaultsWithoutFile(t *testing.T) {
	t.Parallel()
	cfg, err := load(envOf(map[string]string{
		"CONFIG_PATH":     filepath.Join(t.TempDir(), "missing.toml"),
		"JWT_SIGNING_KEY": goodKey,
	}))
	if err != nil {
		t.Fatal(err)
	}
	want := Default()
	want.JWTSigningKey = goodKey
	if cfg != want {
		t.Fatalf("cfg = %+v\nwant %+v", cfg, want)
	}
}

func TestLoad_EnvOverridesFile(t *testing.T) {
	t.Parallel()
	path := writeFile(t, `
bind_addr = "9000"
log_level = "debug"
mongo_url = "mongodb://file:27017"
jwt_signing_key = "`+goodKey+`"
jwt_secure_cookie = true
`)
	cfg, err := load(envOf(map[string]string{
		"CONFIG_PATH":        path,
		"MONGO_URL":          "mongodb://env:27017",
		"JWT_SECURE_COOKIE":  "false",
		"JWT_SESSION_LENGTH": "15",
		"MAX_BODY_BYTES":     "1024",
	}))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.BindAddr != "9000" || cfg.LogLevel != "debug" {
		t.Errorf("file values not applied: %+v", cfg)
	}
	if cfg.MongoURL != "mongodb://env:27017" || cfg.JWTSecureCookie || cfg.JWTSessionLength != 15 || cfg.MaxBodyBytes != 1024 {
		t.Errorf("env overrides not applied: %+v", cfg)
	}
	if cfg.JWTSession() != 15*time.Minute {
		t.Errorf("JWTSession = %v", cfg.JWTSession())
	}
}

func TestLoad_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		env     map[string]string
		file    string
		wantErr string
	}{
		{"missing jwt key", map[string]string{}, "", "jwt_signing_key"},
		{"short jwt key (the old committed one)", map[string]string{"JWT_SIGNING_KEY": "go-todo-app"}, "", "at least 32 bytes"},
		{"bad int env", map[string]string{"JWT_SIGNING_KEY": goodKey, "READ_TIMEOUT": "soon"}, "", "READ_TIMEOUT"},
		{"bad bool env", map[string]string{"JWT_SIGNING_KEY": goodKey, "JWT_SECURE_COOKIE": "maybe"}, "", "JWT_SECURE_COOKIE"},
		{"bad storage", map[string]string{"JWT_SIGNING_KEY": goodKey, "STORAGE": "redis"}, "", "storage"},
		{"bad log level", map[string]string{"JWT_SIGNING_KEY": goodKey, "LOG_LEVEL": "loud"}, "", "log_level"},
		{"weak bcrypt cost", map[string]string{"JWT_SIGNING_KEY": goodKey, "PASSWORD_HASH_COST": "4"}, "", "password_hash_cost"},
		{"zero timeout", map[string]string{"JWT_SIGNING_KEY": goodKey, "IDLE_TIMEOUT": "0"}, "", "idle_timeout"},
		{"malformed toml", map[string]string{"JWT_SIGNING_KEY": goodKey}, "bind_addr = ", "read"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			env := map[string]string{"CONFIG_PATH": filepath.Join(t.TempDir(), "none.toml")}
			if tt.file != "" {
				env["CONFIG_PATH"] = writeFile(t, tt.file)
			}
			for k, v := range tt.env {
				env[k] = v
			}
			_, err := load(envOf(env))
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("err = %v, want it to mention %q", err, tt.wantErr)
			}
		})
	}
}

func TestLoad_MemoryStorageNeedsNoMongo(t *testing.T) {
	t.Parallel()
	_, err := load(envOf(map[string]string{
		"CONFIG_PATH": filepath.Join(t.TempDir(), "none.toml"), "JWT_SIGNING_KEY": goodKey,
		"STORAGE": StorageMemory, "MONGO_URL": "",
	}))
	if err != nil {
		t.Fatal(err)
	}
}

func TestExampleFile_IsValidAndHasNoSecret(t *testing.T) {
	t.Parallel()
	b, err := os.ReadFile("config.example.toml")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "jwt_signing_key =") {
		t.Fatal("config.example.toml must not contain a JWT signing key")
	}
	if _, err := load(envOf(map[string]string{"CONFIG_PATH": "config.example.toml", "JWT_SIGNING_KEY": goodKey})); err != nil {
		t.Fatalf("example config does not load: %v", err)
	}
}

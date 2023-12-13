package apiserver

type Config struct {
	BindAddr         string `toml:"bind_addr"`
	LogLevel         string `toml:"log_level"`
	DatabaseURL      string `toml:"database_url"`
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

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}

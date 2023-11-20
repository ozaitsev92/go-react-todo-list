package apiserver

type Config struct {
	BindAddr        string `toml:"bind_addr"`
	LogLevel        string `toml:"log_level"`
	DatabaseURL     string `toml:"database_url"`
	SessionKey      string `toml:"session_key"`
	WriteTimeout    int    `toml:"write_timeout"`
	ReadTimeout     int    `toml:"read_timeout"`
	IdleTimeout     int    `toml:"idle_timeout"`
	GracefulTimeout int    `toml:"graceful_timeout"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}

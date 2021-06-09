package server

type Config struct {
	Domain        string `toml:"domain"`
	Port          string `toml:"port"`
	LogLevel      string `toml:"log_level"`
	DatabaseURL   string `toml:"database_url"`
	WebCoverPath  string `toml:"web_cover_path"`
	Smtp_host     string `toml:"smtp_host"`
	Smtp_port     string `toml:"smtp_port"`
	Smtp_login    string `toml:"smtp_login"`
	Smtp_password string `toml:"smtp_password"`
	PrimayKey     string `toml:"primary_key"`
}

func NewConfig() *Config {
	return &Config{
		Domain:   "",
		Port:     ":8080",
		LogLevel: "debug",
	}
}

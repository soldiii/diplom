package server

type ServerConfig struct {
	Address     string `toml:"addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
}

func NewServConfig() *ServerConfig {

	return &ServerConfig{
		Address:     ":8080",
		LogLevel:    "debug",
		DatabaseURL: "host=localhost user=postgres password=postgres dbname=supervisor_bd sslmode=disable",
	}

}

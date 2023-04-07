package server

type ServerConfig struct {
	Addr     string `toml:"addr"`
	LogLevel string `toml:"log_level"`
}

func NewServConfig() *ServerConfig {

	return &ServerConfig{
		Addr:     ":8080",
		LogLevel: "debug",
	}

}

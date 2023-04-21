package server

type ServerConfig struct {
	Address  string `toml:"addr"`
	LogLevel string `toml:"log_level"`
}

func NewServConfig() *ServerConfig {

	return &ServerConfig{
		Address:  ":8080",
		LogLevel: "debug",
	}

}

package store

type StoreConfig struct {
	Hostname string `toml:"host"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
	SSLMode  string `toml:"sslmode"`
}

func NewStoreConfig() *StoreConfig {
	return &StoreConfig{
		Hostname: "host",
		User:     "postgres",
		Password: "postgres",
		DBName:   "supervisor_app_bd",
		SSLMode:  "disable",
	}
}

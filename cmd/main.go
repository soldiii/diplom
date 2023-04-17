package main

import (
	"flag"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/soldiii/diplom/internal/server"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "./config/config.toml", "path to config file")
}

func main() {

	logrus.SetFormatter(new(logrus.JSONFormatter))

	flag.Parse()

	servConfig := server.NewServConfig()
	_, err := toml.DecodeFile(configPath, &servConfig)
	if err != nil {
		logrus.Fatalf("Reading or decoding config file error: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading env variables: %s", err.Error())
	}

	serv := server.NewServer(servConfig)

	defer serv.PostgresDB.ClosePostgresDB()

	if err := serv.RunServer(); err != nil {
		logrus.Fatalf("Running server error: %s", err.Error())
	}

}

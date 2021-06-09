package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/whitewolfaster/filin/internal/server"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/server.toml", "path to server config file")
}

func main() {
	flag.Parse()

	config := server.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	server, err := server.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	if err = server.Start(); err != nil {
		log.Fatal(err)
	}
}

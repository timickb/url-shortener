package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/timickb/url-shortener/internal/app/server"
	"gopkg.in/yaml.v3"
)

var (
	configPath string
	storeImpl  string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/server.yaml", "Path to config file")
	flag.StringVar(&storeImpl, "store", "local", "Data storage way")
}

func main() {
	flag.Parse()

	config := server.DefaultConfig()

	configContent, ioErr := ioutil.ReadFile(configPath)

	if ioErr != nil {
		log.Fatal(ioErr)
	}

	if err := yaml.Unmarshal(configContent, &config); err != nil {
		log.Fatal(err)
	}

	config.StoreImpl = storeImpl
	srv := server.NewServer(config)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
